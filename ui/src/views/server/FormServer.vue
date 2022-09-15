<template>
  <ZModal
    :showModal="showModalRef"
    @onCancel="cancel"
    @onOk="submit"
    :title="info.id == 0 ? t('create_remote_server') : t('edit_remote_server')"
    :contentStyle="{ width: '500px' }"
  >
    <Form class="form-server" labelCol="6" wrapperCol="16">
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
        :label="t('server_link')"
        :info="validateInfos.path"
      >
        <input type="text" placeholder="http://127.0.0.1:8085" v-model="modelRef.path" />
      </FormItem>
    </Form>
  </ZModal>
</template>

<script setup lang="ts">
import { useI18n } from "vue-i18n";
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
import Button from "@/components/Button.vue";

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
    ? { id: 0, type: "ztf", lang: "", path: "", name: "" }
    : props.info.value
);

const showModalRef = computed(() => {
  return props.show;
});

const serverInfos = ref([]);

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
  path: [{ required: true, msg: t("pls_input_server_link") }],
});
const { validate, reset, validateInfos } = useForm(modelRef, rulesRef);

const emit = defineEmits<{
  (type: "submit", event: {}): void;
  (type: "cancel", event: {}): void;
}>();

const submit = () => {
  if (validate()) {
    console.log("submit", validate());
    emit("submit", modelRef.value);
  }
};

const clearFormData = () => {
  console.log("clear");
  modelRef.value.name = "";
  modelRef.value.path = "";
  serverInfos.value = [];
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
