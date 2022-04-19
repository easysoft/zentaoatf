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
          <a-select-option value="workspace">{{t('by_workspace')}}</a-select-option>

          <template v-if="currProduct.id">
            <a-select-option value="module">{{t('by_module')}}</a-select-option>
            <a-select-option value="suite">{{t('by_suite')}}</a-select-option>
            <a-select-option value="task">{{t('by_task')}}</a-select-option>
          </template>

        </a-select>

        <a-select
            v-model:value="filerValue"
            @change="selectFilerValue"
            style="width: 180px"
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
          :load-data="onLoadData"
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
      <a-button :disabled="!currProduct.id"
                @click="checkoutCases">{{t('checkout_case')}}</a-button>
      <a-button :disabled="!currProduct.id || checkedKeys.length === 0"
                @click="checkinCases">{{t('checkin_case')}}</a-button>

      <template v-if="currWorkspace.type === 'ztf'">
        <a-button
            v-if="isRunning !== 'true'"
            :disabled="checkedKeys.length === 0"
            @click="execSelected">
          {{t('exec_selected')}}
        </a-button>
        <a-button
            v-if="isRunning === 'true'"
            @click="execStop">
          {{t('stop')}}
        </a-button>
      </template>

      <a-button
          v-if="currWorkspace.type !== 'ztf'"
          @click="toExecUnit">
        {{testToolMap[currWorkspace.type] + t('test')}}
      </a-button>
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
import {setExpandedKeys, getScriptFilters, getExpandedKeys, setScriptFilters} from "@/utils/cache";
import {genWorkspaceToScriptsMap, listFilterItems, getCaseIdsFromReport, syncToZentao} from "../service";
import settings from "@/config/settings";
import {useRouter} from "vue-router";

