<template>
  <ButtonGroup>
    <Button class="rounded border primary-pale" icon="run-all" label="执行所有选中" />
    <Button class="rounded border primary-pale padding-0" iconClass="muted" icon="caret-down" iconSize="1em" style="width: 20px" id="batchRunMenuToggle" />
  </ButtonGroup>
  <DropdownMenu
    :items="[
      {checked: false, title: '执行所有文件'},
      {checked: false, title: '执行上次'},
      {checked: true, title: '执行所有选中', active: true},
      {checked: false, title: '执行当前打开'}
    ]"
    toggle="#batchRunMenuToggle"
    position="bottom-right"
  />
</template>

<script setup lang="ts">
import Button from './Button.vue';
import ButtonGroup from './ButtonGroup.vue';
import DropdownMenu from './DropdownMenu.vue';
</script>
