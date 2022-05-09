<template>
  <ButtonList :gap="0" class="toolbar">
    <slot />
  </ButtonList>
</template>

<script setup lang="ts">
import ButtonList from './ButtonList.vue';
</script>

<style scoped>
.toolbar >>> .btn {
  height: calc(2em);
  font-size: var(--text-size-sm);
}
.toolbar >>> .btn + .btn {
  margin-left: -2px;
}
.toolbar >>> .btn-only-icon {
  width: calc(2em);
}
.toolbar >>> .btn > .icon,
.toolbar >>> .btn > .btn-label {
  opacity: .6;
}
.toolbar >>> .btn.active {
  color: var(--color-primary);
  background-color: var(--color-darken-1);
}
.toolbar >>> .btn:hover > .icon,
.toolbar >>> .btn:focus > .icon,
.toolbar >>> .btn:active > .icon,
.toolbar >>> .btn:hover > .btn-label,
.toolbar >>> .btn:focus > .btn-label,
.toolbar >>> .btn:active > .btn-label {
  opacity: 1;
}
</style>
