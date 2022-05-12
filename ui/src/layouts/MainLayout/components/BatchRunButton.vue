<template>
  <ButtonGroup>
    <Button class="rounded border primary-pale" icon="run-all" label="执行所有选中" />
    <Button class="rounded border primary-pale padding-0" iconClass="muted" icon="caret-down" iconSize="1em" style="width: 20px" id="batchRunMenuToggle" />
  </ButtonGroup>
  <DropdownMenu
    :items="[
      {id: 'all', name: '执行所有文件'},
      {id: 'last', name: '执行上次'},
      {id: 'select', name: '执行所有选中'},
      {id: 'open', name: '执行当前打开'}
    ]"
    keyName="id"
    checkedKey="select"
    activeKey="select"
    :replaceFields="{title: 'name'}"
    toggle="#batchRunMenuToggle"
    position="bottom-right"
    @click="_handleDropDownMenuClick"
  />
</template>

<script setup lang="ts">
import Button from './Button.vue';
import ButtonGroup from './ButtonGroup.vue';
import DropdownMenu from './DropdownMenu.vue';

function _handleDropDownMenuClick(event) {
    console.log('_handleDropDownMenuClick', event);
}
</script>
