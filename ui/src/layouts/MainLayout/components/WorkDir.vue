<template>
  <div class="workdir">
    <Tree 
      :data="treeData" 
      :checkable="checkable" 
      ref="treeRef" 
      @active="selectNode" 
      @check="checkNode" 
      @clickToolbar="onToolbarClicked"
    />

    <ZModal
        :showModal="showModal"
        @onCancel="modalClose"
        @onOk="createNode"
        :title="t('pls_name')"
    >

      <Form labelCol="100px" wrapperCol="60">
        <FormItem name="name" :label="t('name')" :info="validateInfos.name">
          <input v-model="modelRef.name" class="form-control"/>
        </FormItem>
        <FormItem v-if="currentNode.path==''" name="path" :label="t('path')" :info="validateInfos.path">
          <input v-model="modelRef.path" class="form-control"/>
        </FormItem>
        <FormItem v-if="currentNode.path==''" name="type" :label="t('type')" :info="validateInfos.type">
          <select name="type" v-model="modelRef.type" class="form-control">
            <option v-for="item in testTypes" :key="item.value" :value="item.value">{{ item.label }}</option>
          </select>
        </FormItem>
        <FormItem v-if="currentNode.path==''" name="lang" :label="t('default_lang')" :info="validateInfos.lang">
        <select name="type" v-model="modelRef.lang" class="form-control">
            <option v-for="item in langs" :key="item.code" :value="item.code">{{ item.name }}</option>
          </select>
        </FormItem>
      </Form>
    </ZModal>
  </div>
</template>

<script setup lang="ts">
import {useI18n} from "vue-i18n";
import {useStore} from "vuex";
import {ZentaoData} from "@/store/zentao";
import {ScriptData} from "@/views/script/store";
import {expandOneKey, resizeWidth} from "@/utils/dom";
import Form from "./Form.vue";
import FormItem from "./FormItem.vue";
import {useForm} from "@/utils/form";
import Tree from "./Tree.vue";
import ZModal from './Modal.vue';
import notification from "@/utils/notification";
import {unitTestTypesDef, ZentaoCasePrefix, ztfTestTypesDef} from "@/utils/const";

import {computed, defineExpose, onMounted, onUnmounted, ref, watch, getCurrentInstance} from "vue";

import {Modal} from "ant-design-vue";

import bus from "@/utils/eventBus";
import {
  getExpandedKeys,
  getScriptDisplayBy,
  getScriptFilters,
  setExpandedKeys,
  setScriptDisplayBy,
  setScriptFilters
} from "@/utils/cache";
import {
  genWorkspaceToScriptsMap,
  getCaseIdsFromReport,
  getFileNodesUnderParent,
  getNodeMap,
  getSyncFromInfoFromMenu,
  listFilterItems,
  updateNameReq
} from "@/views/script/service";
import settings from "@/config/settings";
import {useRouter} from "vue-router";
import {DropEvent, TreeDragEvent} from "ant-design-vue/es/tree/Tree";
import {isWindows} from "@/utils/comm";
import {ExecStatus} from "@/store/exec";
import debounce from "lodash.debounce";
import throttle from "lodash.debounce";
import {isInArray} from "@/utils/array";
import {PageType} from "@/store/tabs";

const {t} = useI18n();

const zentaoStore = useStore<{ Zentao: ZentaoData }>();
const currSite = computed<any>(() => zentaoStore.state.Zentao.currSite);
const currProduct = computed<any>(() => zentaoStore.state.Zentao.currProduct);

const scriptStore = useStore<{ Script: ScriptData }>();
const currWorkspace = computed<any>(() => scriptStore.state.Script.currWorkspace);

const isWin = isWindows()

const testTypes = ref([...ztfTestTypesDef, ...unitTestTypesDef])
zentaoStore.dispatch('Zentao/fetchLangs')
const langs = computed<any[]>(() => zentaoStore.state.Zentao.langs);

const fromTitle = ref('从禅道同步用例')
const fromVisible = ref(false)

const router = useRouter();
let workspace = router.currentRoute.value.params.workspace as string
workspace = workspace === '-' ? '' : workspace
let seq = router.currentRoute.value.params.seq as string
seq = seq === '-' ? '' : seq
let scope = router.currentRoute.value.params.scope as string
scope = scope === '-' ? '' : scope

