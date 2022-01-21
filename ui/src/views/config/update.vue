<template>
  <a-modal
      :destroy-on-close="true"
      :mask-closable="false"
      title="编辑解析器"
      :visible="visible"
      :onCancel="onCancel"
      width="600px"
  >
    <template #footer>
      <a-button key="back" @click="() => onCancel()">取消</a-button>
      <a-button key="submit" type="primary" :loading="onSubmitLoading" @click="onFinish">提交</a-button>
    </template>

    <a-form :labelCol="{ span: 4 }" :wrapper-col="{span:20}">
      <a-form-item label="脚本语言" v-bind="validateInfos.lang">
        {{languageMap[modelRef.lang]}}
      </a-form-item>
      <a-form-item label="解析器路径" v-bind="validateInfos.val">
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

const useForm = Form.useForm;

interface UpdateFormSetupData {
  modelRef: Ref;
  validateInfos: validateInfos;
  onFinish: () => Promise<void>;
}

export default defineComponent({
  name: 'UpdateForm',
  props: {
    visible: {
      type: Boolean,
      required: true
    },
    values: {
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
  setup(props): UpdateFormSetupData {
    const {t} = useI18n();

    let modelRef = reactive<any>({
      lang: props.values.value.lang || '',
      val: props.values.value.val || '',
    });

    const rulesRef = reactive({
      lang: [{required: true, message: '请输入语言'}],
      val: [{required: true, message: '请输入解析器可执行文件路径'}],
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
      modelRef,
      validateInfos,
      onFinish
    }
  }
})
</script>