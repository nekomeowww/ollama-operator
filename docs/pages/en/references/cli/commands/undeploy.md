# `kollama undeploy`

::: info No need to worry about deleting past [`Model`](/pages/en/references/crd/model) deployments and having to re-download the model image!

Ollama Operator deploys a separate [`StatefulSet`](https://kubernetes.io/zh-cn/docs/concepts/workloads/controllers/statefulset/) resource for storing downloaded Ollama model images and its corresponding storage resource in addition to the normal [`Deployment`](https://kubernetes.io/zh-cn/docs/concepts/workloads/controllers/deployment/) type resource for [`Model`](/pages/en/references/crd/model)s.

Therefore, even if the deployment of a [`Model`](/pages/en/references/crd/model) is deleted, it will not affect the model images that have already been downloaded. They are stored in a separate resource called `ollama-models-store` until manually deleted.

You can check the status of `ollama-models-store` with the following command:

```shell
kubectl describe statefulset ollama-models-store
```

:::

`kollama undeploy` is used to delete the deployment of a [`Model`](/pages/en/references/crd/model).

[[toc]]

## Use cases

### Delete the deployment of a [`Model`](/pages/en/references/crd/model)

```shell
kollama undeploy phi
```

### Delete the deployment of a [`Model`](/pages/en/references/crd/model) in a specific namespace

```shell
kollama undeploy phi --namespace=production
```

## Flags

### `--namespace`

If present, the namespace scope for this CLI request.
