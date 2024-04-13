---
layout: home
sidebar: false

title: Ollama Operator

hero:
  name: Ollama Operator
  text: Large language models, scaled, deployed
  tagline: Yet another operator for running large language models on Kubernetes with ease. Powered by Ollama! 🐫
  image:
    src: /logo.png
    alt: Ollama Operator
  actions:
    - theme: brand
      text: Documentations
      link: /pages/en/guide/overview
    - theme: alt
      text: View on GitHub
      link: https://github.com/nekomeowww/ollama-operator

features:
  - icon: <div i-twemoji:rocket></div>
    title: Launch and chat
    details: Easy to use API, the spec is simple enough to just a few lines of YAML to deploy a model, and you can chat with it right away.
  - icon: <div i-twemoji:ship></div>
    title: Cross-kubernetes
    details: Extend the user experience of Ollama to any Kubernetes cluster, edge or any cloud infrastructure, with the same spec, and chat with it from anywhere.
  - icon: <div i-simple-icons:openai></div>
    title: OpenAI API compatible
    details: Your familiar <code>/v1/chat/completions</code> endpoint is here, with the same request and response format. No need to change your code or switch to another API.
  - icon: <div i-twemoji:parrot></div>
    title: Langchain ready
    details: Power on to function calling, agents, knowldgebase retrieving. Unleash all the power Langchain has out of box with Ollama Operator.
---

### See it in action

<br>

<AsciinemaPlayer src="/demo.cast" />

## Try it out

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
kubectl wait --for=jsonpath='{.status.readyReplicas}'=1 deployment/ollama-operator-controller-manager -n ollama-operator-system
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
