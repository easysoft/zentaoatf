<template>
    <a-modal
      :destroy-on-close="true"
      :mask-closable="false"
      title="新建脚本"
      :visible="visible"
      :onCancel="onCancel"
    >
        <template #footer>
            <a-button key="back" @click="() => onCancel()">取消</a-button>
            <a-button key="submit" type="primary" :loading="onSubmitLoading" @click="onFinish">提交</a-button>
        </template>
        
        <a-form :labelCol="{ span: 4 }" :wrapper-col="{span:20}">
            <a-form-item label="名称" v-bind="validateInfos.name">
                <a-input v-model:value="modelRef.name" placeholder="" />
            </a-form-item>
            <a-form-item label="描述" v-bind="validateInfos.desc">
                <a-input v-model:value="modelRef.desc" placeholder="" />
            </a-form-item>
        </a-form>


    </a-modal>
</template>
<script lang="ts">
import { defineComponent, PropType, reactive } from "vue";
import { useI18n } from "vue-i18n";

import { Props, validateInfos } from 'ant-design-vue/lib/form/useForm';
import { message, Form } from 'ant-design-vue';
const useForm = Form.useForm;

import { Execution } from '../../data.d';

interface CreateFormSetupData {
    modelRef: Omit<Execution, 'id'>;
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
        onCancel: {
            type: Function,
            required: true
        },
        onSubmitLoading: {
            type: Boolean,
            required: true
        },
        onSubmit: {
            type: Function as PropType<(values: Omit<Execution, 'id'>, resetFields: (newValues?: Props | undefined) => void) => void>,
            required: true
        }
    },
    components: {
    },
    setup(props): CreateFormSetupData {

        const { t } = useI18n();

        // 表单值
        const modelRef = reactive<Omit<Execution, 'id'>>({
            name: '',
            desc: '',
        });
        // 表单验证
        const rulesRef = reactive({
            name: [
                {
                    required: true,
                    validator: async (rule: any, value: string) => {
                        if (value === '' || !value) {
                            throw new Error('请输入名称');
                        } else if (value.length > 15) {
                            throw new Error('长度不能大于15个字');
                        }
                    }
                },
            ],
            desc: [],
        });
        // 获取表单内容
        const { resetFields, validate, validateInfos } = useForm(modelRef, rulesRef);
        // 提交
        const onFinish = async () => {           
            try {
                const fieldsValue = await validate<Omit<Execution, 'id'>>();
                props.onSubmit(fieldsValue, resetFields);
            } catch (error) {
                // console.log('error', error);
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