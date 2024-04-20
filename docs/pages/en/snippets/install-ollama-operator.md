1. Install operator.

```shell
kubectl apply -f https://raw.githubusercontent.com/nekomeowww/ollama-operator/main/dist/install.yaml
```

2. Wait for the operator to be ready:

```shell
kubectl wait --for=jsonpath='{.status.replicas}'=2 deployment/ollama-operator-controller-manager -n ollama-operator-system
```
