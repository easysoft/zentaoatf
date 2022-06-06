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
        <input v-if="isElectron" v-model="modelRef.path"
             />
        <Button  v-if="isElectron" @click="selectDir" class="state secondary select-dir-btn">{{t('select')}}</Button>
        <input v-if="!isElectron" v-model="modelRef.path" />
      </FormItem>
      <FormItem name="type" :label="t('type')" :info="validateInfos.type">
        <select name="type" @change="selectType" v-model="modelRef.type">
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
        v-if="modelRef.type === 'ztf'"
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
      <FormItem v-if="showCmd" name="cmd" :label="t('cmd')" :info="validateInfos.cmd">
        <textarea v-model="modelRef.cmd" />
        <div class="t-tips" style="margin-top: 5px;">
          <div>{{ t('tips_test_cmd', {cmd: cmdSample}) }}</div>
        </div>
      </FormItem>
    </Form>
  </ZModal>
</template>

<script setup lang="ts">
import { useI18n } from "vue-i18n";
import { useStore } from "vuex";
import { ZentaoData } from "@/store/zentao";

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
const store = useStore<{ Zentao: ZentaoData }>();
const langs = computed<any[]>(() => store.state.Zentao.langs);
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
const modelRef = ref({} as any);

const showCmd = computed(() => { return modelRef.value.type && modelRef.value.type !== 'ztf' })
const rulesRef = ref({
  name: [{ required: true, msg: t("pls_name") }],
  path: [{ required: true, msg: t("pls_workspace_path") }],
  lang: [{ required: true, msg: t("pls_script_lang") }],
  type: [{ required: true, msg: t("pls_workspace_type") }],
  cmd: [
        {
          trigger: 'blur',
          validator: async (rule: any, value: string) => {
            if (modelRef.value.type !== 'ztf' && (value === '' || !value)) {
              throw new Error(t('pls_cmd'));
            }
          }
        },
      ],
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

const selectDir = () => {
    console.log('selectDir')

    const { ipcRenderer } = window.require('electron')
    ipcRenderer.send(settings.electronMsg, 'selectDir')

    ipcRenderer.on(settings.electronMsgReplay, (event, arg) => {
    console.log(arg)
    modelRef.value.path = arg
    })
}

defineExpose({
  clearFormData,
});
</script>

<style>
.select-dir-btn{
    position: absolute;
    right: 20px;
}
</style>