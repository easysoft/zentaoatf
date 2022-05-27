<template>
  <ZModal
      :showModal="showModalRef"
      @onCancel="close"
      @onOk="submit"
      :title="t('submit_result_to_zentao')"
  >
    <Form labelCol="6" wrapperCol="16">
      <FormItem name="title" :label="t('title')" :info="validateInfos.title">
        <input v-model="modelRef.title" class="form-control" />
      </FormItem>

      <FormItem name="module" :label="t('module')">
        <select name="module" v-model="modelRef.module" class="form-control">
          <option value=""></option>
          <option v-for="item in modules" :key="item.code" :value="item.code+''">
            {{ item.name }}
          </option>
        </select>
      </FormItem>

      <FormItem name="type" :label="t('category')">
        <select name="type" v-model="modelRef.type" class="form-control">
          <option value=""></option>
          <option v-for="item in types" :key="item.code" :value="item.code+''">
            {{ item.name }}
          </option>
        </select>
      </FormItem>

      <FormItem name="openedBuild" :label="t('version')">
        <select name="openedBuild" v-model="modelRef.openedBuild" class="form-control">
          <option value=""></option>
          <option v-for="item in builds" :key="item.code" :value="item.code+''">
            {{ item.name }}
          </option>
        </select>
      </FormItem>

      <FormItem name="severity" :label="t('severity')">
        <select name="severity" v-model="modelRef.severity" class="form-control">
          <option v-for="item in severities" :key="item.code" :value="item.code+''">
            {{ item.name }}
          </option>
        </select>
      </FormItem>

      <FormItem name="priority" :label="t('priority')">
        <select name="priority" v-model="modelRef.pri" class="form-control">
          <option v-for="item in priorities" :key="item.code" :value="item.code+''">
            {{ item.name }}
          </option>
        </select>
      </FormItem>

      <FormItem name="type" :label="t('step')">
        <textarea v-model="modelRef.steps" rows="3" />
      </FormItem>
    </Form>
  </ZModal>
</template>

<script setup lang="ts">
import { useI18n } from "vue-i18n";
import { useStore } from "vuex";
import { ZentaoData } from "@/store/zentao";
import { unitTestTypesDef, ztfTestTypesDef } from "@/utils/const";
import {
  computed,
  defineExpose,
  onMounted,
  withDefaults,
  ref,
  defineProps, reactive,
} from "vue";
import { useForm } from "@/utils/form";
import Form from "./Form.vue";
import FormItem from "./FormItem.vue";
import {queryBugFields, queryTask} from "@/services/zentao";
import {notification} from "ant-design-vue";
import {prepareBugData, submitBugToZentao} from "@/services/bug";

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
const rulesRef = reactive({
  title: [
    { required: true, msg: t('pls_title') },
  ],
});

let modules = ref([])
let types = ref([])
let builds = ref([])
let severities = ref([])
let priorities = ref([])
const getBugData = () => {
  prepareBugData(props.data).then((jsn) => {
    modelRef.value = jsn.data
    modelRef.value.module = ''
    modelRef.value.severity = ''+modelRef.value.severity
    modelRef.value.pri = ''+modelRef.value.pri

    getBugFields()
  })
}
const getBugFields = () => {
  queryBugFields().then((jsn) => {
    modules.value = jsn.data.modules
    types.value = jsn.data.type
    builds.value = jsn.data.build
    severities.value = jsn.data.severity
    priorities.value = jsn.data.pri
  })
}
getBugData()

const { validate, reset, validateInfos } = useForm(modelRef, rulesRef);

const submit = () => {
  console.log('submitBugForm', modelRef.value)
  if (!validate()) {
    return
  }

  const data = Object.assign({
    workspaceId: props.data.workspaceId,
    seq: props.data.seq
  }, modelRef.value)
  data.module = parseInt(data.module)
  data.severity = parseInt(data.severity)
  data.pri = parseInt(data.pri)

  submitBugToZentao(data).then((json) => {
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