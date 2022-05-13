<template>
  <div class="tabs-container relative column single">
    <TabsNav
      class="flex-none"
      :items="items"
      :activeID="activeID"
      :toolbarItems="toolbarItems"
      @click="_handleNavClick"
      @close="_handleNavClose"
    />
    <Button style="position: absolute; top: 50px; right: 0;" class="red" icon="add" @click="_addTestTab">Add test tab</Button>
    <template v-for="tab in tabsList" :key="tab.id">
      <KeepAlive>
        <TabPage v-if="tab.id === activeID" class="flex-auto" :tab="tab" />
      </KeepAlive>
    </template>
  </div>
</template>

<script setup lang="ts">
import {computed, ref} from 'vue';
import {useStore} from 'vuex';
import {TabsData, PageTab} from "@/store/tabs";
import TabsNav, {TabNavItem} from './TabsNav.vue';
import Button from './Button.vue';
import TabPage from './TabPage.vue';
import {useI18n} from "vue-i18n";
const { t } = useI18n();

const store = useStore<{tabs: TabsData}>();

const items = computed<TabNavItem[]>(() => {
    return store.getters['tabs/list'];
});

const toolbarItems = computed(() => {
    return [{
        hint: 'Run',
        icon: 'play'
    }, {
        hint: 'Save',
        icon: 'save'
    }, {
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
</script>

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
