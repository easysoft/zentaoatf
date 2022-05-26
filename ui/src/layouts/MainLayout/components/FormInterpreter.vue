<template>
  <ZModal
    :showModal="showModalRef"
    @onCancel="cancel"
    @onOk="submit"
    :title="t('script_lang')"
  >

    <Form class="form-site" labelCol="6" wrapperCol="16">
        <FormItem name="lang" :label="t('script_lang')" :info="validateInfos.lang">
        <select name="type" v-model="modelRef.lang" class="form-control">
          <option
            v-for="item in languages"
            :key="item.value"
            :value="item.value"
          >
            {{ languageMap[item.value].name }}
          </option>
        </select>
      </FormItem>
<FormItem name="name" :label="t('interpreter_path')" :info="validateInfos.path">
        <input v-model="modelRef.path" class="form-control" @change="selectFile" />
      </FormItem>
<FormItem v-if="!isElectron" name="name" :label="t('interpreter_path')" :info="validateInfos.path">
        <input v-model="modelRef.path" class="form-control" @change="selectFile" />
      </FormItem>

      <a-input v-if="!isElectron" v-model:value="modelRef.path" spellcheck="false"
               @blur="validate('path', { trigger: 'blur' }).catch(() => {})"/>
    </a-form-item>

    <a-form-item :wrapper-col="{ span: wrapperCol.span, offset: labelCol.span }">
      <a-select v-if="interpreterInfos.length > 0" v-model:value="selectedInterpreter" @change="selectInterpreter">
        <a-select-option value="">{{ t('find_to_select', {num: interpreterInfos.length})}}</a-select-option>
        <a-select-option v-for="item in interpreterInfos" :key="item.path" :value="item.path">
          {{ item.info }}
        </a-select-option>
      </a-select>
    </a-form-item>

    <a-form-item :wrapper-col="{ span: wrapperCol.span, offset: labelCol.span }"
                 :class="{'t-dir-right': !isWin}" class="t-right">
      <a-button type="primary" @click.prevent="save" class="t-btn-gap">{{ t('save') }}</a-button> &nbsp;
      <a-button @click="reset" class="t-btn-gap">{{ t('reset') }}</a-button>
    </a-form-item>
      <FormItem name="name" :label="t('name')" :info="validateInfos.name">
        <input v-model="modelRef.name" class="form-control" />
      </FormItem>
      <FormItem name="url" :label="t('zentao_url')" :info="validateInfos.url">
        <input v-model="modelRef.url" class="form-control" />
      </FormItem>
      <FormItem
        name="username"
        :label="t('username')"
        :info="validateInfos.username"
      >
        <input v-model="modelRef.username" class="form-control" />
      </FormItem>
      <FormItem
        name="password"
        :label="t('password')"
        :info="validateInfos.password"
      >
        <input v-model="modelRef.password" class="form-control" />
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
  watch,
} from "vue";
import { useForm } from "@/utils/form";
import Form from "./Form.vue";
import FormItem from "./FormItem.vue";
import { StateType } from "@/views/site/store";
import {getLangSettings} from "@/views/interpreter/service";

export interface FormSiteProps {
  show?: boolean;
  id?: number;
}
const { t } = useI18n();

const languages = ref<any>({})
    const languageMap = ref<any>({})

    const getInterpretersA = async () => {
      const data = await getLangSettings()
      languages.value = data.languages
      languageMap.value = data.languageMap
    }

const props = withDefaults(defineProps<FormSiteProps>(), {
  show: false,
  id: 0
});

const showModalRef = computed(() => {
  return props.show;
});
const store = useStore<{ Site: StateType }>();
const get = async (id: number): Promise<void> => {
  await store.dispatch("Site/get", id);
};

watch(props, () => {
        get(props.id);
})
get(props.id);
const cancel = () => {
  emit("cancel", {});
};

const modelRef = computed(() => store.state.Site.detailResult);
const rulesRef = ref({
  name: [{ required: true, msg: t("pls_name") }],
  url: [
    {
      required: true,
      msg: t("pls_zentao_url"),
    },
    {
      regex:
        /(http?|https):\/\/[-A-Za-z0-9+&@#/%?=~_|!:,.;]+[-A-Za-z0-9+&@#/%=~_|]/i,
      msg: t("pls_zentao_url"),
    },
  ],
  username: [{ required: true, msg: t("pls_username") }],
  password: [{ required: true, msg: t("pls_password") }],
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
.form-site {
  min-width: 500px;
}
</style>