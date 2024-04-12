# Overview

While [Ollama](https://github.com/ollama/ollama) is a powerful tool for running large language models locally, and the user experience of CLI is just the same as using Docker CLI, it's not possible yet to replicate the same user experience on Kubernetes, especially when it comes to running multiple models on the same cluster with loads of resources and configurations.

That's where the Ollama Operator kicks in:

- Install the operator on your Kubernetes cluster
- Apply the needed CRDs
- Create your models
- Wait for the models to be fetched and loaded, that's it!

Thanks to the great works of [lama.cpp](https://github.com/ggerganov/llama.cpp), **no more worries about Python environment, CUDA drivers.**

The journey to large language models, AIGC, localized agents, [🦜🔗 Langchain](https://www.langchain.com/) and more is just a few steps away!

## Features

<div grid="~ cols-[auto_1fr] gap-1" items-start my-1>
  <div h=[1rem]><div i-icon-park-outline:check-one text="green-600" /></div>
  <span>Abilities to run multiple models on the same cluster.</span>
  <div h=[1rem]><div i-icon-park-outline:check-one text="green-600" /></div>
  <span>Compatible with all Ollama models, APIs, and CLI.</span>
  <div h=[1rem]><div i-icon-park-outline:check-one text="green-600" /></div>
  <span>Able to run on <a href="https://kubernetes.io/">general Kubernetes clusters</a>, <a href="https://k3s.io/">K3s clusters</a> (Respberry Pi, TrueNAS SCALE, etc.), <a href="https://kind.sigs.k8s.io/">kind</a>, <a href="https://minikube.sigs.k8s.io/docs/">minikube</a>, etc. You name it!</span>
  <div h=[1rem]><div i-icon-park-outline:check-one text="green-600" /></div>
  <span>Easy to install, uninstall, and upgrade.</span>
  <div h=[1rem]><div i-icon-park-outline:check-one text="green-600" /></div>
  <span>Pull image once, share across the entire node (just like normal images).</span>
  <div h=[1rem]><div i-icon-park-outline:check-one text="green-600" /></div>
  <span>Easy to expose with existing Kubernetes services, ingress, etc.</span>
  <div h=[1rem]><div i-icon-park-outline:check-one text="green-600" /></div>
  <span>Doesn't require any additional dependencies, just Kubernetes</span>
</div>

## Requirements

### Kubernetes cluster

::: tip Do I have to have a complete deployed Kubernetes cluster over cloud or self managed to use Ollama Operator?

In fact, it is not.

For any macOS, Windows device, you just need to install [Docker Desktop](https://www.docker.com/products/docker-desktop/) or the macOS-only [OrbStack](https://), along with the utilities like [kind](https://kind.sigs.k8s.io/) and [minikube](https://minikube.sigs.k8s.io/docs/) tools for running a Kubernetes cluster locally, you can start your own Kubernetes cluster locally.

Kubernetes is not as difficult as you might think, as long as you have Docker and a Kubernetes tool, you can run a Kubernetes cluster locally, and then install the Ollama Operator to run large language models locally.

:::

- Kubernetes
- K3s
- kind
- minikube

### Memory requirements

You should have at least 8 GB of RAM available on your node to run the 7B models, 16 GB to run the 13B models, and 32 GB to run the 33B models.

### Disk requirements

The actual size of downloaded large language models are huge by comparing to the size of general container images.

1. Fast and stable network connection is recommended to download the models.
2. Efficient storage is required to store the models if you want to run models larger than 13B.
