# Deploy through `kollama` CLI

We have a CLI called `kollama` here to simplify the deployment process. It is a simple way to deploy Ollama models to your Kubernetes cluster.

## Getting Started

1. Install the CLI:

```shell
go install github.com/nekomeowww/ollama-operator/cmd/kollama@latest
```

2. Deploy a model:

```shell
kollama deploy phi --expose --node-port 30001
```

That's it.

3. Interact with the model:

```shell
OLLAMA_HOST=<Node ip>:30001 ollama run phi
```

or use the OpenAI API compatible endpoint:

```shell
curl http://<Node ip>:30001/v1/chat/completions -H "Content-Type: application/json" -d '{
  "model": "phi",
  "messages": [
      {
          "role": "user",
          "content": "Hello!"
      }
  ]
}'
```
