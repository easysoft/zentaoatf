<template>
  <Button class="rounded border-light canvas gap-sm" hint="点击查看测试结果">
    <small class="text-muted">上次结果</small>
    <small class="text-yellow">1.3s</small>
    <Icon icon="close-circle" class="text-red space-left" />
    <small class="text-red">4</small>
    <Icon icon="checkmark-circle" class="text-green space-left" />
    <small class="text-green">2</small>
  </Button>
</template>

<script setup lang="ts">
import Button from './Button.vue';
import Icon from './Icon.vue';
</script>
