<template>
  <Panel>
    <template #header>
      <ButtonList :gap="2" class="workdir-panel-nav">
        <Button class="rounded pure padding-sm-h" label="工作目录" labelClass="strong" suffix-icon="caret-down"/>
        <Button class="rounded pure padding-sm-h" label="按套件" suffix-icon="caret-down" />
      </ButtonList>
    </template>
    <template #toolbar-buttons>
      <Button class="rounded pure" hint="折叠所有" icon="subtract-square-multiple" iconSize="1.4em" />
      <Button class="rounded pure" hint="批量选择" icon="select-all-on" />
      <Button class="rounded pure" hint="添加工作目录" icon="folder-add" />
      <Button class="rounded pure" hint="更多操作" icon="more-vert" />
    </template>
    <WorkDir />
  </Panel>
</template>

<script setup lang="ts">
import Panel from './Panel.vue';
import Button from './Button.vue';
import ButtonList from './ButtonList.vue';
import WorkDir from './WorkDir.vue';
</script>

<style scoped>
.workdir-panel-nav {
  margin-left: -6px;
}
</style>
