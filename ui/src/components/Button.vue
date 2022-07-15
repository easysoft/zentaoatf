<template>
  <button
    :class="`btn ${disabled ? 'disabled' : 'state'}${isOnlyIcon ? ' btn-only-icon' : ''}${size ? ` btn-size-${size}` : ''}${active ? ' active' : ''}`"
    type="button"
    :title="hint"
    @click.stop="_handleClick"
  >
    <Icon v-if="icon" class="btn-icon" :class="iconClass ?? (isOnlyIcon ? '' : 'muted')" :icon="icon" :color="iconColor" :size="iconSize ?? '1.2em'" />
    <slot>
      <span class="btn-label" :class="labelClass" v-if="label">{{label}}</span>
    </slot>
    <Icon v-if="suffixIcon" class="btn-suffix-icon" :class="suffixIconClass ?? (isOnlyIcon ? '' : 'muted')" :icon="suffixIcon" :color="suffixIconColor" :size="suffixIconSize" @click="_handleSuffixClick" />
  </button>
</template>

<script setup lang="ts">
import {defineProps, computed, useSlots, defineEmits, useAttrs} from 'vue';
import Icon from './Icon.vue';import {useI18n} from "vue-i18n";

export interface ButtonProps {
    icon?: string,
    iconColor?: string,
    iconSize?: string | number,
    iconClass?: string,
    suffixIcon?: string,
    suffixIconColor?: string,
    suffixIconClass?: string,
    suffixIconSize?: string | number,
    label?: string,
    labelClass?: string,
    size?: '' | 'sm' | 'lg',
    hint?: string,
    disabled?: boolean,
    active?: boolean,
}

const { t } = useI18n();

const props = defineProps<ButtonProps>();

const slots = useSlots();

const isOnlyIcon = computed(() => (!props.label && !slots.default && props.icon && !props.suffixIcon));

const emit = defineEmits<{
    (type: 'click', event: {originalEvent: Event, key: string | number | symbol | null}) : void,
    (type: 'suffixClick', event: {originalEvent: Event, key: string | number | symbol | null}) : void,
    }>();
const attrs = useAttrs();

function _handleClick(originalEvent) {
    if (props.disabled) {
        return;
    }

    const event = {originalEvent, key: attrs['data-key'] as (string | number | symbol | null)};
    emit('click', event);
}

function _handleSuffixClick(originalEvent) {
    if (props.disabled) {
        return;
    }

    const event = {originalEvent, key: attrs['data-key'] as (string | number | symbol | null)};
    emit('suffixClick', event);
}
</script>

<style scoped>
.btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  height: calc(2em + 2px);
  padding: 0 .7em;
  border-style: solid;
  border-width: 1px;
}
.btn-icon {
  margin-right: 0.4em;
  margin-left: -0.2em;
}
.btn-suffix-icon {
  margin-left: 0.4em;
  margin-right: -0.2em;
}
.btn-only-icon {
  padding: 0;
  width: calc(2em + 2px);
}
.btn-only-icon > .btn-icon {
  margin: 0;
}
.btn-size-sm {
  height: calc(2em - 2px);
}
.btn-size-sm.btn-only-icon {
  width: calc(2em - 2px);
}
</style>
