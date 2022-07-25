<template>
  <ZModal
    :showModal="showModalRef"
    @onCancel="cancel"
    @onOk="submit"
    :title="props.workspaceId ? t('edit_workspace') : t('create_workspace')"
    :contentStyle="{width: '600px'}"
  >
    <Form>
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
        v-for="(proxy, index) in proxyList"
        :key="index"
        labelWidth="100px"
        name="proxy_id[]"
        :label="index == 0 ? t('remote_proxy') : ''"
        :helpText="index == proxyList.length-1 ? t('priority_desc') : ''"
      >
        <div class="form-select">
          <select name="type" v-model="proxyList[index]">
            <option :value="0">{{t('local')}}</option>
            <option v-for="item in remoteProxies" :key="item.id" :value="item.id">
              {{ item.path }}
            </option>
          </select>
          <Button
          :disabled="index == 0"
          class="state only-icon"
          icon="arrow-up"
          size="sm"
          @click="changeSort(index, 'up')"
          >
          </Button>
          <Button
          :disabled="index == proxyList.length-1"
          class="state only-icon"
          icon="arrow-down"
          size="sm"
          @click="changeSort(index, 'down')"
          >
          </Button>
          <Button
          class="state only-icon"
          :icon="index ==0 ? 'add' : 'delete'"
          size="sm"
          @click="addProxySelect(index)"
          >
          </Button>
        </div>
      </FormItem>
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
import { get as getWorkspace } from "@/views/workspace/service";
import Icon from '@/components/Icon.vue';

export interface FormWorkspaceProps {
  show?: boolean;
  workspaceId?: number;
}
const { t } = useI18n();
const props = withDefaults(defineProps<FormWorkspaceProps>(), {
  show: false,
  workspaceId: 0,
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

const info = ref({} as any);
const proxyList = ref([0] as any[]);
const loadInfo = async () => {
    if(props.workspaceId === undefined || !props.workspaceId) return;
    await getWorkspace(props.workspaceId).then((json) => {
    if (json.code === 0) {
      info.value = json.data;
      modelRef.value = {
        id: info.value.id,
        name: info.value.name,
        path: info.value.path,
        type: info.value.type,
        lang: info.value.lang,
        cmd: info.value.cmd,
        proxy_id: info.value.proxy_id,
      };
      proxyList.value = info.value.proxies.split(',');
      selectType()
    }
  });
}

loadInfo();

const testTypes = ref([...ztfTestTypesDef, ...unitTestTypesDef]);
const store = useStore<{ Zentao: ZentaoData, proxy: ProxyData }>();
store.dispatch("proxy/fetchProxies");
const langs = computed<any[]>(() => store.state.Zentao.langs);
const remoteProxies = computed<any[]>(() => store.state.proxy.proxies);
const cmdSample = ref('')
const cmdMap = ref(arrToMap(testTypes.value))
const selectType = () => {
    console.log('selectType', modelRef.value.type)

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
const modelRef = ref({type: 'ztf', proxy_id: 0} as any);

const showCmd = computed(() => { return modelRef.value.type && modelRef.value.type !== 'ztf' })
const rulesRef = ref({
  name: [{ required: true, msg: t("pls_name") }],
  path: [{ required: true, msg: t("pls_workspace_path") }],
  lang: [{ required: true, msg: t("pls_script_lang") }],
  type: [{ required: true, msg: t("pls_workspace_type") }],
});

const { validate, reset, validateInfos } = useForm(modelRef, rulesRef);

const emit = defineEmits<{
  (type: "submit", event: {}): void;
  (type: "cancel", event: {}): void;
}>();

const submit = () => {
  modelRef.value.proxies = proxyList.value.join(',');
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

const addProxySelect = (index) => {
    if(index == 0){
        proxyList.value.push(0)
    }else{
        proxyList.value.splice(index, 1)
    }
}

const changeSort = (index, action) => {
    if(action == 'up'){
        if(index == 0){
            return;
        }
        [proxyList.value[index], proxyList.value[index-1]] = [proxyList.value[index-1], proxyList.value[index]]
    }else{
        if(index == proxyList.value.length-1){
            return;
        }
        [proxyList.value[index], proxyList.value[index+1]] = [proxyList.value[index+1], proxyList.value[index]]
    }
}

defineExpose({
  clearFormData,
});
</script>
<style  lang="less" scoped>
.form-select{
    display: flex;
    flex-direction: row;
    align-items: center;
    flex: 1;
    .inner-label{
        word-break: keep-all;
        margin: 0 10px;
    }
    .only-icon{
        outline:none;
        background-color: transparent;
        border: none;
        margin-left: 10px;
    }
}
</style>
