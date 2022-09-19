<template>
  <ZModal
    :showModal="showModalRef"
    @onCancel="cancel"
    @onOk="submit"
    :title="info.id == 0 ? t('create_remote_proxy') : t('edit_remote_proxy')"
    :contentStyle="{ width: '500px' }"
  >
    <Form class="form-proxy" labelCol="6" wrapperCol="16">
      <FormItem
        name="name"
        labelWidth="100px"
        :label="t('name')"
        :info="validateInfos.name"
      >
        <input type="text" v-model="modelRef.name" />
      </FormItem>
      <FormItem
        labelWidth="100px"
        name="path"
        :label="t('proxy_link')"
        :info="validateInfos.path"
        :helpText="t('proxy_desc')"
      >
        <input placeholder="http://127.0.0.1:8085" type="text" v-model="modelRef.path" />
      </FormItem>
    </Form>
  </ZModal>
</template>

<script setup lang="ts">
import { useI18n } from "vue-i18n";
import { useStore } from "vuex";
import {
  computed,
  defineExpose,
  withDefaults,
  ref,
  defineProps,
  defineEmits,
  watch,
} from "vue";

import { useForm } from "@/utils/form";

import Form from "@/components/Form.vue";
import FormItem from "@/components/FormItem.vue";

export interface FormSiteProps {
  show?: boolean;
  info?: any;
}
const { t } = useI18n();

const props = withDefaults(defineProps<FormSiteProps>(), {
  show: false,
  info: {
    id: 0,
    path: "",
    name: "",
  },
});
const info = computed(() =>
  props.info.value == undefined
    ? { id: 0, type: "ztf", lang: "", path: "" }
    : props.info.value
);

const showModalRef = computed(() => {
  return props.show;
});

const cancel = () => {
  emit("cancel", {});
};

const modelRef = ref<any>({
  id: info.value.id,
  name: info.value.name,
  path: info.value.path,
});
const rulesRef = ref({
  name: [{ required: true, msg: t("pls_name") }],
  path: [{ required: true, msg: t("pls_input_proxy_link") }],
});
const { validate, reset, validateInfos } = useForm(modelRef, rulesRef);

const emit = defineEmits<{
  (type: "submit", event: {}): void;
  (type: "cancel", event: {}): void;
}>();

const submit = () => {
  if (validate()) {
    console.log("submit", validate());
    emit("submit", {proxyPath: props.info.proxyPath, ...modelRef.value});
  }
};

const clearFormData = () => {
  console.log("clear");
  modelRef.value.path = "";
  modelRef.value.name = "";
};

defineExpose({
  clearFormData,
});
</script>
<style scoped>
.select-dir-btn {
  width: 60px;
}
</style>
