<template>
  <div class="workdir padding muted">
    <Row :gutter="10">
      <Col width="30px">1</Col>
      <Col :span="3" :offset="4">111</Col>
      <Col :flex="6" :offset="4">wwwwwww</Col>
    </Row>

<!--    <FormItem label="标题" labelCol="50px" wrapperCol="20">
      <a-input />
    </FormItem>-->

    <ScriptTreePage></ScriptTreePage>
  </div>
</template>

<script setup lang="ts">
import ScriptTreePage from "../../../views/script/component/tree.vue";
import {useI18n} from "vue-i18n";
import {useStore} from "vuex";
import {ZentaoData} from "@/store/zentao";
import {computed, onMounted, provide} from "vue";
import {ScriptData} from "@/views/script/store";
import {resizeWidth} from "@/utils/dom";

import Row from "./Row.vue";
import Col from "./Col.vue";
import FormItem from "./FormItem.vue";

const { t } = useI18n();

const zentaoStore = useStore<{ Zentao: ZentaoData }>();
const currSite = computed<any>(() => zentaoStore.state.Zentao.currSite);
const currProduct = computed<any>(() => zentaoStore.state.Zentao.currProduct);

const scriptStore = useStore<{ Script: ScriptData }>();
const currWorkspace = computed<any>(() => scriptStore.state.Script.currWorkspace);

provide('errors', {});

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
