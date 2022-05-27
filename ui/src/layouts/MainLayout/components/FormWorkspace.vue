<template>
  <ZModal
    :showModal="showModalRef"
    @onCancel="cancel"
    @onOk="submit"
    :title="t('create_workspace')"
    :contentStyle="{width: '600px'}"
  >
    <Form labelCol="6" wrapperCol="16">
      <FormItem name="name" :label="t('name')" :info="validateInfos.name">
        <input v-model="modelRef.name" />
      </FormItem>
      <FormItem name="path" :label="t('path')" :info="validateInfos.path">
        <input v-model="modelRef.path" />
      </FormItem>
      <FormItem name="type" :label="t('type')" :info="validateInfos.type">
        <select name="type" v-model="modelRef.type">
          <option
            v-for="item in testTypes"
            :key="item.value"
            :value="item.value"
          >
            {{ item.label }}
          </option>
        </select>
      </FormItem>
      <FormItem
        name="lang"
        :label="t('default_lang')"
        :info="validateInfos.lang"
      >
        <select name="type" v-model="modelRef.lang">
          <option v-for="item in langs" :key="item.code" :value="item.code">
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
  defineEmits,
} from "vue";
import { useForm } from "@/utils/form";
import Form from "./Form.vue";
import FormItem from "./FormItem.vue";

export interface FormWorkspaceProps {
  show?: boolean;
}
const { t } = useI18n();
const props = withDefaults(defineProps<FormWorkspaceProps>(), {
  show: false,
});

const showModalRef = computed(() => {
  return props.show;
});
const testTypes = ref([...ztfTestTypesDef, ...unitTestTypesDef]);
const zentaoStore = useStore<{ Zentao: ZentaoData }>();
const langs = computed<any[]>(() => zentaoStore.state.Zentao.langs);

const cancel = () => {
  emit("cancel", {});
};

const modelRef = ref({});
const rulesRef = ref({
  name: [{ required: true, msg: t("pls_name") }],
  path: [{ required: true, msg: t("pls_workspace_path") }],
  lang: [{ required: true, msg: t("select_ui_lang") }],
  type: [{ required: true, msg: t("pls_workspace_type") }],
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

<style lang="less" scoped>
.workdir {
  height: calc(100vh - 80px);
}
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