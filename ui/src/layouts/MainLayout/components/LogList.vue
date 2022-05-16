<template>
  <div class="log-list padding muted">

    <div id="content" class="content">
      <template v-for="(item, index) in wsMsg.out" :key="index">
        <div v-if="item.msg" class="item"
           :class="['result-'+caseResult[item.info?.key]]">

          <div class="group">
            <Icon v-if="item.info?.key" icon="chevron-right" />
  <!--          <Icon icon="chevron-down" />-->
          </div>

          <div class="sign">
            <Icon icon="circle" />
          </div>

          <div class="time">
            <span>{{ item.time }}</span>
          </div>
          <div class="msg-span">
            <span>{{ item.msg }}</span>
            <span v-if="caseResult[item.info?.key]">
              [ {{ t(caseResult[item.info?.key]) }} ]
            </span>
          </div>
        </div>

        <div v-if="!item.msg">&nbsp;</div>

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
const { t } = useI18n();

const zentaoStore = useStore<{ Zentao: ZentaoData }>();
const currSite = computed<any>(() => zentaoStore.state.Zentao.currSite);
const currProduct = computed<any>(() => zentaoStore.state.Zentao.currProduct);

const websocketStore = useStore<{ WebSocket: WebSocketData }>();
const wsStatus = computed<any>(() => websocketStore.state.WebSocket.connStatus);

const execStore = useStore<{ Exec: ExecStatus }>();
const isRunning = computed<any>(() => execStore.state.Exec.isRunning);

const caseResult = ref({})

// websocket
let wsMsg = reactive({in: '', out: [] as any[]});

const onWebsocketMsgEvent = (data: any) => {
  console.log('WebsocketMsgEvent in ExecLog', data.msg)

  let item = JSON.parse(data.msg) as WsMsg

  if ('isRunning' in wsMsg) {
    console.log('change isRunning to ', item.isRunning)
    execStore.dispatch('Exec/setRunning', item.isRunning)
  }

  item = genExecInfo(item)
  if (item.info && item.info.key && item.info.status) {
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

<style lang="less" scoped>
.log-list {
  .result-pass { color: #68BB8D }
  .result-fail { color: #FC2C25 }

  font-family: HelveticaNeue;
  .content {
    .item {
      .group {
        width: 16px;
        font-size: 13px;
        text-align: center;
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
