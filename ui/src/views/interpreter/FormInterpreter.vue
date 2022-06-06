<template>
  <ZModal
    :showModal="showModalRef"
    @onCancel="cancel"
    @onOk="submit"
    :title="info == undefined ? t('create_interpreter') : t('edit_interpreter')" 
    :contentStyle="{width: '600px'}"
  >
    <Form class="form-interpreter" labelCol="6" wrapperCol="16">
      <FormItem
        name="lang"
        :label="t('script_lang')"
        :info="validateInfos.lang"
      >
        <select
          name="lang"
          v-model="modelRef.lang"
          @change="selectLang"
        >
          <option v-for="item in languages" :key="item" :value="item">
            {{ languageMap[item].name }}
          </option>
        </select>
      </FormItem>
      <FormItem
        v-if="isElectron"
        name="path"
        :label="t('interpreter_path')"
        :info="validateInfos.path"
      >
        <input
          v-model="modelRef.path"
        />
        <Button  v-if="isElectron" @click="selectFile" class="state secondary select-dir-btn">{{t('select')}}</Button>
      </FormItem>
      <FormItem
        v-if="!isElectron"
        name="path"
        :label="t('interpreter_path')"
        :info="validateInfos.path"
      >
        <input v-model="modelRef.path" />
      </FormItem>
      <FormItem
        name="lang"
        v-if="interpreterInfos.length > 0"
        :label="t('script_lang')"
      >
        <select
          name="type"
          v-model="selectedInterpreter"
          @change="selectInterpreter"
        >
          <option value="">
            {{ t("find_to_select", { num: interpreterInfos.length }) }}
          </option>
          <option
            v-for="item in interpreterInfos"
            :key="item.path"
            :value="item.path"
          >
            {{ item.info.length > 50 ? item.info.substring(0,47) + '...' : item.info }}
          </option>
        </select>
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

import settings from "@/config/settings";
import { useForm } from "@/utils/form";
import { getElectron } from "@/utils/comm";

import { StateType } from "@/views/site/store";
import { getLangSettings, getLangInterpreter } from "@/views/interpreter/service";

import Form from "@/components/Form.vue";
import FormItem from "@/components/FormItem.vue";
import Button from "@/components/Button.vue";

export interface FormSiteProps {
  show?: boolean;
  info?: any;
}
const { t } = useI18n();
const isElectron = ref(getElectron());

const props = withDefaults(defineProps<FormSiteProps>(), {
  show: false,
  info: {},
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

const languages = ref<any>({})
const languageMap = ref<any>({})
const interpreterInfos = ref([]);

const getInterpretersA = async () => {
    const data = await getLangSettings();
    languages.value = data.languages;
    languageMap.value = data.languageMap;
};
getInterpretersA();

const store = useStore<{ Site: StateType }>();

const selectedInterpreter = ref("");
const selectInterpreter = async () => {
  console.log("selectInterpreter", selectedInterpreter.value);
  modelRef.value.path = selectedInterpreter.value;
};

const selectLang = async (item) => {
  console.log("selectLang", modelRef.value.lang);

  modelRef.value.path = "";
  selectedInterpreter.value = "";

  if (modelRef.value.lang === "") {
    interpreterInfos.value = [];
    return;
  }

  interpreterInfos.value = await getLangInterpreter(modelRef.value.lang);
  if (interpreterInfos.value == null) interpreterInfos.value = [];
  console.log(interpreterInfos.value);
};

watch(props, () => {
  console.log("watch formInterpreter", props);
  if (props.info.value == undefined) {
    modelRef.value = {
      id: 0,
      lang: "",
      path: "",
    };
    interpreterInfos.value = [];
  } else {
    modelRef.value.id = props.info.value.id;
    modelRef.value.path = props.info.value.path;
    modelRef.value.lang = props.info.value.lang;
    interpreterInfos.value = [];
  }
});

const cancel = () => {
  emit("cancel", {});
};

const modelRef = ref<any>({
  id: 0,
  lang: "",
  path: "",
});
const rulesRef = ref({
  lang: [{ required: true, msg: t("pls_lang") }],
  path: [{ required: true, msg: t("pls_input_interpreter_path") }],
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
  selectedInterpreter.value = "";
  interpreterInfos.value = [];
};

const selectFile = () => {
    console.log('selectFile')

    const { ipcRenderer } = window.require('electron')
    ipcRenderer.send(settings.electronMsg, 'selectFile')

    ipcRenderer.on(settings.electronMsgReplay, (event, arg) => {
    console.log(arg)
    modelRef.value.path = arg
    })
}
defineExpose({
  clearFormData,
});
</script>
