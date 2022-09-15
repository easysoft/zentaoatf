<template>
  <ZModal
    :showModal="showModalRef"
    @onCancel="cancel"
    @onOk="submit"
    :title="info.id == 0 ? t('create_remote_proxy') : t('edit_remote_proxy')"
    :contentStyle="{ width: '500px' }"
  >
    <Form class="form-proxy" labelCol="6" wrapperCol="16">
      <!-- <FormItem 
        name="type" 
        :label="t('type')" 
        :info="validateInfos.type" 
        labelWidth="100px"
        >
        <div class="select">
          <select name="type" v-model="modelRef.type">
            <option
                v-for="item in testTypes"
                :key="item.value"
                :value="item.value"
            >
                {{ item.label }}
            </option>
          </select>
        </div>
      </FormItem>
      <FormItem
        v-if="modelRef.type === 'ztf'"
        name="lang"
        labelWidth="100px"
        :label="t('proxy_lang')"
        :info="validateInfos.lang"
      >
        <div class="select">
          <select name="lang" v-model="modelRef.lang">
            <option v-for="item in langs" :key="item.code" :value="item.code">
              {{ item.name }}
            </option>
          </select>
        </div>
      </FormItem> -->
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

import { StateType } from "@/views/site/store";
import { unitTestTypesDef, ztfTestTypesDef } from "@/utils/const";
import { ZentaoData } from "@/store/zentao";

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
    type: "ztf",
    lang: "",
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

const testTypes = ref([...ztfTestTypesDef, ...unitTestTypesDef]);
const proxyInfos = ref([]);
// const langs = computed<any[]>(() => store.state.Zentao.langs);

const store = useStore<{ Site: StateType; Zentao: ZentaoData }>();

const cancel = () => {
  emit("cancel", {});
};

const modelRef = ref<any>({
  id: info.value.id,
  type: info.value.type,
  lang: info.value.lang,
  name: info.value.name,
  path: info.value.path,
});
const rulesRef = ref({
  name: [{ required: true, msg: t("pls_name") }],
  //   lang: [{ required: true, msg: t("pls_lang") }],
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
    emit("submit", modelRef.value);
  }
};

const clearFormData = () => {
  console.log("clear");
  modelRef.value.path = "";
  modelRef.value.lang = "";
  modelRef.value.name = "";
  proxyInfos.value = [];
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
