<template>
  <div class="script-main">
    <div v-if="currProduct.id" id="main">
      <div id="left">
        <ScriptTreePage></ScriptTreePage>
      </div>

      <div id="splitter-h"></div>

      <div id="right">
        <div class="toolbar">
          <template v-if="scriptCode !== ''">
            <a-button @click="execSingle" type="primary" size="small">{{ t('exec') }}</a-button>

            <a-button @click="extract" type="primary" size="small">{{ t('extract_step') }}</a-button>
          </template>
        </div>

        <div id="right-content">
          <div id="editor-panel">
            <MonacoEditor
                v-if="scriptCode !== ''"
                class="editor"
                :value="scriptCode"
                :language="lang"
                :options="editorOptions"
            />
          </div>

          <div id="splitter-v"></div>

          <div id="logs-panel">
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
        </div>

      </div>
    </div>
    <div v-if="!currProduct.id">
      <a-empty :image="simpleImage"/>
    </div>
  </div>
</template>

<script lang="ts">
import {
  computed,
  ComputedRef,
  defineComponent,
  getCurrentInstance,
  onMounted,
  onUnmounted,
  reactive,
  ref,
  Ref,
  watch
} from "vue";
import {useStore} from "vuex";
import {useI18n} from "vue-i18n";
import { CloseCircleOutlined, CheckOutlined, CloseOutlined} from '@ant-design/icons-vue';

import {ScriptData} from "./store";
import {resizeWidth, resizeHeight, scroll} from "@/utils/dom";
import {Empty, message, notification} from "ant-design-vue";

import {MonacoOptions} from "@/utils/const";
import MonacoEditor from "@/components/Editor/MonacoEditor.vue";
import {ZentaoData} from "@/store/zentao";
import {WebSocket, WsEventName} from "@/services/websocket";
import {getCache} from "@/utils/localCache";
import settings from "@/config/settings";
import {WsMsg} from "@/views/exec/data";
import {genExecInfo} from "@/views/exec/service";
import ScriptTreePage from "./component/tree.vue";

export default defineComponent({
  name: 'ScriptListPage',
  components: {
    ScriptTreePage,
    MonacoEditor,
   CheckOutlined, CloseOutlined, CloseCircleOutlined,
  },
  setup() {
    const { t } = useI18n();

    let selectedNode = {} as any
    let isExpand = ref(false);

    const zentaoStore = useStore<{ zentao: ZentaoData }>();
    const currSite = computed<any>(() => zentaoStore.state.zentao.currSite);
    const currProduct = computed<any>(() => zentaoStore.state.zentao.currProduct);

    let tree = ref(null)

    const scriptStore = useStore<{ script: ScriptData }>();
    let script = computed<any>(() => scriptStore.state.script.detail);
    let scriptCode = ref('')
    let lang = ref('')
    const editorOptions = ref(MonacoOptions)

    watch(script, () => {
      console.log('watch script', script)
      scriptCode.value = script.value.code
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
      console.log('onMounted', tree)
      initWsConn()

      setTimeout(() => {
        resizeWidth('main', 'left', 'splitter-h', 'right', 280, 800)
        resizeHeight('right-content', 'editor-panel', 'splitter-v', 'logs-panel',
            100, 100, 90)
      }, 600)
    })

    const execSingle = () => {
      console.log('exec', selectedNode.props)

      const msg = {act: 'execCase', workspacePath: 'workspacePath', cases: [selectedNode.props.path]}
      console.log('msg', msg)

      wsMsg.out += '\n'
      WebSocket.sentMsg(room, JSON.stringify(msg))
    }

    const extract = () => {
      console.log('extract', selectedNode.props)

      scriptCode.value = ''
      scriptStore.dispatch('script/extractScript', selectedNode.props).then(() => {
        scriptCode.value = script.value.code
        notification.success({
          message: t('extract_success'),
        })
      }).catch(() => {
        notification.error({
          message: t('extract_fail'),
        });
      })
    }

    return {
      t,
      currSite,
      currProduct,

      isExpand,
      extract,
      tree,
      script,
      scriptCode,
      lang,
      editorOptions,
      simpleImage: Empty.PRESENTED_IMAGE_SIMPLE,

      execSingle,
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
      padding: 3px 0 0 3px;

      .toolbar {
        display: flex;
        padding: 0 3px;
        border-bottom: 1px solid #D0D7DE;

        .left {
          flex: 1;
          .ant-select-tree-switcher {
            display: none;
          }
        }

        .right {
          width: 40px;
          text-align: right;
        }

        .ant-btn-link {
          padding: 0px 3px;
          color: #1890ff;
        }
      }

      .tree-panel {
        height: calc(100% - 35px);
        overflow: auto;

        .ant-tree {
          font-size: 16px;
        }
      }
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
        padding: 5px 10px;
        height: 36px;
        text-align: right;

        .ant-btn {
          margin: 0 5px;
        }
      }

      #right-content {
        height: calc(100% - 46px);

        display: flex;
        flex-direction: column;

        #editor-panel {
          flex: 1;

          padding: 0 6px 0 8px;
          overflow: auto;
          background-color: #fff;
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
          background-color: #fff;

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
