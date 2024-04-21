# `kollama undeploy`

::: info 无需担心删除过去的 [`Model`](/pages/zh-CN/references/crd/model) 部署后会需要重新下载模型镜像哦！

Ollama Operator 在为 [`Model`](/pages/zh-CN/references/crd/model) 部署的普通的 [`Deployment`](https://kubernetes.io/zh-cn/docs/concepts/workloads/controllers/deployment/) 类型资源外，还会部署单独的用于存储下载过的模型的 [`StatefulSet`](https://kubernetes.io/zh-cn/docs/concepts/workloads/controllers/statefulset/) 资源和其对应的存储资源。

因此，即使删除了模型的部署，也不会影响到已经下载的模型文件。他们都会存放在名为 `ollama-models-store` 的单独的资源中，直到手动删除。

你可以通过下面的命令来查看 `ollama-models-store` 的状态

```shell
kubectl describe statefulset ollama-models-store
```

:::

`kollama undeploy` 用于删除 [`Model`](/pages/zh-CN/references/crd/model) 的部署。

[[toc]]

## 用例

### 删除 [`Model`](/pages/zh-CN/references/crd/model) 的部署

```shell
kollama undeploy phi
```

### 删除指定命名空间中的 [`Model`](/pages/zh-CN/references/crd/model) 的部署

```shell
kollama undeploy phi --namespace=production
```

## 选项

### `--namespace`

如果配置了该参数，将会在指定的命名空间中删除模型的部署。
