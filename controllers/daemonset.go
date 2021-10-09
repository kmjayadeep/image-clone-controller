package controllers

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// DeploymentController reconciles Deployments
type DaemonSetController struct {
	// client can be used to retrieve objects from the APIServer.
	client client.Client
}

func NewDaemonSetController(client client.Client) reconcile.Reconciler {
	return &DaemonSetController{
		client: client,
	}
}

func (r *DaemonSetController) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	log := log.FromContext(ctx)

	if request.Namespace == "kube-system" {
		log.Info("ignoring the Deployment in kube-system Namespace", "name", request.Name)
		return reconcile.Result{}, nil
	}

	// Fetch the DaemonSet from the cache
	ds := &appsv1.DaemonSet{}
	err := r.client.Get(ctx, request.NamespacedName, ds)
	if errors.IsNotFound(err) {
		log.Error(nil, "Could not find daemonset")
		return reconcile.Result{}, nil
	}

	if err != nil {
		return reconcile.Result{}, fmt.Errorf("could not fetch daemonset: %+v", err)
	}

	log.Info("Reconciling Daemonset", "name", ds.Name)

	podSpec := &ds.Spec.Template.Spec

	changed, err := reconcilePod(ctx, podSpec)

	if err != nil {
		log.Error(err, "unable to reconcile pod spec")
		return reconcile.Result{}, err
	}

	if changed {
		log.Info("daemonset spec updated", "new spec", podSpec)

		err = r.client.Update(ctx, ds)

		if err != nil {
			log.Error(err, "unable to update daemonset")
			return reconcile.Result{}, err
		}

		return reconcile.Result{}, nil
	}

	return reconcile.Result{}, nil
}
