# Contributing to Ollama Operator

### Prerequisites

- `go` version `v1.21.0+`
- `docker` version `17.03+`.
- `kubectl` version `v1.11.3+`.
- Access to a Kubernetes `v1.11.3+` cluster (Either kind, minikube is suitable for local development).

## Development

### Develop with `kind`

We have pre-defined `kind` configurations in the `hack` directory to help you get started with the development environment.

> [!NOTE]
> Install kind if you haven't already. You can install it using the following command:
>
> ```shell
> go install sigs.k8s.io/kind@latest
> ```

To create a `kind` cluster with the configurations defined in `hack/kind-config.yaml`, run the following command:

```sh
kind create cluster --config hack/kind-config.yaml
```

### Test controller locally

#### Build and install CRD to local cluster

```shell
make manifests && make install
```

#### Run the controller locally

```shell
make run
```

#### Create a CRD to test the controller

```shell
kubectl apply -f hack/ollama-model-phi-kind-cluster.yaml
```

#### Verify the resources is reconciling

```shell
kubectl describe models
```

## Deployment

### To Deploy on the cluster

**Build and push your image to the location specified by `IMG`:**

```sh
make docker-build docker-push IMG=<some-registry>/ollama-operator:tag
```

> [!NOTE]
> This image ought to be published in the personal registry you specified.
> And it is required to have access to pull the image from the working environment.
> Make sure you have the proper permission to the registry if the above commands don’t work.

**Install the CRDs into the cluster:**

```sh
make install
```

**Deploy the Manager to the cluster with the image specified by `IMG`:**

```sh
make deploy IMG=<some-registry>/ollama-operator:tag
```

> [!NOTE]
> If you encounter RBAC errors, you may need to grant yourself cluster-admin
> privileges or be logged in as admin.

**Create instances of your solution**
You can apply the samples (examples) from the config/sample:

```sh
kubectl apply -k config/samples/
```

> [!NOTE]
> Ensure that the samples has default values to test it out.

### To Uninstall

**Delete the instances (CRs) from the cluster:**

```sh
kubectl delete -k config/samples/
```

**Delete the APIs(CRDs) from the cluster:**

```sh
make uninstall
```

**UnDeploy the controller from the cluster:**

```sh
make undeploy
```

> [!NOTE]
> Run `make help` for more information on all potential `make` targets

## Project Distribution

### Test your build locally

Developers should test their builds locally before pushing the changes to the repository.

#### Test the build of `kollama` locally

This reporsitory has been configured to build and release by using [GoReleaser](https://goreleaser.com/):

> [!NOTE]
> Install GoReleaser if you haven't already by going through the steps in [GoReleaser Install](https://goreleaser.com/install/)

```shell
goreleaser release --snapshot --clean
```

#### Test the build of `ollama-operator` locally

To build the image for `ollama-operator` locally, run the following command:

```shell
docker buildx build --platform linux/amd64,linux/arm64 .
```

#### Test the release of `kollama` to `krew` locally

Please replace `<Your Tag>` with the tag you want to release.

```shell
docker run -v $(pwd)/.krew.yaml:/tmp/template-file.yaml ghcr.io/rajatjindal/krew-release-bot:v0.0.46 krew-release-bot template --tag <Your Tag> --template-file /tmp/template-file.yaml
```

### Build and distribute the project

Following are the steps to build the installer and distribute this project to users.

1. Build the installer for the image built and published in the registry:

```sh
make build-installer IMG=<some-registry>/ollama-operator:tag
```

NOTE: The makefile target mentioned above generates an 'install.yaml'
file in the dist directory. This file contains all the resources built
with Kustomize, which are necessary to install this project without
its dependencies.

2. Using the installer

Users can just run `kubectl apply -f <URL for YAML BUNDLE>` to install the project, i.e.:

```sh
kubectl apply -f https://raw.githubusercontent.com/<org>/ollama-operator/<tag or branch>/dist/install.yaml
```
