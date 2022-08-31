<template>
  <Panel :title="t('exec_log')" class="log-panel">
    <template #toolbar-buttons>
      <Button class="rounded pure" :hint="t('clear')"
              icon="clear" iconSize="1.4em"
              @click="bus.emit(settings.eventClearWebSocketMsg);"/>
      <Button class="rounded pure" :hint="t('collapse_all')"
              :icon="logContentExpand ? 'subtract-square-multiple' : 'add-square-multiple'" iconSize="1.4em"
              @click="store.commit('global/setLogContentExpand')"/>
      <Button class="rounded pure"
        :hint="logPaneMaximized ? t('restore_panel_size') : t('expand_up')"
        :icon="logPaneMaximized ? 'chevron-down' : 'chevron-up'"
        iconSize="1.5em"
        @click="store.commit('global/setLogPaneResized')" />
    </template>
    <LogList />
  </Panel>
</template>

<script setup lang="ts">
import Panel from '@/components/Panel.vue';
import Button from '@/components/Button.vue';
import LogList from './LogList.vue';
import {useI18n} from "vue-i18n";
import { useStore} from 'vuex';
import { StateType } from '@/store/global'
import {computed} from 'vue';
import bus from "@/utils/eventBus";
import settings from "@/config/settings";
const { t } = useI18n();

const store = useStore<{global: StateType}>()
const logContentExpand = computed<boolean>(() => store.state.global.logContentExpand);
const logPaneMaximized = computed<boolean>(() => store.state.global.logPaneMaximized);

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
