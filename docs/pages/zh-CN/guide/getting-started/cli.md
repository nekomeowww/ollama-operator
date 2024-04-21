# 通过 `kollama` CLI 进行部署

我们有一个名为 `kollama` 的 CLI 来简化部署过程。这是一种将 Ollama 模型部署到 Kubernetes 集群的简单方法。

## 开始操作

1. 通过 `go install` 安装 `kollama` CLI：

```shell
go install github.com/nekomeowww/ollama-operator/cmd/kollama@latest
```

> 要了解 `kollama` CLI 支持的命令，请参阅 [`kollama`](/pages/zh-CN/references/cli/)。

2. 部署 Ollama 模型 CRD 到 Kubernetes 集群：

```shell
kollama deploy phi --expose --node-port 30101
```

> 有关 `deploy` 命令的更多信息，请参阅 [`kollama deploy`](/pages/zh-CN/references/cli/commands/deploy)。

3. 开始与模型进行交互吧：

```shell
OLLAMA_HOST=<节点 IP>:30101 ollama run phi
```

或者使用 `curl` 连接到与 OpenAI API 兼容的接口：

```shell
curl http://<节点 IP>:30101/v1/chat/completions -H "Content-Type: application/json" -d '{
  "model": "phi",
  "messages": [
      {
          "role": "user",
          "content": "Hello!"
      }
  ]
}'
```

