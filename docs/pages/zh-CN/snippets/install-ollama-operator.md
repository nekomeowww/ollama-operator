1. 安装 Operator.

```shell
kubectl apply -f https://raw.githubusercontent.com/nekomeowww/ollama-operator/main/dist/install.yaml
```

2. 等待 Operator 就绪：

```shell
kubectl wait --for=jsonpath='{.status.readyReplicas}'=1 deployment/ollama-operator-controller-manager -n ollama-operator-system
```
