<template>
  <div class="left-pannel-contain">
    <div :class="checkable && checkedKeys.length ? 'workdir-with-btn' : 'workdir'">
      <Tree
        :data="treeData"
        :checkable="checkable"
        ref="treeRef"
        @active="selectNode"
        @rightClick="onRightClick"
        @check="checkNode"
        @clickToolbar="onToolbarClicked"
        @collapse="expandNode"
        :defaultCollapsedMap="collapsedMap"
        :defaultCollapsed="true"
      />
      <FormNode :show="showModal" @submit="createNode" @cancel="modalClose" :path="currentNode.path" :name="currentNode.title" ref="formNode" />
    </div>
    <Button
      v-if="checkable && checkedKeys.length"
      class="rounded border primary-pale run-selected" icon="run-all"
      :label="t('exec_selected')"
      @click="execSelected"
     />

    <div v-if="contextNode.id && rightVisible" :style="menuStyle">
      <TreeContextMenu :treeNode="contextNode" :clipboardData="clipboardData" :onMenuClick="menuClick" :siteId="currSite.id"/>
    </div>
    <FormSyncFromZentao v-if="showSyncFromZentaoModal"
      :show="showSyncFromZentaoModal"
      @submit="syncFromZentaoSubmit"
      @cancel="showSyncFromZentaoModal = !showSyncFromZentaoModal"
      :workspaceId="syncFromZentaoWorkspaceId"
      ref="syncFromZentaoRef"
    />
    <FormWorkspace
      v-if="showWorkspaceModal"
      :show="showWorkspaceModal"
      @submit="createWorkSpace"
      @cancel="modalWorkspaceClose"
      ref="formWorkspace"
      :workspaceId="currentNode.workspaceId"
     />
  </div>
</template>

<script setup lang="ts">
import { useI18n } from "vue-i18n";
import { useStore } from "vuex";
import { StateType as GlobalData } from "@/store/global";
import { ZentaoData } from "@/store/zentao";
import { ScriptData } from "@/views/script/store";
import { WorkspaceData } from "@/store/workspace";
import {getContextMenuStyle, resizeWidth} from "@/utils/dom";
import Tree from "@/components/Tree.vue";
import notification from "@/utils/notification";
import { computed, defineExpose, onMounted, onUnmounted, ref, watch, onBeforeUnmount } from "vue";
import Button from '@/components/Button.vue';
import TreeContextMenu from './TreeContextMenu.vue';
import FormSyncFromZentao  from "./FormSyncFromZentao.vue";
import FormWorkspace from "@/views/workspace/FormWorkspace.vue";

import bus from "@/utils/eventBus";
import {getExpandedKeys, getScriptDisplayBy, getScriptFilters, setExpandedKeys } from "@/utils/cache";
import {getCaseIdsFromReport, getNodeMap, listFilterItems, getFileNodesUnderParent, genWorkspaceToScriptsMap } from "@/views/script/service";
import { useRouter } from "vue-router";
import { isWindows } from "@/utils/comm";
import debounce from "lodash.debounce";
import Modal from "@/utils/modal"
import FormNode from "./FormNode.vue";
import settings from "@/config/settings";

const { t } = useI18n();

const store = useStore<{ global: GlobalData, Zentao: ZentaoData, Script: ScriptData, Workspace: WorkspaceData }>();
const global = computed<any>(() => store.state.global.tabIdToWorkspaceIdMap);
const currSite = computed<any>(() => store.state.Zentao.currSite);
const currProduct = computed<any>(() => store.state.Zentao.currProduct);

const currWorkspace = computed<any>(() => store.state.Script.currWorkspace);

const isWin = isWindows()

store.dispatch('Zentao/fetchLangs')
const langs = computed<any[]>(() => store.state.Zentao.langs);

const router = useRouter();
let workspace = router.currentRoute.value.params.workspace as string
workspace = workspace === '-' ? '' : workspace
let seq = router.currentRoute.value.params.seq as string
seq = seq === '-' ? '' : seq
let scope = router.currentRoute.value.params.scope as string
scope = scope === '-' ? '' : scope

