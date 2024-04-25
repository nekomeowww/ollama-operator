# `kollama deploy`

`kollama deploy` command is used to deploy a [`Model`](/pages/en/references/crd/model) to the Kubernetes cluster. It's basic a wrapper and utility CLI to interact with the Ollama Operator by manipulating CRD resources.

[[toc]]

## Use cases

### Deploy model that lives on [registry.ollama.ai](https://registry.ollama.ai)

```shell
kollama deploy phi
```

### Deploy to a specific namespace

```shell
kollama deploy phi --namespace=production
```

### Deploy [`Model`](/pages/en/references/crd/model) that lives on a custom registry

```shell
kollama deploy phi --image=registry.example.com/library/phi:latest
```

### Deploy [`Model`](/pages/en/references/crd/model) with exposed NodePort service for external access

```shell
kollama deploy phi --expose
```

### Deploy [`Model`](/pages/en/references/crd/model) with exposed LoadBalancer service for external access

```shell
kollama deploy phi --expose --service-type=LoadBalancer
```

### Deploy [`Model`](/pages/en/references/crd/model) with resources limits

The following example deploys the `phi` model with CPU limit to `1` and memory limit to `1Gi`.

```shell
kollama deploy phi --limit=cpu=1 --limit=memory=1Gi
```

## Flags

### `--namespace`

If present, the namespace scope for this CLI request.

### `--image`

Default: `registry.ollama.ai/library/<model name>:latest`

```shell
kollama deploy phi --image=registry.ollama.ai/library/phi:latest
```

Model image to deploy.

- If not specified, the [`Model`](/pages/en/references/crd/model) name will be used as the image name (will be pulled from `registry.ollama.ai/library/<model name>` by default if no registry is specified). For example, if the [`Model`](/pages/en/references/crd/model) name is `phi`, the image name will be `registry.ollama.ai/library/phi:latest`.
- If not specified, the tag will be latest.

### `--limit` (supports multiple flags)

> Multiple limits can be specified by using the flag multiple times.

Resource limits for the deployed [`Model`](/pages/en/references/crd/model). This is useful for clusters that don't have a large enough number of resources, or if you want to deploy multiple [`Models`](/pages/en/references/crd/model) in a cluster with limited resources.

::: tip For resource limits on NVIDIA, AMD GPUs...

In Kubernetes, any GPU resource follows this pattern for resources labels:

```yaml
resources:
  limits:
    gpu-vendor.example/example-gpu: 1 # requesting 1 GPU
```

Using `nvidia.com/gpu` allows you to limit the number of NVIDIA GPUs, therefore when using `kollama deploy` you can use `--limit nvidia.com/gpu=1` to specify the number of NVIDIA GPUs as `1`:

```shell
kollama deploy phi --limit=nvidia.com/gpu=1
```

this is what it may looks like in the YAML configuration file:


```yaml
resources:
  limits:
    nvidia.com/gpu: 1 # requesting 1 GPU # [!code focus]
```

