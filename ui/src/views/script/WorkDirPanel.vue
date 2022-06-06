<template>
  <Panel class="workdir-panel">
    <template #header>
      <ButtonList :gap="2" class="workdir-panel-nav">
        <Button id="displayByMenuToggle"
                :label="te('by_' + displayBy) ? t('by_' + displayBy) : t('by_workspace')"
                labelClass="strong"
                class="rounded pure padding-sm-h"
                suffix-icon="caret-down"/>

        <DropdownMenu :items="displayTypes"
                  :checkedKey="displayBy"
                  @click="onDisplayByChanged"
                  toggle="#displayByMenuToggle">
        </DropdownMenu>

        <Button id="displayByFilter"
                :label="te('by_' + filerType) ? t('by_' + filerType) : t('by_workspace')"
                labelClass="strong"
                class="rounded pure padding-sm-h"
                suffix-icon="caret-down"/>

    <div class="dropdownMenu-container">

        <DropdownMenu class="childMenu" toggle="#displayByFilter"
                  id="parentMenu"
                  :items="FilterTyles"
                  :checkedKey="filerType"
                  @click="onFilterTypeChanged"
                  :hideOnClickMenu="false"
                  >
        </DropdownMenu>
        <DropdownMenu class="childMenu" toggle="#parentMenu"
                  :items="filerItems"
                  :checkedKey="filerValue"
                  keyName="value"
                  @click="onFilterValueChanged"
                  :hideOnClickMenu="true"
                  :replace-fields="replaceFields"
                  triggerEvent="click"
                  >
        </DropdownMenu>

    </div>
      </ButtonList>
    </template>

    <template #toolbar-buttons>
      <Button class="rounded pure" :hint="t('create_workspace')" @click="showModal=!showModal" icon="folder-add" />
      <Button class="rounded pure" :hint="t('batch_select')" icon="select-all-on" @click="_handleBatchSelectBtnClick" :active="workDirRef?.isCheckable" />
      <Button @click="_handleToggleAllBtnClick" class="rounded pure" :hint="workDirRef?.isAllCollapsed ? t('collapse') : t('expand_all')" :icon="workDirRef?.isAllCollapsed ? 'add-square-multiple' : 'subtract-square-multiple'" iconSize="1.4em" />
<!--      <Button class="rounded pure" :hint="t('more_actions')" icon="more-vert" />-->
    </template>

    <WorkDir ref="workDirRef" />
    <FormWorkspace
      :show="showModal"
      @submit="createWorkSpace"
      @cancel="modalClose"
      ref="formWorkspace"
     />
  </Panel>
</template>

<script setup lang="ts">
import {computed, defineComponent, onMounted, onUnmounted, ref, watch} from "vue";

import Panel from '@/components/Panel.vue';
import Button from '@/components/Button.vue';
import ButtonList from '@/components/ButtonList.vue';
import WorkDir from '@/views/script/WorkDir.vue';
import DropdownMenu from '@/components/DropdownMenu.vue';
import {getScriptDisplayBy, setScriptDisplayBy, setScriptFilters, getScriptFilters} from "@/utils/cache";
import {useStore} from "vuex";
import {ZentaoData} from "@/store/zentao";
import debounce from "lodash.debounce";
import {useI18n} from "vue-i18n";
import {ScriptData} from "@/views/script/store";
import {
  genWorkspaceToScriptsMap,
  listFilterItems,
  getCaseIdsFromReport,
  getSyncFromInfoFromMenu, getNodeMap, getFileNodesUnderParent, updateNameReq
} from "@/views/script/service";
import FormWorkspace from "@/views/workspace/FormWorkspace.vue";
import notification from "@/utils/notification";

const { t, locale, te } = useI18n();

const store = useStore<{ Zentao: ZentaoData, Script: ScriptData }>();
const currSite = computed<any>(() => store.state.Zentao.currSite);
const currProduct = computed<any>(() => store.state.Zentao.currProduct);

