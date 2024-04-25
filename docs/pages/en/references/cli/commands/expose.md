# `kollama expose`

`kolamma expose` is a command to expose a [`Model`](/pages/en/references/crd/model) as a service.

[[toc]]

## Use cases

### Expose a [`Model`](/pages/en/references/crd/model) service

```shell
kollama expose phi
```

### Expose a [`Model`](/pages/en/references/crd/model) service in a specific namespace

```shell
kollama expose phi --namespace=production
```

### Expose the service as a [`LoadBalancer`](https://kubernetes.io/docs/concepts/services-networking/service/#loadbalancer) type

```shell
kollama expose phi --service-type=LoadBalancer
```

## Flags

### `--namespace`

If present, the namespace scope for this CLI request.

### `--service-type`

```shell
kollama deploy phi --expose --service-type=NodePort
```

Default: `NodePort`

Type of the [`Service`](https://kubernetes.io/docs/concepts/services-networking/service/) to expose the model. **Only valid when [`--expose`](#expose) is specified.**

If not specified, the service will be exposed as [`NodePort`](https://kubernetes.io/docs/concepts/services-networking/service/#type-nodeport).

::: tip To understand how many Services are associated to models...

```shell
kubectl get svc --selector ollama.ayaka.io/type=model
```

:::

Use [`LoadBalancer`](https://kubernetes.io/docs/concepts/services-networking/service/#loadbalancer) to expose the service as [`LoadBalancer`](https://kubernetes.io/docs/concepts/services-networking/service/#loadbalancer).

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
