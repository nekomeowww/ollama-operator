---
layout: home
sidebar: false

title: Ollama Operator

hero:
  name: Ollama Operator
  text: Large language models, scaled, deployed
  tagline: Yet another operator for running large language models on Kubernetes with ease. Powered by Ollama! 🐫
  image:
    src: /logo.png
    alt: Ollama Operator
  actions:
    - theme: brand
      text: Getting Started
      link: /pages/en/guide/getting-started/
    - theme: alt
      text: View on GitHub
      link: https://github.com/nekomeowww/ollama-operator

features:
  - icon: <div i-twemoji:rocket></div>
    title: Launch and chat
    details: Easy to use API, the spec is simple enough to just a few lines of YAML to deploy a model, and you can chat with it right away.
  - icon: <div i-twemoji:ship></div>
    title: Cross-kubernetes
    details: Extend the user experience of Ollama to any Kubernetes cluster, edge or any cloud infrastructure, with the same spec, and chat with it from anywhere.
  - icon: <div i-simple-icons:openai></div>
    title: OpenAI API compatible
    details: Your familiar <code>/v1/chat/completions</code> endpoint is here, with the same request and response format. No need to change your code or switch to another API.
  - icon: <div i-twemoji:parrot></div>
    title: Langchain ready
    details: Power on to function calling, agents, knowldgebase retrieving. Unleash all the power Langchain has out of box with Ollama Operator.
---

<script setup>
import { NuAsciinemaPlayer } from '@nolebase/ui-asciinema'
</script>

### Getting started

<GettingStartedBlocksEn />

### See it in action

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
