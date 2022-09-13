<template>
  <div class="tabs-container relative column single">
    <TabsNav
        class="flex-none"
        :items="items"
        :activeID="activeID"
        :toolbarItems="toolbarItems"
        @click="_handleNavClick"
        @close="_handleNavClose"
        @clickToolbar="onToolbarClick"
    />
    <template v-for="tab in tabsList" :key="tab.id">
      <KeepAlive>
        <TabPage v-if="tab.id === activeID" class="flex-auto" :tab="tab" ref="tabsRef" />
      </KeepAlive>
    </template>
    <BugsModal
      v-if="showBugsModal"
      @cancel="showBugsModal=!showBugsModal"
      :caseIds="caseIds"
    />
  </div>
</template>

<script setup lang="ts">
import {computed, onMounted, ref, watch} from 'vue';
import {useStore} from 'vuex';
import {PageTab, TabsData} from "@/store/tabs";
import TabsNav, {TabNavItem} from './TabsNav.vue';
import TabPage from './TabPage.vue';
import {useI18n} from "vue-i18n";
import bus from "@/utils/eventBus";
import settings from "@/config/settings";
import {ScriptData} from "@/views/script/store";
import {loadBugs} from "@/services/bug";
import { ZentaoData } from "@/store/zentao";
import { StateType as GlobalData } from "@/store/global";

import BugsModal from '@/views/result/BugsModal.vue';

const {t} = useI18n();

const store = useStore<{ global: GlobalData, tabs: TabsData, Script: ScriptData, Zentao: ZentaoData }>();
const global = computed<any>(() => store.state.global.tabIdToWorkspaceIdMap);
const saveBtnAbled = computed<any>(() => store.state.Script.currentCodeChanged);
const currWorkspace = computed<any>(() => store.state.Script.currWorkspace);
const bugMap = computed<any>(() => store.state.Zentao.bugMap);
const showBugsModal = ref(false)
const caseIds = ref([] as number[])

const items = computed<TabNavItem[]>(() => {
  return store.getters['tabs/list'];
});

let toolbarItemArr = [
  {
    key: 'run',
    hint: 'Run',
    icon: 'play'
  },
  {
    key: 'save',
    hint: 'Save',
    icon: 'save',
    disabled: !saveBtnAbled.value
  }
];
const toolbarItems = ref([] as any[]);
// const tabPageRef = ref<InstanceType<typeof TabPage> | null>(null)
const tabsRef = ref<InstanceType<typeof TabPage>[] | null>(null)

const tabsList = computed(() => {
  return store.getters['tabs/list'];
});

const activeID = computed((): string => {
  return store.state.tabs.activeID;
});

const treeDataMap = computed<any>(() => store.state.Script.treeDataMap);

watch(activeID, () => {
  console.log('watch activeID', activeID.value)
  if (activeID.value.indexOf('script-') > -1) {
    toolbarItems.value = toolbarItemArr
    const path = activeID.value.replace('script-', '');
    store.dispatch('Result/getStatistic', treeDataMap.value[path])
  } else if (activeID.value.indexOf('code-') > -1) {
    toolbarItems.value = [toolbarItemArr[1]]
  } else {
    toolbarItems.value = []
  }
  store.dispatch("Zentao/fetchBugs", {});
  showBugBtn()
}, {deep: true})

watch(bugMap, () => {
  console.log('watch bugMap', bugMap.value)
  showBugBtn()
}, {deep: true})

watch(saveBtnAbled, () => {
  toolbarItemArr = [
    {
    key: 'run',
    hint: 'Run',
    icon: 'play'
    },
    {
    key: 'save',
    hint: 'Save',
    icon: 'save',
    disabled: !saveBtnAbled.value
    }
  ];
  if (activeID.value.indexOf('script-') > -1) {
    toolbarItems.value = toolbarItemArr
  } else if (activeID.value.indexOf('code-') > -1) {
    toolbarItems.value = [toolbarItemArr[1]]
  } else {
    toolbarItems.value = []
  }
}, {deep: true})
const showBugBtn = () => {
  if (activeID.value.indexOf('script-') === 0) {
    const path = activeID.value.replace('script-', '')
    const relateBugs = bugMap.value[treeDataMap.value[path].caseId] || []
    const hasBugBtn = toolbarItems.value.filter(item => item.key === 'bugs').length > 0;
    if(relateBugs.length > 0 && !hasBugBtn) {
        toolbarItems.value.push({
            key: 'bugs',
            hint: 'Bugs',
            icon: 'bug'
        })
    }else if(relateBugs.length === 0 && hasBugBtn) {
        toolbarItems.value = toolbarItems.value.filter(item => item.key !== 'bugs')
    }
    caseIds.value = [treeDataMap.value[path].caseId]
  }
}

const testTabIDRef = ref(0);

const onToolbarClick = (e) => {
  console.log('onToolbarClick', e.key, activeID.value)

  switch (e.key) {
    case 'run': {
      let path = activeID.value
      if (path.indexOf('script-') === 0) {
        path = path.replace('script-', '')
      }

      const workspaceId = global.value[activeID.value]
      bus.emit(settings.eventExec,
         {execType: 'ztf', scripts: [{ path: path, workspaceId: workspaceId }]});
      break;
    }
    case 'save': {
      if (tabsRef.value) {
        bus.emit(settings.eventScriptSave, {path: activeID.value})
      }
      break;
    }
    case 'bugs': {
      showBugsModal.value = true;
      break;
    }
  }
}

function _handleNavClick(item) {
  console.log('_handleNavClick', item);
  store.dispatch('tabs/open', item);
  store.dispatch('Script/updateCurrentCodeChanged', item.changed)
}

function _handleNavClose(item) {
  console.log('_handleNavClose', item);
  store.dispatch('tabs/close', item);
}

if (process.env?.NODE_ENV === 'development') {
  onMounted(() => {
    Object.assign(window, {
      $openPage: (tab: string | PageTab): void => {
        store.dispatch('tabs/open', typeof tab === 'string' ? {id: tab, type: tab, title: tab} : tab);
      }
    });
  });
}
</script>

<style lang="less">
#content {
  .category-output { color: #95a5a6 }
  .category-run { color: #1890ff }
  .category-result { color: #68BB8D }
  .category-error { color: #FC2C25 }

  &.category-run {
    .category-output { display: none !important; }
  }
  &.category-result {
    .category-output, .category-run { display: none !important; }
  }
  &.category-error {
    .category-output, .category-run, .category-result { display: none !important; }
  }

  &.status-pass {
    .category-error { display: none !important; }
  }
  &.status-fail {
    .category-output, .category-run, .category-result { display: none !important; }
  }

  .item {
    display: flex;

    .no-span {
      width: 60px;
      text-align: right;
    }

    .msg-span {
      flex: 1;
    }
  }
}
</style>

<style scoped>
.tabs-container {
  position: relative;
}
</style>
