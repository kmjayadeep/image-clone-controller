package main

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// reconcileReplicaSet reconciles ReplicaSets
type reconcileReplicaSet struct {
	// client can be used to retrieve objects from the APIServer.
	client client.Client
}

func (r *reconcileReplicaSet) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	log := log.FromContext(ctx)

	// Fetch the ReplicaSet from the cache
	deploy := &appsv1.Deployment{}
	err := r.client.Get(ctx, request.NamespacedName, deploy)
	if errors.IsNotFound(err) {
		log.Error(nil, "Could not find deployment")
		return reconcile.Result{}, nil
	}

	if err != nil {
		return reconcile.Result{}, fmt.Errorf("could not fetch Deployment: %+v", err)
	}

	// Print the ReplicaSet
	log.Info("Reconciling Deployment", "name", deploy.Name)

	return reconcile.Result{}, nil
}
