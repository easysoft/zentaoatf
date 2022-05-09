<template>
  <div class="panel">
    <slot name="heading">
      <header
        class="panel-heading state"
        :class="headerClass"
        @click="state.collapsed = !state.collapsed"
      >
        <slot name="header">
          <div class="title" :class="titleClass ?? 'strong'">{{title}}</div>
        </slot>
        <slot name="toolbar">
          <Toolbar v-if="toolbar || $slots['toolbar-buttons']" :buttons="toolbar">
            <slot name="toolbar-buttons"></slot>
          </Toolbar>
        </slot>
      </header>
    </slot>
    <template v-if="!state.collapsed">
      <slot name="body">
        <div class="panel-body" :class="bodyClass">
          <slot></slot>
        </div>
      </slot>
      <slot name="footer"></slot>
    </template>
  </div>
</template>

<script setup lang="ts">
import {defineProps, reactive} from 'vue';
import Button, {ButtonProps} from './Button.vue';
import Toolbar from './Toolbar.vue';

const props = defineProps<{
    title?: string,
    defaultCollapsed?: boolean,
    collapsable?: boolean,
    toolbar?: ButtonProps[],
    headerClass?: string,
    titleClass?: string,
    bodyClass?: string,
}>();

const state = reactive({collapsed: !!props.defaultCollapsed});
</script>

<style scoped>
.panel-heading {
  height: calc(2em + 4px);
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: space-between;
  padding: 0 var(--space-base);
}
.panel-heading >>> .toolbar {
  margin-right: -6px;
}
.panel-body {
  overflow: auto;
  overflow: overlay;
  scroll-behavior: smooth;
}
</style>
