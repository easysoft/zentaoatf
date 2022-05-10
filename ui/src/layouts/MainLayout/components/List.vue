<template>
  <div class="list" :class="{compact, divider}">
    <slot />
  </div>
</template>

<script setup lang="ts">
import { defineProps } from 'vue';
import ListItem, {ListItemProps} from './ListItem.vue';

defineProps<{
    compact?: boolean,
    divider?: boolean,
    items?: ListItemProps[]
}>();

</script>
