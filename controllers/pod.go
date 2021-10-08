package controllers

import (
	"context"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const BACKUP_REGISTRY = "repo.treescale.com/kmjayadeep"

func reconcilePod(ctx context.Context, podSpec *corev1.PodSpec) (bool, error) {
	log := log.FromContext(ctx)
	log.Info("reconciling pod")

	for i, c := range podSpec.Containers {
		image := c.Image
		log.Info("checking container image", "name", image)

		// check if image is from backup registry
		if strings.HasPrefix(image, BACKUP_REGISTRY) {

		}

		podSpec.Containers[i].Image = image
	}

	return false, nil
}
