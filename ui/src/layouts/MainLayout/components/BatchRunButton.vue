<template>
  <ButtonGroup>
    <Button class="rounded border primary-pale" icon="run-all" label="执行所有选中" />
    <Button class="rounded border primary-pale padding-0" iconClass="muted" icon="caret-down" iconSize="1em" style="width: 20px" />
  </ButtonGroup>
</template>

<script setup lang="ts">
import Button from './Button.vue';
import ButtonGroup from './ButtonGroup.vue';
</script>
