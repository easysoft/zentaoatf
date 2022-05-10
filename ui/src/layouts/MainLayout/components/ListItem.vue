<template>
  <div
    class="list-item"
    :class="{disabled, divider, state: !disabled, compact, active}"
    @click="disabled ? null: _handleClick"
  >
    <slot name="leading" />
    <Icon
      class="list-item-icon"
      v-if="icon"
      :icon="icon"
      :class="iconClass"
      :color="iconColor"
    />
    <slot name="content">
      <div class="list-item-content">
        <div
          v-if="title"
          class="list-item-title"
          :class="titleClass"
        >
          {{title}}
        </div>
        <div
          v-if="subtitle"
          class="list-item-subtitle"
          :class="subtitleClass"
        >
          {{subtitle}}
        </div>
        <slot />
      </div>
    </slot>
    <div v-if="trailingText"
      class="list-item-trailing-text"
      :class="trailingTextClass"
    >{{trailingText}}</div>
    <Icon
      v-if="trailingIcon"
      class="list-item-trailing-icon"
      :icon="trailingIcon"
      :class="trailingIconClass"
      :color="trailingIconColor"
    />
    <Icon
      v-if="trailingAngle"
      icon="chevron-right"
      class="list-item-trailing-angle"
    />
    <slot name="trailing" />
  </div>
</template>

<script setup lang="ts">
import { defineProps, defineEmits } from 'vue';
import Icon from './Icon.vue';

export interface ListItemProps {
    disabled?: boolean,
    active?: boolean,
    divider?: boolean,
    compact?: boolean,
    icon?: string,
    iconColor?: string,
    iconSize?: string | number,
    iconClass?: string,
    trailingIcon?: string,
    trailingIconColor?: string,
    trailingIconClass?: string,
    trailingIconSize?: string | number,
    title?: string,
    titleClass?: string,
    subtitle?: string,
    subtitleClass?: string,
    trailingText?: string,
    trailingTextClass?: string,
    trailingAngle?: boolean,
    url?: string,
    click?: (event: Event) => void
}

const props = defineProps<ListItemProps>();

const emit = defineEmits<{(event: 'click', e: Event) : void}>();

function _handleClick(event) {
    if (props.disabled) {
        return;
    }

    if (typeof props.url === 'string' && props.url.length) {
        window.open(props.url);
    }
    if (typeof props.click === 'function') {
        props.click(event);
    }

    emit('click', event);
}
</script>

<style scoped>
.list-item {
  display: flex;
  flex-direction: row;
  align-items: center;
  padding: var(--space-sm) var(--space-base);
  gap: var(--space-base);
}
.compact > .list-item,
.list-item.compact {
  gap: var(--space-sm);
}
.divider > .list-item + .list-item,
.list-item.divider + .list-item {
  border-top: 1px solid var(--color-darken-1)
}
.list-item-subtitle,
.compact > .list-item .list-item-title,
.list-item.compact .list-item-title {
  font-size: 0.9230769231em;
}
.list-item-content {
  flex: auto;
}
</style>
