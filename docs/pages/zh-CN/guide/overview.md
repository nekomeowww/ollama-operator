# 概览

即便 [Ollama](https://github.com/ollama/ollama) 已经是一个强大的用于在本地运行大型语言模型的工具，并且 CLI 的用户体验与使用 Docker CLI 相同，但可惜的是，目前还无法在 Kubernetes 上直接复刻相同的用户体验，特别是同一集群上在运行多个模型时，涉及大量资源和配置。

这就是 Ollama Operator 发挥作用的地方：

- 在您的 Kubernetes 集群上安装 operator
- 应用所需的 CRDs
- 创建您的模型
- 等待模型被获取和加载，就是这样！

多亏了 [lama.cpp](https://github.com/ggerganov/llama.cpp) 的出色工作，**不再担心 Python 环境、CUDA 驱动程序**。
通往大型语言模型、AIGC、本地化代理、[🦜🔗 Langchain](https://www.langchain.com/) 等的旅程只需几步之遥！

## 能力

<div grid="~ cols-[auto_1fr] gap-1" items-start my-1>
  <div h=[1rem]><div i-icon-park-outline:check-one text="green-600" /></div>
  <span>在同一集群上运行多个模型的能力</span>
  <div h=[1rem]><div i-icon-park-outline:check-one text="green-600" /></div>
  <span>与所有 Ollama 模型、API 和 CLI 兼容</span>
  <div h=[1rem]><div i-icon-park-outline:check-one text="green-600" /></div>
  <span>可以在 <a href="https://kubernetes.io/">常规 Kubernetes 集群</a>、<a href="https://k3s.io/">K3s 集群</a> (Respberry Pi（树莓派），TrueNAS SCALE，等等), <a href="https://kind.sigs.k8s.io/">kind</a>, <a href="https://minikube.sigs.k8s.io/docs/">minikube</a> 上运行</span>
  <div h=[1rem]><div i-icon-park-outline:check-one text="green-600" /></div>
  <span>易于安装、卸载和升级</span>
  <div h=[1rem]><div i-icon-park-outline:check-one text="green-600" /></div>
  <span>一次拉取，全节点共享（就像普通镜像一样）</span>
  <div h=[1rem]><div i-icon-park-outline:check-one text="green-600" /></div>
  <span>易于与现有的 Kubernetes 服务、Ingress，微服务网关等结合使用</span>
  <div h=[1rem]><div i-icon-park-outline:check-one text="green-600" /></div>
  <span>除去 Kubernetes 以外，什么都不需要配置</span>
</div>

## 需求

### Kubernetes 集群

::: tip 我必须要有一整套云上或者自部署的 Kubernetes 集群才能用 Ollama Operator 吗？

其实并不是，对于任意的 macOS，Windows 设备而言，只需要安装了 [Docker Desktop](https://www.docker.com/products/docker-desktop/) 或者 macOS 独享的 [OrbStack](https://orbstack.dev/)，配合用于在本地运行 Kubernetes 集群的 [kind](https://kind.sigs.k8s.io/) 和 [minikube](https://minikube.sigs.k8s.io/docs/) 工具即可在本地启动一个自己的 Kubernetes 集群。

Kubernetes 并没有想象中那么难，只要有 Docker 和一个 Kubernetes 工具，就可以在本地运行 Kubernetes 集群，然后安装 Ollama Operator，就可以在本地运行大型语言模型了。

:::

- Kubernetes
- K3s
- kind
- minikube

### 内存需求

要运行 7B 机型，节点上至少应有 8GB 内存；要运行 13B 机型，节点上至少应有 16GB 内存；要运行 33B 机型，节点上至少应有 32GB 内存。


### 磁盘需求

与一般容器镜像的大小相比，下载的大型语言模型的实际大小非常大。

1. 建议使用快速稳定的网络连接下载模型。
2. 如果要运行大于 13B 的模型，则需要高效的存储设备来存储模型。
