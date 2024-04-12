# 架构设计

Ollama Operator 会根据 CRD 的定义，创建两个主要组件：

1. **模型推理服务**：模型推断服务器是一个简单的 API 服务器，用于运行模型并为模型的 API 提供服务。它在 Kubernetes 集群中作为 `Deployment` 创建。
2. **模型镜像托管服务**：我们需要一项服务来存储下载过的模型镜像并重复使用它们，而不是在模型推理服务创建的时候每次都下载自己的模型，因此将创建一个 `StatefulSet` 和一个 `PersistentVolumeClaim` 。实际的文件将会被存储在动态调配的 `PersistentVolume` 中。

::: info 为什么会需要额外的 **模型镜像托管服务**？
虽然 Ollama 的 [`Modelfile`](https://github.com/ollama/ollama/blob/main/docs/modelfile.md) 创建的模型镜像是有效的 OCI 格式映像，但由于镜像内的 `contentType` 值和 `Modelfile` 镜像的整体结构与一般的容器映像不兼容，因此无法直接使用一般容器运行时（containerd，docker）运行模型。因此，我们需要在 Kubernetes 集群上持久化**模型镜像托管服务**的独立 Service/Deployment，以保存和缓存之前下载的模型镜像。

虽说是一个额外的服务，但是实际上我们并不会创建一个 Docker Registry 或者 Harbor 这样大体量的 Registry 服务器，而是简单地用 `ollama serve` 这个内置命令来启动一个简单服务，在每次的模型推理服务创建时，都会有单独的镜像拉取请求直接发送到这个服务上。
:::

它创建的具体资源及其之间的关系如下图所示：

 <picture>
 <source
   srcset="/architecture-theme-night.png"
   media="(prefers-color-scheme: dark)"
 />
 <source
   srcset="/architecture-theme-day.png"
   media="(prefers-color-scheme: light), (prefers-color-scheme: no-preference)"
 />
 <img src="/architecture-theme-day.png" />
</picture>
