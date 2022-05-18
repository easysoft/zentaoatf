<template>
  <div class="tab-page-exec-unit">
    <Form labelCol="180px" wrapperCol="60">
      <FormItem name="cmd" :label="t('test_cmd')" :info="validateInfos.cmd">
        <input v-model="modelRef.cmd" @keydown="keydown"/>
      </FormItem>

      <FormItem v-if="currProduct.id" name="submitResult" :label="t('submit_result_to_zentao')">
        <input v-model="modelRef.submitResult" type="checkbox">
      </FormItem>

      <FormItem>
        <a-button :disabled="isRunning === 'true' || !modelRef.cmd" @click="start" type="primary"
                  class="t-btn-gap">
          {{ t('exec') }}
        </a-button>
        <a-button v-if="isRunning === 'true'" @click="stop" class="t-btn-gap">
          {{ t('stop') }}
        </a-button>
      </FormItem>

      <FormItem>
        <span class="t-tips">{{ t('cmd_nav') }}</span>
      </FormItem>
    </Form>
  </div>
</template>

<script setup lang="ts">
import {withDefaults, defineProps, computed, ref} from "vue";
import { PageTab } from "@/store/tabs";
import {useI18n} from "vue-i18n";
import {useStore} from "vuex";
import {ZentaoData} from "@/store/zentao";
import {useForm} from "@/utils/form";
const { t } = useI18n();

import Form from "./Form.vue";
import FormItem from "./FormItem.vue";

const props = withDefaults(defineProps<{
  tab: PageTab
}>(), {})

const zentaoStore = useStore<{ Zentao: ZentaoData }>();
const currProduct = computed<any>(() => zentaoStore.state.Zentao.currProduct);

const workspaceId = computed<any>(() => props.tab.data.workspaceId)
const workspaceType = computed<any>(() => props.tab.data.workspaceType)

const modelRef = ref({})
const isRunning = ref(false)

const rulesRef = ref({
  cmd: [
    {required: true, msg: 'Please input test command.'},
  ],
})
const {validate, reset, validateInfos} = useForm(modelRef, rulesRef);

const start = () => {
  console.log('start')

  console.log(modelRef.value)
}
const stop = () => {
  console.log('stop')
}

console.log(workspaceId, workspaceType)

</script>

<style lang="less" scoped>
.tab-page-exec-unit {
  padding: 16px;
}
</style>
