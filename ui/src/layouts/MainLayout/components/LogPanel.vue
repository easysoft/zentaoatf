<template>
  <Panel :title="t('exec_log')" class="log-panel">
    <template #toolbar-buttons>
      <Button class="rounded pure" :hint="t('more_actions')" icon="subtract-square-multiple" iconSize="1.4em" />
      <Button class="rounded pure" :hint="t('collapse_all')" icon="chevron-up" iconSize="1.5em" />
      <Button class="rounded pure" :hint="t('expand_up')" icon="more-vert" />
    </template>

    <LogList />

  </Panel>
</template>

<script setup lang="ts">
import Panel from './Panel.vue';
import Button from './Button.vue';
import WorkDir from './WorkDir.vue';
import LogList from './LogList.vue';
import {useI18n} from "vue-i18n";
const { t } = useI18n();



</script>

<style lang="less">
.log-panel {
  height: 100%;

  .panel-heading {

  }

  .panel-body {
    height: calc(100% - 30px);
    overflow-y: auto;
  }
}
</style>
