# Getting Started

::: tip Don't have an existing Kubernetes cluster?

Run the following commands to create a new Kubernetes cluster with `kind`:

::: code-group

```shell [macOS]
brew install --cask docker
brew install docker kind kubectl
wget https://raw.githubusercontent.com/nekomeowww/ollama-operator/main/hack/kind-config.yaml
kind create cluster --config kind-config.yaml
```

```shell [Windows]
Invoke-WebRequest  -OutFile "./Docker Desktop Installer.exe"
Start-Process 'Docker Desktop Installer.exe' -Wait install
start /w "" "Docker Desktop Installer.exe" install

scoop install docker kubectl go
go install sigs.k8s.io/kind@latest
wget https://raw.githubusercontent.com/nekomeowww/ollama-operator/main/hack/kind-config.yaml
kind create cluster --config kind-config.yaml
```

```shell [Linux]
# refer to Install Docker Engine on Debian | Docker Docs https://docs.docker.com/engine/install/debian/
# and Install and Set Up kubectl on Linux | Kubernetes https://kubernetes.io/docs/tasks/tools/install-kubectl-linux/
```

:::

1. Install operator.

```shell
kubectl apply -f https://raw.githubusercontent.com/nekomeowww/ollama-operator/main/dist/install.yaml
```

2. Wait for the operator to be ready:

```shell
kubectl wait --for=jsonpath='{.status.replicas}'=2 deployment/ollama-operator-controller-manager -n ollama-operator-system
```

3. Create one `Model` CRD to rule them all.

::: tip What is CRD?

CRD stands for Custom Resource Definition for Kubernetes, which extends the Kubernetes API by allowing users to customize resource types.

Cervices called Operator and Controller manages these custom resources to deploy, manage, and monitor applications in a Kubernetes cluster.

Ollama Operator manages the deployment and operation of large language models through a CRD with version number `ollama.ayaka.io/v1` and type `Model`.

```yaml
apiVersion: ollama.ayaka.io/v1 # [!code focus]
kind: Model # [!code focus]
metadata:
  name: phi
spec:
  image: phi
```

:::

::: warning Working with `kind`?

The default provisioned `StorageClass` in `kind` is `standard`, and will only work with `ReadWriteOnce` access mode, therefore if you would need to run the operator with `kind`, you should specify `persistentVolume` with `accessMode: ReadWriteOnce` in the `Model` CRD:

```yaml
apiVersion: ollama.ayaka.io/v1
kind: Model
metadata:
  name: phi
spec:
  image: phi
  persistentVolume: # [!code focus]
    accessMode: ReadWriteOnce # [!code focus]
```

:::

Copy the following command to create a phi model CRD:

```shell
cat <<EOF >> ollama-model-phi.yaml
apiVersion: ollama.ayaka.io/v1
kind: Model
metadata:
  name: phi
spec:
  image: phi
  persistentVolume:
    accessMode: ReadWriteOnce
EOF
```

or you can create your own file:

::: code-group

```yaml [ollama-model-phi.yaml]
apiVersion: ollama.ayaka.io/v1 # [!code ++]
kind: Model # [!code ++]
metadata: # [!code ++]
  name: phi # [!code ++]
spec: # [!code ++]
  image: phi # [!code ++]
```

:::

4. Apply the `Model` CRD to your Kubernetes cluster:

```shell
kubectl apply -f ollama-model-phi.yaml
```

5. Wait for the model to be ready:

```shell
kubectl wait --for=jsonpath='{.status.readyReplicas}'=1 deployment/ollama-model-phi
```

6. Ready! Now let's forward the ports to access the model:

```shell
kubectl port-forward svc/ollama-model-phi ollama
```

7. Interact with the model:

```shell
ollama run phi
```

or use the OpenAI API compatible endpoint:

```shell
curl http://localhost:11434/v1/chat/completions -H "Content-Type: application/json" -d '{
  "model": "phi",
  "messages": [
      {
          "role": "user",
          "content": "Hello!"
      }
  ]
}'
```
