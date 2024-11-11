1. 安装 Operator.

```shell
kubectl apply \
  --server-side=true \
  -f https://raw.githubusercontent.com/nekomeowww/ollama-operator/v0.10.1/dist/install.yaml
```

2. 等待 Operator 就绪：

```shell
kubectl wait \
   -n ollama-operator-system \
  --for=jsonpath='{.status.readyReplicas}'=1 deployment/ollama-operator-controller-manager
```
