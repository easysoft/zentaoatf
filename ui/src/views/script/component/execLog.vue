<template>
  <div id="exec-log-main">
    <div v-if="wsStatus === 'success'" class="ws-status" :class="wsStatus">
      <CheckOutlined />
      <span class="text">{{t('ws_conn_success')}}</span>
      <span @click="hideWsStatus" class="icon-close"><CloseCircleOutlined /></span>
    </div>

    <div v-if="wsStatus === 'fail'" class="ws-status" :class="wsStatus">
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
          <a-select-option value="result">显示结果日志</a-select-option>
          <a-select-option value="run">显示执行日志</a-select-option>
          <a-select-option value="output">显示详细日志</a-select-option>
        </a-select>

        <a-select
            v-model:value="logStatus"
            size="small"
        >
          <a-select-option value="">所有结果</a-select-option>
          <a-select-option value="pass">成功结果</a-select-option>
          <a-select-option value="fail">失败结果</a-select-option>
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
import {WebSocket, WsEventName} from "@/services/websocket";
import {getCache} from "@/utils/localCache";
import settings from "@/config/settings";
import {WsMsg} from "@/views/exec/data";
import {genExecInfo} from "@/views/exec/service";
import bus from "@/utils/eventBus";
import {logLevelMap} from "@/utils/const";

export default defineComponent({
  name: 'ScriptExecLogPage',
  components: {
   CheckOutlined, CloseOutlined, CloseCircleOutlined,
  },
  setup() {
    const { t } = useI18n();

    const zentaoStore = useStore<{ zentao: ZentaoData }>();
    const currSite = computed<any>(() => zentaoStore.state.zentao.currSite);
    const currProduct = computed<any>(() => zentaoStore.state.zentao.currProduct);

    const scriptStore = useStore<{ script: ScriptData }>();
    let script = computed<any>(() => scriptStore.state.script.detail);

    const logLevel = ref('result')
    const logStatus = ref('')

    watch(currProduct, () => {
      console.log('watch currProduct', currProduct)
    }, {deep: true})

    // websocket
    let init = true;
    let isRunning = ref('false');
    let wsMsg = reactive({in: '', out: [] as any[]});

    let room = ''
    getCache(settings.currWorkspace).then((token) => {
      room = token || ''
    })

    const {proxy} = getCurrentInstance() as any;
    WebSocket.init(proxy)

    let wsStatus = ref('success')

    if (init) {
      proxy.$sub(WsEventName, (data) => {
        console.log('---', data[0].msg);
        const jsn = JSON.parse(data[0].msg) as WsMsg

        if (jsn.conn) { // ws connection status updated
          wsStatus.value = jsn.conn
          return
        }

        if ('isRunning' in jsn) {
          isRunning.value = jsn.isRunning
        }

        const msg = genExecInfo(jsn)
        if (msg) {
          wsMsg.out.push(msg)
        }
        scroll('content')
      });
      init = false;
    }

    const initWsConn = (): void => {
      console.log("initWsConn")
      getCache(settings.currWorkspace).then (
          (workspacePath) => {
            const msg = {act: 'init', workspacePath: workspacePath}
            console.log('msg', msg)
            WebSocket.sentMsg(room, JSON.stringify(msg))
          }
      )
    }
    const hideWsStatus = (): void => {
      wsStatus.value = ''
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

      initWsConn()
    })
    onBeforeUnmount( () => {
      bus.off(settings.eventExec, exec);
    })

    const exec = (scripts: any) => {
      console.log('exec', scripts)

      const workspaceIds = [] as number[]
      const mp = {}
      scripts.forEach((item) => {
        if (!mp[item.workspaceId]) {
          mp[item.workspaceId] = []
          workspaceIds.push(item.workspaceId)
        }

        mp[item.workspaceId].push(item.path)
      })

      const sets = [] as any[]
      workspaceIds.forEach((workspaceId) => {
        const set = {workspaceId: workspaceId, cases: mp[workspaceId]}
        sets.push(set)
      })

      const msg = {act: 'execCase', testSets: sets}
      console.log('msg', msg)

      // wsMsg.out.push({msg: ''})
      WebSocket.sentMsg(room, JSON.stringify(msg))
    }

    const checkoutCases = () => {
      console.log('checkoutCases')
    }
    const checkinCases = () => {
      console.log('checkinCases')
    }

    return {
      t,
      currSite,
      currProduct,

      logLevel,
      logStatus,
      script,
      exec,
      checkoutCases,
      stop,
      hideWsStatus,
      addCount,
      wsMsg,
      wsStatus,
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
