<template>
  <ButtonGroup class="space-left">
    <Button class="rounded border lighten-16" icon="zentao" iconColor="var(--color-blue)" iconClass="off-off" label="当前禅道站点" suffix-icon="caret-down" id="siteMenuToggle" />
    <Button class="rounded border lighten-16" icon="cube" label="当前产品" suffix-icon="caret-down" id="productMenuToggle" />
  </ButtonGroup>
  <DropdownMenu
    :items="[
      {checked: false, title: 'product 1'},
      {checked: false, title: 'product 2'},
      {checked: true, title: 'product 3', active: true},
      {checked: false, title: 'product 4'}
    ]"
    toggle="#productMenuToggle"
  />
  <DropdownMenu toggle="#siteMenuToggle">
    <div>custom site menu</div>
  </DropdownMenu>
</template>

<script setup lang="ts">
import Button from './Button.vue';
import ButtonGroup from './ButtonGroup.vue';
import DropdownMenu from './DropdownMenu.vue';
</script>
