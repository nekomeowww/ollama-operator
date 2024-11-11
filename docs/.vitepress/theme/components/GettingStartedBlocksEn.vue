<template>
  <TitleBlockContainerGroup>
    <TitleBlockContainer
      href="/pages/en/guide/getting-started/cli"
      class="[&_.title-block-container]:bg-zinc-200 [&_.title-block-container]:dark:bg-zinc-600 [&_.title-block-container]:text-zinc-900 [&_.title-block-container]:dark:text-zinc-100 text-$vp-c-brand-1!"
    >
      <template #title>
        <div flex-1 flex items-center gap-2>
          <div class="title-block-container" transition="all ease-in-out duration-300" p-2 rounded-lg w-10 h-10 items-center flex justify-center>
            <div i-icon-park-outline:bank-card-two />
          </div>
          <span font-semibold>Deploy models through Kollama CLI</span>
        </div>
      </template>
      <div flex="~ col" h-full gap-2>
        <div flex-1>
          <p m="0!" pt="2" text="zinc-600 dark:zinc-300">
            Dislike YAML?
          </p>
          <p m="0!" text="zinc-600 dark:zinc-300">
            Faster and better UX?
          </p>
          <p m="0!" text="zinc-600 dark:zinc-300">
            No worries, kollama here to rescue!
          </p>
        </div>
        <pre class="title-block-body whitespace-pre-line m-0" transition="all ease-in-out duration-300" font-mono p-2 bg="zinc-200/50 dark:zinc-700/80" rounded-lg>
          <span text="zinc-400"># <span>as regular binary cli</span></span>
          $ <span text="blue-800 dark:blue-300">kollama </span><span text="green-800 dark:green-300">deploy phi</span>
        </pre>
        <span items-center>or</span>
        <pre class="title-block-body whitespace-pre-line m-0" transition="all ease-in-out duration-300" font-mono p-2 bg="zinc-200/50 dark:zinc-700/80" rounded-lg>
          <span text="zinc-400"># <span>as kubectl plugin</span></span>
          $ <span text="blue-800 dark:blue-300">kubectl </span><span text="green-800 dark:green-300">ollama deploy phi</span>
        </pre>
      </div>
    </TitleBlockContainer>
    <TitleBlockContainer
      href="/pages/en/guide/getting-started/crd"
      class="[&_.title-block-container]:bg-zinc-200 [&_.title-block-container]:dark:bg-zinc-600 [&_.title-block-container]:text-zinc-900 [&_.title-block-container]:dark:text-zinc-100 text-$vp-c-brand-1!"
    >
      <template #title>
        <div flex-1 flex items-center gap-2>
          <div class="title-block-container" transition="all ease-in-out duration-300" p-2 rounded-lg w-10 h-10 items-center flex justify-center>
            <div i-icon-park-outline:file-code />
          </div>
          <span font-semibold>Deploy models through CRD</span>
        </div>
      </template>
      <div flex="~ col" h-full>
        <div flex-1>
          <p m="0!" pt="2" text="zinc-600 dark:zinc-300">
            Fine-grained control over parameters?
          </p>
          <p m="0!" text="zinc-600 dark:zinc-300">
            GitOps and CI/CD?
          </p>
          <p m="0!" text="zinc-600 dark:zinc-300">
            CRD is simple enough with 6 lines!
          </p>
        </div>
        <pre class="title-block-body whitespace-pre-line m-0" transition="all ease-in-out duration-300" font-mono p-2 bg="zinc-200/50 dark:zinc-700/80" rounded-lg>
          <span text="blue-800 dark:blue-300">apiVersion: </span><span text="green-800 dark:green-300">ollama.ayaka.io/v1</span>
          <span text="blue-800 dark:blue-300">kind: </span><span text="green-800 dark:green-300">Model</span>
          <span text="blue-800 dark:blue-300">metadata:</span>
          <span text="blue-800 dark:blue-300">&nbsp;&nbsp;name: </span><span text="green-800 dark:green-300">phi</span>
          <span text="blue-800 dark:blue-300">spec:</span>
          <span text="blue-800 dark:blue-300">&nbsp;&nbsp;image: </span><span text="green-800 dark:green-300">phi</span>
        </pre>
      </div>
    </TitleBlockContainer>
  </TitleBlockContainerGroup>
</template>
