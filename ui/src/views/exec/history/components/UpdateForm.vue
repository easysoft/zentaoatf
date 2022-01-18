<template>
    <a-modal
      :destroy-on-close="true"
      :mask-closable="false"
      title="编辑脚本"
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

interface UpdateFormSetupData {
    modelRef: Omit<Execution, 'id'>;
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
            type: Object as PropType<Partial<Execution>>,
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
            type: Function as PropType<(values: Execution, resetFields: (newValues?: Props | undefined) => void) => void>,
            required: true
        }
    },
    components: {
    },
    setup(props): UpdateFormSetupData {
        const { t } = useI18n();

        // 表单值
        const modelRef = reactive<Execution>({
            id: props.values.id || 0,
            name: props.values.name || '',
            desc: props.values.desc || '',
        });
        // 表单验证
        const rulesRef = reactive({
            id: [],
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
                const fieldsValue = await validate<Execution>();
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