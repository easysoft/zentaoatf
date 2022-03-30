<template>
  <div class="tree-main">
    <div class="toolbar">
      <div class="left">
        <a-select
            v-model:value="filerType"
            @change="selectFilerType"
            style="width: 120px"
        >
          <a-select-option value=""></a-select-option>
          <a-select-option value="workspace">按目录</a-select-option>
          <a-select-option value="module">按模块</a-select-option>
          <a-select-option value="suite">按套件</a-select-option>
          <a-select-option value="task">按任务</a-select-option>
        </a-select>

        <a-select
            v-model:value="filerValue"
            @change="selectFilerValue"
            style="width: 200px"
            :dropdownMatchSelectWidth="false"
        >
          <a-select-option value=""></a-select-option>
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
          v-if="!treeDataEmpty"
          :tree-data="treeData"
          v-model:expandedKeys="expandedKeys"
          v-model:selectedKeys="selectedKeys"
          v-model:checkedKeys="checkedKeys"

          :replace-fields="replaceFields"
          show-icon
          checkable
          @expand="expandNode"
          @select="selectNode"
          @check="checkNode"
      >
        <template #icon="slotProps">
          <DatabaseOutlined v-if="slotProps.type==='workspace'" />
          <FolderOutlined v-if="slotProps.type==='dir' && !slotProps.expanded" />
          <FolderOpenOutlined v-if="slotProps.type==='dir' && slotProps.expanded" />

          <FileOutlined v-if="slotProps.type==='file'" />
        </template>
      </a-tree>

      <a-empty v-if="treeDataEmpty" :image="simpleImage"/>
    </div>

    <div class="actions">
      <a-button @click="checkoutCases">导出禅道用例</a-button>
      <a-button :disabled="checkedKeys.length === 0" @click="checkinCases">更新禅道用例</a-button>
      <a-button :disabled="checkedKeys.length === 0" @click="execSelected">执行选中</a-button>
    </div>

    <a-modal
        :title="fromTitle"
        v-if="fromVisible"
        :visible="true"
        :onCancel="onCancel"
        width="800px"
        :destroy-on-close="true"
        :mask-closable="false"
        :footer="null"
    >
      <SyncFromZentao
          :onClose="onSave"
      />
    </a-modal>

  </div>
</template>

<script lang="ts">
import {
  computed,
  defineComponent,
  onMounted,
  ref,
  watch
} from "vue";
import {useStore} from "vuex";
import {useI18n} from "vue-i18n";

import {ScriptData} from "../store";
import {Empty, message, notification} from "ant-design-vue";
import { DatabaseOutlined, FolderOutlined, FolderOpenOutlined, FileOutlined} from '@ant-design/icons-vue';

import bus from "@/utils/eventBus";
import {ZentaoData} from "@/store/zentao";
import {cacheExpandedKeys, getScriptFilters, retrieveExpandedKeys, setScriptFilters} from "@/utils/cache";
import {genWorkspaceToScriptsMap, listFilterItems, syncToZentao} from "../service";
import settings from "@/config/settings";
import {getCaseIdsFromReport} from "../service";
import {useRouter} from "vue-router";

import SyncFromZentao from "./syncFromZentao.vue"

