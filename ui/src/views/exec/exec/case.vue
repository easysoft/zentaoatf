<template>
  <div class="indexlayout-main-conent">

    <a-card :bordered="false">
      <template #title>
        {{t('exec')}}{{t('case')}}
      </template>

      <template #extra>
        <div class="opt">
          <a-button v-if="isRunning == 'false'" :disabled="checkedKeys.length == 0" @click="exec" type="primary">{{ t('exec') }}</a-button>
          <a-button v-if="isRunning == 'true'" @click="stop" type="primary">{{t('stop')}}</a-button>

          <a-button @click="back" type="link">{{ t('back') }}</a-button>
        </div>
      </template>

      <div id="main">
        <div id="left">
          <div class="toolbar">
            <div class="left"></div>
            <div class="right">
              <a-button @click="expandAll" type="link">
                <span v-if="!isExpand">{{ t('expand_all') }}</span>
                <span v-if="isExpand">{{ t('collapse_all') }}</span>
              </a-button>
            </div>
          </div>

          <div class="tree-panel">
            <a-tree
                ref="tree"
                :tree-data="treeData"
                :replace-fields="replaceFields"
                checkable
                show-icon
                @expand="expandNode"
                @select="selectNode"
                @check="checkNode"

                v-model:selectedKeys="selectedKeys"
                v-model:checkedKeys="checkedKeys"
                v-model:expandedKeys="expandedKeys"
            >
              <template #icon="slotProps">
                <icon-svg v-if="slotProps.isDir" type="folder-outlined"></icon-svg>
                <icon-svg v-if="!slotProps.isDir" type="file-outlined"></icon-svg>
              </template>
            </a-tree>
          </div>
        </div>

        <div id="resize"></div>

        <div id="content">
          <div v-if="wsStatus === 'success'" class="ws-status" :class="wsStatus">
            <icon-svg type="pass-outline" />
            <span class="text">{{t('ws_conn_success')}}</span>
            <span @click="hideWsStatus" class="icon-close"><icon-svg type="close" /></span>
          </div>
          <div v-if="wsStatus === 'fail'" class="ws-status" :class="wsStatus">
            <icon-svg type="fail-outline" />
            <span class="text">{{t('ws_conn_success')}}</span>
            <span @click="hideWsStatus" class="icon-close"><icon-svg type="close" /></span>
          </div>

          <div id="logs" class="logs" :class="{ 'with-status': wsStatus }">
            <span v-html="wsMsg.out"></span>
          </div>
        </div>
      </div>

    </a-card>

  </div>
</template>

<script lang="ts">
import {computed, ComputedRef, defineComponent, getCurrentInstance, onMounted, reactive, Ref, ref, watch} from "vue";
import {Form, notification} from "ant-design-vue";
import {useRouter} from "vue-router";
import {useStore} from "vuex";
import {WorkspaceData} from "@/store/workspace";
import IconSvg from "@/components/IconSvg";
import {execCase} from "@/views/exec/exec/service";
import {getCache} from "@/utils/localCache";
import settings from "@/config/settings";
import {WebSocket, WsEventName} from "@/services/websocket";
import {PrefixSpace, resizeWidth, scroll} from "@/utils/dom";
import {genExecInfo, get, getCaseIdsFromReport} from "@/views/exec/service";
import {WsMsg} from "@/views/exec/data";
import throttle from "lodash.throttle";
import {useI18n} from "vue-i18n";

const useForm = Form.useForm;

interface ExecCasePageSetupData {
  t: (key: string | number) => string;
  model: any
  seq: string

  treeData: ComputedRef<any[]>;
  replaceFields: any,
  expandNode: (expandedKeys: string[], e: any) => void;
  selectNode: (keys, e) => void;
  checkNode: (keys, e) => void;
  isExpand: Ref<boolean>;
  expandAll: (e) => void;
  selectedKeys: Ref<string[]>
  checkedKeys: Ref<string[]>
  expandedKeys: Ref<string[]>
  tree: Ref;

  wsMsg: any,
  wsStatus: Ref<string>,
  exec: (keys) => void;
  stop: (keys) => void;
  isRunning: Ref<string>;
  hideWsStatus: () => void;
  back: () => void;
}

