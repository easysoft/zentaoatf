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
            <a-select-option value="suite">{{t('by_suite')}}</a-select-option>
            <a-select-option value="task">{{t('by_task')}}</a-select-option>
          </template>

        </a-select>

        <a-select
            v-model:value="filerValue"
            @change="selectFilerValue"
            style="width: 120px"
            :dropdownMatchSelectWidth="false"
        >
          <a-select-option value=""></a-select-option>
          <a-select-option v-for="item in filerItems" :key="item.value" :value="item.value">
            {{item.label}}
          </a-select-option>
        </a-select>
      </div>

      <div class="right">
        <a-button @click="expandAllOrNot" type="link">
          <span v-if="!isExpand">{{ t('expand_all') }}</span>
          <span v-if="isExpand">{{ t('collapse_all') }}</span>
        </a-button>

        <a-button @click="showCheckboxOrNot" type="link">
          <span v-if="!showCheckbox">{{ t('show_checkbox') }}</span>
          <span v-if="showCheckbox">{{ t('hide_checkbox') }}</span>
        </a-button>

        <a-button @click="onDisplayBy" type="link">
          <span v-if="displayBy === 'workspace'">按模块</span>
          <span v-if="displayBy === 'module'">按目录</span>
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
          :checkable="showCheckbox"
          @expand="expandNode"
          @select="selectNode"
          @check="checkNode"

          draggable
          @rightClick="onRightClick"
          @dragenter="onDragEnter"
          @drop="onDrop"
      >

        <template #title="slotProps">
          <span v-if="!slotProps.isEdit">
            <span :class="[{'no-script': noScript(slotProps.path)}]">
              {{slotProps.title}}
            </span>
          </span>

          <span v-else class="name-editor">
            <a-input v-model:value="editedData[slotProps.path]"
                     @keyup.enter=updateName(slotProps.path)
                     @click.stop
                     class="edit-input"/>

            <span class="btns">
              <CheckOutlined @click.stop="updateName(slotProps.path)"/>
              <CloseOutlined @click.stop="cancelUpdate(slotProps.path)"/>
            </span>
          </span>
        </template>
      </a-tree>

      <div v-if="contextNode.path" :style="menuStyle">
        <TreeContextMenu
            :displayBy="displayBy"
            :treeNode="contextNode"
            :onSubmit="menuClick"/>
      </div>

      <a-empty v-if="treeDataEmpty" :image="simpleImage"/>
    </div>

    <div class="actions">
      <template v-if="currWorkspace.type === 'ztf'" >
        <span class="btn-wrapper">
          <a-button :hidden="!(currProduct.id)" @click="checkoutCases">
            {{t('checkout_case')}}
          </a-button>
        </span>

        <span class="btn-wrapper">
          <a-button :hidden="!(currProduct.id && checkedKeys.length > 0)" @click="checkinCases">
            {{t('checkin_case')}}
          </a-button>
        </span>

        <span class="btn-wrapper">
          <a-button :hidden="!(checkedKeys.length > 0 && isRunning !== 'true')"
              @click="execSelected">
            {{t('exec_selected')}}
          </a-button>

          <a-button :hidden="!(checkedKeys.length > 0 && isRunning === 'true')"
              @click="execStop">
            {{t('stop')}}
          </a-button>
        </span>
      </template>

      <span v-if="currWorkspace.type !== 'ztf'" class="btn-wrapper">
        <a-button @click="toExecUnit">
          {{testToolMap[currWorkspace.type] + t('test')}}
        </a-button>
      </span>
    </div>

    <NameForm
        v-if="nameFormVisible"
        :onSubmit="createNode"
        :onCancel="() => nameFormVisible = false"
    />

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
  onMounted, onUnmounted,
  ref,
  watch
} from "vue";
import {useStore} from "vuex";
import {useI18n} from "vue-i18n";

import {ScriptData} from "../store";
import {Empty, Modal, notification} from "ant-design-vue";

import bus from "@/utils/eventBus";
import {ZentaoData} from "@/store/zentao";
import {
  setExpandedKeys,
  getScriptFilters,
  getExpandedKeys,
  setScriptFilters,
  getScriptDisplayBy,
  setScriptDisplayBy
} from "@/utils/cache";
import {
  genWorkspaceToScriptsMap,
  listFilterItems,
  getCaseIdsFromReport,
  getSyncFromInfoFromMenu, getNodeMap, getFileNodesUnderParent, updateNameReq
} from "../service";
import settings from "@/config/settings";
import {useRouter} from "vue-router";
import {DropEvent, TreeDragEvent} from "ant-design-vue/es/tree/Tree";
import {CloseOutlined, CheckOutlined} from "@ant-design/icons-vue";