export default defineComponent({
  name: 'ScriptTreePage',
  components: {
    SyncFromZentao,
    DatabaseOutlined, FolderOutlined, FolderOpenOutlined, FileOutlined,
  },
  props: {
  },

  setup(props) {
    const { t } = useI18n();

    const fromTitle = ref('从禅道同步用例')
    const fromVisible = ref(false)

    const router = useRouter();
    let workspace = router.currentRoute.value.params.workspace  as string
    workspace = workspace === '-' ? '' : workspace
    let seq = router.currentRoute.value.params.seq  as string
    seq = seq === '-' ? '' : seq
    let scope = router.currentRoute.value.params.scope as string
    scope = scope === '-' ? '' : scope

    const zentaoStore = useStore<{ zentao: ZentaoData }>();
    const currSite = computed<any>(() => zentaoStore.state.zentao.currSite);
    const currProduct = computed<any>(() => zentaoStore.state.zentao.currProduct);

    const store = useStore<{ script: ScriptData }>();
    const treeData = computed<any>(() => store.state.script.list);
    const treeDataEmpty = computed<boolean>(() => !(treeData.value.length > 0 &&
        treeData.value[0] && treeData.value[0].children))

    const filerType = ref('')
    const filerValue = ref('')

    const scriptStore = useStore<{ script: ScriptData }>();

    const selectCasesFromReport = async (): Promise<void> => {
      if (!seq) return

      getCaseIdsFromReport(workspace, seq, scope).then((json) => {
         checkedKeys.value = json.data
      })
    }
    selectCasesFromReport()

    watch(currProduct, () => {
      console.log('watch currProduct', currProduct.value.id)
      initData()
    }, {deep: true})

    watch(treeData, (currConfig) => {
      console.log('watch treeData', treeData.value)

      getNodeMap(treeData.value[0])

      retrieveExpandedKeys().then(async keys => {
        console.log('keys')
        if (keys) expandedKeys.value = keys

        if (!expandedKeys.value || expandedKeys.value.length === 0) {
          getOpenKeys(treeData.value[0], false) // expend first level folder
          await cacheExpandedKeys(expandedKeys.value)
        }
      })
    }, {deep: true})

    let filerItems = ref([] as any)

    const loadScripts = async () => {
      console.log(`=== filerType: ${filerType.value}, filerValue: ${filerValue.value}`)

      const params = {filerType: filerType.value, filerValue: filerValue.value} as any
      store.dispatch('script/listScript', params)
    }
    loadScripts()

    const loadFilterItems = async () => {
      if (!filerType.value) {
        const data = await getScriptFilters()
        filerType.value = data.by
      }

      filerValue.value = ''
      if (!filerType.value) {
        filerItems.value = []
        return
      }

      const result = await listFilterItems(filerType.value)
      filerItems.value = result.data
    }

    const initData = async () => {
      console.log('init')
      loadFilterItems()
      loadScripts()
    }
    initData()

    const selectFilerType = async (val) => {
      console.log('selectFilerType', val)
      await setScriptFilters(val, '')
      await loadFilterItems()
      await loadScripts()
    }
    const selectFilerValue = async (val) => {
      console.log('selectFilerValue', val)
      await loadScripts()
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
    getOpenKeys(treeData.value[0], false)

    const selectedKeys = ref<string[]>([])
    const checkedKeys = ref<string[]>([])

    const replaceFields = {
      key: 'path',
    };

    let isExpand = ref(false);

    let tree = ref(null)

    onMounted(() => {
      console.log('onMounted', tree)
    })

    const execSelected = () => {
      console.log('execSelected')

      // cancel selecting any nodes
      selectedKeys.value = []
      scriptStore.dispatch('script/getScript', null)

      const leafNodes = getLeafNodes()
      bus.emit(settings.eventExec, leafNodes);
    }

    const checkoutCases = () => {
      console.log('checkoutCases')
      fromVisible.value = true
    }
    const onSave = async () => {
      console.log('onSave')
      fromVisible.value = false
      loadScripts()
    }
    const onCancel = async () => {
      console.log('onCancel')
      fromVisible.value = false
    }

    const checkinCases = () => {
      console.log('checkinCases')

      const leafNodes = getLeafNodes()
      const sets = genWorkspaceToScriptsMap(leafNodes)
      console.log('sets', sets)

      syncToZentao(sets).then((json) => {
        if (json.code === 0) {
          notification.success({
            message: `同步成功`,
          });
        } else {
          notification.error({
            message: `同步失败`,
            description: json.msg,
          });
        }
      })
    }

    const selectNode = (selectedKeys, e) => {
      console.log('selectNode', selectedKeys)

      let data = null

      if (e.selectedNodes.length > 0) {
        data = e.selectedNodes[0].props
      }

      scriptStore.dispatch('script/getScript', data)
    }
    const checkNode = () => {
      console.log('checkNode', checkedKeys)
    }

    const expandNode = (keys: string[], e: any) => {
      console.log('expandNode', keys)
      cacheExpandedKeys(expandedKeys.value)
    }
    const expandAll = (e) => {
      console.log('expandAll')
      isExpand.value = !isExpand.value

      expandedKeys.value = []
      if (isExpand.value) getOpenKeys(treeData.value[0], true)

      cacheExpandedKeys(expandedKeys.value)
    }

    const nodeMap = {}
    const getNodeMap = (node): void => {
      if (!node) return

      nodeMap[node.path] = node
      if (node.children) {
        node.children.forEach(c => {
          getNodeMap(c)
        })
      }
    }
    const getLeafNodes = (): string[] => {
      console.log('nodeMap', nodeMap)
      let arr = [] as string[]
      checkedKeys.value.forEach(k => {
        if (nodeMap[k].type === 'file') {
          arr.push(nodeMap[k])
        }
      })
      return arr
    }

    return {
      t,

      currSite,
      currProduct,
      treeData,
      treeDataEmpty,
      filerItems,

      filerType,
      filerValue,
      selectFilerType,
      selectFilerValue,

      replaceFields,
      expandNode,
      selectNode,
      checkNode,
      execSelected,
      checkoutCases,
      checkinCases,
      isExpand,
      expandAll,
      tree,
      expandedKeys,
      selectedKeys,
      checkedKeys,
      simpleImage: Empty.PRESENTED_IMAGE_SIMPLE,

      fromTitle,
      fromVisible,
      onSave,
      onCancel,
    }
  }

})
</script>

<style lang="less" scoped>

.tree-main {
  width: 100%;
  height: 100%;
  padding: 0;
  display: flex;
  flex-direction: column;

  .toolbar {
    display: flex;
    padding: 4px 3px;
    height: 40px;
    //border-bottom: 1px solid #D0D7DE;

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
    flex: 1;
    overflow: auto;

    .ant-tree {
      font-size: 16px;
    }
  }

  .actions {
    padding: 4px;
    height: 40px;
    text-align: center;

    .ant-btn {
      margin: 0 5px;
    }
  }
}

</style>

<style lang="less">

</style>
