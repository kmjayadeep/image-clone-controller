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
type DeploymentController struct {
	// client can be used to retrieve objects from the APIServer.
	client client.Client
}

func NewDeploymentController(client client.Client) reconcile.Reconciler{
  return &DeploymentController{
    client: client,
  }
}

func (r *DeploymentController) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
  log := log.FromContext(ctx)

  if request.Namespace == "kube-system" {
    log.Info("ignoring the Deployment in kube-system Namespace", "name", request.Name)
    return reconcile.Result{}, nil
  }

	// Fetch the Deployment from the cache
	deploy := &appsv1.Deployment{}
	err := r.client.Get(ctx, request.NamespacedName, deploy)
	if errors.IsNotFound(err) {
		log.Error(nil, "Could not find deployment")
		return reconcile.Result{}, nil
	}

	if err != nil {
		return reconcile.Result{}, fmt.Errorf("could not fetch Deployment: %+v", err)
	}

	log.Info("Reconciling Deployment", "name", deploy.Name)

	return reconcile.Result{}, nil
}