import SyncFromZentao from "./syncFromZentao.vue"
import {isWindows} from "@/utils/comm";
import {testToolMap} from "@/utils/testing";
import {ExecStatus} from "@/store/exec";
import debounce from "lodash.debounce";
import throttle from "lodash.debounce";
import {expandOneKey} from "@/utils/dom";
import TreeContextMenu from "./treeContextMenu.vue"
import {ZentaoCasePrefix} from "@/utils/const";

import NameForm from "./nodeName.vue";
import {isInArray} from "@/utils/array";

export default defineComponent({
  name: 'ScriptTreePage',
  components: {
    CloseOutlined, CheckOutlined,
    TreeContextMenu, NameForm, SyncFromZentao,
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
      onTreeDataChanged()
    }, {deep: true})

    let filerItems = ref([] as any)

    const loadScripts = async () => {
      console.log(`loadScripts should be executed only once`)
      console.log(`filerType: ${filerType.value}, filerValue: ${filerValue.value}`)

      const params = {displayBy: displayBy.value, filerType: filerType.value, filerValue: filerValue.value} as any
      store.dispatch('Script/listScript', params)
    }

    // 异步加载
    const onLoadData = async (treeNode: any) => {
      console.log('onLoadData')
      await store.dispatch('Script/loadChildren', treeNode)
    }

    const onTreeDataChanged =async () => {
      getNodeMapCall()

      getExpandedKeys(currSite.value.id, currProduct.value.id).then(async cachedKeys => {
        console.log('cachedKeys', currSite.value.id, currProduct.value.id)

        if (cachedKeys) expandedKeys.value = cachedKeys

        if (!cachedKeys || cachedKeys.length === 0) {
          getOpenKeys(treeData.value[0], false) // expend first level folder
          await setExpandedKeys(currSite.value.id, currProduct.value.id, expandedKeys.value)
        }
      })
    }

    // display
    const loadDisplayBy = async () => {
      displayBy.value = await getScriptDisplayBy(currSite.value.id, currProduct.value.id)
    }

    // filters
    const loadFilterItems = async () => {
      const data = await getScriptFilters(displayBy.value, currSite.value.id, currProduct.value.id)

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
      await setScriptFilters(displayBy.value, currSite.value.id, currProduct.value.id, type, '')

      await loadFilterItems()
      await loadScripts()
    }
    const selectFilerValue = async (val) => {
      console.log('selectFilerValue', val)
      await setScriptFilters(displayBy.value, currSite.value.id, currProduct.value.id, filerType.value, val)

      await loadScripts()
    }

    const initData = debounce(async () => {
      console.log('init')
      if (!currSite.value.id) return

      await loadDisplayBy()
      await loadFilterItems()
      await loadScripts()
    }, 50)

    // only do it when switch from another pages, otherwise will called by watching currProduct method.
    if (filerValue.value.length === 0) initData()

    const expandedKeys = ref<string[]>([]);
    const getOpenKeys = (treeNode: any, openAll: boolean) => { // expand top one level if openAll is false
      if (!treeNode) return
      expandedKeys.value.push(treeNode.path)

      if (treeNode.children && openAll) {
        treeNode.children.forEach((item, index) => {
          getOpenKeys(item, openAll)
        })
      }

      console.log('keys', expandedKeys.value)
    }

    const selectedKeys = ref<string[]>([])
    const checkedKeys = ref<string[]>([])

    const replaceFields = {
      key: 'path',
    };

    let isExpand = ref(false);
    let showCheckbox = ref(false)
    let displayBy = ref('workspace')

    let tree = ref(null)

    onMounted(() => {
      console.log('onMounted', tree)
      document.addEventListener("click", clearMenu)
    })
    onUnmounted(() => {
      document.removeEventListener("click", clearMenu)
    })

    const execSelected = () => {
      console.log('execSelected')

      selectNothing()

      const leafNodes = getCheckedFileNodes()
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
    }
    const onCancel = async () => {
      console.log('onCancel')
      fromVisible.value = false
    }

    const checkinCases = () => {
      console.log('checkinCases')

      const fileNodes = getCheckedFileNodes()
      const workspaceWithScripts = genWorkspaceToScriptsMap(fileNodes)

      scriptStore.dispatch('Script/syncToZentao', workspaceWithScripts).then((resp => {
        if (resp.code === 0) {
          notification.success({message: t('sync_success')});
        } else {
          notification.error({message: t('sync_fail'), description: resp.data.msg});
        }
      }))
    }

    const selectNode = (keys) => {
      console.log('selectNode', selectedKeys.value)

      const node = treeDataMap[selectedKeys.value[0]]

      if (node.workspaceType !== 'ztf') checkNothing()

      scriptStore.dispatch('Script/getScript', node)
      if(node.type === 'file'){
          store.dispatch('tabs/open', {
                id: node.path,
                title: node.title,
                changed: false,
                type: 'script',
                data: node.path
          });
      }

      scriptStore.dispatch('Script/changeWorkspace',
          {id: node.workspaceId, type: node.workspaceType})
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
    const expandAllOrNot = (e) => {
      console.log('expandAllOrNot')
      isExpand.value = !isExpand.value

      expandedKeys.value = []
      if (isExpand.value) getOpenKeys(treeData.value[0], true)
      else getOpenKeys(treeData.value[0], false)

      setExpandedKeys(currSite.value.id, currProduct.value.id, expandedKeys.value)
    }

    const showCheckboxOrNot = (e) => {
      console.log('showCheckboxOrNot')
      showCheckbox.value = !showCheckbox.value
    }

    const onDisplayBy = () => {
      console.log('onDisplayBy')
      if (displayBy.value === 'workspace') displayBy.value = 'module'
      else displayBy.value = 'workspace'

      setScriptDisplayBy(displayBy.value, currSite.value.id, currProduct.value.id)

      loadScripts()
    }

    // context menu
    let rightVisible = false
    let contextNode = ref({} as any)
    let menuStyle = ref({} as any)
    const editedData = ref<any>({})
    const nameFormVisible = ref(false)

    const treeDataMap = {}
    const getNodeMapCall = throttle(async () => {
      getNodeMap(treeData.value[0], treeDataMap)
    }, 300)

    let rightClickedNode = {} as any
    let createAct = ''

    const onRightClick = (e) => {
      console.log('onRightClick', e.node.dataRef)
      const {event, node} = e

      rightClickedNode = node.dataRef

      const y = event.currentTarget.getBoundingClientRect().top
      const x = event.currentTarget.getBoundingClientRect().right

      contextNode.value = {
        pageX: x,
        pageY: y,
        path: rightClickedNode.path,
        title: rightClickedNode.title,
        type: rightClickedNode.type,
        workspaceId: rightClickedNode.workspaceId,
        // parentId: node.dataRef.parentId
      }

      menuStyle.value = {
        position: 'fixed',
        maxHeight: 40,
        textAlign: 'center',
        left: `${x + 10}px`,
        top: `${y + 6}px`
        // display: 'flex',
        // flexDirection: 'row'
      }
    }

    const updateName = (path) => {
      const title = editedData.value[path]
      console.log('updateName', path, title)
      updateNameReq(path, title).then((json) => {
        if (json.code === 0) {
          treeDataMap[path].title = title
          treeDataMap[path].isEdit = false
        }
      })
    }
    const cancelUpdate = (path) => {
      console.log('cancelUpdate', path)
      treeDataMap[path].isEdit = false
    }

    const createNode = (model) => {
      const arr = createAct.split('_')
      const mode = arr[1]
      const type = arr[2]

      scriptStore.dispatch('Script/createScript', {
        name: model.name, mode: mode, type: type, target: rightClickedNode.path,
        workspaceId: rightClickedNode.workspaceId, productId: currProduct.value.id,
      }).then((result) => {
        if (result) {
          notification.success({message: t('create_success')});
          nameFormVisible.value = false

          if (mode == 'child') {
            expandedKeys.value.push(rightClickedNode.path)
          }
          if (type === 'dir') {
            expandedKeys.value.push(result)
          }
          setExpandedKeys(currSite.value.id, currProduct.value.id, expandedKeys.value)

        } else {
          notification.error({message: t('create_fail')});
        }
      })
    }

    const menuClick = (act: string, node: any) => {
      console.log('menuClick', act, node)
      createAct = act
      rightClickedNode = node

      if (isInArray(act, ['add_brother_node', 'add_child_node', 'add_brother_dir', 'add_child_dir'])) {
        nameFormVisible.value = true
        return

      } else if (act === 'rename') {
        selectedKeys.value = [rightClickedNode.path]
        selectNode(selectedKeys.value)
        editedData.value[rightClickedNode.path] = treeDataMap[rightClickedNode.path].title
        treeDataMap[rightClickedNode.path].isEdit = true
        return

      } else if (act === 'remove') {
        console.log('remove ', rightClickedNode)
        const typ = rightClickedNode.type === 'file' ? t('script') : t('dir')
        Modal.confirm({
          title: '删除项目',
          content: t('confirm_delete', {
            name: rightClickedNode.title,
            typ: typ,
          }),
          okText: t('confirm'),
          cancelText: t('cancel'),
          onOk: async () => {
            removeNode(rightClickedNode.path)
          }
        });
        return

      } else if (act === 'sync_from_zentao') {
        const node = treeDataMap[rightClickedNode.path]
        const data = getSyncFromInfoFromMenu(rightClickedNode.path, node)
        data.workspaceId = treeDataMap[rightClickedNode.path].workspaceId

        scriptStore.dispatch('Script/syncFromZentao', data).then((resp => {
          if (resp.code === 0) {
            notification.success({message: t('sync_success')});
          } else {
            notification.error({message: t('sync_fail'), description: resp.data.msg});
          }
        }))

        return
      } else if (act === 'sync_to_zentao') {
        const node = treeDataMap[rightClickedNode.path]

        const fileNodes = getFileNodesUnderParent(node)
        const workspaceWithScripts = genWorkspaceToScriptsMap(fileNodes)

        scriptStore.dispatch('Script/syncToZentao', workspaceWithScripts).then((resp => {
          if (resp.code === 0) {
            notification.success({message: t('sync_success')});
          } else {
            notification.error({message: t('sync_fail'), description: resp.data.msg});
          }
        }))

        return
      }

      const arr = act.split('_')
      addNode(arr[1], arr[2])

      clearMenu()
    }

    const clearMenu = () => {
      console.log('clearMenu')
      contextNode.value = ref({})
    }
    const addNode = (mode, type) => {
      console.log('addNode', rightClickedNode.path)
      store.dispatch('Interface/createInterface',
          {mode: mode, type: type, target: rightClickedNode.path, name: type === 'dir' ? '新目录' : '新接口'})
          .then((newNode) => {
            console.log('newNode', newNode)
            selectedKeys.value = [newNode.id] // select new node
            expandOneKey(treeDataMap, newNode.parentId, expandedKeys.value) // expend new node
          })
    }
    const removeNode = (path) => {
      store.dispatch('Script/deleteScript', path);
    }

    const onDragEnter = (info: TreeDragEvent) => {
      console.log('onDragEnter', info);
    };

    const onDrop = (info: DropEvent) => {
      console.log('onDrop', info);

      const dragKey = info.dragNode.eventKey;
      const dropKey = info.node.eventKey;
      let dropPos = info.dropPosition > 1 ? 1 : info.dropPosition;
      if (!treeDataMap[dropKey].isDir && dropPos === 0) dropPos = 1
      console.log(dragKey, dropKey, dropPos);

      store.dispatch('Interface/moveInterface', {dragKey: dragKey, dropKey: dropKey, dropPos: dropPos});
    }

    const getCheckedFileNodes = (): string[] => {
      console.log('getCheckedFileNodes')

      let arr = [] as string[]
      checkedKeys.value.forEach(k => {
        if (treeDataMap[k].type === 'file') {
          arr.push(treeDataMap[k])
        }
      })
      return arr
    }

    const noScript = (str) => {
      if (str.indexOf(ZentaoCasePrefix) == 0) {
        return true
      }
      return false
    }

    return {
      t,
      isWin,

      currSite,
      currProduct,
      treeData,
      currWorkspace,
      testToolMap,
      nameFormVisible,
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
      showCheckbox,
      displayBy,
      expandAllOrNot,
      showCheckboxOrNot,
      onDisplayBy,
      tree,
      expandedKeys,
      selectedKeys,
      checkedKeys,
      simpleImage: Empty.PRESENTED_IMAGE_SIMPLE,

      fromTitle,
      fromVisible,
      onSave,
      onCancel,
      noScript,

      rightVisible,
      contextNode,
      menuStyle,
      treeDataMap,
      editedData,
      rightClickedNode,
      updateName,
      cancelUpdate,
      onRightClick,
      menuClick,
      clearMenu,
      createNode,
      removeNode,
      onDragEnter,
      onDrop,
    }
  }

})
</script>

<style lang="less">
.tree-main {
  .tree-panel {
    margin-top: -8px;
    margin-left: -10px;

    .ant-tree > li > span {
      display: none !important;
    }

    .ant-tree {
      .ant-tree-node-content-wrapper {
        display: inline-block;
        height: 26px;
        margin: 0;
        padding: 0;
      }
      .ant-tree-title {
        display: inline-block;
        height: 26px;
        margin: -2px 0 0 0;
        padding: 0;

        .name-editor {
          display: inline-block;
          height: 26px;
          margin: 0;
          padding: 0;

          .edit-input {
            height: 26px;
            line-height: 26px;

            margin: 0 3px 0 0;
            padding: 0 5px;
          }
        }
      }
    }
  }

  .no-script {
    color: #a9aeb4;
  }
}
</style>

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
      width: 170px;
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
  }

  .actions {
    padding: 4px;
    height: 40px;
    text-align: center;

    .btn-wrapper {
      display: inline-block;
      width: 100px;

      .ant-btn {
        width: 96px;
        margin: 0 5px;
        padding: 3px 3px;
      }
    }
  }
}

</style>

<style lang="less">

</style>
