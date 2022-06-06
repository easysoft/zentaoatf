<template>
  <div class="tabs-nav">
    <nav class="tabs-nav-list scrollbar-hover" ref="navListRef">
      <template v-for="item in items" :key="item.id">
        <div
          class="tabs-nav-item"
          :class="{active: item.id === activeID, changed: item.changed}"
          @click.stop="emit('click', item)"
        >
          <Icon :icon="_getItemIcon(item)" />
          <span class="tabs-nav-title" :title="item.title">{{item.titleFunc ? item.titleFunc() : item.title}}</span>
          <Icon v-if="item.readonly" icon="lock-closed" class="muted" />
          <div class="tabs-nav-close state rounded">
            <Icon icon="close" @click.stop="emit('close', item)" />
          </div>
        </div>
      </template>
    </nav>
    <Toolbar
      v-if="toolbarItems"
      class="tabs-nav-toolbar"
      :buttons="toolbarItems"
      @click="emit('clickToolbar', $event)"
    />
  </div>
</template>

<script setup lang="ts">
import { defineProps, defineEmits, watch, nextTick, ref } from 'vue';
import Icon from './Icon.vue';
import Toolbar, {ToolbarItemProps} from './Toolbar.vue';
import {useI18n} from "vue-i18n";

export interface TabNavItem {
  id: string;
    title: string;
    titleFunc: TitleFunc;
    icon?: string;
    type?: 'script' | 'result' | 'settings' | 'sites';
    changed?: boolean;
    readonly?: boolean;
}

export interface TitleFunc {():string;}

const { t } = useI18n();

const navListRef = ref<HTMLDivElement>();

const props = defineProps<{
    items: TabNavItem[],
    activeID?: string,
    toolbarItems?: ToolbarItemProps[],
}>();


const emit = defineEmits<{
    (type: 'click', event: TabNavItem) : void,
    (type: 'close', event: TabNavItem) : void,
    (type: 'clickToolbar', event: any) : void,
}>();

function _getItemIcon(item: TabNavItem) {
    if (item.icon) {
        return item.icon;
    }

    if (item.type) {
        if (item.type === 'script') {
            return 'file-text';
        }
        if (item.type === 'result') {
            return 'text-box-search';
        }
        if (item.type === 'settings') {
            return 'settings'
        }
        if (item.type === 'sites') {
            return 'globe'
        }
    }
    return 'file';
}

watch(() => props.activeID, () => {
    nextTick(() => {
        if (!navListRef.value) {
            return;
        }
        const activeNavItem = navListRef.value.querySelector('.tabs-nav-item.active');
        if (activeNavItem) {
            if (typeof activeNavItem['scrollIntoViewIfNeeded'] === 'function') {
                activeNavItem['scrollIntoViewIfNeeded']({behavior: 'smooth'});
            } else {
                activeNavItem.scrollIntoView({behavior: 'smooth'});
            }
        }
    });
});
</script>

<style scoped>
.tabs-nav {
  display: flex;
  flex-direction: row;
  flex-wrap: nowrap;
  justify-content: space-between;
}
.tabs-nav-list {
  display: flex;
  flex-direction: row;
  align-items: center;
  flex-wrap: nowrap;
  overflow-x: auto;
  overflow-x: overlay;
  height: var(--tabs-nav-height, 30px);
}
.tabs-nav-list::-webkit-scrollbar {
  height: 4px;
}
.tabs-nav-item {
  display: flex;
  align-items: center;
  flex-direction: row;
  height: var(--tabs-nav-height, 30px);
  padding-left: var(--space-base);
  gap: var(--space-sm);
  padding-right: var(--tabs-nav-height, 30px);
  background-color: var(--color-darken-1);
  position: relative;
  cursor: pointer;
}
.tabs-nav-item + .tabs-nav-item {
  border-left: 1px solid var(--color-canvas);
}
.tabs-nav-item:hover {
  background-color: var(--color-darken-2);
}
.tabs-nav-item.active {
  background-color: var(--color-canvas);;
}

.tabs-nav-title {
  display: block;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 120px;
  opacity: .7;
}
.tabs-nav-item.active .tabs-nav-title,
.tabs-nav-item:hover .tabs-nav-title {
  opacity: 1;
}
.tabs-nav-close {
  position: absolute;
  width: calc(var(--tabs-nav-height, 30px) - 2 * var(--space-sm));
  height: calc(var(--tabs-nav-height, 30px) - 2 * var(--space-sm));
  top: var(--space-sm);
  right: var(--space-sm);
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0;
  transition: opacity .2s;
}
.tabs-nav-item:hover .tabs-nav-close {
  opacity: .7;
}
.tabs-nav-close:hover {
  opacity: 1;
}
.tabs-nav-item.changed::before {
  content: ' ';
  display: block;
  position: absolute;
  width: calc(var(--tabs-nav-height, 30px) / 3);
  height: calc(var(--tabs-nav-height, 30px) / 3);
  background-color: var(--color-primary);
  right: calc(var(--tabs-nav-height, 30px) / 3);
  top: calc(var(--tabs-nav-height, 30px) / 3);
  border-radius: 50%;
}
.tabs-nav-item.changed:hover::before {
  opacity: 0;
}
.tabs-nav-toolbar {
  padding-right: var(--space-sm);
  padding-left: var(--space-sm);
}
</style>