const filerType = ref('')
const filerValue = ref('')
const showModal = ref(false)
const toolbarAction = ref('')
const currentNode = ref({} as any) // parent node for create node
const collapsedMap = ref({} as any)
const checkedKeys = ref<string[]>([])
const showSyncFromZentaoModal = ref(false);
const syncFromZentaoWorkspaceId = ref(0);

onMounted(() => {
  console.log('onMounted')
  initData();
  setTimeout(() => {
    resizeWidth('main', 'left', 'splitter-h', 'right', 380, 800)
  }, 600)
  bus.on(settings.eventWebSocketMsg, onWebsocketMsgEvent);
})

onBeforeUnmount( () => {
  console.log('onBeforeUnmount')
  bus.off(settings.eventWebSocketMsg, onWebsocketMsgEvent);
})

const onWebsocketMsgEvent = (data: any) => {
  console.log('WebsocketMsgEvent in WatchFile', data.msg)

  let item = JSON.parse(data.msg)
  if(item.category == 'watch'){
    loadScripts()
  }
}

const onToolbarClicked = (e) => {
  const node = e.node == undefined ? treeDataMap.value[''] : treeDataMap.value[e.node.id]
  store.dispatch('Script/changeWorkspace',
    { id: node.workspaceId, type: node.workspaceType })

  currentNode.value = node;
  switch (e.event.key) {
    case 'runTest':
      runTest(currentNode);
      break;
    case 'createFile':
    case 'createWorkspace':
    case 'createDir':
        currentNode.value = {}
      showModal.value = true;
      toolbarAction.value = e.event.key;
      break;
    case 'editWorkspace':
      showWorkspaceModal.value = true;
      break;
    case 'deleteWorkspace':
      Modal.confirm({
          title: t('delete'),
          content: t('confirm_to_delete_workspace', { p: node.title }),
          showOkBtn: true
        },
        {
          "onOk": () => {
            store.dispatch('Workspace/removeWorkspace', node.path)
              .then((response) => {
              if (response) {
                notification.success({ message: t('delete_success') });
                loadScripts()
              }
            })
          }
        }
      );
    break;
    case 'runScript':
      console.log('run script', currentNode.value);
      bus.emit(settings.eventExec,
        {execType: currentNode.value.workspaceType === 'ztf' ? 'ztf' : 'unit', scripts: currentNode.value.isLeaf ? [currentNode.value] : currentNode.value.children});
      break;
    case 'checkinCase':
      console.log('checkin case', currentNode.value);
      break;
    case 'checkoutCase':
      console.log('checkout case', currentNode.value);
      break;
  }
}

const runTest = (node) => {
  console.log('runTest', node.value)

  store.dispatch('tabs/open', {
    id: 'workspace-' + node.value.workspaceId,
    title: node.value.title,
    type: 'execUnit',
    changed: false,
    data: {
      workspaceId: node.value.workspaceId,
      workspaceType: node.value.workspaceType,
    }
  });
}

const modalClose = () => {
  showModal.value = false;
}

const treeRef = ref<{ isAllCollapsed: () => boolean, toggleAllCollapsed: () => void }>();

let treeData = computed<any>(() => store.state.Script.list);

const checkable = ref(false);

function toggleCheckable(toggle?: boolean) {
  if (toggle === undefined) {
    toggle = !checkable.value;
  }
  checkable.value = toggle;
}


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
}, { deep: true })

watch(treeData, () => {
  console.log('watch treeData', treeData.value)
  onTreeDataChanged()
}, { deep: true })

let filerItems = ref([] as any)

const loadScripts = async () => {
  console.log(`loadScripts should be executed only once`)
  console.log(`filerType: ${filerType.value}, filerValue: ${filerValue.value}`)

  const params = { displayBy: displayBy.value, filerType: filerType.value, filerValue: filerValue.value } as any
  store.dispatch('Script/listScript', params)
}

