<template>
  <div class="log-list scrollbar-y">
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
import bus from "@/utils/eventBus";
import settings from "@/config/settings";
import {ZentaoData} from "@/store/zentao";
import {WebSocket} from "@/services/websocket";

import Icon from '@/components/Icon.vue';
import {momentTime} from "@/utils/datetime";
import {isInArray} from "@/utils/array";
import {StateType as GlobalStateType} from "@/store/global";
const { t } = useI18n();

const store = useStore<{global: GlobalStateType, Zentao: ZentaoData, WebSocket: WebSocketData, Exec: ExecStatus}>();
const logContentExpand = computed<boolean>(() => store.state.global.logContentExpand);

const currSite = computed<any>(() => store.state.Zentao.currSite);
const currProduct = computed<any>(() => store.state.Zentao.currProduct);
const wsStatus = computed<any>(() => store.state.WebSocket.connStatus);
const isRunning = computed<any>(() => store.state.Exec.isRunning);

const cachedExecData = ref({})
const caseCount = ref(1)
const caseResult = ref({})
const caseDetail = ref({})

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

const onWebsocketMsgEvent = (data: any) => {
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

  wsMsg.out.push(item)
  scroll('content')
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

const exec = (data: any) => {
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
  let msg = {}
  if (execType === 'ztf') {
    const scripts = data.scripts
    const sets = genWorkspaceToScriptsMap(scripts)
    msg = {act: 'execCase', testSets: sets}

  } else if (execType === 'unit') {
    const set = {workspaceId: data.id, workspaceType: data.type, cmd: data.cmd, submitResult: data.submitResult}

    msg = {act: 'execUnit', testSets: [set]}

  } else if (execType === 'stop') {
    msg = {act: 'execStop'}
  }

  console.log('exec testing', msg)
  WebSocket.sentMsg(settings.webSocketRoom, JSON.stringify(msg))
}

const logLevel = ref('result')
const logStatus = ref('')

</script>

<style lang="less">
.log-list {
  .result-pass {
    color: var(--color-green)
  }

  .result-fail {
    color: var(--color-red)
  }

  .content {
    white-space: normal;
    .item {
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
          position: relative;
          top: 2px;
        }
      }
      .sign {
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