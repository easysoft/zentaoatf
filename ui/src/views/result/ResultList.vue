<template>
  <List compact class="result-list">
    <ListItem
      v-for="item in results"
      :key="item.no"
      class="result-list-item"
      :title="item.displayName"
      titleClass="text-ellipsis"
      :trailingText="momentUnixFormat(item.startTime, 'HH:mm')"
      trailingTextClass="muted small"
      @click="showDetail(item)"
    >

      <template #leading>
        <Icon v-if="item.fail" icon="close-circle" class="flex-none text-red" />
        <Icon v-else icon="checkmark-circle" class="flex-none text-green" />
      </template>
      <template #trailing>
        <Button
          :hint="t('re_exec')"
          icon="refresh"
          class="pure rounded text-primary result-list-item-btn"
          @click="refreshExec(item)"
        />
      </template>
    </ListItem>
  </List>
</template>

<script setup lang="ts">
import {useStore} from "vuex";
import {computed} from "vue";
import {useI18n} from "vue-i18n";
import List from '@/components/List.vue';
import ListItem from '@/components/ListItem.vue';
import Icon from '@/components/Icon.vue';
import {StateType} from "@/views/result/store";
import {momentUnixFormat} from "@/utils/datetime";
import {ZentaoData} from "@/store/zentao";
import Button from '@/components/Button.vue';
import bus from "@/utils/eventBus";
import settings from "@/config/settings";
import useResultList from '@/hooks/use-result-list';

const { t } = useI18n();

const {results} = useResultList();
const store = useStore<{ Zentao: ZentaoData, Result: StateType }>();
const report = computed<any>(() => store.state.Result.detailResult);

const refreshExec = async (item): Promise<void> => {
  await store.dispatch('Result/get', {workspaceId: item.workspaceId, seq: item.seq});
  const testType = report.value.testType;
  if (testType === "func") {
    const getCaseIdsInReport = (reportVal) => {
      const allCases: object[] = [];
      const failedCases: object[] = [];

      reportVal.funcResult.forEach(cs => {
        const item = {path: cs.path, workspaceId: reportVal.workspaceId}
        allCases.push(item)
        if (cs.status === 'fail') failedCases.push(item)
      })

      return {all: allCases, fail: failedCases}
    }
    const caseMap = getCaseIdsInReport(report.value)
    const cases = caseMap['all']
    bus.emit(settings.eventExec, {execType: 'ztf', scripts: cases});

  } else if (testType === "unit") {
    const data = {
      execType: 'unit',
      cmd: report.value.testCommand,
      id: report.value.workspaceId,
      type: report.value.workspaceType,
      submitResult: report.value.submitResult,
    }
    console.log(data)
    bus.emit(settings.eventExec, data);
  }
};

const showDetail = (item) => {
    store.dispatch('tabs/open', {
    id: 'result-' + item.no,
    title: item.displayName,
    type: 'result',
    data: {seq: item.seq, workspaceId: item.workspaceId}
  });
}
</script>

<style scoped>
.result-list-item-btn {display: none}
.result-list-item {min-height: 28px; padding-top: 0; padding-bottom: 0;}
.result-list-item:hover .result-list-item-btn {display: flex;}
.result-list-item:hover:deep(.list-item-trailing-text) {display: none;}
</style>
