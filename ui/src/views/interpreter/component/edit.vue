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
    </a-form-item>

    <a-form-item :wrapper-col="{ span: wrapperCol.span, offset: labelCol.span }">
      <a-select v-if="interpreterInfos.length > 0" v-model:value="selectedInterpreter" @change="selectInterpreter">
        <a-select-option value="">找到 {{interpreterInfos.length}} 项可选</a-select-option>
        <a-select-option v-for="item in interpreterInfos" :key="item.path" :value="item.path">
          {{ item.info }}
        </a-select-option>
      </a-select>
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
  setup(props) {
    const {t} = useI18n();
    const isElectron = ref(!!window.require)

    const languages = ref<any>({})
    const languageMap = ref<any>({})

    const getInterpretersA = async () => {
      const data = await getLangSettings()
      languages.value = data.languages
      languageMap.value = data.languageMap
    }
    getInterpretersA()

    const interpreterInfos = ref([])

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

      modelRef.value.path = ''
      selectedInterpreter.value = ''

      if (item === '') {
        interpreterInfos.value = []
        return
      }

      interpreterInfos.value = await getLangInterpreter(item)
      console.log(interpreterInfos.value)
    }
    const selectedInterpreter = ref('')
    const selectInterpreter= async (item) => {
      console.log('selectInterpreter', item)
      modelRef.value.path = item
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

      interpreterInfos,
      validate,
      validateInfos,
      modelRef,
      selectLang,
      selectInterpreter,
      selectedInterpreter,
      selectDir,
      save,
      reset,

      languages,
      languageMap,
      labelCol: {span: 6},
      wrapperCol: {span: 18},
    }

  }
})
</script>