<template>
  <div id="exec-log-main">
    <div v-if="showStatus && wsStatus === 'success'" class="ws-status" :class="wsStatus">
      <CheckOutlined />
      <span class="text">{{t('ws_conn_success')}}</span>
      <span @click="hideWsStatus" class="icon-close"><CloseCircleOutlined /></span>
    </div>

    <div v-if="showStatus && wsStatus === 'fail'" class="ws-status" :class="wsStatus">
      <CloseOutlined />
      <span class="text">{{t('ws_conn_success')}}</span>
      <span @click="hideWsStatus" class="icon-close"><CloseCircleOutlined /></span>
    </div>

    <div id="logs" class="logs" :class="{ 'with-status': wsStatus }">
      <div class="log-filters">
        <a-select
            v-model:value="logLevel"
            size="small"
        >
          <a-select-option value="result">{{ t('show_result_log') }}</a-select-option>
          <a-select-option value="run">{{ t('show_exec_log') }}</a-select-option>
          <a-select-option value="output">{{ t('show_detail_log') }}</a-select-option>
        </a-select>

        <a-select
            v-model:value="logStatus"
            size="small"
        >
          <a-select-option value="">{{ t('all_result') }}</a-select-option>
          <a-select-option value="pass">{{ t('pass_result') }}</a-select-option>
          <a-select-option value="fail">{{ t('fail_result') }}</a-select-option>
        </a-select>

      </div>

      <div id="content" :class="['status-'+logStatus, 'category-'+logLevel]">
        {{ void (count = 1) }}
        <div v-for="(item, index) in wsMsg.out" :key="index" :class="['category-'+item.category, 'item']">
          <div class="no-span">
            <span v-if="item.msg">{{ count }}. </span>
          </div>
          <div class="msg-span">
            <span v-if="item.msg">{{ item.msg }}</span>
            <span v-if="!item.msg">&nbsp;</span>
          </div>

          <template v-if="addCount(item)"> {{ void (count++) }}</template>
        </div>
      </div>
    </div>

  </div>
</template>

<script lang="ts">
import {
  computed,
  defineComponent,
  getCurrentInstance, onBeforeUnmount,
  onMounted,
  reactive,
  ref,
  watch
} from "vue";
import {useStore} from "vuex";
import {useI18n} from "vue-i18n";
import { CloseCircleOutlined, CheckOutlined, CloseOutlined} from '@ant-design/icons-vue';

import {ScriptData} from "../store";
import {scroll} from "@/utils/dom";
import {ZentaoData} from "@/store/zentao";
import {WebSocket} from "@/services/websocket";
import {getCache} from "@/utils/localCache";
import settings from "@/config/settings";
import bus from "@/utils/eventBus";
import {logLevelMap} from "@/utils/const";
import {WsMsg} from "@/types/data";
import {genExecInfo, genWorkspaceToScriptsMap} from "../service";
import {ExecStatus} from "@/store/exec";
import {WebSocketData} from "@/store/websoket";

