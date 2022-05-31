<template>
  <div class="workdir">
    <Tree :data="treeData" :checkable="checkable" ref="treeRef" @active="selectNode" @check="checkNode"
      @clickToolbar="onToolbarClicked" />
    <FormNode :show="showModal" @submit="createNode" @cancel="modalClose" ref="formNode" />
  </div>
</template>

<script setup lang="ts">
import { useI18n } from "vue-i18n";
import { useStore } from "vuex";
import { ZentaoData } from "@/store/zentao";
import { ScriptData } from "@/views/script/store";
import { WorkspaceData } from "@/store/workspace";
import { expandOneKey, resizeWidth } from "@/utils/dom";
import Form from "./Form.vue";
import FormItem from "./FormItem.vue";
import { useForm } from "@/utils/form";
import Tree from "./Tree.vue";
import notification from "@/utils/notification";
import { computed, defineExpose, onMounted, onUnmounted, ref, watch } from "vue";

import bus from "@/utils/eventBus";
import {
  getExpandedKeys,
  getScriptDisplayBy,
  getScriptFilters,
  setExpandedKeys,
} from "@/utils/cache";
import {
  getCaseIdsFromReport,
  getNodeMap,
  listFilterItems,
} from "@/views/script/service";
import { useRouter } from "vue-router";
import { isWindows } from "@/utils/comm";
import debounce from "lodash.debounce";
import throttle from "lodash.debounce";
import Modal from "@/utils/modal"
import FormNode from "./FormNode.vue";

const { t } = useI18n();

const zentaoStore = useStore<{ Zentao: ZentaoData }>();
const currSite = computed<any>(() => zentaoStore.state.Zentao.currSite);
const currProduct = computed<any>(() => zentaoStore.state.Zentao.currProduct);

const scriptStore = useStore<{ Script: ScriptData }>();
const currWorkspace = computed<any>(() => scriptStore.state.Script.currWorkspace);

const workspaceStore = useStore<{ Workspace: WorkspaceData }>();
const isWin = isWindows()

zentaoStore.dispatch('Zentao/fetchLangs')
const langs = computed<any[]>(() => zentaoStore.state.Zentao.langs);

const router = useRouter();
let workspace = router.currentRoute.value.params.workspace as string
workspace = workspace === '-' ? '' : workspace
let seq = router.currentRoute.value.params.seq as string
seq = seq === '-' ? '' : seq
let scope = router.currentRoute.value.params.scope as string
scope = scope === '-' ? '' : scope

const store = useStore<{ Script: ScriptData }>();
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
    { id: node.workspaceId, type: node.workspaceType })

  currentNode.value = node;
  if (e.event.key == 'runTest') {
    runTest(currentNode);
  } else if (e.event.key == 'createFile' || e.event.key == 'createWorkspace') {
    showModal.value = true;
  } else if (e.event.key === 'deleteWorkspace') {
    Modal.confirm({
      title: t('delete'),
      content: t('confirm_to_delete_workspace', { p: node.title }),
      showOkBtn: true
    },
      {
        "onOk": () => {
          workspaceStore.dispatch('Workspace/removeWorkspace', node.path)
            .then((response) => {
              if (response) {
                notification.success({ message: t('delete_success') });
                loadScripts()
              }
            })
        }
      }
    )
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
}, { deep: true })

watch(treeData, (currConfig) => {
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
  console.log('selectNode', activeNode.activeID)

  const node = treeDataMap[activeNode.activeID]
  if (node.workspaceType !== 'ztf') checkNothing()

  scriptStore.dispatch('Script/getScript', node)
  if (node.type === 'file') {
    store.dispatch('tabs/open', {
      id: 'workspace-' + node.path,
      title: node.title,
      changed: false,
      type: 'script',
      data: node.path
    });
  }

  scriptStore.dispatch('Script/changeWorkspace',
    { id: node.workspaceId, type: node.workspaceType })
}

const checkNode = (checkedKeys) => {
  console.log('checkNode', checkedKeys.checked)
  //   scriptStore.dispatch('Script/changeWorkspace',
  //       {id: e.node.dataRef.workspaceId, type: e.node.dataRef.workspaceType})
}

const checkNothing = () => {
  checkedKeys.value = []
}

let contextNode = ref({} as any)
let menuStyle = ref({} as any)
const editedData = ref<any>({})
const nameFormVisible = ref(false)

const treeDataMap = {}
const getNodeMapCall = throttle(async () => {
  treeData.value.forEach(item => {
    getNodeMap(item, treeDataMap)
  })
}, 300)

let rightClickedNode = {} as any
let createAct = ''

const formNode = ref({} as any)
const createNode = (formData) => {
  const mode = 'child';
  let type = 'dir';
  if (currentNode.value.isLeaf) {
    type = 'node';
  }
  scriptStore.dispatch('Script/createScript', {
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
    } else {
      notification.error({ message: t('create_fail') });
    }
  })
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
</style>
