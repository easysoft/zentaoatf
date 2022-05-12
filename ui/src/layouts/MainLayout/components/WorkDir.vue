<template>
  <div class="workdir padding muted">
    <Row type="flex">
      <Column :width="30">1</Column>
      <Column :span="12" :offset="4">wwwwwww</Column>
    </Row>

    <ScriptTreePage></ScriptTreePage>
  </div>
</template>

<script setup lang="ts">
import ScriptTreePage from "../../../views/script/component/tree.vue";
import {useI18n} from "vue-i18n";
import {useStore} from "vuex";
import {ZentaoData} from "@/store/zentao";
import {computed, onMounted} from "vue";
import {ScriptData} from "@/views/script/store";
import {resizeWidth} from "@/utils/dom";

import Row from "./Row.vue";
import Column from "./Column.vue";

const { t } = useI18n();

const zentaoStore = useStore<{ Zentao: ZentaoData }>();
const currSite = computed<any>(() => zentaoStore.state.Zentao.currSite);
const currProduct = computed<any>(() => zentaoStore.state.Zentao.currProduct);

const scriptStore = useStore<{ Script: ScriptData }>();
const currWorkspace = computed<any>(() => scriptStore.state.Script.currWorkspace);

onMounted(() => {
  console.log('onMounted')
  setTimeout(() => {
    resizeWidth('main', 'left', 'splitter-h', 'right', 380, 800)
  }, 600)
})

</script>

<style lang="less" scoped>
.workdir {

}
</style>
