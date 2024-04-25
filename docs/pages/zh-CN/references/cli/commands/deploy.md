# `kollama deploy`

`kollama deploy` 命令用于将 [`Model`](/pages/zh-CN/references/crd/model) 部署到 Kubernetes 集群。它是通过操作 CRD 资源与 Ollama Operator 交互的基本封装和工具类型 CLI。

[[toc]]

## 用例

### 部署存储于 [registry.ollama.ai](https://registry.ollama.ai) 镜像仓库上的模型

```shell
kollama deploy phi
```

### 部署到特定命名空间（namespace）

```shell
kollama deploy phi --namespace=production
```

### 部署存储于自定义镜像仓库上的模型

```shell
kollama deploy phi --image=registry.example.com/library/phi:latest
```

### 部署带有暴露 NodePort 服务的 [`Model`](/pages/zh-CN/references/crd/model) 以供外部访问

```shell
kollama deploy phi --expose
```

::: tip 了解 [`Model`](/pages/zh-CN/references/crd/model) 使用的 NodePort 端口号...

```shell
kubectl get svc --selector model.ollama.ayaka.io/name=<model name> -o json | jq ".spec.ports[0].nodePort"
```

:::

::: tip 配置期望分配的 NodePort 端口号...

```shell
kollama deploy phi --expose --node-port=30000
```

:::

### 部署带有暴露 LoadBalancer 服务的 [`Model`](/pages/zh-CN/references/crd/model) 以供外部访问

```shell
kollama deploy phi --expose --service-type=LoadBalancer
```

### 部署有着资源限制的 [`Model`](/pages/zh-CN/references/crd/model)

下面的示例部署了 `phi` 模型，并限制 CPU 使用率为 `1` 个核心，内存使用量为 `1Gi`。

```shell
kollama deploy phi --limit=cpu=1 --limit=memory=1Gi
```

## 选项

### `--namespace`

如果配置了该参数，将会在指定的命名空间中部署 [`Model`](/pages/zh-CN/references/crd/model)。

### `--image`

默认：`registry.ollama.ai/library/<model name>:latest`

```shell
kollama deploy phi --image=registry.ollama.ai/library/phi:latest
```

要部署的模型镜像。

- 如果未指定，将使用 [`Model`](/pages/zh-CN/references/crd/model) 名称作为镜像名称（如果未指定镜像仓库（Registry），这个时候会默认从 `registry.ollama.ai/library/<model name>` 拉取）。例如，如果 [`Model`](/pages/zh-CN/references/crd/model) 名称是 `phi`，最终获取的镜像名称将是 `registry.ollama.ai/library/phi:latest`。
- 如果没有指定，标签将会使用 `latest` 的。

### `--limit`（支持多次使用）

> 多次使用该选项可指定多个资源限制。

为即将部署的 [`Model`](/pages/zh-CN/references/crd/model) 指定资源限制。这对于没有足够多资源的集群，或者是希望在有限资源的集群中部署多个 [`Model`](/pages/zh-CN/references/crd/model) 是非常有用的。

::: tip 对于 NVIDIA、AMD GPU 的资源限制...

在 Kubernetes 中，任何 GPU 资源都遵循这个格式：

```yaml
resources:
  limits:
    gpu-vendor.example/example-gpu: 1 # requesting 1 GPU
```

使用 `nvidia.com/gpu` 可以限制 NVIDIA GPU 的数量，因此，在使用 `kollama deploy` 时，你可以使用 `--limit nvidia.com/gpu=1` 来指定 NVIDIA GPU 的数量为 `1`：

```shell
kollama deploy phi --limit=nvidia.com/gpu=1
```

```yaml
resources:
  limits:
    nvidia.com/gpu: 1 # requesting 1 GPU # [!code focus]
```

