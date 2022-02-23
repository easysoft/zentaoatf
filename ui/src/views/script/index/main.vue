<template>
  <div v-if="!currProject.path">
    <a-empty :image="simpleImage" :description="t('pls_create_project')"/>
  </div>

  <div v-if="currProject.path" class="indexlayout-main-conent">
    <div v-if="currProject.type === 'unit'" class="panel">
      {{ t('no_script_for_unittest') }}
    </div>

    <div id="main" v-if="currProject.type === 'func'">
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
              show-icon
              @expand="expandNode"
              @select="selectNode"
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
        <div class="toolbar">
          <a-button @click="extract" type="primary">{{ t('extract_step') }}</a-button>
        </div>
        <div class="panel">
          <highlightjs v-if="scriptCode" autodetect :code="scriptCode" />
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import {computed, ComputedRef, defineComponent, onMounted, ref, Ref, watch} from "vue";
import {useStore} from "vuex";
import IconSvg from "@/components/IconSvg";
import {ProjectData} from "@/store/project";
import {ScriptData} from "../store";
import {resizeWidth} from "@/utils/dom";
import {Empty, message, notification} from "ant-design-vue";
import {useI18n} from "vue-i18n";
import 'highlight.js/lib/common';
import hljsVuePlugin from "@highlightjs/vue-plugin";
import 'highlight.js/styles/googlecode.css'

interface ListScriptPageSetupData {
  t: (key: string | number) => string;
  currProject: ComputedRef;
  treeData: ComputedRef<any[]>;
  replaceFields: any,
  tree: Ref;
  script: ComputedRef
  scriptCode: Ref<string>;

  isExpand: Ref<boolean>;
  expandNode: (expandedKeys: string[], e: any) => void;
  selectNode: (keys, e) => void;

  expandAll: (e) => void;
  extract: () => void;
  expandedKeys: Ref<string[]>
  simpleImage: any
}

export default defineComponent({
  name: 'ScriptListPage',
  components: {
    IconSvg,
    highlightjs: hljsVuePlugin.component
  },
  setup(): ListScriptPageSetupData {
    const { t } = useI18n();

    const replaceFields = {
      key: 'path',
    };

    let selectedNode = {} as any
    let isExpand = ref(false);
    const store = useStore<{ project: ProjectData }>();
    const currProject = computed<any>(() => store.state.project.currProject);
    const treeData = computed<any>(() => store.state.project.scriptTree);
    const expandedKeys = ref<string[]>([]);

    const getOpenKeys = (treeNode, isAll) => {
      if (!treeNode) return
      expandedKeys.value.push(treeNode.path)
      if (treeNode.children && isAll) {
        treeNode.children.forEach((item, index) => {
          getOpenKeys(item, isAll)
        })
      }
    }

    expandedKeys.value = []
    getOpenKeys(treeData.value[0], false)
    watch(treeData, (currConfig) => {
      expandedKeys.value = []
      getOpenKeys(treeData.value[0], false)
    })

    let tree = ref(null)

    const storeScript = useStore<{ script: ScriptData }>();
    let script = computed<any>(() => storeScript.state.script.detail);
    let scriptCode = ref('')

    onMounted(() => {
      console.log('onMounted', tree)
      resizeWidth('main', 'left', 'resize', 'content', 280, 800)
    })

    const expandNode = (keys: string[], e: any) => {
      console.log('expandNode', keys[0], e)
    }

    const selectNode = (selectedKeys, e) => {
      console.log('selectNode', e.selectedNodes)
      if (e.selectedNodes.length === 0) return
      selectedNode = e.selectedNodes[0]

      scriptCode.value = ''
      storeScript.dispatch('script/getScript', selectedNode.props).then(() => {
        console.log('===', script)
        scriptCode.value = script.value.code
      })
    }

    const expandAll = (e) => {
      console.log('expandAll')
      isExpand.value = !isExpand.value

      expandedKeys.value = []
      if (isExpand.value) {
        getOpenKeys(treeData.value[0], true)
      }
    }
    const extract = () => {
      console.log('extract', selectedNode.props)
      storeScript.dispatch('script/extractScript', selectedNode.props).then(() => {
        notification.success({
          message: `提取注释为测试步骤和验证点成功！`,
        })
      }).catch(() => {
        notification.error({
          message: `提取注释为测试步骤和验证点失败！`,
        });
      })
    }

    return {
      t,
      currProject,
      treeData,
      replaceFields,
      expandNode,
      selectNode,
      isExpand,
      expandAll,
      extract,
      tree,
      expandedKeys,
      script,
      scriptCode,
      simpleImage: Empty.PRESENTED_IMAGE_SIMPLE,
    }
  }

})
</script>

<style lang="less">
.indexlayout-main-conent {
  margin: 0px;
  height: 100%;

  .panel {
    padding: 20px;
  }

  #main {
    display: flex;
    height: 100%;
    width: 100%;

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
        height: calc(100% - 35px);
        overflow: auto;

        .ant-tree {
          font-size: 16px;
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

    #content {
      width: 80%;
      height: 100%;
      flex: 1;

      .toolbar {
        padding: 5px 10px;
        height: 46px;
        text-align: right;
      }

      .panel {
        padding: 0 16px 0 12px;
        height: calc(100% - 46px);
        overflow: auto;
      }
    }
  }
}
</style>
