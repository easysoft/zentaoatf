<template>
  <ButtonList :gap="0" class="btn-group">
    <slot />
  </ButtonList>
</template>

<script setup lang="ts">
import ButtonList from './ButtonList.vue';
</script>

<style scoped>
.btn-group >>> .btn + .btn {
  margin-left: -1px;
  position: relative;
}
.btn-group >>> .btn:hover,
.btn-group >>> .btn:focus,
.btn-group >>> .btn:active {
  z-index: 1;
}
.btn-group >>> .btn:not(:last-child) {
  border-bottom-right-radius: 0;
  border-top-right-radius: 0;
}
.btn-group >>> .btn:not(:first-child) {
  border-bottom-left-radius: 0;
  border-top-left-radius: 0;
}
</style>
