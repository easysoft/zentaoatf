<template>
  <ButtonList :gap="0" class="btn-group">
    <slot />
  </ButtonList>
</template>

<script setup lang="ts">
import ButtonList from './ButtonList.vue';
</script>

<style scoped>
.btn-group :deep(.btn + .btn) {
  margin-left: -1px;
  position: relative;
}
.btn-group :deep(.btn:hover),
.btn-group :deep(.btn:focus),
.btn-group :deep(.btn:active) {
  z-index: 1;
}
.btn-group :deep(.btn:not(:last-child)) {
  border-bottom-right-radius: 0;
  border-top-right-radius: 0;
}
.btn-group :deep(.btn:not(:first-child)) {
  border-bottom-left-radius: 0;
  border-top-left-radius: 0;
}
</style>
