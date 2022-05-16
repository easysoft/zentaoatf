<template>
  <ButtonList
    class="toolbar"
    :gap="gap"
    :buttons="items"
    :defaultBtnClass="defaultBtnClass"
    :defaultIconSize="defaultIconSize"
    @click="$emit('click', $event)"
  >
    <slot />
  </ButtonList>
</template>

<script setup lang="ts">
import { defineProps, withDefaults } from 'vue';
import ButtonList, {ButtonListItemProps} from './ButtonList.vue';
import {useI18n} from "vue-i18n";
const { t } = useI18n();

export type ToolbarItemProps = ButtonListItemProps;

withDefaults(defineProps<{
    items?: ToolbarItemProps[] | Record<string, any>[],
    defaultBtnClass?: string,
    defaultIconSize?: string | number,
    gap?: number
}>(), {
    defaultBtnClass: 'rounded pure',
    defaultIconSize: '1.4em',
    gap: 0
});

</script>

<style scoped>
.toolbar :deep(.btn) {
  height: calc(2em);
  font-size: var(--text-size-sm);
}
.toolbar :deep(.btn + .btn) {
  margin-left: -2px;
}
.toolbar :deep(.btn-only-icon) {
  width: calc(2em);
}
.toolbar :deep(.btn > .icon) ,
.toolbar :deep(.btn > .btn-label) {
  opacity: .6;
}
.toolbar :deep(.btn.active) {
  color: var(--color-primary);
  background-color: var(--color-darken-1);
}
.toolbar :deep(.btn:hover > .icon) ,
.toolbar :deep(.btn:focus > .icon) ,
.toolbar :deep(.btn:active > .icon) ,
.toolbar :deep(.btn:hover > .btn-label) ,
.toolbar :deep(.btn:focus > .btn-label) ,
.toolbar :deep(.btn:active > .btn-label) {
  opacity: 1;
}
</style>
