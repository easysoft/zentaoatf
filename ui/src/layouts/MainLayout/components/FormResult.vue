<template>
  <ZModal
    :showModal="showModalRef"
    @onCancel="close"
    @onOk="submit"
    :title="t('submit_result_to_zentao')"
  >
    <Form labelCol="6" wrapperCol="16">
      <FormItem name="taskId" :label="t('task')">
        <select name="taskId" v-model="modelRef.taskId" class="form-control">
          <option value="0"></option>
          <option v-for="item in tasks" :key="item.id" :value="item.id">
            {{ item.name }}
          </option>
        </select>
      </FormItem>
    </Form>
  </ZModal>
</template>

<script setup lang="ts">
import { useI18n } from "vue-i18n";
import { useStore } from "vuex";
import { ZentaoData } from "@/store/zentao";
import { ScriptData } from "@/views/script/store";
import { unitTestTypesDef, ztfTestTypesDef } from "@/utils/const";
import {
  computed,
  defineExpose,
  onMounted,
  withDefaults,
  ref,
  defineProps,
  defineEmits, reactive, PropType,
} from "vue";
import { useForm } from "@/utils/form";
import Form from "./Form.vue";
import FormItem from "./FormItem.vue";
import {queryTask} from "@/services/zentao";
import {submitResultToZentao} from "@/views/result/service";
import {notification} from "ant-design-vue";

export interface FormWorkspaceProps {
  data?: any;
  show?: boolean;
  finish?: Function;
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

const modelRef = ref({taskId: 0});
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
  }, modelRef.value)
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
.form-control {
  width: 100%;
  color: #495057;
  background-color: #fff;
  border: 1px solid #ced4da;
  border-radius: 0.25rem;
  transition: border-color 0.15s ease-in-out, box-shadow 0.15s ease-in-out;
}
.z-form-item-label {
  font-weight: 400;
  color: #212529;
  text-align: left;
  box-sizing: border-box;
  display: inline-block;
  position: relative;
  width: 100%;
  padding-right: 15px;
  padding-left: 15px;
  padding-top: calc(0.375rem + 1px);
  padding-bottom: calc(0.375rem + 1px);
  margin-bottom: 0;
  line-height: 1.5;
}
.z-form-item {
  display: flex;
  align-items: center;
}
.form-control:focus {
  color: #495057;
  background-color: #fff;
  border-color: #80bdff;
  outline: 0;
  box-shadow: 0 0 0 0.2rem rgb(0 123 255 / 25%);
}
</style>