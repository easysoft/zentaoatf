<template>
  <div class="dropdown">
    <div
      class="dropdown-toggle"
      :class="toggleClass"
      :style="toggleStyle"
    >
      <slot />
    </div>
    <DropdownMenu
      :toggle="toggle"
      :class="menuClass"
      :style="menuStyle"
      :position="menuPosition"
      :items="items"
      :checkedKey="checkedKey"
      :activeKey="activeKey"
      :keyName="keyName"
      :replaceFields="replaceFields"
      :listClass="listClass"
      :listCompact="listCompact"
      :listDivider="listDivider"
    >
      <slot name="menu" />
    </DropdownMenu>
  </div>
</template>

<script setup lang="ts">
import { defineProps } from 'vue';
import DropdownMenu, {DropdownMenuPosition} from './DropdownMenu.vue';
import {ListItemProps, ListItemKey} from './ListItem.vue';

defineProps<{
    toggle?: object | string,
    items?: ListItemProps[] | Record<string, any>[],
    keyName?: string,
    checkedKey?: ListItemKey,
    activeKey?: ListItemKey,
    replaceFields?: Record<string, string>, // {title: 'name'}
    listClass?: string,
    listCompact?: boolean,
    listDivider?: boolean,
    showOnHover?: boolean,
    toggleClass?: string,
    toggleStyle?: object | string,
    menuClass?: string,
    menuStyle?: object | string,
    menuPosition?: DropdownMenuPosition
}>();
</script>

<style scoped>
</style>
