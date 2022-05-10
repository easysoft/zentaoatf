<template>
  <Icon class="icon" :icon="icon" :width="size ?? '1.2em'" :height="size ?? '1.2em'" :color="color" />
</template>

<script setup lang="ts">
import { Icon } from '@iconify/vue/dist/offline';
import { defineProps } from 'vue';
import './icons';

defineProps<{
    icon: string,
    size?: number | string,
    color?: string,
}>();
</script>
