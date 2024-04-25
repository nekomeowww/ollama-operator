# Supported models

The out of box supported models are listed below.
You can use them directly by specifying the model image in the `Model` CRD.

> [!TIP] Not just these
> By the power of [`Modelfile`](https://github.com/ollama/ollama/blob/main/docs/modelfile.md) backed by Ollama, you can create and bundle any of your own model. **As long as it's a GGUF formatted model.**

| Model                                                              | Parameters | Size  | Model image         | Full model image URL                           | Multi-modal | Uncensored |
| ------------------------------------------------------------------ | ---------- | ----- | ------------------- | ---------------------------------------------- | ----------- | ---------- |
| [Phi-3 Mini](https://ollama.com/library/phi3)                      | 3.8B       | 2.3GB | `phi3`              | `registry.ollama.ai/library/phi3`              |             |            |
| [Llama 3](https://ollama.com/library/llama3)                       | 8B         | 4.7GB | `llama3`            | `registry.ollama.ai/library/llama3`            |             |            |
| [Dolphin Llama 3](https://ollama.com/library/dolphin-llama3)       | 8B         | 4.7GB | `dolphin-llama3`    | `registry.ollama.ai/dolphin-llama3`            |             | ✅          |
| [WizardLM-2](https://ollama.com/library/wizardlm2)                 | 7B         | 4.1GB | `wizardlm2`         | `registry.ollama.ai/library/wizardlm2`         |             |            |
| [Llama 2](https://ollama.com/library/llama2)                       | 7B         | 3.8GB | `llama2`            | `registry.ollama.ai/library/llama2`            |             |            |
| [Mistral](https://ollama.com/library/mistral)                      | 7B         | 4.1GB | `mistral`           | `registry.ollama.ai/library/mistral`           |             |            |
| [Mixtral 8x7B](https://ollama.com/library/mixtral:8x7b)            | 8x7B       | 26GB  | `mixtral:8x7b`      | `registry.ollama.ai/library/mixtral:8x7b`      |             |            |
| [Mixtral 8x22B](https://ollama.com/library/mixtral:8x22b)          | 8x22B      | 80GB  | `mixtral:8x22b`     | `registry.ollama.ai/library/mixtral:8x22b`     |             |            |
| [Command R](https://ollama.com/library/command-r)                  | 35B        | 20GB  | `command-r`         | `registry.ollama.ai/library/command-r`         |             |            |
| [Command R Plus](https://ollama.com/library/command-r-plus)        | 104B       | 59GB  | `command-r-plus`    | `registry.ollama.ai/library/command-r-plus`    |             |            |
| [Dolphin Phi](https://ollama.com/library/dolphin-phi)              | 2.7B       | 1.6GB | `dolphin-phi`       | `registry.ollama.ai/library/dolphin-phi`       |             | ✅          |
| [Phi-2](https://ollama.com/library/phi)                            | 2.7B       | 1.7GB | `phi`               | `registry.ollama.ai/library/phi`               |             |            |
| [Neural Chat](https://ollama.com/library/neural-chat)              | 7B         | 4.1GB | `neural-chat`       | `registry.ollama.ai/library/neural-chat`       |             |            |
| [Starling](https://ollama.com/library/starling-lm)                 | 7B         | 4.1GB | `starling-lm`       | `registry.ollama.ai/library/starling-lm`       |             |            |
| [Code Llama](https://ollama.com/library/codellama)                 | 7B         | 3.8GB | `codellama`         | `registry.ollama.ai/library/codellama`         |             |            |
| [Llama 2 Uncensored](https://ollama.com/library/llama2-uncensored) | 7B         | 3.8GB | `llama2-uncensored` | `registry.ollama.ai/library/llama2-uncensored` |             | ✅          |
| [Llama 2 13B](https://ollama.com/library/llama2)                   | 13B        | 7.3GB | `llama2:13b`        | `registry.ollama.ai/library/llama2:13b`        |             |            |
| [Llama 2 70B](https://ollama.com/library/llama2)                   | 70B        | 39GB  | `llama2:70b`        | `registry.ollama.ai/library/llama2:70b`        |             |            |
| Orca Mini                                                          | 3B         | 1.9GB | `orca-mini`         | `registry.ollama.ai/library/orca-mini`         |             |            |
| Vicuna                                                             | 7B         | 3.8GB | `vicuna`            | `registry.ollama.ai/library/vicuna`            |             |            |
| LLaVA                                                              | 7B         | 4.5GB | `llava`             | `registry.ollama.ai/library/llava`             | ✅           |            |
| Gemma 2B                                                           | 2B         | 1.4GB | `gemma:2b`          | `registry.ollama.ai/library/gemma:2b`          |             |            |
| Gemma 7B                                                           | 7B         | 4.8GB | `gemma:7b`          | `registry.ollama.ai/library/gemma:7b`          |             |            |

Full list of available images can be found at [Ollama Library](https://ollama.com/library).
