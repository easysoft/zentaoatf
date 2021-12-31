<template>
    <a-modal
      :destroy-on-close="true"
      :mask-closable="false"
      title="编辑"
      :visible="visible"
      :onCancel="onCancel"
    >
        <template #footer>
            <a-button key="back" @click="() => onCancel()">取消</a-button>
            <a-button key="submit" type="primary" :loading="onSubmitLoading" @click="onFinish">提交</a-button>
        </template>
        
        <a-form :labelCol="{ span: 4 }" :wrapper-col="{span:20}">
            <a-form-item label="位置" v-bind="validateInfos.type">
                <TypeSelect v-model:value="modelRef.type" placeholder="请选择" />
            </a-form-item>
            <a-form-item label="名称" v-bind="validateInfos.name">
                <a-input v-model:value="modelRef.name" placeholder="请输入名称" />
            </a-form-item>
            <a-form-item label="网址" v-bind="validateInfos.href">
                <a-input v-model:value="modelRef.href" placeholder="请输入网址" />
            </a-form-item>

            <a-form-item label="备注" v-bind="validateInfos.desc">
                <a-input v-model:value="modelRef.desc" placeholder="请输入备注" />
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

import TypeSelect from './TypeSelect.vue';
import { TableListItem } from "../data.d";

interface UpdateFormSetupData {
    modelRef: Omit<TableListItem, 'id'>;
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
            type: Object as PropType<Partial<TableListItem>>,
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
            type: Function as PropType<(values: TableListItem, resetFields: (newValues?: Props | undefined) => void) => void>,
            required: true
        }
    },
    components: {
        TypeSelect
    },
    setup(props): UpdateFormSetupData {

        const { t } = useI18n();

        // 表单值
        const modelRef = reactive<TableListItem>({
            id: props.values.id || 0,
            name: props.values.name || '',
            desc: props.values.desc || '',
            href: props.values.href || '',
            type: props.values.type || ''
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
            href: [
                {
                    required: true,
                    validator: async (rule: any, value: string) => {
                        if (value === '' || !value) {
                            throw new Error('请输入网址');
                        } else if (!/^(https?:)/.test(value)) {
                            throw new Error('请输入正确的网址');
                        }
                    },
                },
            ],
            type: [
                {
                    required: true,
                    message: '请选择'
                }
            ]         
        });
        // 获取表单内容
        const { resetFields, validate, validateInfos } = useForm(modelRef, rulesRef);
        // 提交
        const onFinish = async () => {           
            try {
                const fieldsValue = await validate<TableListItem>();
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