> [有关配合 `nvidia/k8s-device-plugin` 使用资源标签的文档](https://github.com/NVIDIA/k8s-device-plugin?tab=readme-ov-file#enabling-gpu-support-in-kubernetes)

使用 `amd.com/gpu` 可以限制 AMD GPU 的数量，在使用 `kollama deploy` 时，你可以使用 `--limit amd.com/gpu=1` 来指定 AMD GPU 的数量为 `1`。

```shell
kollama deploy phi --limit=amd.com/gpu=1
```

最终会渲染为：

```yaml
resources:
  limits:
    amd.com/gpu: 1 # requesting a GPU  # [!code focus]
```

> [关于配合 `ROCm/k8s-device-plugin` 使用 Label 的 YAML 配置文件的示例](https://github.com/ROCm/k8s-device-plugin/blob/4607bf06b700e53803d566e0bf9555f773f0b4f1/example/pod/alexnet-gpu.yaml)

你可以在这里阅读更多：[调度 GPUs | Kubernetes](https://kubernetes.io/zh-cn/docs/tasks/manage-gpus/scheduling-gpus/)

:::

::: details 我已经部署过 [`Model`](/pages/zh-CN/references/crd/model)，但是我想要更改资源限制...

当然可以，用 [`kubectl set resources`](https://kubernetes.io/zh-cn/docs/reference/kubectl/generated/kubectl_set/kubectl_set_resources/) 命令来可以更改资源限制：

```shell
kubectl set resources deployment -l model.ollama.ayaka.io/name=<model name> --limits cpu=4
```

改内存限制：

```shell
kubectl set resources deployment -l model.ollama.ayaka.io/name=<model name> --limits memory=8Gi
```

:::

格式是 `<resource>=<quantity>`.

比如：`--limit=cpu=1` `--limit=memory=1Gi`.


### `--storage-class`

```shell
kollama deploy phi --storage-class=standard
```

要给 [`Model`](/pages/zh-CN/references/crd/model) 部署服务关联的 [`PersistentVolumeClaim`](https://kubernetes.io/zh-cn/docs/concepts/storage/persistent-volumes/#persistentvolumeclaims) 使用的 [`StorageClass`](https://kubernetes.io/zh-cn/docs/concepts/storage/storage-classes/#storageclass-objects)。

如果没有指定，则会使用 [默认的 `StorageClass`](https://kubernetes.io/zh-cn/docs/concepts/storage/storage-classes/#default-storageclass)。

### `--pv-access-mode`

```shell
kollama deploy phi --pv-access-mode=ReadWriteMany
```

用于给 Ollama Operator 所创建的 image store 的 [`StatefulSet`](https://kubernetes.io/zh-cn/docs/concepts/workloads/controllers/statefulset/) 资源所关联的 [`PersistentVolume`](https://kubernetes.io/zh-cn/docs/concepts/storage/persistent-volumes/#introduction) 所使用的 [访问模式（Access mode）](https://kubernetes.io/zh-cn/docs/concepts/storage/persistent-volumes/#access-modes)。

如果未指定，则默认使用 [`ReadWriteOnce`](https://kubernetes.io/zh-cn/docs/concepts/storage/persistent-volumes/#access-modes) 作为[访问模式（Access mode）](https://kubernetes.io/zh-cn/docs/concepts/storage/persistent-volumes/#access-modes)的值。

如果你会将 [`Model`](/pages/zh-CN/references/crd/model) 部署到 Ollama Operator 默认支持的 [kind](https://kind.sigs.k8s.io/) 和 [k3s](https://k3s.io/) 集群，你应该保持其当前的默认值 [`ReadWriteOnce`](https://kubernetes.io/zh-cn/docs/concepts/storage/persistent-volumes/#access-modes)。有且仅有当你部署到一个自托管的集群的时候，且 [`StorageClass`](https://kubernetes.io/zh-cn/docs/concepts/storage/storage-classes/#storageclass-objects) 支持时，就可以指定访问模式为 [`ReadWriteMany`](https://kubernetes.io/zh-cn/docs/concepts/storage/persistent-volumes/#access-modes)。

### `--expose`

默认：`false`

```shell
kollama deploy phi --expose
```

是否通过服务公开 [`Model`](/pages/zh-CN/references/crd/model) 所暴露的接口供外部访问，使其方便与 [`Model`](/pages/zh-CN/references/crd/model) 交互。

::: info 其实创建 Model 资源时，也会创建一个 `ClusterIP` 类型的服务

在没有指定 `--expose` 情况下，在为 [`Model`](/pages/zh-CN/references/crd/model) 创建资源时，Ollama Operator 也默认将为 [`Model`](/pages/zh-CN/references/crd/model) 创建一个用于集群中其他服务内部直接请求到 [`Model`](/pages/zh-CN/references/crd/model) 的关联服务，方便其他服务直接与之进行集成，其类型为 `ClusterIP`，名称与 [`Model`](/pages/zh-CN/references/crd/model) 所关联的 Deployment 的名称（也就是 `ollama-model-<model name>`）相同。

:::

默认情况下，[`--expose`](#expose) 参数包含在内时，将会创建一个类型为 [`NodePort`](https://kubernetes.io/zh-cn/docs/concepts/services-networking/service/#type-nodeport) 的 [`Service`](https://kubernetes.io/zh-cn/docs/concepts/services-networking/service/)。

你可以使用 [`--service-type`](#service-type) 参数加上 `LoadBalancer` 值（也就是 `--service-type=LoadBalancer`）来创建一个 [`LoadBalancer`](https://kubernetes.io/zh-cn/docs/concepts/services-networking/service/#loadbalancer) 类型的 [`Service`](https://kubernetes.io/zh-cn/docs/concepts/services-networking/service/)。

### `--service-type`

```shell
kollama deploy phi --expose --service-type=NodePort
```

默认：`NodePort`

暴露所部署的 [`Model`](/pages/zh-CN/references/crd/model) 服务时所连带创建的 [`Service`](https://kubernetes.io/zh-cn/docs/concepts/services-networking/service/) 类型。**有且仅有当 [`--expose`](#expose) 被指定时，该参数才会生效。**

如果没有指定该参数，那么创建时，将会默认使用 [`NodePort`](https://kubernetes.io/zh-cn/docs/concepts/services-networking/service/#type-nodeport) 类型的服务。

::: tip 了解有多少服务与 [`Model`](/pages/zh-CN/references/crd/model) 相关联...

```shell
kubectl get svc --selector ollama.ayaka.io/type=model
```

:::

可以指定为 [`LoadBalancer`](https://kubernetes.io/zh-cn/docs/concepts/services-networking/service/#loadbalancer) 来暴露一个 [`LoadBalancer`](https://kubernetes.io/zh-cn/docs/concepts/services-networking/service/#loadbalancer) 类型的服务。

### `--service-name`

```shell
kollama deploy phi --expose --service-name=phi-svc-nodeport
```

默认：`ollama-model-<model name>-<service type>`

暴露 [`Model`](/pages/zh-CN/references/crd/model) 所使用的 [`Service`](https://kubernetes.io/zh-cn/docs/concepts/services-networking/service/) 的名称。

如果未指定，则将使用 [`Model`](/pages/zh-CN/references/crd/model) 名称作为服务名称，并将 `-nodeport` 作为 [`NodePort`](https://kubernetes.io/zh-cn/docs/concepts/services-networking/service/#type-nodeport) 类型的服务的后缀。

### `--node-port`

```shell
kollama deploy phi --expose --service-type=NodePort --node-port=30000
```

默认：随机协商的端口

::: tip 了解 [`Model`](/pages/zh-CN/references/crd/model) 使用的 `NodePort` 端口号...

```shell
kubectl get svc --selector model.ollama.ayaka.io/name=<model name> -o json | jq ".spec.ports[0].nodePort"
```

:::

::: warning 并不可以随便指定端口号哦！

有这么几个限制是存在的：

1. 默认情况下 `30000-32767` 是 Kubernetes 集群中的 `NodePort` 端口范围。如果你想要使用这个范围之外的端口，你需要在集群中配置 `--service-node-port-range` 参数。
2. 你不能使用已经被其他服务占用的端口号。

有关自己选择端口号的更多信息，请参考 [Kubernetes 官方文档有关 `nodePort` 的章节](https://kubernetes.io/zh-cn/docs/concepts/services-networking/service/#nodeport-custom-port)。

:::

暴露 [`Model`](/pages/zh-CN/references/crd/model) 为 [`NodePort`](https://kubernetes.io/zh-cn/docs/concepts/services-networking/service/#type-nodeport) 类型的服务时所使用的指定的 [`nodePort`](https://kubernetes.io/zh-cn/docs/concepts/services-networking/service/#nodeport-custom-port) 端口号。

如果没有指定，将分配一个随机端口。该参数有且仅有在 [`--expose`](#expose) 指定，或者 [`--service-type`](#service-type) 设置为 `NodePort` 时才有效。