const onTreeDataChanged = async () => {
  getNodeMapCall()

  getExpandedKeys(currSite.value.id, currProduct.value.id).then(async cachedKeys => {
    console.log('cachedKeys', currSite.value.id, currProduct.value.id)

    if (cachedKeys) expandedKeys.value = cachedKeys
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
    treeNode.children.forEach((item) => {
      getOpenKeys(item, openAll)
    })
  }

  console.log('keys', expandedKeys.value)
}

watch(expandedKeys, () => {
    console.log('watch expandedKeys')
    for (let treeDataKey in treeDataMap.value) {
        collapsedMap.value[treeDataKey] = expandedKeys.value.indexOf(treeDataKey) !== -1 ? false : true
    }
}, { deep: true })

let isExpand = ref(false);
let showCheckbox = ref(false)
let displayBy = ref('workspace')

let tree = ref(null)

onMounted(() => {
  console.log('onMounted', tree)
})
onUnmounted(() => {
  console.log('onUnmounted', tree)
})

const selectNode = (activeNode) => {
  console.log('selectNode', activeNode.activeID, global.value)

  const node = treeDataMap.value[activeNode.activeID]
  console.log('node id', node.caseId)
  if (node.workspaceType !== 'ztf') checkNothing()

  store.dispatch('Script/getScript', node)
  store.dispatch('Result/getStatistic', node)
  if (node.type === 'file') {
    const tabId = node.workspaceType === 'ztf' && node.path.indexOf('.exp') !== node.path.length - 4
        ? 'script-' + node.path : 'code-' + node.path
    global.value[tabId] = node.workspaceId

    store.dispatch('tabs/open', {
      id: tabId,
      title: node.title,
      changed: false,
      type: 'script',
      data: node.path
    });
  }

  store.dispatch('Script/changeWorkspace',
    { id: node.workspaceId, type: node.workspaceType })
}

const checkNode = (keys) => {
  console.log('checkNode', keys.checked)
  store.dispatch('Script/setCheckedNodes', keys.checked)
  let checkedKeysTmp:string[] = [];
  for(let checkedKey in keys.checked){
    if(keys.checked[checkedKey] === true){
        checkedKeysTmp.push(checkedKey)
    }
  }
  checkedKeys.value = checkedKeysTmp;
}

const checkNothing = () => {
  checkedKeys.value = []
}

const execSelected = () => {
    let arr = [] as string[]
    checkedKeys.value.forEach(checkedKey => {
      if (treeDataMap.value[checkedKey]?.type === 'file') {
        arr.push(treeDataMap.value[checkedKey])
      }
    })
    bus.emit(settings.eventExec, { execType: 'ztf', scripts: arr });
}

const nameFormVisible = ref(false)

const treeDataMap = computed<any>(() => store.state.Script.treeDataMap);
const getNodeMapCall = debounce(async () => {
  treeData.value.forEach(item => {
    getNodeMap(item, treeDataMap.value)
  })
}, 300)

let rightClickedNode = {} as any

const formNode = ref({} as any)
const createNode = (formData) => {
  if(formData.path != ""){
    store.dispatch('Script/renameScript', formData)
    formNode.value.clearFormData()
    showModal.value = false;
    return;
  }
  const mode = 'child';
  let type = 'dir';
  if(toolbarAction.value === 'createFile') type = 'node'
  store.dispatch('Script/createScript', {
    name: formData.name, mode: mode, type: type, target: currentNode.value.path,
    workspaceId: currentNode.value.workspaceId, productId: currProduct.value.id,
  }).then((result) => {
    if (result) {
      formNode.value.clearFormData()
      showModal.value = false;
      notification.success({ message: t('create_success') });
      nameFormVisible.value = false

      if (mode == 'child') {
        expandedKeys.value.push(rightClickedNode.path)
      }
      if (type === 'dir') {
        expandedKeys.value.push(result)
      }
      setExpandedKeys(currSite.value.id, currProduct.value.id, expandedKeys.value)
    }
  })
}

const expandNode = (expandedKeysMap) => {
    console.log('expandNode', expandedKeysMap.collapsed)
    let expandedKeysTmp:string[] = [];
    for(let key in expandedKeysMap.collapsed){
        if(expandedKeysMap.collapsed[key]){
            expandedKeys.value.forEach((item, index) => {
                if(item === key){
                    expandedKeys.value.splice(index, 1)
                }
            })
        }else{
            expandedKeysTmp.push(key)
        }
    }
    expandedKeys.value = expandedKeysTmp;
    console.log('expandkeys', expandedKeys.value)
    setExpandedKeys(currSite.value.id, currProduct.value.id, expandedKeys.value)
}

let menuStyle = ref({} as any)
let contextNode = ref({} as any)
const clipboardAction = ref('')
const clipboardData = ref({} as any)

const rightVisible = ref(false)

const onRightClick = (e) => {
  console.log('onRightClick', e)
  const {event, node} = e

  const contextNodeData = treeDataMap.value[node.id]
  contextNode.value = {
    id: contextNodeData.id,
    title: contextNodeData.title,
    type: contextNodeData.type,
    isLeaf: contextNodeData.isLeaf,
    workspaceId: contextNodeData.workspaceId,
    workspaceType: contextNodeData.workspaceType,
  }

  menuStyle.value = getContextMenuStyle(event.currentTarget.getBoundingClientRect().right, event.currentTarget.getBoundingClientRect().top, 260)

  rightVisible.value = true
}

const menuClick = (menuKey: string, targetId: number) => {
  const contextNodeData = treeDataMap.value[targetId]
  console.log('menuClick', menuKey, targetId, contextNodeData)

  if(menuKey === 'exec'){
    execScript(contextNodeData)
  } else if (menuKey === 'copy' || menuKey === 'cut') {
    clipboardAction.value = menuKey
    clipboardData.value = contextNodeData

  } else if (menuKey === 'paste') {
    console.log(clipboardData.value)
    const data = {
      srcKey: clipboardData.value.id,
      srcType: clipboardData.value.type,
      srcWorkspaceId: clipboardData.value.workspaceId,
      distKey: contextNodeData.id,
      distType: contextNodeData.type,
      distWorkspaceId: contextNodeData.workspaceId,
      action: clipboardAction.value,
    }
    store.dispatch('Script/pasteScript', data)

  } else if (menuKey === 'delete') {
    Modal.confirm({
      title: t("confirm_delete", {
        name: contextNodeData.title,
        typ: t("node"),
      }),
      okText: t("confirm"),
      cancelText: t("cancel"),
      onOk: async () => {
        store.dispatch('Script/deleteScript', contextNodeData.id)
      },
    });
} else if (menuKey === 'rename') {
    showModal.value = true;
    currentNode.value = contextNodeData;
  } else if(menuKey == 'sync-from-zentao'){
    syncFromZentao(contextNodeData)
  } else if(menuKey === 'sync-to-zentao'){
    checkinCases(contextNodeData)
  } else {
    clipboardAction.value = ''
    clipboardData.value = {}

    if (menuKey === 'open-in-explore') {
      const { ipcRenderer } = window.require('electron')
      ipcRenderer.send(settings.electronMsg, {action: 'openInExplore', path: contextNodeData.path})

    } else if (menuKey === 'open-in-terminal') {
      const { ipcRenderer } = window.require('electron')

      ipcRenderer.send(settings.electronMsg, {action: 'openInTerminal', path: contextNodeData.path})

    }
  }

  clearMenu()
}
const checkinCases = (node) => {
  if(node.workspaceType == 'ztf'){
    console.log('checkinCases')
    const fileNodes = getFileNodesUnderParent(node)
    const workspaceWithScripts = genWorkspaceToScriptsMap(fileNodes)
    store.dispatch('Script/syncToZentao', workspaceWithScripts).then((resp => {
    if (resp.code === 0) {
        notification.success({message:
              t('sync_success', {
                success: resp.data.success,
                ignore: resp.data.total - resp.data.success
              }
          )});
    } else {
        notification.error({message: t('sync_fail'), description: resp.data.msg});
    }
    }))
  }
}
const syncFromZentao = (node) => {
    if(node.workspaceType == 'ztf'){
      if(node.type == 'workspace'){
        showSyncFromZentaoModal.value = true;
        syncFromZentaoWorkspaceId.value = node.workspaceId;
      }else if(node.type == 'dir'){
        checkoutCases(node.workspaceId, node)
      }else if(node.type == 'file'){
        checkout(node.workspaceId, node.caseId, node.path)
      }else if(node.type == 'module'){
        checkoutFromModule(node.workspaceId, node)
      }
    }
}
const checkoutCases = (workspaceId, node) => {
    if(node.children == undefined || node.children.length == 0){
        return;
    }
    node.children.forEach(item => {
        if(item.type == 'dir'){
            checkoutCases(workspaceId, item)
        }else if(item.type == 'file' && item.caseId){
            checkout(workspaceId, item.caseId, item.path, false)
        }
    });
    notification.success({
        message: t('sync_from_zentao_success', {
                success: node.children.length,
              }),
      });
}
const checkoutFromModule = (workspaceId, node) => {
    if(node.children == undefined || node.children.length == 0){
        return;
    }
    console.log('checkout from module', workspaceId, node.children[0].moduleId)
    const data = {moduleId: node.children[0].moduleId, workspaceId: workspaceId}
    store.dispatch('Script/syncFromZentao', data).then((resp => {
    if (resp.code === 0) {
      notification.success({
        message: t('sync_from_zentao_success', {
            success: resp.data.length,
        }),
      });
    } else {
        notification.error({
          message: resp.data.msg,
        });
    }
    }))
}
const checkout = (workspaceId, caseId, path, successNotice = true) => {
    console.log('checkout', workspaceId, caseId, path)
    const data = {caseId: caseId, workspaceId: workspaceId, casePath: path}
    store.dispatch('Script/syncFromZentao', data).then((resp => {
    if (resp.code === 0) {
      successNotice && notification.success({
        message: t('sync_from_zentao_success', {
          success: 1,
        }),
      });
    } else {
        notification.error({
          message: resp.data.msg,
        });
    }
    }))
}
const syncFromZentaoRef = ref({} as any)
const syncFromZentaoSubmit = (model) => {
  store.dispatch("Script/syncFromZentao", model).then((resp) => {
    if (resp.code === 0) {
      notification.success({
        message: t("sync_from_zentao_success", {success: resp.data == undefined? 0 : resp.data.length, ignore:0}),
      });
      showSyncFromZentaoModal.value = false;
      syncFromZentaoRef.value.clearFormData()
    } else {
      notification.error({
        message: resp.data.msg,
      });
    }
  });
}
const execScript = (node) => {
  if(node.workspaceType !== 'ztf'){
    runTest(ref(node));
  }else{
    bus.emit(settings.eventExec,
        {execType: 'ztf', scripts: node.type === 'file' ? [node] : node.children});
  }
}
const clearMenu = () => {
  console.log('clearMenu')
  contextNode.value = ref(null)
}

const showWorkspaceModal = ref(false)
const formWorkspace = ref({} as any)
const createWorkSpace = (formData) => {
    store.dispatch('Workspace/save', formData).then((response) => {
        if (response) {
            formWorkspace.value.clearFormData()
            notification.success({message: t('save_success')});
            showWorkspaceModal.value = false;
            loadScripts()
        }
    })
};
const modalWorkspaceClose = () => {
  showWorkspaceModal.value = false;
}

defineExpose({
  get isCheckable() {
    return checkable.value;
  },
  get isAllCollapsed() {
    return treeRef.value?.isAllCollapsed();
  },
  toggleAllCollapsed() {
    return treeRef.value?.toggleAllCollapsed();
  },
  toggleCheckable,
  onToolbarClicked,
  loadScripts
});

onMounted(() => {
  console.log('onMounted')
  document.addEventListener("click", clearMenu)
})
onUnmounted(() => {
  document.removeEventListener("click", clearMenu)
})

</script>

<style lang="less" scoped>
.left-pannel-contain{
  text-align: center;
  .workdir {
      height: calc(100vh - 80px);
      overflow: auto;
      text-align: left;
  }
  .workdir-with-btn {
      height: calc(100vh - 120px);
      overflow: auto;
      text-align: left;
  }
  .run-selected{
    max-width: 100px;
    margin: auto;
    text-align: center;
    margin-top: 10px;
  }
}
</style>
