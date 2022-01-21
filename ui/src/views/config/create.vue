<template>
    <a-modal
      :destroy-on-close="true"
      :mask-closable="false"
      title="新建解析器"
      :visible="visible"
      :onCancel="onCancel"
      width="600px"
    >
        <template #footer>
          <a-button key="submit" type="primary" :loading="onSubmitLoading" @click="onFinish">提交</a-button>
          <a-button key="back" @click="onCancel">取消</a-button>
        </template>
        
        <a-form :labelCol="{ span: 4 }" :wrapper-col="{span:20}">
            <a-form-item label="脚本语言" v-bind="validateInfos.lang">
              <a-select v-model:value="modelRef.lang">
                <a-select-option key="" value="">&nbsp;</a-select-option>
                <a-select-option v-for="item in languages" :key="item" :value="item">{{languageMap[item]}}</a-select-option>
              </a-select>
            </a-form-item>
            <a-form-item label="解析器路径" v-bind="validateInfos.value">
                <a-input v-model:value="modelRef.val" placeholder="" />
            </a-form-item>
        </a-form>

    </a-modal>
</template>

<script lang="ts">
import { defineComponent, PropType, reactive, Ref } from "vue";
import { useI18n } from "vue-i18n";

import { Props, validateInfos } from 'ant-design-vue/lib/form/useForm';
import { message, Form } from 'ant-design-vue';
const useForm = Form.useForm;

interface CreateInterpreterFormSetupData {
  modelRef: Ref
  validateInfos: validateInfos;
  onFinish: () => Promise<void>;
}

export default defineComponent({
    name: 'CreateForm',
    props: {
        visible: {
            type: Boolean,
            required: true
        },
      languages: {
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
    components: {
    },
    setup(props): CreateInterpreterFormSetupData {
        const { t } = useI18n();

        const modelRef = reactive({lang: '', value: ''} as any)
        const rulesRef = reactive({
          lang: [ { required: true, message: '请输入语言' } ],
          val: [ { required: true, message: '请输入解析器可执行文件路径' } ],
        });

        const { resetFields, validate, validateInfos } = useForm(modelRef, rulesRef);

        const onFinish = async () => {
          validate()
            .then(() => {
              props.onSubmit(modelRef, resetFields);
            })
        }
        
        return {
          modelRef,
          validateInfos,
          onFinish
        }

    }
})
</script>