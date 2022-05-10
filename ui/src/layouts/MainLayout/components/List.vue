<template>
  <div class="list" :class="{compact, divider}">
    <template v-if="items">
      <ListItem
        v-for="({key, ...btnProps}) in itemList"
        :key="key"
        v-bind="btnProps"
      />
    </template>
    <slot />
  </div>
</template>

<script setup lang="ts">
import { defineProps, computed } from 'vue';
import ListItem, { ListItemProps } from './ListItem.vue';

const props = defineProps<{
    compact?: boolean,
    divider?: boolean,
    items?: ListItemProps[]
}>();

const itemList = computed(() => {
    if (!props.items) {
        return null;
    }
    return props.items.map((x, i) => ({
        key: i,
        ...x
    }));
});

</script>
