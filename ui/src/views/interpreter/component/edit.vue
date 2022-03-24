<template>
  <a-form :label-col="labelCol" :wrapper-col="wrapperCol">
    <a-form-item :label="t('script_lang')" v-bind="validateInfos.lang">
      <a-select v-model:value="modelRef.lang" @change="selectLang">
        <a-select-option key="" value="">&nbsp;</a-select-option>
        <a-select-option v-for="item in languages" :key="item" :value="item">{{ languageMap[item].name }}</a-select-option>
      </a-select>
    </a-form-item>

    <a-form-item :label="t('interpreter_path')" v-bind="validateInfos.path">
      <a-input-search v-if="isElectron" v-model:value="modelRef.path"
                      @search="selectDir" spellcheck="false"
                      @blur="validate('path', { trigger: 'blur' }).catch(() => {})">
        <template #enterButton>
          <a-button>选择</a-button>
        </template>
      </a-input-search>

      <a-input v-if="!isElectron" v-model:value="modelRef.path" spellcheck="false"
               @blur="validate('path', { trigger: 'blur' }).catch(() => {})"/>

      <div v-if="languageSettings?.path" class="t-italic">
        <div>当前{{languageMap[modelRef.lang]?.name}}运行环境：</div>
        <div>{{languageSettings.path}}</div>
        <div>{{languageSettings.info}}</div>
      </div>

      <div v-if="!languageSettings?.path && interpreterPath" class="t-italic">
        <div>{{languageMap[modelRef.lang]?.name}}可执行文件的路径类似：</div>
        <div>{{interpreterPath}}</div>
      </div>

    </a-form-item>

    <a-form-item :wrapper-col="{ span: wrapperCol.span, offset: labelCol.span }">
      <a-button type="primary" @click.prevent="save">{{ t('save') }}</a-button> &nbsp;
      <a-button style="margin-left: 10px" @click="reset">{{ t('reset') }}</a-button>
    </a-form-item>

  </a-form>
</template>

<script lang="ts">
import {defineComponent, reactive, ref, Ref, PropType, computed, ComputedRef} from "vue";
import {useI18n} from "vue-i18n";

import {validateInfos} from 'ant-design-vue/lib/form/useForm';
import {Form} from 'ant-design-vue';
import {getLangInterpreter, saveInterpreter} from "@/views/interpreter/service";
import {getLangSettings} from "../service";

const useForm = Form.useForm;

interface EditInterpreterFormSetupData {
  t: (key: string | number) => string;
  validate: any
  validateInfos: validateInfos;
  selectLang: (v) => void
  selectDir: () => void
  save: () => Promise<void>;
  reset: () => Promise<void>;

  interpreterPath: Ref<string>
  modelRef: Ref;
  languages: Ref<[]>,
  languageMap: Ref,
  languageSettings: Ref

  isElectron: Ref<boolean>;
  labelCol: any
  wrapperCol: any
}

export default defineComponent({
  name: 'EditInterpreterForm',
  props: {
    model: {
      type: Object as PropType<any>,
      required: true
    },

    onClose: {
      type: Function,
      required: true
    },
  },
  components: {},
  setup(props): EditInterpreterFormSetupData {
    const {t} = useI18n();
    const isElectron = ref(!!window.require)

    const languages = ref<any>({})
    const languageMap = ref<any>({})
    const languageSettings = ref({})

    const getInterpretersA = async () => {
      const data = await getLangSettings()
      languages.value = data.languages
      languageMap.value = data.languageMap
    }
    getInterpretersA()

    const interpreterPath = computed(() => {
      return languageMap.value[modelRef.value.lang]?.interpreter
    })

    let modelRef = ref<any>({
      id: props.model.value.id,
      lang: props.model.value.lang || '',
      path: props.model.value.path || '',
    });

    const rulesRef = reactive({
      lang: [{required: true, message: t('pls_input_lang')}],
      path: [{required: true, message: t('pls_input_interpreter_path')}],
    });

    const {resetFields, validate, validateInfos} = useForm(modelRef, rulesRef);

    const selectLang = async (item) => {
      console.log('selectLang', item)
      languageSettings.value = await getLangInterpreter(item)
    }

    const selectDir = () => {
      console.log('selectDir')

      if (isElectron.value) {
        const {dialog} = window.require('@electron/remote');
        dialog.showOpenDialog({
          properties: ['openDirectory']
        }).then(result => {
          if (result.filePaths && result.filePaths.length > 0) {
            modelRef.value.path = result.filePaths[0]
          }
        }).catch(err => {
          console.log(err)
        })
      }
    }

    const save = async () => {
      validate()
        .then(() => {
          saveInterpreter(modelRef.value).then((json) => {
            if (json.code === 0) {
              props.onClose()
            }
          })
        })
    }

    const reset = async () => {
      resetFields()
    }

    return {
      t,
      isElectron,

      interpreterPath,
      validate,
      validateInfos,
      modelRef,
      selectLang,
      selectDir,
      save,
      reset,

      languages,
      languageMap,
      languageSettings,
      labelCol: {span: 6},
      wrapperCol: {span: 18},
    }

  }
})
</script>