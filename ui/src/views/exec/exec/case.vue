<template>
  <div class="indexlayout-main-conent">
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
          <div class="toolbar">{{isRunning}}
            <a-button :disabled="isRunning" @click="exec" type="primary">执行</a-button>
            <a-button @click="stop">停止</a-button>
            <a-button @click="back" type="link">返回</a-button>
          </div>
          <div class="panel">
            <pre id="logs">{{ wsMsg.out }}</pre>
          </div>
        </div>
      </div>
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
import {resizeWidth} from "@/utils/dom";

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
  isRunning: Ref<boolean>;
  back: () => void;
}

export default defineComponent({
  name: 'ExecutionListPage',
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

    let init = true;
    let isRunning = ref(false);
    let wsMsg = reactive({in: '', out: ''});

    let room: string | null = ''
    getCache(settings.currProject).then((token) => {
      room = token
    })
    const scroll = () => {
      const elem = document.getElementById('logs')
      if (elem) {
        setTimeout(function(){
          elem.scrollTop = elem.scrollHeight + 100;
        },200);
      }
    }

    const {proxy} = getCurrentInstance() as any;
    WebSocket.init(proxy)

    let i = 0
    if (init) {
      proxy.$sub(WsEventName, (data) => {
        console.log(data[0].msg);
        const jsn = JSON.parse(data[0].msg)
        if ('isRunning' in jsn) {
          isRunning.value = jsn.isRunning
        }
        let msg = jsn.msg
        msg = msg.replace(/^"+/,'').replace(/"+$/,'')

        wsMsg.out = wsMsg.out + i++ + '. ' + msg + '\n';
        scroll()
      });
      init = false;
    }

    onMounted(() => {
      console.log('onMounted', tree)
      resizeWidth('main', 'left', 'resize', 'content', 280, 800)
      initWsConn()
    })

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

    const exec = (): void => {
      console.log("exec")
      if (checkedKeys.value.length == 0) {
        wsMsg.out += '请选择用例执行。\n'
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
      router.push(`/execution/history`)
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
      width: 80%;
      height: 100%;
      flex: 1;

      .toolbar {
        padding: 5px 10px;
        height: 46px;
        text-align: right;

        .ant-btn {
          margin-left: 8px;
        }
      }

      .panel {
        padding: 0 16px;
        height: calc(100% - 50px);
        overflow: auto;

        #logs {
          margin: 0;
          padding: 0;
          height: calc(100% - 10px);
          width: 100%;
          overflow-y: auto;
          white-space: pre-wrap;
          word-wrap: break-word;
        }
      }
    }
  }
}
</style>