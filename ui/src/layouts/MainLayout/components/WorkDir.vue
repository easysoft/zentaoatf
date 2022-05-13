<template>
  <div class="workdir padding muted">
<!--    <Row :gutter="10">
      <Col width="30px">1</Col>
      <Col :span="3" :offset="4">111</Col>
      <Col :flex="6" :offset="4">wwwwwww</Col>
    </Row>-->

    <Form labelCol="50px" wrapperCol="60">
      <FormItem name="name" label="标题" :info="validateInfos.name">
        <input v-model="modelRef.name" />

        &nbsp;&nbsp;&nbsp;
        <select>
          <option value="1" > Miner </option>
        </select>
        &nbsp;&nbsp;&nbsp;
        <input name="sex" checked type="radio" value="female"/><label>女</label>
        &nbsp;&nbsp;&nbsp;
        <input name='subject' type="checkbox" checked="checked" value="English"/><label>英语</label>

      </FormItem>

      <FormItem size="small">
        <button @click="submit" type="button">提交</button>
      </FormItem>
    </Form>

    <ScriptTreePage></ScriptTreePage>
  </div>
</template>

<script setup lang="ts">
import ScriptTreePage from "../../../views/script/component/tree.vue";
import {useI18n} from "vue-i18n";
import {useStore} from "vuex";
import {ZentaoData} from "@/store/zentao";
import {computed, onMounted, provide, ref} from "vue";
import {ScriptData} from "@/views/script/store";
import {resizeWidth} from "@/utils/dom";

import Row from "./Row.vue";
import Col from "./Col.vue";
import Form from "./Form.vue";
import FormItem from "./FormItem.vue";
import {useForm} from "@/utils/form";

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

const modelRef = ref({})
const rulesRef = ref({
  name: [
    {required: true, msg: 'Please input name.'},
  ],
})

const { validate, reset, validateInfos } = useForm(modelRef, rulesRef);

const submit = () => {
  console.log('submit')

  validate()
  console.log(validateInfos)
}

</script>

<style lang="less" scoped>
.workdir {

}
</style>
