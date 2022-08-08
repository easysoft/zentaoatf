<template>
<div>
  <Panel v-show="showStatistices" :title="t('exec_statistics')">
    <ResultStatistic :path="scriptPath" />
  </Panel>
  <Panel :title="t('test_result')">
    <ResultList />
  </Panel>
</div>
</template>

<script setup lang="ts">
import Panel from '@/components/Panel.vue';
import ResultList from '@/views/result/ResultList.vue';
import ResultStatistic from '@/views/result/ResultStatistic.vue';
import {useI18n} from "vue-i18n";
import {computed} from "vue";
import {useStore} from "vuex";
import {TabsData} from "@/store/tabs";

const { t } = useI18n();

const store = useStore<{ tabs: TabsData }>();
const showStatistices = computed((): boolean => {
  return store.state.tabs.activeID.indexOf('script-') > -1
});
const scriptPath = computed((): string => {
  if(store.state.tabs.activeID.indexOf('script-') > -1){
    return store.state.tabs.activeID.replace('script-', '');
  }
  return ""
});

</script>

<style scoped>
</style>