export default defineComponent({
  name: 'ExecutionCasePage',
  components: {
    IconSvg
  },
  setup(): ExecCasePageSetupData {
    const { t } = useI18n();

    const router = useRouter();
    let seq = router.currentRoute.value.params.seq  as string
    seq = seq === '-' ? '' : seq
    let scope = router.currentRoute.value.params.scope as string
    scope = scope === '-' ? '' : scope

    const model = {}
    const replaceFields = {
      key: 'path',
    };

    let isExpand = ref(false);
    const store = useStore<{ workspace: WorkspaceData }>();
    const treeData = computed<any>(() => store.state.workspace.scriptTree);
    const expandedKeys = ref<string[]>([]);
    const selectedKeys = ref<string[]>([]);
    const checkedKeys = ref<string[]>([])

    const selectCasesFromReport = async (): Promise<void> => {
      if (!seq) return

      get(seq).then((json) => {
        setTimeout(() => { // wait tree init completed
          checkedKeys.value = getCaseIdsFromReport(json.data, scope)
        }, 300)
      })
    }
    selectCasesFromReport()

    const getOpenKeys = throttle((treeNode, isAll) => {
      if (!treeNode) return

      expandedKeys.value.push(treeNode.path)
      if (treeNode.children && isAll) {
        treeNode.children.forEach((item, index) => {
          getOpenKeys(item, isAll)
        })
      }
    }, 600)

    const nodeTypeMap = {}
    const getNodeTypeMap = (node): void => {
      if (!node) return

      nodeTypeMap[node.path] = !node.isDir
      if (node.children) {
        node.children.forEach(c => {
          getNodeTypeMap(c)
        })
      }
    }
    const getLeafNodeKeys = (): string[] => {
      let arr = [] as string[]
      checkedKeys.value.forEach(k => {
        if (nodeTypeMap[k]) {
          arr.push(k)
        }
      })
      return arr
    }

    // getOpenKeys(treeData.value[0], false)
    // getNodeTypeMap(treeData.value[0])
    watch(treeData, (currConfig) => {
      expandedKeys.value = []
      getOpenKeys(treeData.value[0], false)
      getNodeTypeMap(treeData.value[0])
    })

    let tree = ref(null)
    const expandNode = (keys: string[], e: any) => {
      console.log('expandNode', keys[0], e)
    }
    const selectNode = (keys, e) => {
      console.log('selectNode', selectedKeys)
    }
    const checkNode = (keys, e) => {
      console.log('checkNode', checkedKeys)
    }

    const expandAll = (e) => {
      console.log('expandAll')
      isExpand.value = !isExpand.value

      expandedKeys.value = []
      if (isExpand.value) {
        getOpenKeys(treeData.value[0], true)
      }
    }

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
      resizeWidth('main', 'left', 'resize', 'content', 280, 800)
      initWsConn()
    })

    const exec = (): void => {
      console.log("exec")
      if (checkedKeys.value.length == 0) {
        wsMsg.out += t('pls_select_case') + '\n'
        return
      }

      getCache(settings.currWorkspace).then (
          (workspacePath) => {
            const leafNodeKeys = getLeafNodeKeys()
            const msg = {act: 'execCase', workspacePath: workspacePath, cases: leafNodeKeys}
            console.log('msg', msg)

            wsMsg.out += '\n'
            WebSocket.sentMsg(room, JSON.stringify(msg))
          }
      )
    }

    const stop = (): void => {
      console.log("stop")
      getCache(settings.currWorkspace).then (
          (workspacePath) => {
            const msg = {act: 'execStop', workspacePath: workspacePath}
            console.log('msg', msg)
            WebSocket.sentMsg(room, JSON.stringify(msg))
          }
      )
    }

    const back = (): void => {
      router.push(`/exec/history`)
    }

    return {
      t,
      model,
      seq,

      treeData,
      tree,
      replaceFields,
      expandNode,
      selectNode,
      checkNode,
      isExpand,
      expandAll,
      expandedKeys,
      selectedKeys,
      checkedKeys,
      back,

      wsMsg,
      wsStatus,
      isRunning,
      exec,
      stop,
      hideWsStatus,
    }
  }

})
</script>

<style lang="less" scoped>
.indexlayout-main-conent {
  margin: 0px;
  height: 100%;
  #main {
    display: flex;
    height: 100%;

    #left {
      width: 380px;
      height: 100%;
      padding: 3px;

      .toolbar {
        display: flex;
        padding: 0 3px;
        border-bottom: 1px solid #D0D7DE;

        .left {
          flex: 1;
        }

        .right {
          width: 70px;
          text-align: right;
        }

        .ant-btn-link {
          padding: 0px 3px;
          color: #1890ff;
        }
      }

      .tree-panel {
        height: calc(100% - 33px);
        overflow: auto;
      }
    }

    #resize {
      width: 2px;
      height: 100%;
      background-color: #D0D7DE;
      cursor: ew-resize;

      &.active {
        background-color: #a9aeb4;
      }
    }

    #content {
      flex: 1;
      height: 100%;
      padding: 0px;

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
</style>