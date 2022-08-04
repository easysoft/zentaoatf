<template>
  <div class="panel">
    <slot name="heading">
      <header
        class="panel-heading"
        :class="[headerClass, collapsable ? 'state' : '']"
        @click="collapsable ? _toggle : null"
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
    <template v-if="!state.collapsed || !collapsable">
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
import {ButtonProps} from './Button.vue';
import Toolbar from './Toolbar.vue';
import {useI18n} from "vue-i18n";
const { t } = useI18n();

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

function _toggle(toggle) {
    if (toggle === undefined) {
        toggle = !state.collapsed;
    }
    state.collapsed = toggle;
}
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
.panel-heading :deep(.toolbar) {
  margin-right: -6px;
}
.panel-body {
  overflow: auto;
  overflow: overlay;
  scroll-behavior: smooth;
}
</style>
