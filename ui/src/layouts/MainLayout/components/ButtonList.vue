<template>
  <div class="btn-list" :style="`gap:${gap ?? 8}px`">
    <template v-if="buttonPropsList">
      <Button
        v-for="({key, ...btnProps}) in buttonPropsList"
        :key="key"
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
    buttons?: ButtonProps[] | Record<string, any>[],
    replaceFields?: Record<string, string>,
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
    return props.buttons.map((x, i) => {
        if (props.replaceFields && Button.props) {
            return Object.keys(Button.props).reduce((item, propName) => {
                const replacePropName = props.replaceFields ? props.replaceFields[propName] : null;
                item[propName] = x[typeof replacePropName === 'string' ? replacePropName : propName];
                return item;
            }, {key: i})
        }
        return {
            key: i,
            'class': props.defaultBtnClass,
            iconClass: props.defaultIconClass,
            suffixIconClass: props.defaultSuffixIconClass,
            labelClass: props.defaultLabelIconClass,
            iconSize: props.defaultIconSize,
            suffixIconSize: props.defaultSuffixIconSize,
            ...x
        };
    });
});
</script>

<style scoped>
.btn-list {
  display: flex;
  flex-direction: row;
  align-items: center;
}
</style>
