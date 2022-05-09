<template>
  <div class="btn-list" :style="`gap:${gap ?? 8}px`">
    <template v-if="buttonPropsList">
      <Button
        v-for="({name, ...btnProps}) in buttonPropsList"
        :key="name"
        v-bind="btnProps"
      />
    </template>
    <slot />
  </div>
</template>

<script setup lang="ts">
import {defineProps, computed, useSlots} from 'vue';
import Button, {ButtonProps} from './Button.vue';

const props = defineProps<{
    buttons?: ButtonProps[],
    defaultBtnClass?: string,
    defaultIconClass?: string,
    defaultIconSize?: string | number,
    defaultSuffixIconClass?: string,
    defaultSuffixIconSize?: string | number,
    defaultLabelIconClass?: string,
    gap?: number
}>();

const buttonPropsList = computed(() => {
    if (!props.buttons) {
        return null;
    }
    return props.buttons.map((x, i) => ({
        name: i,
        'class': props.defaultBtnClass,
        iconClass: props.defaultIconClass,
        suffixIconClass: props.defaultSuffixIconClass,
        labelClass: props.defaultLabelIconClass,
        iconSize: props.defaultIconSize,
        suffixIconSize: props.defaultSuffixIconSize,
        ...x
    }));
});
</script>

<style scoped>
.btn-list {
  display: flex;
  flex-direction: row;
  align-items: center;
}
</style>
