# 通过 CRD 部署

## 部署模型

1. 创建一个 [`Model`](/pages/zh-CN/references/crd/model) 类型的 CRD 资源

::: tip 什么是 CRD？

CRD 是 Kubernetes 的自定义资源定义（Custom Resource Definition）的缩写，它允许用户自定义资源类型，从而扩展 Kubernetes API。

名为 Operator 的服务可以管理这些自定义资源，以便在 Kubernetes 集群中部署、管理和监控应用程序。

Ollama Operator 就是通过版本号为 `ollama.ayaka.io/v1`，类型为 `Model` 的 CRD 来管理大型语言模型的部署和运行的。

```yaml
apiVersion: ollama.ayaka.io/v1 # [!code focus]
kind: Model # [!code focus]
metadata:
  name: phi
spec:
  image: phi
```

:::

::: warning 使用了 `kind` 作为集群吗？

`kind` 默认配置的 `StorageClass` 是 `standard`，并且仅适用于 `ReadWriteOnce` 访问模式，因此，如果您需要使用 `kind` 运行这个 Operator 并部署模型，您应该在 `Model` CRD 中使用 `accessMode：ReadWriteOnce` 指定 `persistentVolume`：

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

复制以下命令以创建一个名为 phi 的 [`Model` CRD](/pages/zh-CN/references/crd/model)：

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

或者您可以创建自己的文件：

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

1. 将 [`Model` CRD](/pages/zh-CN/references/crd/model) 应用到 Kubernetes 集群：

```shell
kubectl apply -f ollama-model-phi.yaml
```

3. 等待 [`Model` CRD](/pages/zh-CN/references/crd/model) 就绪：

```shell
kubectl wait --for=jsonpath='{.status.readyReplicas}'=1 deployment/ollama-model-phi
```

4. 准备就绪！现在让我们转发访问模型的端口到本地：

```shell
kubectl port-forward svc/ollama-model-phi ollama
```

5. 直接与模型交互：

```shell
ollama run phi
```

或者使用 `curl` 连接到与 OpenAI API 兼容的接口：

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
