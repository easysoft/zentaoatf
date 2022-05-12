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
</style>
