<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue';
import 'asciinema-player/dist/bundle/asciinema-player.css';

const props = defineProps<{
  src: string;
}>()

const playerRef = ref<HTMLDivElement>();
const asciinemaPlayer = ref<any>();

onMounted(async () => {
  if (!window)
    return
  if (!document)
    return
  if (!playerRef.value)
    return

  const AsciinemaPlayer = await import('asciinema-player')
  asciinemaPlayer.value = AsciinemaPlayer.create(props.src, playerRef.value, {
    loop: true,
    autoPlay: true,
    cols: 100,
    rows: 20,
    speed: 3,
  })
})

onUnmounted(() => {
  if (!window)
    return
  if (!document)
    return
  if (asciinemaPlayer.value)
    asciinemaPlayer.value.dispose()
})
</script>

<template>
  <div ref="playerRef" w-full rounded-xl overflow-hidden></div>
</template>
