<template>
  <ZModal
      :showModal="showModalRef"
      @onCancel="close"
      @onOk="submit"
      :title="t('submit_bug_to_zentao')"
      :contentStyle="{width: '500px'}"
  >
    <Form>
      <FormItem name="title" :label="t('title')" :info="validateInfos.title">
        <input type="text" v-model="modelRef.title" />
      </FormItem>

      <FormItem name="module" :label="t('module')">
        <div class="select">
          <select name="module" v-model="modelRef.module">
            <option value=""></option>
            <option v-for="item in modules" :key="item.code" :value="item.code+''">
              {{ item.name }}
            </option>
          </select>
        </div>
      </FormItem>

      <FormItem name="type" :label="t('category')">
        <div class="select">
          <select name="type" v-model="modelRef.type">
            <option value=""></option>
            <option v-for="item in types" :key="item.code" :value="item.code+''">
              {{ item.name }}
            </option>
          </select>
        </div>
      </FormItem>

      <FormItem name="openedBuild" :label="t('version')">
        <div class="select">
          <select name="openedBuild" v-model="modelRef.openedBuild">
            <option value=""></option>
            <option v-for="item in builds" :key="item.code" :value="item.code+''">
              {{ item.name }}
            </option>
          </select>
        </div>
      </FormItem>

      <FormItem name="severity" :label="t('severity')">
        <div class="select">
          <select name="severity" v-model="modelRef.severity">
            <option v-for="item in severities" :key="item.code" :value="item.code+''">
              {{ item.name }}
            </option>
          </select>
        </div>
      </FormItem>

      <FormItem name="priority" :label="t('priority')">
        <div class="select">
          <select name="priority" v-model="modelRef.pri">
            <option v-for="item in priorities" :key="item.code" :value="item.code+''">
              {{ item.name }}
            </option>
          </select>
        </div>
      </FormItem>

      <FormItem name="type" :label="t('step')">
        <textarea v-model="modelRef.steps" rows="6" />
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
import Form from "@/components/Form.vue";
import FormItem from "@/components/FormItem.vue";
import {queryBugFields, queryTask} from "@/services/zentao";
import notification from "@/utils/notification";
import {prepareBugData, submitBugToZentao} from "@/services/bug";

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
  }, modelRef.value) as any
  data.module = parseInt(data.module)
  data.severity = parseInt(data.severity)
  data.pri = parseInt(data.pri)
  if (!Array.isArray(data.openedBuild)) data.openedBuild = [data.openedBuild]

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
