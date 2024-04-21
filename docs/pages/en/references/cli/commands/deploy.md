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
kubectl get svc --selector model.ollama.ayaka.io
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
