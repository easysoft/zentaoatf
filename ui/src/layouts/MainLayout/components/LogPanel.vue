<template>
  <Panel :title="t('exec_log')" class="log-panel">
    <template #toolbar-buttons>
      <Button class="rounded pure" :hint="t('collapse_all')"
              :icon="logContentExpand ? 'subtract-square-multiple' : 'add-square-multiple'" iconSize="1.4em"
              @click="globalStore.commit('global/setLogContentExpand')"/>
      <Button class="rounded pure" 
        :hint="logPaneMaximized ? t('restore_panel_size') : t('expand_up')"
        :icon="logPaneMaximized ? 'chevron-down' : 'chevron-up'"
        iconSize="1.5em"
        @click="globalStore.commit('global/setLogPaneResized')" />
		<Button class="rounded pure" :hint="t('more_actions')" icon="more-vert" />
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
import { useStore} from 'vuex';
import { StateType } from '@/store/global'
import {computed} from 'vue';
const { t } = useI18n();

const globalStore = useStore<{global: StateType}>()
const logContentExpand = computed<boolean>(() => globalStore.state.global.logContentExpand);
const logPaneMaximized = computed<boolean>(() => globalStore.state.global.logPaneMaximized);

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
