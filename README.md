# Image clone controller

image clone controller is a kubernetes controller to safe guard against
the risk of container images disappearing from public registry

# Deployment

The yaml files for deploying the controller are located in `deploy`
folder.

1. Build the controller and push it into private registry
```
docker build -t repo.treescale.com/kmjayadeep/image-clone-controller:latest .
docker push repo.treescale.com/kmjayadeep/image-clone-controller:latest
```
1. Edit the `controller.yaml` and make sure the controller image is
   pointing to your private registry
1. Deploy a kubernetes secret with the name `registry-secret` which
   contains credentials to access the container registry
1. Apply the manifests
```
kubectl apply -f deploy/
```

# Testing

Deploy a sample deployment in default namespace

```
kubectl create deploy nginx --image nginx
```

Watch the deployment being updated with new image

```
kubectl get po -w
```

verify that the image is updated
```
❯ k describe po -l app=nginx | grep Image
    Image:          repo.treescale.com/kmjayadeep/nginx:latest
    Image ID:       docker.io/library/nginx@sha256:06e4235e95299b1d6d595c5ef4c41a9b12641f6683136c18394b858967cd1506
```

# Demo

watch the short demo on `asciinema` : <https://asciinema.org/a/441210>
