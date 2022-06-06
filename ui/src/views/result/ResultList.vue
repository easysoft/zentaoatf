<template>
  <div class="result-list">
    <List compact divider>
    <div v-for="item, index in models" :key="index" :class="'list-item-container ' + (item.checked==1?'checked':'')" @click="showDetail(item)" @mouseenter="changeControlIcon($event, index)" @mouseleave="changeControlIcon($event, index)">
        <ListItem
          icon="checkmark-circle"
          class="inline-left"
          iconClass="text-green"
          v-if="item.fail==0"
          :title="item.total != 1 ? item.workspaceName + '(' + item.total + ')' : item.testScriptName"
          trailingTextClass="muted small"
        >
        </ListItem>
        <ListItem
          icon="close-circle"
          class="inline-left"
          iconClass="text-red"
          v-else
          :title="item.total != 1 ? item.workspaceName + '(' + item.total + ')' : item.testScriptName"
          trailingTextClass="muted small"
        />
        <span v-if="item.checked == 0 || item.checked == undefined">{{momentTime(item.startTime, 'hh:mm')}}</span>
        <div v-else>
            <Toolbar
                class="tree-node-toolbar"
                :buttons="toolbarItems"
                @click="clickToolbar($event, item)"
            />
        </div>
    </div>
    </List>
  </div>
</template>

<script setup lang="ts">
import List from '@/components/List.vue';
import ListItem from '@/components/ListItem.vue';
import {StateType} from "@/views/result/store";
import {PaginationConfig, QueryParams} from "@/types/data";

import {useI18n} from "vue-i18n";
import {useRouter} from "vue-router";
import {useStore} from "vuex";
import {computed, onMounted, watch, ref} from "vue";
import {momentUnixDefFormat} from "@/utils/datetime";
import {ZentaoData} from "@/store/zentao";
import Toolbar, { ToolbarItemProps } from '@/components/Toolbar.vue';
import bus from "@/utils/eventBus";
import settings from "@/config/settings";

const { t } = useI18n();
const router = useRouter();

const momentTime = momentUnixDefFormat
const toolbarItems = [
    {hint: 're_exec', icon: 'refresh', key: 'refreshExec'},
    { hint: 'show_detail_log', icon: 'file', key: 'showDetail'},
];
const store = useStore<{ Zentao: ZentaoData, Result: StateType }>();
const models = computed<any[]>(() => store.state.Result.queryResult.result)

const pagination = computed<PaginationConfig>(() => store.state.Result.queryResult.pagination);
const queryParams = ref<QueryParams>({
    keywords: '', enabled: '1', page: pagination.value.page, pageSize: pagination.value.pageSize
});
const currProduct = computed<any>(() => store.state.Zentao.currProduct);

const list = (page: number) => {
    store.dispatch('Result/list', {
    keywords: queryParams.value.keywords,
    enabled: queryParams.value.enabled,
    pageSize: pagination.value.pageSize,
    page: page});
}
list(1);

watch(currProduct, () => {
  list(1);
}, { deep: true })

const clickToolbar = (e, item) => {
    console.log(e, item)
    if(e.key == 'refreshExec'){
        refreshExec(item)
    }else if(e.key == 'showDetail'){
        showDetail(item)
    }
}
const report = computed<any>(() => store.state.Result.detailResult);
const get = async (workspaceId, seq): Promise<void> => {
    await store.dispatch('Result/get', {workspaceId: workspaceId, seq: seq});
}

const refreshExec = async (item): Promise<void> => {
  await get(item.workspaceId, item.seq)
  const testType = report.value.testType;
  if (testType === "func") {
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

const showDetail = (item) => {
    store.dispatch('tabs/open', {
    id: 'result-' + item.no,
    title: item.total != 1 ? item.workspaceName + '(' + item.total + ')' : item.testScriptName,
    type: 'result',
    data: {seq:item.seq, workspaceId: item.workspaceId}
  });
}

const changeControlIcon = (e, index) => {
    for(let i=0; i < models.value.length; i++){
        if(i == index){
            models.value[index].checked = !models.value[index].checked;
        }else{
            models.value[i].checked = false;
        }
    }
}

onMounted(() => {
    console.log("onMounted")
})

</script>

<style scoped>
.list-item-container{
    display: flex;
    align-items: center;
    flex-direction: row;
    justify-content: space-between;
    width: 100%;
}

.icon{
    margin-left: 8px;
    cursor: pointer;
}

.checked{
    background-color: #E2E5E9;
}

.inline-left{
    min-width: 80%;
}
.result-list{
    padding-right: 20px;
}
</style>
