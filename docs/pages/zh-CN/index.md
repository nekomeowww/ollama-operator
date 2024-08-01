---
layout: home
sidebar: false

title: Ollama Operator

hero:
  name: Ollama Operator
  text: 大语言模型，伸缩自如，轻松部署
  tagline: 一个在 Kubernetes 上让部署和运行大型语言模型变得轻松简单的 Operator，由 Ollama 强力驱动 🐫
  image:
    src: /logo.png
    alt: Ollama Operator
  actions:
    - theme: brand
      text: 快速开始
      link: /pages/zh-CN/guide/getting-started/
    - theme: alt
      text: 在 GitHub 上查看
      link: https://github.com/nekomeowww/ollama-operator

features:
  - icon: <div i-twemoji:rocket></div>
    title: 简单易用
    details: 易于使用的 API，足够简单的 CRD 规格，只需几行 YAML 定义即可部署一个模型，然后立即与之交互。
  - icon: <div i-twemoji:ship></div>
    title: 兼容各种 Kubernetes
    details: 将 Ollama 的用户体验扩展到任何 Kubernetes 集群、边缘或任何云基础设施，使用相同的 CRD API，从任何地方与之交互。
  - icon: <div i-simple-icons:openai></div>
    title: 兼容 OpenAI API
    details: 您熟悉的 <code>/v1/chat/completions</code> 接口就在这里，具有相同的请求和响应格式。无需更改代码或切换到其他 API。
  - icon: <div i-twemoji:parrot></div>
    title: 随时对接 Langchain
    details: 强大的功能调用、代理、知识库检索。使用 Ollama Operator，释放 Langchain 开箱即用的所有功能。
---

<script setup>
import { NuAsciinemaPlayer } from '@nolebase/ui-asciinema'
</script>

### 开始

<GettingStartedBlocksZhCn />

### 一睹为快

<br>

<div w-full rounded-xl overflow-hidden>
  <NuAsciinemaPlayer
    src="/demo.cast"
    :loop="true"
    :autoPlay="true"
    :rows="20"
    :speed="3"
    w-full
  />
</div>
