<template>
  <div id="log-list" class="log-list scrollbar-y">
    <pre id="content" class="content">
      <template v-for="(item, index) in wsMsg.out" :key="index">
        {{ void (info = item.info) }}
        {{ void (csKey = info?.key) }}

        <code class="item small"
             :class="[
                 csKey && caseDetail[csKey] ? 'show-detail' : '',

                 csKey ? 'case-item' : '',
                 info?.status === 'start' ? 'case-start' : '',
                 info?.status === 'start' ? 'result-'+caseResult[csKey] : '',

                 info?.status === 'start-task' ? 'strong' : ''
             ]">

          <div class="group">
            <template v-if="info?.status === 'start'">
              <span @click="showDetail(item.info?.key)" class="link state center">
                <Icon v-if="!caseDetail[csKey]" icon="chevron-right" />
                <Icon v-if="caseDetail[csKey]" icon="chevron-down" />
              </span>
            </template>
          </div>

          <div class="sign">
            <Icon v-if="item.msg" icon="circle" />
            <span v-else>&nbsp;</span>
          </div>

          <div class="time">
            <span>{{ item.time }}</span>
          </div>
          <div class="msg-span">
            <span v-html="item.msg"></span>
            <span v-if="info?.status === 'start' && caseResult[csKey]">
              [ {{ t(caseResult[csKey]) }} ]
            </span>
          </div>
        </code>

      </template>
    </pre>
  </div>
</template>

<script setup lang="ts">
import {useI18n} from "vue-i18n";
import {WsMsg} from "@/types/data";
import {genExecInfo, genWorkspaceToScriptsMap} from "@/views/script/service";
import {scroll} from "@/utils/dom";
import {computed, onBeforeUnmount, onMounted, reactive, ref, watch} from "vue";
import {useStore} from "vuex";
import {WebSocketData} from "@/store/websoket";
import {ExecStatus} from "@/store/exec";
import {ProxyData} from "@/store/proxy";
import {WorkspaceData} from "@/store/workspace";
import bus from "@/utils/eventBus";
import settings from "@/config/settings";
import {ZentaoData} from "@/store/zentao";
import {WebSocket} from "@/services/websocket";

import Icon from '@/components/Icon.vue';
import {momentTime} from "@/utils/datetime";
import {isInArray} from "@/utils/array";
import {StateType as GlobalStateType} from "@/store/global";
import { get as getWorkspace, uploadToProxy } from "@/views/workspace/service";
import {mvLog} from "@/views/result/service";
import { key } from "localforage";
const { t } = useI18n();

const store = useStore<{global: GlobalStateType, Zentao: ZentaoData, WebSocket: WebSocketData, Exec: ExecStatus, proxy: ProxyData, workspace: WorkspaceData}>();
const logContentExpand = computed<boolean>(() => store.state.global.logContentExpand);

store.dispatch("proxy/fetchProxies");
const currSite = computed<any>(() => store.state.Zentao.currSite);
const currProduct = computed<any>(() => store.state.Zentao.currProduct);
const wsStatus = computed<any>(() => store.state.WebSocket.connStatus);
const isRunning = computed<any>(() => store.state.Exec.isRunning);
const proxyMap = computed<any>(() => store.state.proxy.proxyMap);
const currentWorkspace = ref({} as any);

const cachedExecData = ref({})
const caseCount = ref(1)
const caseResult = ref({})
const caseDetail = ref({})
const realPathMap = ref({})
const lastWsMsg = ref({})

// websocket
let wsMsg = reactive({in: '', out: [] as any[]});

const showDetail = (key) => {
  console.log('showDetail', key)
  caseDetail.value[key] = !caseDetail.value[key]
}

const hideAllDetailOrNot = (val) => {
  console.log('hideAllDetailOrNot')
  Object.keys(caseDetail.value).forEach((key => {
    caseDetail.value[key] = val;
  }))
  console.log(caseDetail.value)
}

watch(logContentExpand, () => {
  console.log('watch logContentExpand', logContentExpand.value)
  hideAllDetailOrNot(logContentExpand.value)
}, {deep: true})

const onWebsocketMsgEvent = async (data: any) => {
  console.log('WebsocketMsgEvent in ExecLog', data.msg)

  let item = JSON.parse(data.msg) as WsMsg

  if ('isRunning' in wsMsg) {
    console.log('change isRunning to ', item.isRunning)
    store.dispatch('Exec/setRunning', item.isRunning)
  }

  if (item.info?.status === 'start') {
    const key = item.info.key + '-' + caseCount.value
    caseDetail.value[key] = logContentExpand.value
  }

  item = genExecInfo(item, caseCount.value)
  if (item.info && item.info.key && isInArray(item.info.status, ['pass', 'fail', 'skip'])) { // set case result
    store.dispatch('Result/list', {
        keywords: '',
        enabled: 1,
        pageSize: 10,
        page: 1
        });
    caseResult.value[item.info.key] = item.info.status
  }

  Object.keys(realPathMap.value).forEach(key => {
    item.msg = item.msg.replace(key, realPathMap.value[key])
  })

  if(currentWorkspace.value.proxy_id > 0){
    await downloadLog(item);
    if( item.msg == ""){
      return;
    }
  }
  if(JSON.stringify(lastWsMsg.value) === JSON.stringify(item)){
    return;
  }
  lastWsMsg.value = item;
  wsMsg.out.push(item)
  scroll('log-list')
}

