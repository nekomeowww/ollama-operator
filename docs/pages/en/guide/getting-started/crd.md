# Deploy models through CRD

## Deploy the model

1. Create one [`Model` CRD](/pages/en/references/crd/model) to rule them all.

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

Copy the following command to create a phi [`Model` CRD](/pages/en/references/crd/model):

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

2. Apply the [`Model` CRD](/pages/en/references/crd/model) to your Kubernetes cluster:

```shell
kubectl apply -f ollama-model-phi.yaml
```

1. Wait for the [`Model`](/pages/en/references/crd/model) to be ready:

```shell
kubectl wait --for=jsonpath='{.status.readyReplicas}'=1 deployment/ollama-model-phi
```

4. Ready! Now let's forward the ports to access the model:

```shell
kubectl port-forward svc/ollama-model-phi ollama
```

5. Interact with the model:

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
