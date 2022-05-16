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

    <Button style="position: absolute; top: 50px; right: 0;" class="red" icon="add" @click="_addTestTab">Add test tab
    </Button>

    <template v-for="tab in tabsList" :key="tab.id">
      <KeepAlive>
        <TabPage v-if="tab.id === activeID" class="flex-auto" :tab="tab"/>
      </KeepAlive>
    </template>
  </div>
</template>

<script setup lang="ts">
import {computed, onMounted, ref} from 'vue';
import {useStore} from 'vuex';
import {PageTab, TabsData} from "@/store/tabs";
import TabsNav, {TabNavItem} from './TabsNav.vue';
import Button from './Button.vue';
import TabPage from './TabPage.vue';
import {useI18n} from "vue-i18n";
import bus from "@/utils/eventBus";
import settings from "@/config/settings";
import {ScriptData} from "@/views/script/store";

const {t} = useI18n();

const store = useStore<{ tabs: TabsData }>();

const scriptStore = useStore<{ Script: ScriptData }>();
const currWorkspace = computed<any>(() => scriptStore.state.Script.currWorkspace);

const items = computed<TabNavItem[]>(() => {
  return store.getters['tabs/list'];
});

const toolbarItems = computed(() => {
  return [{
    key: 'run',
    hint: 'Run',
    icon: 'play'
  }, {
    key: 'save',
    hint: 'Save',
    icon: 'save'
  }, {
    key: 'run',
    hint: 'More actions',
    icon: 'more-vert'
  }];
});

const tabsList = computed(() => {
  return store.getters['tabs/list'];
});

const activeID = computed(() => {
  return store.state.tabs.activeID;
});

const testTabIDRef = ref(0);

const onToolbarClick = (e) => {
  console.log('onToolbarClick', e.key, activeID.value)

  if (e.key === 'run') {
    bus.emit(settings.eventExec, {execType: 'ztf',
      scripts: [
          { path: activeID.value, workspaceId: currWorkspace.value.id }]});
  }
}

function _handleNavClick(item) {
  console.log('_handleNavClick', item);
  store.dispatch('tabs/open', item);
}

function _handleNavClose(item) {
  console.log('_handleNavClose', item);
  store.dispatch('tabs/close', item);
}

function _addTestTab() {
  testTabIDRef.value++;
  store.dispatch('tabs/open', {
    id: `testTab-${testTabIDRef.value}`,
    title: `TestTab ${testTabIDRef.value}`,
    changed: Math.random() > 0.5,
    type: ['script', 'sites', 'settings', 'result', ''][Math.floor(Math.random() * 5)],
    data: Math.random()
  });
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


<!--
<template>
  <div class="script-tabs padding muted">
    <ZtfScriptPage v-if="currWorkspace?.type === 'ztf'"></ZtfScriptPage>
    <UnitScriptPage v-if="currWorkspace?.type !== 'ztf'"></UnitScriptPage>
  </div>
</template>

<script setup lang="ts">
import {computed, defineComponent, onMounted, onUnmounted, ref, watch} from "vue";

import ZtfScriptPage from "../../../views/script/component/ztf.vue"
import UnitScriptPage from "../../../views/script/component/unit.vue"
import {useI18n} from "vue-i18n";
import {useStore} from "vuex";
import {ZentaoData} from "@/store/zentao";
import {ScriptData} from "@/views/script/store";

const { t } = useI18n();

const zentaoStore = useStore<{ Zentao: ZentaoData }>();
const currSite = computed<any>(() => zentaoStore.state.Zentao.currSite);
const currProduct = computed<any>(() => zentaoStore.state.Zentao.currProduct);

const scriptStore = useStore<{ Script: ScriptData }>();
const currWorkspace = computed<any>(() => scriptStore.state.Script.currWorkspace);

</script>

<style lang="less" >
.script-tabs {
  height: 100%;

  .monaco-editor {
    padding: 10px 0;
  }
}
</style> -->
