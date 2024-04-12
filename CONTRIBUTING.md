# Contributing to Ollama Operator

### Prerequisites

- `go` version `v1.21.0+`
- `docker` version `17.03+`.
- `kubectl` version `v1.11.3+`.
- Access to a Kubernetes `v1.11.3+` cluster (Either kind, minikube is suitable for local development).

### Develop with `kind`

We have pre-defined `kind` configurations in the `hack` directory to help you get started with the development environment.

```sh
kind create cluster --config hack/kind.yaml
```

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

Users can just run kubectl apply -f <URL for YAML BUNDLE> to install the project, i.e.:

```sh
kubectl apply -f https://raw.githubusercontent.com/<org>/ollama-operator/<tag or branch>/dist/install.yaml
```
