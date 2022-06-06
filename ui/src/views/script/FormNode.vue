<template>
  <ZModal
    :showModal="showModalRef"
    @onCancel="cancel"
    @onOk="submit"
    :title="t('create')"
    :contentStyle="{width: '400px'}"
  >
    <Form labelCol="6" wrapperCol="16">
      <FormItem name="name" :label="t('name')" :info="validateInfos.name">
        <input type="text" v-model="modelRef.name" />
      </FormItem>
    </Form>
  </ZModal>
</template>

<script setup lang="ts">
import { useI18n } from "vue-i18n";
import {
  computed,
  defineExpose,
  onMounted,
  withDefaults,
  ref,
  defineProps,
  defineEmits,
  watch,
} from "vue";
import { useForm } from "@/utils/form";
import Form from "@/components/Form.vue";
import FormItem from "@/components/FormItem.vue";

export interface FormWorkspaceProps {
  show?: boolean;
}
const { t } = useI18n();
const props = withDefaults(defineProps<FormWorkspaceProps>(), {
  show: false,
});

watch(props, () => {
    if(!props.show){
        setTimeout(() => {
            validateInfos.value = {};
        }, 200);
    }
})

const showModalRef = computed(() => {
  return props.show;
});

const cancel = () => {
  emit("cancel", {});
};

const modelRef = ref({});
const rulesRef = ref({
  name: [{ required: true, msg: t("pls_name") }],
});

const { validate, reset, validateInfos } = useForm(modelRef, rulesRef);

const emit = defineEmits<{
  (type: "submit", event: {}): void;
  (type: "cancel", event: {}): void;
}>();

const submit = () => {
  if (validate()) {
    emit("submit", modelRef.value);
  }
};

const clearFormData = () => {
  modelRef.value = {};
};

defineExpose({
  clearFormData,
});
</script>
