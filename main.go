package main

import (
	"os"

  "github.com/kmjayadeep/image-clone-controller/controllers"
	appsv1 "k8s.io/api/apps/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

func init() {
	log.SetLogger(zap.New())
}

func main() {
	entryLog := log.Log.WithName("entrypoint")

	// Setup a Manager
	entryLog.Info("setting up manager")
	mgr, err := manager.New(config.GetConfigOrDie(), manager.Options{})
	if err != nil {
		entryLog.Error(err, "unable to set up overall controller manager")
		os.Exit(1)
	}

	// Setup a new controller to reconcile ReplicaSets
	entryLog.Info("Setting up deployment controller")
	deployCtl, err := controller.New("deployment-controller", mgr, controller.Options{
		Reconciler: controllers.NewDeploymentController(mgr.GetClient()),
	})
	if err != nil {
		entryLog.Error(err, "unable to set up deployment controller")
		os.Exit(1)
	}

	// Watch Deployments and enqueue Deployment object key
	if err := deployCtl.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForObject{}); err != nil {
		entryLog.Error(err, "unable to watch Deployments")
		os.Exit(1)
	}

	entryLog.Info("Setting up daemonset controller")
	dsCtl, err := controller.New("daemonset-controller", mgr, controller.Options{
		Reconciler: controllers.NewDaemonSetController(mgr.GetClient()),
	})
	if err != nil {
		entryLog.Error(err, "unable to set up daemonset controller")
		os.Exit(1)
	}

	// Watch Daemonsets and enqueue Daemonset object key
  if err := dsCtl.Watch(&source.Kind{Type: &appsv1.DaemonSet{}}, &handler.EnqueueRequestForObject{}); err != nil {
    entryLog.Error(err, "unable to watch Daemonsets")
    os.Exit(1)
  }

	entryLog.Info("starting manager")
	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		entryLog.Error(err, "unable to run manager")
		os.Exit(1)
	}
}
