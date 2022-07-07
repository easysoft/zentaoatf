<template>
  <ZModal
    :showModal="showModalRef"
    @onCancel="cancel"
    @onOk="submit"
    :title="t('create_workspace')"
    :contentStyle="{width: '600px'}"
  >
    <Form labelCol="6" wrapperCol="16">
      <FormItem labelWidth="100px" name="name" :label="t('name')" :info="validateInfos.name">
        <input type="text" v-model="modelRef.name" />
      </FormItem>
      <FormItem labelWidth="100px" name="path" :label="t('path')" :info="validateInfos.path">
        <input type="text" v-if="isElectron" v-model="modelRef.path" />
        <Button  v-if="isElectron" @click="selectDir('path')" class="state secondary flex-none rounded">{{t('select')}}</Button>
        <input type="text" v-if="!isElectron" v-model="modelRef.path" />
      </FormItem>
      <FormItem labelWidth="100px" name="type" :label="t('type')" :info="validateInfos.type">
        <div class="select">
          <select name="type" @change="selectType" v-model="modelRef.type">
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
        labelWidth="100px"
        v-if="modelRef.type === 'ztf'"
        name="lang"
        :label="t('default_lang')"
        :info="validateInfos.lang"
      >
        <div class="select">
          <select name="type" v-model="modelRef.lang">
            <option v-for="item in langs" :key="item.code" :value="item.code">
              {{ item.name }}
            </option>
          </select>
        </div>
      </FormItem>
      <FormItem labelWidth="100px" v-if="showCmd" name="cmd" :label="t('cmd')" :info="validateInfos.cmd" :helpText="t('tips_test_cmd', {cmd: cmdSample})">
        <textarea v-model="modelRef.cmd" />
      </FormItem>
      <FormItem
        labelWidth="100px"
        name="proxy_id"
        :label="t('remote_proxy')"
        :info="validateInfos.proxy_id"
      >
        <div class="select">
          <select name="type" @change="updateRules" v-model="modelRef.proxy_id">
            <option :value="0">{{t('local')}}</option>
            <option v-for="item in remoteProxies" :key="item.id" :value="item.id">
              {{ item.path }}
            </option>
          </select>
        </div>
      </FormItem>
      <div v-if="modelRef.proxy_id">
        <FormItem
        labelWidth="100px"
        name="auth_type"
        :label="t('git_auth_type')"
        :info="validateInfos.auth_type"
        >
        <div class="select">
            <select name="auth_type" @change="updateRules" v-model="modelRef.auth_type">
            <option value="ssh">ssh</option>
            <option value="password">password</option>
            </select>
        </div>
        </FormItem>
        <FormItem
        labelWidth="100px"
        v-if="modelRef.auth_type === 'password'"
        name="username"
        :label="t('username')"
        :info="validateInfos.username"
        >
        <input type="text" v-model="modelRef.username" />
        </FormItem>
        <FormItem 
          v-if="modelRef.auth_type === 'ssh'" 
          labelWidth="100px" 
          name="rsa_key" 
          :label="t('rsa_path')" 
          :info="validateInfos.rsa_key"
          >
          <input type="text" v-if="isElectron" v-model="modelRef.rsa_key" />
          <Button  v-if="isElectron" @click="selectDir('rsa_key')" class="state secondary flex-none rounded">{{t('select')}}</Button>
          <input type="text" v-if="!isElectron" v-model="modelRef.rsa_key" />
        </FormItem>
        <FormItem
        labelWidth="100px"
        name="password"
        :label="t('password')"
        :info="validateInfos.password"
        >
        <input type="password" v-model="modelRef.password" />
        </FormItem>
      </div>
    </Form>
  </ZModal>
</template>

<script setup lang="ts">
import { useI18n } from "vue-i18n";
import { useStore } from "vuex";
import { ZentaoData } from "@/store/zentao";
import { ProxyData } from "@/store/proxy";

import { unitTestTypesDef, ztfTestTypesDef } from "@/utils/const";
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
import {arrToMap} from "@/utils/array";
import settings from "@/config/settings";
import Button from "@/components/Button.vue";

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

const testTypes = ref([...ztfTestTypesDef, ...unitTestTypesDef]);
const store = useStore<{ Zentao: ZentaoData, proxy: ProxyData }>();
store.dispatch("proxy/fetchProxies");
const langs = computed<any[]>(() => store.state.Zentao.langs);
const remoteProxies = computed<any[]>(() => store.state.proxy.proxies);
const cmdSample = ref('')
const cmdMap = ref(arrToMap(testTypes.value))
const selectType = () => {
    console.log('selectType')

    if (modelRef.value.type !== 'ztf') {
        cmdSample.value = cmdMap.value[modelRef.value.type].cmd
        modelRef.value.cmd = cmdSample.value.split('product_id')[1].trim()
        rulesRef.value.lang = [{ required: false, msg: t("pls_script_lang") }]
    }else{
        rulesRef.value.lang = [{ required: true, msg: t("pls_script_lang") }]
    }
}

const cancel = () => {
  emit("cancel", {});
};

const isElectron = ref(!!window.require)
const modelRef = ref({type: 'ztf', proxy_id: 0, auth_type: 'ssh'} as any);

const showCmd = computed(() => { return modelRef.value.type && modelRef.value.type !== 'ztf' })
const rulesRef = ref({
  name: [{ required: true, msg: t("pls_name") }],
  path: [{ required: true, msg: t("pls_workspace_path") }],
  lang: [{ required: true, msg: t("pls_script_lang") }],
  type: [{ required: true, msg: t("pls_workspace_type") }],
  password: [],
  username: [],
  rsa_key: [],
});

const updateRules = () => {
    if (modelRef.value.proxy_id && modelRef.value.auth_type === 'password') {
        rulesRef.value.username = [{ required: true, msg: t("pls_username") }];
        rulesRef.value.password = [{ required: true, msg: t("pls_password") }];
    } else if(modelRef.value.proxy_id && modelRef.value.auth_type === 'ssh') {
        rulesRef.value.rsa_key = [{ required: true, msg: t("pls_rsa") }]
    }
}

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
  modelRef.value = {type: 'ztf'};
};

const selectDir = (key) => {
    console.log('selectDir')

    const { ipcRenderer } = window.require('electron')
    ipcRenderer.send(settings.electronMsg, 'selectDir')

    ipcRenderer.on(settings.electronMsgReplay, (event, arg) => {
    console.log(arg)
    modelRef.value[key] = arg
    })
}

defineExpose({
  clearFormData,
});
</script>
