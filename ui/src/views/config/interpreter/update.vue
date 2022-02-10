<template>
  <a-modal
      :destroy-on-close="true"
      :mask-closable="false"
      :title="t('edit_interpreter')"
      :visible="visible"
      :onCancel="onCancel"
      width="600px"
  >
    <template #footer>
      <a-button key="submit" type="primary" :loading="onSubmitLoading" @click="onFinish">{{t('save')}}</a-button>
      <a-button key="back" @click="() => onCancel()">{{t('cancel')}}</a-button>
    </template>

    <a-form :labelCol="{ span: 4 }" :wrapper-col="{span:20}">
      <a-form-item :label="t('script_lang')" v-bind="validateInfos.lang">
        {{languageMap[modelRef.lang]}}
      </a-form-item>
      <a-form-item :label="t('interpreter_path')" v-bind="validateInfos.val">
        <a-input v-model:value="modelRef.val" placeholder=""/>
      </a-form-item>
    </a-form>
  </a-modal>
</template>
<script lang="ts">
import {defineComponent, onMounted, PropType, reactive, Ref} from "vue";
import {useI18n} from "vue-i18n";

import {Props, validateInfos} from 'ant-design-vue/lib/form/useForm';
import {Form, message} from 'ant-design-vue';
import {Interpreter} from "@/views/config/data";

const useForm = Form.useForm;

interface UpdateInterpreterFormSetupData {
  t: (key: string | number) => string;
  modelRef: Ref<Interpreter>;
  validateInfos: validateInfos;
  onFinish: () => Promise<void>;
}

export default defineComponent({
  name: 'UpdateInterpreterForm',
  props: {
    visible: {
      type: Boolean,
      required: true
    },
    model: {
      type: Object as PropType<any>,
      required: true
    },
    languageMap: {
      required: true
    },
    onCancel: {
      type: Function,
      required: true
    },
    onSubmitLoading: {
      type: Boolean,
      required: true
    },
    onSubmit: {
      type: Function as PropType<(values: any, resetFields: (newValues?: Props | undefined) => void) => void>,
      required: true
    }
  },
  components: {},
  setup(props): UpdateInterpreterFormSetupData {
    const {t} = useI18n();

    let modelRef = reactive<any>({
      lang: props.model.value.lang || '',
      val: props.model.value.val || '',
    });

    const rulesRef = reactive({
      lang: [{required: true, message: 'pls_input_lang'}],
      val: [{required: true, message: t('pls_input_interpreter_path')}],
    });

    const {resetFields, validate, validateInfos} = useForm(modelRef, rulesRef);
    const onFinish = async () => {
      try {
        const fieldsValue = await validate<any>();
        props.onSubmit(fieldsValue, resetFields);
      } catch (error) {
        message.warning(t('app.global.form.validatefields.catch'));
      }
    };

    return {
      t,
      modelRef,
      validateInfos,
      onFinish
    }
  }
})
</script>