export default defineComponent({
  name: 'ScriptExecLogPage',
  components: {
   CheckOutlined, CloseOutlined, CloseCircleOutlined,
  },
  setup() {
    const { t } = useI18n();

    const websocketStore = useStore<{ WebSocket: WebSocketData }>();
    const wsStatus = computed<any>(() => websocketStore.state.WebSocket.connStatus);

    const zentaoStore = useStore<{ Zentao: ZentaoData }>();
    const currSite = computed<any>(() => zentaoStore.state.Zentao.currSite);
    const currProduct = computed<any>(() => zentaoStore.state.Zentao.currProduct);

    const execStore = useStore<{ Exec: ExecStatus }>();
    const isRunning = computed<any>(() => execStore.state.Exec.isRunning);

    const scriptStore = useStore<{ Script: ScriptData }>();
    const script = computed<any>(() => scriptStore.state.Script.detail);

    const logLevel = ref('result')
    const logStatus = ref('')

    watch(currProduct, () => {
      console.log('watch currProduct', currProduct)
    }, {deep: true})

    // websocket
    let wsMsg = reactive({in: '', out: [] as any[]});

    const getWsMsg = (data: any) => {
      console.log('OnWebSocketEvent in ExecLog', data.msg)

      const jsn = JSON.parse(data.msg) as WsMsg

      if (jsn.conn) { // update connection status
        return
      }

      if ('isRunning' in jsn) {
        execStore.dispatch('Exec/setRunning', jsn.isRunning)
      }

      const msg = genExecInfo(jsn)
      if (msg) {
        wsMsg.out.push(msg)
      }
      scroll('content')
    }
    let init = true;
    if (init) {
      bus.on(settings.eventWebSocket, getWsMsg);
      init = false;
    }

    const showStatus = ref(true)
    const hideWsStatus = (): void => {
      showStatus.value = false
    }

    const addCount = (item): boolean => {
      let ret = false

      if (!item.msg) return false

      const levelCode = logLevelMap[logLevel.value].code
      const itemCode = logLevelMap[item.category].code

      if (!item.category || !logLevelMap[item.category]) {
        ret = true
      } else {
        ret = levelCode <= itemCode
      }

      if (logStatus.value === 'pass') {
        console.log('===', item.category)

        if (item.category === 'result') {
          return true
        } else {
          return false
        }
      } else if (logStatus.value === 'fail') {
        if (item.category === 'error') {
          return true
        } else {
          return false
        }
      }

      return ret
    }

    onMounted(() => {
      console.log('onMounted')
      bus.on(settings.eventExec, exec);
    })
    onBeforeUnmount( () => {
      bus.off(settings.eventExec, exec);
      bus.off(settings.eventWebSocket, getWsMsg);
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

      console.log('msg', msg)
      WebSocket.sentMsg(settings.webSocketRoom, JSON.stringify(msg))
    }

    return {
      t,
      currSite,
      currProduct,

      logLevel,
      logStatus,
      script,
      exec,
      stop,
      hideWsStatus,
      addCount,
      wsMsg,
      wsStatus,
      showStatus,
      isRunning,
    }
  }

})
</script>

<style lang="less">
#content {
  .category-output { color: #95a5a6 }
  .category-run { color: #1890ff }
  .category-result { color: #68BB8D }
  .category-error { color: #FC2C25 }

  &.category-run {
    .category-output { display: none !important; }
  }
  &.category-result {
    .category-output, .category-run { display: none !important; }
  }
  &.category-error {
    .category-output, .category-run, .category-result { display: none !important; }
  }

  &.status-pass {
    .category-error { display: none !important; }
  }
  &.status-fail {
    .category-output, .category-run, .category-result { display: none !important; }
  }

  .item {
    display: flex;

    .no-span {
      width: 60px;
      text-align: right;
    }

    .msg-span {
      flex: 1;
    }
  }
}
</style>

<style lang="less" scoped>
#exec-log-main {
  height: 100%;

  .ws-status {
    padding-left: 8px;
    height: 44px;
    line-height: 44px;
    color: #333333;

    &.success {
      background-color: #DAF7E9;
      svg {
        color: #DAF7E9;
      }
    }
    &.error {
      background-color: #FFD6D0;
      svg {
        color: #FFD6D0;
      }
    }

    .text {
      display: inline-block;
      margin-left: 5px;
    }
    .icon-close {
      position: absolute;
      padding: 5px;
      line-height: 34px;
      right: 15px;
      cursor: pointer;
      svg {
        font-size: 8px;
        color: #333333;
      }
    }
  }

  #logs {
    margin: 0;
    padding: 10px;
    position: relative;
    width: 100%;
    overflow-y: auto;
    white-space: pre-wrap;
    word-wrap: break-word;
    font-family:monospace;

    height: 100%;
    &.with-status {
      height: calc(100% - 45px);
    }

    .log-filters {
      position: absolute;
      top: 5px;
      right: 5px;
      width: 230px;
      .ant-select {
        margin: 0 3px;
      }
    }

    #content {
      height: 100%;
      overflow-y: auto;

      .item {
        .no-span {
          display: inline-block;
          text-align: right;
          min-width: 50px;
          text-align: right;
        }
      }

    }
  }
}
</style>

<style lang="less">

</style>
