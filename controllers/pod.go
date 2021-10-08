package controllers

import (
	"context"
	"strings"

	"github.com/google/go-containerregistry/pkg/name"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/log"
	// "github.com/google/go-containerregistry/pkg/authn"
	// "github.com/google/go-containerregistry/pkg/v1/remote"
)

const BACKUP_REGISTRY = "repo.treescale.com"
const BACKUP_REGISTRY_ORG = "kmjayadeep"

func reconcilePod(ctx context.Context, podSpec *corev1.PodSpec) (bool, error) {
	log := log.FromContext(ctx)
	log.Info("reconciling pod")
	changed := false

	for i, c := range podSpec.Containers {
		image := c.Image
		log.Info("checking container image", "name", image)

		ref, err := name.NewTag(image)
		if err != nil {
			return false, err
		}

		registry := ref.RegistryStr()

		// container is not using the backup image
		if registry != BACKUP_REGISTRY {
			repo := ref.RepositoryStr()
			tag := ref.TagStr()
			log.Info("container using public image", "registry", registry, "repo", repo, "tag", tag)

			backupImage := getBackupImage(repo, tag)

			podSpec.Containers[i].Image = backupImage
			changed = true
		}

	}

	return changed, nil
}

func getBackupImage(repo, tag string) string {
	parts := strings.Split(repo, "/")
	return BACKUP_REGISTRY + "/" + BACKUP_REGISTRY_ORG + "/" + parts[len(parts)-1] + ":" + tag
}
