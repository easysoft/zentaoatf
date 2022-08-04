<template>
  <ZModal
    :showModal="showModalRef"
    @onCancel="cancel"
    @onOk="submit"
    :title="t('sync-from-zentao')"
    :contentStyle="{ width: '600px' }"
  >
    <Form>
      <FormItem labelWidth="140px" :label="t('module')">
        <select v-model="modelRef.moduleId">
          <option key="" value="">&nbsp;</option>
          <option v-for="item in modules" :key="item.id" :value="item.id">
            <span v-html="item.name"></span>
          </option>
        </select>
      </FormItem>

      <FormItem labelWidth="140px" :label="t('suite')">
        <select v-model="modelRef.suiteId">
          <option key="" value="">&nbsp;</option>
          <option v-for="item in suites" :key="item.id" :value="item.id">
            {{ item.name }}
          </option>
        </select>
      </FormItem>

      <FormItem labelWidth="140px" :label="t('task')">
        <select v-model="modelRef.taskId">
          <option key="" value="">&nbsp;</option>
          <option v-for="item in tasks" :key="item.id" :value="item.id">
            {{ item.name }}
          </option>
        </select>
      </FormItem>

      <FormItem labelWidth="140px" :label="t('lang')">
        <select v-model="modelRef.lang">
          <option v-for="item in langs" :key="item.code" :value="item.code">
            {{ item.name }}
          </option>
        </select>
      </FormItem>

      <FormItem labelWidth="140px" :label="t('save_by_module')">
        <Switch v-model="modelRef.saveByModule" />
      </FormItem>

      <FormItem labelWidth="140px" :label="t('independent_expect')">
        <Switch v-model="modelRef.independentFile" />
      </FormItem>
    </Form>
  </ZModal>
</template>

<script setup lang="ts">
import {
  ref,
  computed,
  defineExpose,
  watch,
  withDefaults,
  defineProps,
  defineEmits,
} from "vue";
import { useI18n } from "vue-i18n";
import { useStore } from "vuex";
import { WorkspaceData } from "../workspace/store";
import { isWindows } from "@/utils/comm";
import { get as getWorkspace } from "@/views/workspace/service";
import Form from "@/components/Form.vue";
import FormItem from "@/components/FormItem.vue";
import { useForm } from "@/utils/form";
import Switch from "@/components/Switch.vue";
import { ZentaoData } from "@/store/zentao";

export interface FormWorkspaceProps {
  show?: boolean;
  workspaceId?: number;
}
const { t } = useI18n();
const isWin = isWindows();
const disabled = ref(false);
const props = withDefaults(defineProps<FormWorkspaceProps>(), {
  show: false,
});
watch(props, () => {
  if(!props.show) disabled.value = false;
  modelRef.value.workspaceId = props.workspaceId;
  selectWorkspace();
});
const modelRef = ref({
  workspaceId: props.workspaceId,
  lang: 'php',
  independentFile: false,
} as any);
const rulesRef = ref({
  workspaceId: [{ required: true, msg: t("pls_select_workspace") }],
  lang: [{ required: true, msg: t("pls_script_lang") }],
});
const { validate, reset, validateInfos } = useForm(modelRef, rulesRef);
watch(props, () => {
  if (!props.show) {
    setTimeout(() => {
      validateInfos.value = {};
    }, 200);
  }
});

const showModalRef = computed(() => {
  return props.show;
});

const emit = defineEmits<{
  (type: "submit", event: {}): void;
  (type: "cancel", event: {}): void;
}>();

const cancel = () => {
  emit("cancel", {});
};

const submit = () => {
  if(disabled.value) {
    return;
  }
  disabled.value = true;
  console.log("syncFromZentaoSubmit", console.log(modelRef.value));
  if (validate()) {
    emit("submit", modelRef.value);
  }
};

const clearFormData = () => {
  modelRef.value = {};
};

const store = useStore<{ Workspace: WorkspaceData, Zentao: ZentaoData }>();
const currSite = computed<any>(() => store.state.Zentao.currSite);
const currProduct = computed<any>(() => store.state.Zentao.currProduct);
const langs = computed<any[]>(() => store.state.Zentao.langs);
const modules = computed<any[]>(() => store.state.Zentao.modules);
const suites = computed<any[]>(() => store.state.Zentao.suites);
const tasks = computed<any[]>(() => store.state.Zentao.tasks);
const fetchData = () => {
  if(currSite.value.id == undefined || currSite.value.id <= 1
      || currProduct.value.id == undefined || currProduct.value.id <= 0) return;
  store.dispatch("Zentao/fetchModules", {
    siteId: currSite.value.id,
    productId: currProduct.value.id,
  });
  store.dispatch("Zentao/fetchSuites", {
    siteId: currSite.value.id,
    productId: currProduct.value.id,
  });
  store.dispatch("Zentao/fetchTasks", {
    siteId: currSite.value.id,
    productId: currProduct.value.id,
  });
};
fetchData();

watch(
  currProduct,
  () => {
    fetchData();
  },
  { deep: true }
);

watch(
  currSite,
  () => {
    fetchData();
  },
  { deep: true }
);

const selectWorkspace = () => {
  if (0 == modelRef.value.workspaceId) {
    modelRef.value.lang = "";
    return;
  }
  getWorkspace(parseInt(modelRef.value.workspaceId)).then((json) => {
    if (json.code === 0) {
      modelRef.value.lang = json.data.lang;
    }
  });
};
selectWorkspace();

defineExpose({
  clearFormData,
});
</script>