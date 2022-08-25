<template>
  <ZModal
    id="syncToZentaoModal"
    :showModal="showModalRef"
    @onCancel="close"
    @onOk="submit"
    :title="t('submit_result_to_zentao')"
    :contentStyle="{width: '500px'}"
  >
    <Form>
      <FormItem name="taskId" :label="t('pls_select_task')" labelWidth="120px">
        <div class="select">
          <select name="taskId" v-model="modelRef.taskId">
            <option value="0"></option>
            <option v-for="item in tasks" :key="item.id" :value="item.id">
              {{ item.name }}
            </option>
          </select>
        </div>
      </FormItem>

      <FormItem v-if="!modelRef.taskId" name="name" :label="t('or_input_task_name')" labelWidth="120px">
        <input type="text" v-model="modelRef.name" />
      </FormItem>
    </Form>
  </ZModal>
</template>

<script setup lang="ts">
import { useI18n } from "vue-i18n";
import { useStore } from "vuex";
import { ZentaoData } from "@/store/zentao";
import { unitTestTypesDef, ztfTestTypesDef } from "@/utils/const";
import {computed, defineExpose, withDefaults, ref, defineProps } from "vue";
import { useForm } from "@/utils/form";
import Form from "@/components/Form.vue";
import FormItem from "@/components/FormItem.vue";
import {queryTask} from "@/services/zentao";
import {submitResultToZentao} from "@/views/result/service";
import notification from "@/utils/notification";

export interface FormWorkspaceProps {
  data?: any;
  show?: boolean;
  finish: Function;
}
const { t } = useI18n();
const props = withDefaults(defineProps<FormWorkspaceProps>(), {
  data: {},
  show: false,
});

const showModalRef = computed(() => {
  return props.show;
});
const testTypes = ref([...ztfTestTypesDef, ...unitTestTypesDef]);
const store = useStore<{ Zentao: ZentaoData }>();
const currProduct = computed<any>(() => store.state.Zentao.currProduct);

const modelRef = ref({taskId: 0} as any);
const rulesRef = ref({});

let tasks = ref([])
const listTask = () => {
  queryTask(currProduct.value.id).then((jsn) => {
    tasks.value = jsn.data
  })
}
listTask()

const { validate, reset, validateInfos } = useForm(modelRef, rulesRef);

const submit = () => {
  console.log('submitResultForm', modelRef.value)

  if (!validate()) {
    return
  }

  const data = Object.assign({
    workspaceId: props.data.workspaceId,
    seq: props.data.seq,
  }, modelRef.value) as any
  data.taskId = parseInt(data.taskId)

  console.log('data', data)

  submitResultToZentao(data).then((json) => {
    console.log('json', json)
    if (json.code === 0) {
      notification.success({
        message: t('submit_success'),
      });
      close()
    } else {
      notification.error({
        message: t('submit_failed'),
        description: json.msg,
      });
    }
  })
}
const close = () => {
  props.finish()
};

const clearFormData = () => {
  modelRef.value = {};
};

defineExpose({
  clearFormData,
});
</script>

<style lang="less" scoped>

</style>