const execStore = useStore<{ Exec: ExecStatus }>();
const isRunning = computed<any>(() => execStore.state.Exec.isRunning);

const store = useStore<{ Script: ScriptData }>();
const treeDataEmpty = computed<boolean>(() => !(treeData.value.length > 0 &&
    treeData.value[0] && treeData.value[0].children))

const filerType = ref('')
const filerValue = ref('')
const showModal = ref(false)
const currentNode = ref({}) // parent node for create node

onMounted(() => {
  console.log('onMounted')
  initData();
  setTimeout(() => {
    resizeWidth('main', 'left', 'splitter-h', 'right', 380, 800)
  }, 600)
})

const onToolbarClicked = (e) => {
  const node = e.node == undefined ? treeDataMap[''] : treeDataMap[e.node.id]
  scriptStore.dispatch('Script/changeWorkspace',
      {id: node.workspaceId, type: node.workspaceType})

  currentNode.value = node;
  if(e.event == undefined){
      e.event = {key : 'createWorkspace'};// create workspace
  }
  if (e.event.key == 'runTest') {
    runTest(currentNode);
  } else if (e.event.key == 'createFile' || e.event.key == 'createWorkspace') {
    showModal.value = true;
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

const modelRef = ref({})
const rulesRef = ref({
  name: [
    {required: true, msg: t('pls_name')},
  ],
  path: [
    {required: true, msg: t('pls_workspace_path')},
  ],
  lang: [
    {required: true, msg: t('select_ui_lang')},
  ],
  type: [
    {required: true, msg: t('pls_workspace_type')},
  ],
})

const {validate, reset, validateInfos} = useForm(modelRef, rulesRef);

const treeRef = ref<{ isAllCollapsed: () => boolean, toggleAllCollapsed: () => void }>();

let treeData = computed<any>(() => scriptStore.state.Script.list);

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

const onTreeDataChanged = async () => {
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

const selectNode = (activeNode) => {
  console.log('selectNode', activeNode.activeID)

  const node = treeDataMap[activeNode.activeID]
  if (node.workspaceType !== 'ztf') checkNothing()

  scriptStore.dispatch('Script/getScript', node)
  if (node.type === 'file') {
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

const checkNode = (checkedKeys) => {
  console.log('checkNode', checkedKeys.checked)
  scriptStore.dispatch('Script/setCheckedNodes', checkedKeys.checked)
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

const createWorkSpace = () => {
  if(validate()){
  store.dispatch('Workspace/save', modelRef.value).then((response) => {
    if (response) {
      modelRef.value = {};
      notification.success({message: t('save_success')});
      showModal.value = false;
    }
  })
  }

};


const createNode = () => {
  if (currentNode.value.path == '') {
    createWorkSpace();
    return;
  }
  const mode = 'child';
  let type = 'dir';
  if (currentNode.value.isLeaf) {
    type = 'node';
  }
  scriptStore.dispatch('Script/createScript', {
    name: modelRef.value.name, mode: mode, type: type, target: currentNode.value.path,
    workspaceId: currentNode.value.workspaceId, productId: currProduct.value.id,
  }).then((result) => {
    if (result) {
      showModal.value = false;
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

  store.dispatch('Script/moveScript', {dragKey: dragKey, dropKey: dropKey, dropPos: dropPos});
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
  onToolbarClicked
});
</script>

<style lang="less" scoped>
.workdir {
  height: calc(100vh - 80px);
}
.form-control{
    width: 100%;
    color: #495057;
    background-color: #fff;
    border: 1px solid #ced4da;
    border-radius: .25rem;
    transition: border-color .15s ease-in-out,box-shadow .15s ease-in-out;
}
.z-form-item-label{
    font-weight: 400;
    color: #212529;
    text-align: left;
    box-sizing: border-box;
    display: inline-block;
    position: relative;
    width: 100%;
    padding-right: 15px;
    padding-left: 15px;
    flex: 0 0 16.666667%;
    max-width: 16.666667%;
    padding-top: calc(.375rem + 1px);
    padding-bottom: calc(.375rem + 1px);
    margin-bottom: 0;
    line-height: 1.5;
}
.z-form-item{
    display: flex;
    align-items: center;
}
.form-control:focus {
    color: #495057;
    background-color: #fff;
    border-color: #80bdff;
    outline: 0;
    box-shadow: 0 0 0 0.2rem rgb(0 123 255 / 25%);
}
</style>
