<script setup lang="ts">
const props = defineProps<{
  title?: string
  href?: string
}>()
</script>

<template>
  <a
    w-full
    p-4
    rounded-lg
    transition="all ease-in-out duration-300"
    cursor-pointer
    flex="~ col"
    decoration="none!"
    bg="zinc-50 dark:zinc-800"
    border="~ solid 2 transparent hover:zinc-200 dark:hover:zinc-600"
    :href="props.href"
  >
    <div flex="1">
      <div flex="~" items-center>
        <div v-if="props.title" text-lg font-semibold flex-1>
          <span>{{ props.title }}</span>
        </div>
        <template v-else>
          <slot name="title" />
        </template>
      </div>
    </div>
    <slot />
  </a>
</template>
