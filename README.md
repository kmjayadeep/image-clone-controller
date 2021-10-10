# Image clone controller

image clone controller is a kubernetes controller to safe guard against
the risk of container images disappearing from public registry

# Deployment

The yaml files for deploying the controller are located in `deploy`
folder.

1. Edit the `controller.yaml` and make sure the controller image is
   pointing to your private registry
1. Deploy a kubernetes secret with the name `registry-secret` which
   contains credentials to access the container registry
1. Apply the manifests
