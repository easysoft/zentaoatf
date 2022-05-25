<template>
  <div class="tab-page-exec-unit">
    <Form labelCol="6" wrapperCol="26">
      <FormItem name="cmd" :label="t('test_cmd')" :info="validateInfos.cmd">
        <input v-model="modelRef.cmd" @keydown="keydown"/>
      </FormItem>

      <FormItem v-if="currProduct.id" name="submitResult" :label="t('submit_result_to_zentao')">
        <input v-model="modelRef.submitResult" type="checkbox">
      </FormItem>

      <FormItem>
        <Button :disabled="isRunning === 'true' || !modelRef.cmd" @click="start" class="rounded">
          {{ t('exec') }}
        </Button>
        <Button v-if="isRunning === 'true'" @click="stop" class="rounded pure">
          {{ t('stop') }}
        </Button>
      </FormItem>

      <FormItem>
        <span class="t-tips">{{ t('cmd_nav') }}</span>
      </FormItem>
    </Form>
  </div>
</template>

<script setup lang="ts">
import {withDefaults, defineProps, computed, ref, watch} from "vue";
import { PageTab } from "@/store/tabs";
import {useI18n} from "vue-i18n";
import {useStore} from "vuex";
import {ZentaoData} from "@/store/zentao";
import Button from './Button.vue';
import {useForm} from "@/utils/form";
const { t } = useI18n();

import Form from "./Form.vue";
import FormItem from "./FormItem.vue";
import {get} from "@/views/workspace/service";
import {ScriptData} from "@/views/script/store";
import {getCmdHistories, setCmdHistories} from "@/utils/cache";
import {WorkspaceData} from "@/store/workspace";
import bus from "@/utils/eventBus";
import settings from "@/config/settings";

const props = withDefaults(defineProps<{
  tab: PageTab
}>(), {})

const zentaoStore = useStore<{ Zentao: ZentaoData }>();
const currProduct = computed<any>(() => zentaoStore.state.Zentao.currProduct);

const scriptStore = useStore<{ Script: ScriptData }>();
const currWorkspace = computed<any>(() => scriptStore.state.Script.currWorkspace);

const workspaceId = computed<any>(() => props.tab.data.workspaceId)
const workspaceType = computed<any>(() => props.tab.data.workspaceType)

const modelRef = ref({} as any)
const isRunning = ref(false)

const histories = ref([] as any[])
const historyIndex = ref(0)

const rulesRef = ref({
  cmd: [
    {required: true, msg: 'Please input test command.'},
  ],
})
const {validate, reset, validateInfos} = useForm(modelRef, rulesRef);

const loadWorkspaceCmd = async () => {
  get(currWorkspace.value.id).then((json) => {
    modelRef.value = Object.assign({cmd: json.data.cmd}, currWorkspace.value)
  })
}

const loadCmdHistories = async () => {
  histories.value = await getCmdHistories(currWorkspace.value.id)
  historyIndex.value = histories.value? histories.value.length : 0
}

if (currWorkspace.value.id > 0) {
  loadCmdHistories()
  loadWorkspaceCmd()
}

watch(currWorkspace, () => {
  console.log('watch currWorkspace', currWorkspace)

  loadCmdHistories()
  loadWorkspaceCmd()
}, {deep: true})

if (workspaceId.value > 0 && workspaceId.value !== currWorkspace.value.id) {
  scriptStore.dispatch('Script/changeWorkspace',
      {id: workspaceId.value, type: workspaceType.value})
}

const start = () => {
  console.log('start exec unit test', modelRef.value)

  addHistory()

  const data = Object.assign({execType: 'unit'}, modelRef.value)
  bus.emit(settings.eventExec, data);
}

const stop = () => {
  console.log('stop')
}

const addHistory = () => {
  if (modelRef.value.cmd !== histories.value[histories.value.length - 1])
    histories.value.push(modelRef.value.cmd)
  if (histories.value.length > 10)
    histories.value = histories.value.slice(histories.value.length - 10)

  setCmdHistories(currWorkspace.value.id, histories.value)
  historyIndex.value = histories.value.length
}
const keydown = (e) => {
  console.log('keydown', e.code)

  if (e.code === 'ArrowUp') {
    if (historyIndex.value > 0) historyIndex.value--
    modelRef.value.cmd = histories.value[historyIndex.value]
  } else if (e.code === 'ArrowDown') {
    if (historyIndex.value < histories.value.length - 1) historyIndex.value++
    modelRef.value.cmd = histories.value[historyIndex.value]
  }
}

console.log(workspaceId, workspaceType)

</script>

<style lang="less" scoped>
.tab-page-exec-unit {
  padding: 16px;
}
</style>
