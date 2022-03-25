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
      <span v-html="wsMsg.out"></span>
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
import {resizeWidth, resizeHeight, scroll} from "@/utils/dom";
import {ZentaoData} from "@/store/zentao";
import {WebSocket, WsEventName} from "@/services/websocket";
import {getCache} from "@/utils/localCache";
import settings from "@/config/settings";
import {WsMsg} from "@/views/exec/data";
import {genExecInfo} from "@/views/exec/service";
import bus from "@/utils/eventBus";

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

    watch(currProduct, () => {
      console.log('watch currProduct', currProduct)
    }, {deep: true})

    // websocket
    let init = true;
    let isRunning = ref('false');
    let wsMsg = reactive({in: '', out: ''});

    let room = ''
    getCache(settings.currWorkspace).then((token) => {
      room = token || ''
    })

    const {proxy} = getCurrentInstance() as any;
    WebSocket.init(proxy)

    let wsStatus = ref('')
    let i = 1
    if (init) {
      proxy.$sub(WsEventName, (data) => {
        console.log(data[0].msg);
        const jsn = JSON.parse(data[0].msg) as WsMsg

        if (jsn.conn) { // ws connection status updated
          wsStatus.value = jsn.conn
          return
        }

        if ('isRunning' in jsn) {
          isRunning.value = jsn.isRunning
        }

        wsMsg.out += genExecInfo(jsn, i)
        i++
        scroll('logs')
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

      wsMsg.out += '\n'
      // WebSocket.sentMsg(room, JSON.stringify(msg))
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

      script,
      exec,
      checkoutCases,
      stop,
      hideWsStatus,
      wsMsg,
      wsStatus,
      isRunning,
    }
  }

})
</script>

<style lang="less" scoped>
.script-main {
  margin: 0px;
  height: 100%;

  #main {
    display: flex;
    height: 100%;
    width: 100%;

    #left {
      width: 380px;

      height: 100%;
      padding: 0;
    }

    #splitter-h {
      width: 2px;
      height: 100%;
      background-color: #D0D7DE;
      cursor: ew-resize;

      &.active {
        background-color: #a9aeb4;
      }
    }

    #right {
      flex: 1;
      height: 100%;

      .toolbar {
        padding: 4px 10px;
        height: 40px;
        text-align: right;

        .ant-btn {
          margin: 0 5px;
        }
      }

      #right-content {
        height: calc(100% - 50px);

        display: flex;
        flex-direction: column;

        #editor-panel {
          flex: 1;

          padding: 0 6px 0 8px;
          overflow: auto;
        }

        #splitter-v {
          width: 100%;
          height: 2px;
          background-color: #D0D7DE;
          cursor: ns-resize;

          &.active {
            background-color: #a9aeb4;
          }
        }

        #logs-panel {
          height: 160px;

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
            width: 100%;
            overflow-y: auto;
            white-space: pre-wrap;
            word-wrap: break-word;
            font-family:monospace;

            height: 100%;
            &.with-status {
              height: calc(100% - 45px);
            }
          }
        }
      }
    }
  }
}
</style>

<style lang="less">
.monaco-editor {
  padding: 10px 0;
}
</style>