let displayBy = ref('workspace')
let isExpand = ref(false);
const filerType = ref('')
const filerValue = ref('')

const displayTypes = ref([
  {key: 'workspace', title: t('by_workspace')},
  {key: 'module', title: t('by_module')},
])

const setDisplayTypes = () => {
  displayTypes.value = [
    {key: 'workspace', title: t('by_workspace')},
    {key: 'module', title: t('by_module')},
  ]
}

const FilterTyles = ref([
  {key: 'suite', title: t('by_suite')},
  {key: 'task', title: t('by_task')},
])

const setFilterTypes = () => {
  FilterTyles.value = [
    {key: 'suite', title: t('by_suite')},
  {key: 'task', title: t('by_task')},
  ]
}
watch(
  locale,
  () => {
    setFilterTypes();
    setDisplayTypes();
  },
  { deep: true }
);

const loadDisplayBy = async () => {
  displayBy.value = await getScriptDisplayBy(currSite.value.id, currProduct.value.id)
}

const initData = debounce(async () => {
  console.log('init')
  if (!currSite.value.id) return

  await loadDisplayBy()
  await loadFilterItems()
  await loadScripts()
}, 50)

const replaceFields = {
      key: 'value',
      title: 'label',
    };

let filerItems = ref([] as any)
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

watch(currProduct, () => {
  console.log('watch currProduct', currProduct.value.id)
  initData()
}, {deep: true})

// only do it when switch from another pages, otherwise will called by watching currProduct method.
if (filerValue.value.length === 0) initData()

const onDisplayByChanged = (item) => {
  console.log('onDisplayBy')
  displayBy.value = item.key

  setScriptDisplayBy(displayBy.value, currSite.value.id, currProduct.value.id)

  loadScripts()
}

const onFilterTypeChanged = async (item) => {
  console.log('onFilterTypeChanged')
  filerType.value = item.key
  await setScriptFilters(displayBy.value, currSite.value.id, currProduct.value.id, filerType.value, filerValue.value)
  loadFilterItems();
}

const onFilterValueChanged = async (item) => {
  console.log('onFilterValueChanged', displayBy.value, filerType.value, item.key)
  filerValue.value = item.key

  await setScriptFilters(displayBy.value, currSite.value.id, currProduct.value.id, filerType.value, filerValue.value)
  await loadScripts()
}

const loadScripts = async () => {
  console.log(`filerType: ${filerType.value}, filerValue: ${filerValue.value}`)

  const params = {displayBy: displayBy.value, filerType: filerType.value, filerValue: filerValue.value} as any
  store.dispatch('Script/listScript', params)
}

const workDirRef = ref<{toggleCheckable: () => void, isCheckable: boolean, toggleAllCollapsed: () => void, isAllCollapsed: boolean}>();

const showModal = ref(false)
const modalClose = () => {
  showModal.value = false;
}
const formWorkspace = ref({} as any)
const createWorkSpace = (formData) => {
    store.dispatch('Workspace/save', formData).then((response) => {
        if (response) {
            formWorkspace.value.clearFormData()
            notification.success({message: t('save_success')});
            showModal.value = false;
            workDirRef.value.loadScripts()
        }
    })
};

function _handleBatchSelectBtnClick() {
    if (workDirRef.value) {
        workDirRef.value.toggleCheckable();
    }
}

function _handleToggleAllBtnClick() {
    if (workDirRef.value) {
        workDirRef.value.toggleAllCollapsed();
    }
}

</script>

<style lang="less" >
.workdir-panel {
  height: 100%;

  .panel-heading {
    .workdir-panel-nav {
      margin-left: -6px;
    }
  }

  .panel-body {
    height: calc(100% - 30px);
  }
}
.dropdownMenu-container {
    display: flex;
    position: fixed;
    z-index: 100;
    opacity: 1;
    top: 69px;
    left: 71px;
    min-width: none !important;
}
.childMenu{
    position: unset!important;
    min-width: 0!important;
}
</style>
