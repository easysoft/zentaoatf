<template>
  <ZModal
    :showModal="showModalRef"
    @onCancel="cancel"
    @onOk="submit"
    :title="info.id == 0 ? t('create_interpreter') : t('edit_interpreter')"
    :contentStyle="{width: '500px'}"
  >
    <Form class="form-interpreter">
      <FormItem
        v-if="props.proxyId > 0"
        labelWidth="100px"
        name="proxy_id"
        :label="t('copy_from')"
        :info="validateInfos.proxy_id"
      >
        <select
          name="proxy_id"
          v-model="modelRef.proxy_id"
          @change="selectProxy"
        >
          <option :value="-1">{{t('local')}}</option>
          <option v-for="item in remoteProxies" :key="item.id" :value="item.id">
            {{ item.name }}
          </option>
        </select>
      </FormItem>
      <FormItem
        labelWidth="100px"
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
        labelWidth="100px"
        v-if="isElectron"
        name="path"
        :label="t('interpreter_path')"
        :info="validateInfos.path"
      >
        <input type="text"
          v-model="modelRef.path"
        />
        <Button  v-if="isElectron" @click="selectFile" class="state secondary select-dir-btn">{{t('select')}}</Button>
      </FormItem>
      <FormItem
        labelWidth="100px"
        v-if="!isElectron"
        name="path"
        :label="t('interpreter_path')"
        :info="validateInfos.path"
      >
        <input type="text" v-model="modelRef.path" />
      </FormItem>
      <FormItem
        labelWidth="100px"
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
import { ProxyData } from "@/store/proxy";

import Form from "@/components/Form.vue";
import FormItem from "@/components/FormItem.vue";
import Button from "@/components/Button.vue";

export interface FormSiteProps {
  show?: boolean;
  info?: any;
  proxyPath?: string;
  proxyId?: number;
}
const { t } = useI18n();
const isElectron = ref(getElectron());

const props = withDefaults(defineProps<FormSiteProps>(), {
  show: false,
  proxyPath: 'local',
  info: ref({
    id: 0,
    lang: "",
    path: "",
  }),
});
const info = computed(() => {
  return props.info.value == undefined ? {id: 0,lang: "",path: ""} : props.info.value;
});
const showModalRef = computed(() => {
  return props.show;
});

const languages = ref<any>({})
const languageMap = ref<any>({})
const interpreterInfos = ref([]);

const getInterpretersA = async () => {
    const data = await getLangSettings(props.proxyPath);
    languages.value = data.languages;
    languageMap.value = data.languageMap;
};
getInterpretersA();

const store = useStore<{ Site: StateType, proxy: ProxyData }>();

const remoteProxies = computed<any[]>(() => store.state.proxy.proxies.filter(item => {
    if(item.id != props.proxyId){
        return true;
    }
    return false
}));
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

  interpreterInfos.value = await getLangInterpreter(modelRef.value.lang, props.proxyPath);
  if (interpreterInfos.value == null) interpreterInfos.value = [];
  console.log(interpreterInfos.value);
};

const cancel = () => {
  emit("cancel", {});
};

const modelRef = ref<any>({
  id: info.value.id,
  lang: info.value.lang,
  path: info.value.path,
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

const selectProxy = () => {
    if(modelRef.value.proxy_id == -1){
        rulesRef.value.lang = [];
        rulesRef.value.path = [];
    }else{
        rulesRef.value.lang = [{ required: true, msg: t("pls_lang") }];
        rulesRef.value.path = [{ required: true, msg: t("pls_input_interpreter_path") }];
        remoteProxies.value.forEach(item => {
            if(item.id == modelRef.value.proxy_id){
                modelRef.value.proxyPath = item.path;
            }
        })
    }
}
defineExpose({
  clearFormData,
});
</script>
<style scoped>
.select-dir-btn{
  width: 60px;
}
</style>
