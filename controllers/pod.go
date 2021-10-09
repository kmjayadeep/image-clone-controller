package controllers

import (
	"context"
	"strings"

	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const BACKUP_REGISTRY = "repo.treescale.com"
const BACKUP_REGISTRY_ORG = "kmjayadeep"

func reconcilePod(ctx context.Context, podSpec *corev1.PodSpec) (bool, error) {
	log := log.FromContext(ctx)
	log.Info("reconciling pod")
	changed := false

	log.Info("checking containers")
	for i, c := range podSpec.Containers {
		is, err := isBackupImage(ctx, c.Image)
		if err != nil {
			return false, err
		}

		if !is {
			continue
		}

		image, err := backupImage(ctx, c.Image)
		podSpec.Containers[i].Image = image
		changed = true
	}

	log.Info("checking init containers")
	for i, c := range podSpec.InitContainers {
		is, err := isBackupImage(ctx, c.Image)
		if err != nil {
			return false, err
		}

		if !is {
			continue
		}

		image, err := backupImage(ctx, c.Image)
		podSpec.InitContainers[i].Image = image
		changed = true
	}

	log.Info("checking ephemeral containers")
	for i, c := range podSpec.EphemeralContainers {
		is, err := isBackupImage(ctx, c.Image)
		if err != nil {
			return false, err
		}

		if !is {
			continue
		}

		image, err := backupImage(ctx, c.Image)
		podSpec.EphemeralContainers[i].Image = image
		changed = true
	}

	return changed, nil
}

// backup current image to backup registry and return new image url
func backupImage(ctx context.Context, image string) (string, error) {
	log := log.FromContext(ctx)
	ref, err := name.NewTag(image)
	if err != nil {
		return "", err
	}

	registry := ref.RegistryStr()
	repo := ref.RepositoryStr()
	tag := ref.TagStr()
	log.Info("container using public image", "registry", registry, "repo", repo, "tag", tag)

	backupImage := getBackupImage(repo, tag)

	crane.Tag(image, backupImage)

	img, err := remote.Image(ref)
	if err != nil {
		log.Error(err, "Unable to parse image")
		return "", err
	}

	err = crane.Push(img, backupImage)
	if err != nil {
		log.Error(err, "Unable to push image")
		return "", err
	}
	log.Info("image pushed to backup registry")
	return backupImage, nil
}

// check if the image is already taken from backup registry
func isBackupImage(ctx context.Context, image string) (bool, error) {
	log := log.FromContext(ctx)
	log.Info("checking container image", "name", image)

	ref, err := name.NewTag(image)
	if err != nil {
		log.Error(err, "unable to parse tag")
		return false, err
	}

	registry := ref.RegistryStr()

	return registry != BACKUP_REGISTRY, nil
}

// get backup registry image name from original
func getBackupImage(repo, tag string) string {
	parts := strings.Split(repo, "/")
	return BACKUP_REGISTRY + "/" + BACKUP_REGISTRY_ORG + "/" + parts[len(parts)-1] + ":" + tag
}