import SyncFromZentao from "./syncFromZentao.vue"
import {isWindows} from "@/utils/comm";
import {testToolMap} from "@/utils/testing";
import {ExecStatus} from "@/store/exec";

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
    const isWin = isWindows()

    const fromTitle = ref('从禅道同步用例')
    const fromVisible = ref(false)

    const router = useRouter();
    let workspace = router.currentRoute.value.params.workspace  as string
    workspace = workspace === '-' ? '' : workspace
    let seq = router.currentRoute.value.params.seq  as string
    seq = seq === '-' ? '' : seq
    let scope = router.currentRoute.value.params.scope as string
    scope = scope === '-' ? '' : scope

    const zentaoStore = useStore<{ Zentao: ZentaoData }>();
    const currSite = computed<any>(() => zentaoStore.state.Zentao.currSite);
    const currProduct = computed<any>(() => zentaoStore.state.Zentao.currProduct);

    const execStore = useStore<{ Exec: ExecStatus }>();
    const isRunning = computed<any>(() => execStore.state.Exec.isRunning);

    const store = useStore<{ Script: ScriptData }>();
    const treeData = computed<any>(() => store.state.Script.list);
    const currWorkspace = computed<any>(() => store.state.Script.currWorkspace);
    const treeDataEmpty = computed<boolean>(() => !(treeData.value.length > 0 &&
        treeData.value[0] && treeData.value[0].children))

    const filerType = ref('')
    const filerValue = ref('')

    const scriptStore = useStore<{ script: ScriptData }>();

    const selectCasesFromReport = async (): Promise<void> => {
      if (!seq) return

      getCaseIdsFromReport(workspace, seq, scope).then((json) => {
         checkedKeys.value = json.data
        router.push(`/script/index`) // remove the params of re-test
      })
    }
    selectCasesFromReport()

    watch(currProduct, () => {
      console.log('watch currProduct', currProduct.value.id)
      initData()
    }, {deep: true})

    watch(treeData, (currConfig) => {
      console.log('watch treeData', treeData.value)

      treeData.value[0].title = t('test_script')
      getNodeMap(treeData.value[0])

      getExpandedKeys(currSite.value.id, currProduct.value.id).then(async keys => {
        console.log('keys')
        if (keys) expandedKeys.value = keys

        if (!expandedKeys.value || expandedKeys.value.length === 0) {
          getOpenKeys(treeData.value[0], false) // expend first level folder
          await setExpandedKeys(currSite.value.id, currProduct.value.id, expandedKeys.value)
        }
      })
    }, {deep: true})

    let filerItems = ref([] as any)

    const loadScripts = async () => {
      console.log(`=== filerType: ${filerType.value}, filerValue: ${filerValue.value}`)

      const params = {filerType: filerType.value, filerValue: filerValue.value} as any
      store.dispatch('Script/listScript', params)
    }

    // 异步加载
    const onLoadData = async (treeNode: any) => {
      console.log('onLoadData')
      await store.dispatch('Script/loadChildren', treeNode)
    }

    // filters
    const loadFilterItems = async () => {
      const data = await getScriptFilters(currSite.value.id, currProduct.value.id)

      if (!filerType.value) {
        filerType.value = data.by
      }
      filerValue.value = data.val

      if (!currProduct.value.id && filerType.value !== 'workspace') {
        filerType.value = 'workspace'
        filerValue.value = ''
      }

      if (filerType.value) {
        const result = await listFilterItems(filerType.value)
        filerItems.value = result.data

        let found = false
        if (filerItems.value) {
          filerItems.value.forEach((item) => {
            // console.log(`${filerValue.value}, ${item.value}`)
            if (filerValue.value === item.value) found = true
          })
        }

        if (!found) filerValue.value = ''
      }
    }
    const selectFilerType = async (type) => {
      console.log('selectFilerType', type)
      await setScriptFilters(currSite.value.id, currProduct.value.id, type, '')

      await loadFilterItems()
      await loadScripts()
    }
    const selectFilerValue = async (val) => {
      console.log('selectFilerValue', val)
      await setScriptFilters(currSite.value.id, currProduct.value.id, filerType.value, val)

      await loadScripts()
    }

    const initData = async () => {
      console.log('init')
      await loadFilterItems()
      await loadScripts()
    }
    // only do it when switch from another pages, otherwise will called by watching currProduct method.
    if (filerValue.value.length === 0) initData()

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

      selectNothing()

      const leafNodes = getLeafNodes()
      bus.emit(settings.eventExec, {execType: 'ztf', scripts: leafNodes});
    }
    const execStop = () => {
      console.log('execStop')
      const data = Object.assign({execType: 'stop'})
      bus.emit(settings.eventExec, data);
    }

    const toExecUnit = () => {
      console.log('toExecUnit')
      selectNothing()
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
            message: t('sync_success'),
          });
        } else {
          notification.error({
            message: t('sync_fail'),
            description: json.msg,
          });
        }
      })
    }

    const selectNode = (selectedKeys, e) => {
      console.log('selectNode', e.node.dataRef.workspaceId)

      if (e.node.dataRef.workspaceType !== 'ztf') checkNothing()

      scriptStore.dispatch('Script/getScript', e.node.dataRef)

      scriptStore.dispatch('Script/changeWorkspace',
          {id: e.node.dataRef.workspaceId, type: e.node.dataRef.workspaceType})
    }

    const checkNode = (checkedKeys, e) => {
      console.log('checkNode', e)
      selectNothing()
      scriptStore.dispatch('Script/changeWorkspace',
          {id: e.node.dataRef.workspaceId, type: e.node.dataRef.workspaceType})
    }
    const selectNothing = () => {
      selectedKeys.value = []
      scriptStore.dispatch('Script/getScript', null)
    }
    const checkNothing = () => {
      checkedKeys.value = []
    }

    const expandNode = (keys: string[], e: any) => {
      console.log('expandNode')
      setExpandedKeys(currSite.value.id, currProduct.value.id, expandedKeys.value)
    }
    const expandAll = (e) => {
      console.log('expandAll')
      isExpand.value = !isExpand.value

      expandedKeys.value = []
      if (isExpand.value) getOpenKeys(treeData.value[0], true)

      setExpandedKeys(currSite.value.id, currProduct.value.id, expandedKeys.value)
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
      isWin,

      currSite,
      currProduct,
      treeData,
      currWorkspace,
      testToolMap,
      treeDataEmpty,
      filerItems,

      filerType,
      filerValue,
      selectFilerType,
      selectFilerValue,

      replaceFields,
      onLoadData,
      expandNode,
      selectNode,
      checkNode,
      isRunning,
      execSelected,
      execStop,
      toExecUnit,
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
      width: 60px;
      text-align: center;
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
      padding: 4px 6px;
    }
  }
}

</style>

<style lang="less">

</style>
