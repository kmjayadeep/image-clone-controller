package controllers

import (
	"sigs.k8s.io/controller-runtime/pkg/log"

	corev1 "k8s.io/api/core/v1"
)

func reconcilePod(podSpec *corev1.PodSpec) (bool, error) {
	log := log.Log.WithName("entrypoint")
	return false, nil
}
