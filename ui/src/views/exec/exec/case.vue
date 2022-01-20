<template>
  <div class="indexlayout-main-conent">

    <a-card :bordered="false">
      <template #title>
        执行套件
      </template>
      <template #extra>
        <div class="opt">
          <a-button v-if="isRunning == 'false'" @click="exec" type="primary">执行</a-button>
          <a-button v-if="isRunning == 'true'" @click="stop" type="primary">停止</a-button>

          <a-button @click="back" type="link">返回</a-button>
        </div>
      </template>

      <div id="main">
        <div id="left">
          <div class="toolbar">
            <div class="left"></div>
            <div class="right">
              <a-button @click="expandAll" type="link">
                <span v-if="!isExpand">展开全部</span>
                <span v-if="isExpand">收缩全部</span>
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
          <div id="logs">
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
import {ProjectData} from "@/store/project";
import IconSvg from "@/components/IconSvg";
import {execCase} from "@/views/exec/exec/service";
import {getCache} from "@/utils/localCache";
import settings from "@/config/settings";
import {WebSocket, WsEventName} from "@/services/websocket";
import {PrefixSpace, resizeWidth, scroll, SetWidth} from "@/utils/dom";
import {genExecInfo} from "@/views/exec/service";

const useForm = Form.useForm;

interface ExecCasePageSetupData {
  model: any

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
  exec: (keys) => void;
  stop: (keys) => void;
  isRunning: Ref<string>;
  back: () => void;
}

export default defineComponent({
  name: 'ExecutionCasePage',
  components: {
    IconSvg
  },
  setup(): ExecCasePageSetupData {
    const router = useRouter();
    const model = {}

    const replaceFields = {
      key: 'path',
    };

    let isExpand = ref(false);
    const store = useStore<{ project: ProjectData }>();
    const treeData = computed<any>(() => store.state.project.scriptTree);
    const expandedKeys = ref<string[]>([]);
    const selectedKeys = ref<string[]>([]);
    const checkedKeys = ref<string[]>([]);

    const getOpenKeys = (treeNode, isAll) => {
      if (!treeNode) return

      expandedKeys.value.push(treeNode.path)
      if (treeNode.children && isAll) {
        treeNode.children.forEach((item, index) => {
          getOpenKeys(item, isAll)
        })
      }
    }

    getOpenKeys(treeData.value[0], false)
    watch(treeData, (currConfig) => {
      expandedKeys.value = []
      getOpenKeys(treeData.value[0], false)
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

    let room: string | null = ''
    getCache(settings.currProject).then((token) => {
      room = token
    })

    const {proxy} = getCurrentInstance() as any;
    WebSocket.init(proxy)

    let i = 1
    if (init) {
      proxy.$sub(WsEventName, (data) => {
        console.log(data[0].msg);
        const jsn = JSON.parse(data[0].msg)

        if ('isRunning' in jsn) {
          isRunning.value = jsn.isRunning
        }

        wsMsg.out += genExecInfo(jsn, i)
        i++
        scroll('logs')
      });
      init = false;
    }

    onMounted(() => {
      console.log('onMounted', tree)
      resizeWidth('main', 'left', 'resize', 'content', 280, 800)
      initWsConn()
    })

    const exec = (): void => {
      console.log("exec")
      if (checkedKeys.value.length == 0) {
        wsMsg.out += '请选择需要执行的用例。\n'
        return
      }

      getCache(settings.currProject).then (
          (projectPath) => {
            const msg = {act: 'execCase', projectPath: projectPath, cases: checkedKeys.value}
            console.log('msg', msg)

            wsMsg.out += '\n'
            WebSocket.sentMsg(room, JSON.stringify(msg))
          }
      )
    }
    const stop = (): void => {
      console.log("stop")
      getCache(settings.currProject).then (
          (projectPath) => {
            const msg = {act: 'execStop', projectPath: projectPath}
            console.log('msg', msg)
            WebSocket.sentMsg(room, JSON.stringify(msg))
          }
      )
    }
    const initWsConn = (): void => {
      console.log("initWsConn")
      getCache(settings.currProject).then (
          (projectPath) => {
            const msg = {act: 'init', projectPath: projectPath}
            console.log('msg', msg)
            WebSocket.sentMsg(room, JSON.stringify(msg))
          }
      )
    }

    const back = (): void => {
      router.push(`/exec/history`)
    }

    return {
      model,

      treeData,
      replaceFields,
      expandNode,
      selectNode,
      checkNode,
      isExpand,
      expandAll,
      tree,
      expandedKeys,
      selectedKeys,
      checkedKeys,

      wsMsg,

      exec,
      stop,
      isRunning,
      back,
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
      padding: 16px;
      overflow: auto;

      #logs {
        margin: 0;
        padding: 0;
        height: calc(100% - 10px);
        width: 100%;
        overflow-y: auto;
        white-space: pre-wrap;
        word-wrap: break-word;
        font-family:monospace;
      }
    }
  }
}
</style>