const downloadLog = async (item) => {
  const msg = item.msg;
  if(msg.indexOf('Report') !== 0 && msg.indexOf('报告') !== 0){
    return;
  }
  let logPath = '';
  if(msg.indexOf('Report') === 0){
    logPath = msg.replace('Report', '')
    logPath = logPath.replace('.', '')
  }else if(msg.indexOf('报告') === 0){
    logPath = msg.replace('报告', '')
    logPath = logPath.replace('。', '')
  }
  mvLog({file: logPath, workspaceId: currentWorkspace.value.id}).then(resp => {
    if(resp.code === 0){
      item.msg = item.msg.replace(logPath, resp.data)
      let emptMsg = {...item}
      emptMsg.msg = '';
      emptMsg.isRunning = false;
      emptMsg.time = '';
      wsMsg.out.push(emptMsg)
      store.dispatch('Result/list', {
        keywords: '',
        enabled: 1,
        pageSize: 10,
        page: 1
        });
    }
  })
}

onMounted(() => {
  console.log('onMounted')
  bus.on(settings.eventExec, exec);
  bus.on(settings.eventWebSocketMsg, onWebsocketMsgEvent);
})
onBeforeUnmount( () => {
  console.log('onBeforeUnmount')
  bus.off(settings.eventExec, exec);
  bus.off(settings.eventWebSocketMsg, onWebsocketMsgEvent);
})

const exec = async (data: any) => {
  console.log('exec', data)

  let execType = data.execType
  if (execType === 'previous') {
    data = cachedExecData.value
    execType = data.execType
  } else {
    cachedExecData.value = data
  }

  if (execType === 'ztf' && (!data.scripts || data.scripts.length === 0)) {
    const msgCancel = {
      msg: `<span class="strong">`+t('case_num_empty')+`</span>`,
      time: momentTime(new Date())}
    wsMsg.out.push(msgCancel)
    return
  }

  caseCount.value++
  let msg = {} as any
  let workspaceId = 0;
  if (execType === 'ztf') {
    const scripts = data.scripts
    const sets = genWorkspaceToScriptsMap(scripts)
    if(!workspaceId){
        if(scripts.length > 0){
            workspaceId = scripts[0].workspaceId
        }
    }
    console.log('===', sets)
    msg = {act: 'execCase', testSets: sets}

  } else if (execType === 'unit') {
    const set = {workspaceId: data.id, workspaceType: data.type, cmd: data.cmd, submitResult: data.submitResult, name: data.name}
    workspaceId = workspaceId == 0 ? data.id : workspaceId;

    msg = {act: 'execUnit', testSets: [set]}

  } else if (execType === 'stop') {
    msg = {act: 'execStop'}
  }

  console.log('exec testing', msg)
  currentWorkspace.value = {};
  let workspaceInfo = {};
  if(workspaceId>0){
    workspaceInfo = await getWorkspace(workspaceId)
  }
  if (msg.testSets !== undefined && workspaceInfo.data != undefined && workspaceInfo.data.proxy_id > 0) {
    currentWorkspace.value = workspaceInfo.data;
    const resp = await uploadToProxy(msg.testSets);
    const testSetsMap = resp.data;
    realPathMap.value = testSetsMap;
    let casesMap = {};
    var keys = Object.keys(testSetsMap);
    keys.forEach((val) => {
        casesMap[testSetsMap[val]] = val;
    });
    msg.testSets.forEach((set, setIndex) => {
        if(set.cases != undefined){
            set.cases.forEach((casePath, caseIndex) => {
                msg.testSets[setIndex].cases[caseIndex] = casesMap[casePath];
            }
        );
        }
    })
  }
  WebSocket.sentMsg(settings.webSocketRoom, JSON.stringify(msg), proxyMap.value[workspaceInfo?.data?.proxy_id])
}

const logLevel = ref('result')
const logStatus = ref('')

</script>

<style lang="less">
.log-list {
  height: 100%;
  .result-pass {
    color: var(--color-green)
  }
  .result-fail {
    color: var(--color-red)
  }
  .result-skip {
    color: var(--color-gray)
  }

  .content {
    white-space: normal;
    .item {
      line-height: 18px;
      &.case-item {
        &:not(.case-start) {
          display: none !important;
        }
      }
      &.show-detail:not(.case-start) {
        display: flex !important;
      }
      &:hover {
        background-color: var(--color-darken-1);
      }

      .group {
        width: 16px;
        font-size: 13px;
        text-align: center;
        .link {
          //position: relative;
          //top: 2px;
        }
      }
      .sign {
        margin: auto;

        width: 20px;
        font-size: 6px;
        text-align: center;
      }
      .time {
        width: 80px;
      }
    }
  }
}
</style>
