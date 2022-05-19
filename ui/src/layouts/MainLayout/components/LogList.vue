<template>
  <div class="log-list padding muted">

    <div id="content" class="content">
      <template v-for="(item, index) in wsMsg.out" :key="index">
        {{ void (info = item.info) }}
        {{ void (csKey = info?.key) }}

        <div class="item"
             :class="[
                 csKey && caseDetail[csKey] ? 'show-detail' : '',

                 csKey ? 'case-item' : '',
                 info?.status === 'start' ? 'case-start' : '',
                 info?.status === 'start' ? 'result-'+caseResult[csKey] : '',

                 info?.status === 'start-task' ? 'z-border' : ''
             ]">

          <div class="group">
            <template v-if="info?.status === 'start'">
              <span @click="showDetail(item.info?.key)" class="link">
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
        </div>

      </template>
    </div>

  </div>
</template>

<script setup lang="ts">
import {useI18n} from "vue-i18n";
import {WsMsg} from "@/types/data";
import {genExecInfo, genWorkspaceToScriptsMap} from "@/views/script/service";
import {scroll} from "@/utils/dom";
import {computed, onBeforeUnmount, onMounted, reactive, ref} from "vue";
import {useStore} from "vuex";
import {WebSocketData} from "@/store/websoket";
import {ExecStatus} from "@/store/exec";
import bus from "@/utils/eventBus";
import settings from "@/config/settings";
import {ZentaoData} from "@/store/zentao";
import {WebSocket} from "@/services/websocket";

import Icon from './Icon.vue';
import {momentTime} from "@/utils/datetime";
import {isInArray} from "@/utils/array";
const { t } = useI18n();

const zentaoStore = useStore<{ Zentao: ZentaoData }>();
const currSite = computed<any>(() => zentaoStore.state.Zentao.currSite);
const currProduct = computed<any>(() => zentaoStore.state.Zentao.currProduct);

const websocketStore = useStore<{ WebSocket: WebSocketData }>();
const wsStatus = computed<any>(() => websocketStore.state.WebSocket.connStatus);

const execStore = useStore<{ Exec: ExecStatus }>();
const isRunning = computed<any>(() => execStore.state.Exec.isRunning);

const caseCount = ref(1)
const caseResult = ref({})
const caseDetail = ref({})

// websocket
let wsMsg = reactive({in: '', out: [] as any[]});

const showDetail = (key) => {
  console.log('showDetail', key)
  caseDetail.value[key] = !caseDetail.value[key]
}

const onWebsocketMsgEvent = (data: any) => {
  console.log('WebsocketMsgEvent in ExecLog', data.msg)

  let item = JSON.parse(data.msg) as WsMsg

  if ('isRunning' in wsMsg) {
    console.log('change isRunning to ', item.isRunning)
    execStore.dispatch('Exec/setRunning', item.isRunning)
  }

  item = genExecInfo(item, caseCount.value)
  if (item.info && item.info.key && isInArray(item.info.status, ['pass', 'fail', 'skip'])) { // set case result
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
  caseCount.value++

  const execType = data.execType

  let msg = {}
  if (execType === 'ztf') {
    const scripts = data.scripts
    const sets = genWorkspaceToScriptsMap(scripts)
    msg = {act: 'execCase', testSets: sets}

  } else if (execType === 'unit') {
    const set = {workspaceId: data.id, workspaceType: data.type, cmd: data.cmd, submitResult: data.submitResult,
      productId: currProduct.value.id}
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
    color: #68BB8D
  }

  .result-fail {
    color: #FC2C25
  }
}
</style>

<style lang="less">
.log-list {
  height: 100%;
  font-family: HelveticaNeue;
  .content {
    height: 100%;
    overflow-y: auto;
    .item {
      &.case-item {
        &.case-start {
        }
        &:not(.case-start) {
          display: none !important;
        }
      }
      &.show-detail:not(.case-start) {
        display: flex !important;
      }

      .group {
        width: 16px;
        font-size: 13px;
        text-align: center;
        .link {
          cursor: pointer;
        }
      }
      .sign {
        width: 30px;
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
