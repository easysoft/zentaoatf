<template>
  <div class="indexlayout-main-conent">
    <div id="main">
      <div id="left">
        <div class="toolbar">
          <div class="left">
            <a-select
                ref="select"
                v-model:value="filerType"
                @change="selectFilerType"
                style="width: 90px"
                :bordered="false"
            >
              <a-select-option value="workspace">按目录</a-select-option>
              <a-select-option value="module">按模块</a-select-option>
              <a-select-option value="suite">按套件</a-select-option>
              <a-select-option value="task">按任务</a-select-option>
            </a-select>

            <a-select
                ref="select"
                v-model:value="filerValue"
                @change="selectFilerValue"
                style="width: 200px"
                :bordered="false"
                :dropdownMatchSelectWidth="false"
            >
              <a-select-option v-for="item in filerItems" :key="item.value" :value="item.value">
                {{item.label}}
              </a-select-option>
            </a-select>
          </div>
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
              <DatabaseOutlined v-if="slotProps.type==='workspace'" />
              <FolderOutlined v-if="slotProps.type==='dir' && !slotProps.expanded" />
              <FolderOpenOutlined v-if="slotProps.type==='dir' && slotProps.expanded" />

              <FileOutlined v-if="slotProps.type==='file'" />
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
          <MonacoEditor
              v-if="scriptCode !== ''"
              class="editor"
              :value="scriptCode"
              :language="lang"
              :options="editorOptions"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import {computed, ComputedRef, defineComponent, onMounted, onUnmounted, ref, Ref, watch} from "vue";
import {useStore} from "vuex";
import {useI18n} from "vue-i18n";
import { CodeOutlined, DatabaseOutlined, FolderOutlined, FolderOpenOutlined, FileOutlined} from '@ant-design/icons-vue';

import {ScriptData} from "./store";
import {resizeWidth} from "@/utils/dom";
import {Empty, message, notification} from "ant-design-vue";

import {MonacoOptions} from "@/utils/const";
import MonacoEditor from "@/components/Editor/MonacoEditor.vue";
import {ZentaoData} from "@/store/zentao";
import {cacheExpandedKeys, getScriptFilters, retrieveExpandedKeys, setScriptFilters} from "@/utils/cache";
import {listFilterItems} from "@/views/script/service";
import _ from "lodash";

interface ListScriptPageSetupData {
  t: (key: string | number) => string;
  currSite: ComputedRef;
  currProduct: ComputedRef;
  treeData: ComputedRef<any[]>;
  replaceFields: any,
  tree: Ref;

  filerItems: Ref<[]>
  filerType: Ref
  filerValue: Ref
  selectFilerType: (v) => void;
  selectFilerValue: (v) => void;

  script: ComputedRef
  scriptCode: Ref<string>;
  lang: Ref<string>;
  editorOptions: Ref

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
    DatabaseOutlined, FolderOutlined, FolderOpenOutlined, FileOutlined,
    MonacoEditor,
  },
  setup(): ListScriptPageSetupData {
    const { t } = useI18n();

    const replaceFields = {
      key: 'path',
    };

    let selectedNode = {} as any
    let isExpand = ref(false);

    const zentaoStore = useStore<{ zentao: ZentaoData }>();
    let scriptLoaded = computed<any>(() => zentaoStore.state.zentao.scriptLoaded);
    const currSite = computed<any>(() => zentaoStore.state.zentao.currSite);
    const currProduct = computed<any>(() => zentaoStore.state.zentao.currProduct);
    const treeData = computed<any>(() => zentaoStore.state.zentao.testScripts);

    watch(treeData, (currConfig) => {
      console.log('watch treeData', treeData.value)
      retrieveExpandedKeys().then(async keys => {
        console.log('keys', keys, expandedKeys.value)
        if (keys) expandedKeys.value = keys

        if (!expandedKeys.value || expandedKeys.value.length === 0) {
          getOpenKeys(treeData.value[0], false) // expend first level folder
          await cacheExpandedKeys(expandedKeys.value)
        }
      })
    }, {deep: true})

    console.log(`treeData loaded ${scriptLoaded.value}`, treeData.value.length)
    if (!scriptLoaded.value) { // switch to current page
      zentaoStore.dispatch('zentao/fetchSitesAndProducts', {needLoadScript: true})
    }

    let filerItems = ref([] as any)
    const filerType = ref('')
    const filerValue = ref('')

    const loadFilterItems = async () => {
      const filter = await getScriptFilters()
      filerType.value = filter.by
      filerValue.value = filter.val
      listFilterItems(filerType.value).then((data) => {
        filerItems.value = data.data
      })
    }
    loadFilterItems()

    const selectFilerType = async (val) => {
      console.log('selectFilerType', val)
      await setScriptFilters(val, null)
      await loadFilterItems()
    }
    const selectFilerValue = async (val) => {
      console.log('selectFilerValue', val)
      await setScriptFilters(filerType.value, val)
    }

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

    let tree = ref(null)

    const storeScript = useStore<{ script: ScriptData }>();
    let script = computed<any>(() => storeScript.state.script.detail);
    let scriptCode = ref('')
    let lang = ref('')
    const editorOptions = ref(MonacoOptions)

    onMounted(() => {
      console.log('onMounted', tree)
      resizeWidth('main', 'left', 'resize', 'content', 280, 800)
    })

    onUnmounted(() => {
      console.log('onUnmounted', tree)
    })
    const selectNode = (selectedKeys, e) => {
      console.log('selectNode', e.selectedNodes)
      if (e.selectedNodes.length === 0 || e.selectedNodes[0].props.isDir) {
        scriptCode.value = ''
        return
      }
      selectedNode = e.selectedNodes[0]

      scriptCode.value = ''
      storeScript.dispatch('script/getScript', selectedNode.props).then(() => {
        scriptCode.value = script.value.code
        lang.value = script.value.lang
      })
    }

    const expandNode = (keys: string[], e: any) => {
      console.log('expandNode', expandedKeys.value)
      cacheExpandedKeys(expandedKeys.value)
    }
    const expandAll = (e) => {
      isExpand.value = !isExpand.value

      expandedKeys.value = []
      if (isExpand.value) getOpenKeys(treeData.value[0], true)

      console.log('expandAll', expandedKeys.value)
      cacheExpandedKeys(expandedKeys.value)
    }

    const extract = () => {
      console.log('extract', selectedNode.props)

      scriptCode.value = ''
      storeScript.dispatch('script/extractScript', selectedNode.props).then(() => {
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
      treeData,
      filerItems,

      filerType,
      filerValue,
      selectFilerType,
      selectFilerValue,

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
      lang,
      editorOptions,
      simpleImage: Empty.PRESENTED_IMAGE_SIMPLE,
    }
  }

})
</script>

<style lang="less" scoped>
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
      padding: 3px 0 0 3px;

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

<style lang="less">
.monaco-editor {
  padding: 10px 0;
}
</style>
