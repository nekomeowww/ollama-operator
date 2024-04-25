# `kollama expose`

`kollama expose` 用于暴露 [`Model`](/pages/zh-CN/references/crd/model) 服务。

[[toc]]

## 用例

### 暴露 [`Model`](/pages/zh-CN/references/crd/model) 服务

```shell
kollama expose phi
```

### 暴露指定命名空间中的 [`Model`](/pages/zh-CN/references/crd/model) 服务

```shell
kollama expose phi --namespace=production
```

### 暴露为 [`LoadBalancer`](https://kubernetes.io/zh-cn/docs/concepts/services-networking/service/#loadbalancer) 类型的服务

```shell
kollama expose phi --service-type=LoadBalancer
```

## 选项

### `--namespace`

如果配置了该参数，将会在指定的命名空间中操作 [`Model`](/pages/zh-CN/references/crd/model)。

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