> [Documentation on using resource labels with `nvidia/k8s-device-plugin`](https://github.com/NVIDIA/k8s-device-plugin?tab=readme-ov-file#enabling-gpu-support-in-kubernetes)

Using `amd.com/gpu` allows you to limit the number of AMD GPUs, therefore when using `kollama deploy` you can use `--limit amd.com/gpu=1` to specify the number of AMD GPUs as `1`.

```shell
kollama deploy phi --limit=amd.com/gpu=1
```

this is what it may looks like in the YAML configuration file:

```yaml
resources:
  limits:
    amd.com/gpu: 1 # requesting a GPU  # [!code focus]
```

> [Example YAML manifest of labels with `ROCm/k8s-device-plugin`](https://github.com/ROCm/k8s-device-plugin/blob/4607bf06b700e53803d566e0bf9555f773f0b4f1/example/pod/alexnet-gpu.yaml)

Your can read more here: [Schedule GPUs | Kubernetes](https://kubernetes.io/docs/tasks/manage-gpus/scheduling-gpus/)

:::

::: details I have deployed [`Model`](/pages/en/references/crd/model), but I want to change the resource limit...

Of course you can, with the [`kubectl set resources`](https://kubernetes.io/zh-cn/docs/reference/kubectl/generated/kubectl_set/kubectl_set_resources/) command, you can change the resource limit:

```shell
kubectl set resources deployment -l model.ollama.ayaka.io/name=<model name> --limits cpu=4
```

For memory limits:

```shell
kubectl set resources deployment -l model.ollama.ayaka.io/name=<model name> --limits memory=8Gi
```

:::

The format is `<resource>=<quantity>`.

For example: `--limit=cpu=1` `--limit=memory=1Gi`.

### `--storage-class`

```shell
kollama deploy phi --storage-class=standard
```

[`StorageClass`](https://kubernetes.io/docs/concepts/storage/storage-classes/#storageclass-objects) to use for the [`Model`](/pages/en/references/crd/model)'s associated [`PersistentVolumeClaim`](https://kubernetes.io/docs/concepts/storage/persistent-volumes/#persistentvolumeclaims).

If not specified, the [default `StorageClass`](https://kubernetes.io/docs/concepts/storage/storage-classes/#default-storageclass) will be used.

### `--pv-access-mode`

```shell
kollama deploy phi --pv-access-mode=ReadWriteMany
```

[Access mode](https://kubernetes.io/docs/concepts/storage/persistent-volumes/#access-modes) for Ollama Operator created image store (to cache pulled images)'s [`StatefulSet`](https://kubernetes.io/docs/concepts/workloads/controllers/statefulset/) resource associated [`PersistentVolume`](https://kubernetes.io/docs/concepts/storage/persistent-volumes/#introduction).

If not specified, the access mode will be [`ReadWriteOnce`](https://kubernetes.io/docs/concepts/storage/persistent-volumes/#access-modes).

If you are deploying [`Model`](/pages/en/references/crd/model)s into default deployed [kind](https://kind.sigs.k8s.io/) and [k3s](https://k3s.io/) clusters, you should keep it as [`ReadWriteOnce`](https://kubernetes.io/docs/concepts/storage/persistent-volumes/#access-modes). If you are deploying [`Model`](/pages/en/references/crd/model)s into a custom cluster, you can set it to [`ReadWriteMany`](https://kubernetes.io/docs/concepts/storage/persistent-volumes/#access-modes) if [`StorageClass`](https://kubernetes.io/docs/concepts/storage/storage-classes/#storageclass-objects) supports it.

### `--expose`

Default: `false`

```shell
kollama deploy phi --expose
```

Whether to expose the [`Model`](/pages/en/references/crd/model) through a Service for external access and makes it easy to interact with the [`Model`](/pages/en/references/crd/model).

::: info Actually, when creating a Model resource, a `ClusterIP` type service will be created

At the case where users didn't supply either `--expose` flag, Ollama Operator will create a associated service for the [`Model`](/pages/en/references/crd/model) with the type of `ClusterIP` with the same name as the corresponding Deployment by default, and the service will be used for internal communication between the [`Model`](/pages/en/references/crd/model) and other services in the cluster.

:::

By default, [`--expose`](#expose) will create a [`NodePort`](https://kubernetes.io/docs/concepts/services-networking/service/#type-nodeport) service.

Use `--expose=LoadBalancer` to create a [`LoadBalancer`](https://kubernetes.io/docs/concepts/services-networking/service/#loadbalancer) service.

### `--service-type`

```shell
kollama deploy phi --expose --service-type=NodePort
```

Default: `NodePort`

Type of the [`Service`](https://kubernetes.io/docs/concepts/services-networking/service/) to expose the [`Model`](/pages/en/references/crd/model). **Only valid when [`--expose`](#expose) is specified.**

If not specified, the service will be exposed as [`NodePort`](https://kubernetes.io/docs/concepts/services-networking/service/#type-nodeport).

::: tip To understand how many Services are associated to [`Model`](/pages/en/references/crd/model)...

```shell
kubectl get svc --selector ollama.ayaka.io/type=model
```

:::

Use [`LoadBalancer`](https://kubernetes.io/docs/concepts/services-networking/service/#loadbalancer) to expose the service as [`LoadBalancer`](https://kubernetes.io/docs/concepts/services-networking/service/#loadbalancer).

### `--service-name`

```shell
kollama deploy phi --expose --service-name=phi-svc-nodeport
```

Default: `ollama-model-<model name>-<service type>`

Name of the [`Service`](https://kubernetes.io/docs/concepts/services-networking/service/) to expose the [`Model`](/pages/en/references/crd/model).

If not specified, the [`Model`](/pages/en/references/crd/model) name will be used as the service name with `-nodeport` as the suffix for [`NodePort`](https://kubernetes.io/docs/concepts/services-networking/service/#type-nodeport).

### `--node-port`

```shell
kollama deploy phi --expose --service-type=NodePort --node-port=30000
```

Default: Random port

::: tip To understand what NodePort is used for the [`Model`](/pages/en/references/crd/model)...

```shell
kubectl get svc --selector model.ollama.ayaka.io/name=<model name> -o json | jq ".spec.ports[0].nodePort"
```

:::

::: warning You can't simply specify a port number!

There are several restrictions:

1. By default, `30000-32767` is the `NodePort` port range in the Kubernetes cluster. If you want to use ports outside this range, you need to configure the `--service-node-port-range` parameter for the cluster.
2. You can't use the port number already occupied by other services.

For more information about choosing your own port number, please refer to [Chapter of Kubernetes Official Document about `nodePort`](https://kubernetes.io/docs/concepts/services-networking/service/#nodeport-custom-port).

:::

[`nodePort`](https://kubernetes.io/docs/concepts/services-networking/service/#nodeport-custom-port) to expose the [`Model`](/pages/en/references/crd/model).

If not specified, a random port will be assigned. Only valid when [`--expose`](#expose) is specified, and [`--service-type`](#service-type) is set to NodePort.
