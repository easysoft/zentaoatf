<template>
    <div class="indexlayout-main-conent">
        <a-form  :wrapper-col="{span:24}">
            <a-card :bordered="false"  title="基础信息" style="margin-bottom: 20px">
                <a-row :gutter="16">
                    <a-col :lg="8" :md="12" :sm="24">
                        <a-form-item label="标题：" v-bind="validateInfos.title">
                            <a-input v-model:value="modelRef.title" placeholder="请输入" />
                        </a-form-item>
                    </a-col>
                    <a-col :lg="8" :md="12" :sm="24">
                        <a-form-item label="起止日期" v-bind="validateInfos.date">
                            <a-range-picker v-model:value="modelRef.date" style="width:100%" />
                        </a-form-item>
                    </a-col>
                    <a-col :lg="8" :md="12" :sm="24">
                        <a-form-item label="下拉选择" v-bind="validateInfos.select">
                            <a-select v-model:value="modelRef.select"  placeholder="请选择" allowClear>
                                <a-select-option value="1">select1</a-select-option>
                                <a-select-option value="2">select2</a-select-option>
                                <a-select-option value="3">select3</a-select-option>
                            </a-select>
                        </a-form-item>
                    </a-col>
                    <a-col :lg="8" :md="12" :sm="24">
                        <a-form-item label="单选按钮1">
                            <a-radio-group  v-model:value="modelRef.radio1" >
                                <a-radio value="1">item 1</a-radio>
                                <a-radio value="2">item 2</a-radio>
                                <a-radio value="3">item 3</a-radio>
                            </a-radio-group>
                        </a-form-item>
                    </a-col>
                </a-row>                
            </a-card>

            <a-card :bordered="false"  title="拓展信息" style="margin-bottom: 20px">
                <a-row :gutter="16">
                    <a-col :lg="8" :md="12" :sm="24">
                        <a-form-item label="单选按钮2"  v-bind="validateInfos.radio2">
                            <a-radio-group  v-model:value="modelRef.radio2" >
                                <a-radio-button value="1">item 1</a-radio-button>
                                <a-radio-button value="2">item 2</a-radio-button>
                                <a-radio-button value="2">item 3</a-radio-button>
                            </a-radio-group>
                        </a-form-item>
                    </a-col>
                    <a-col :lg="8" :md="12" :sm="24">
                        <a-form-item label="复选框" v-bind="validateInfos.checkbox">
                            <a-checkbox-group v-model:value="modelRef.checkbox">
                                <a-checkbox value="1" name="type">
                                Online
                                </a-checkbox>
                                <a-checkbox value="2" name="type">
                                Promotion
                                </a-checkbox>
                                <a-checkbox value="3" name="type">
                                Offline
                                </a-checkbox>
                            </a-checkbox-group>
                        </a-form-item>
                    </a-col>
                    <a-col :lg="8" :md="12" :sm="24">
                        <a-form-item label="备注" v-bind="validateInfos.remark">
                            <a-textarea v-model:value="modelRef.remark" :rows="1" />
                        </a-form-item>
                    </a-col>
                </a-row>
            </a-card>

            <a-card :bordered="false"  title="表格信息">
                <a-form-item label="" v-bind="validateInfos.users">
                    <TableForm v-model:value="modelRef.users" />
                </a-form-item>
            </a-card>




            <FooterToolbar class="text-align-right">
                <a-button type="primary" @click="handleSubmit" :loading="submitLoading">提交</a-button>
                <a-button @click="resetFields" style="margin-left: 10px;">重置</a-button>  
            </FooterToolbar>


        </a-form>
    </div>
</template>
<script lang="ts">
import { defineComponent, reactive, Ref, ref } from "vue";
import { useStore } from "vuex";

import { Props, validateInfos } from 'ant-design-vue/lib/form/useForm';
import { message, Form } from 'ant-design-vue';
const useForm = Form.useForm;

import { FormDataType } from "./data.d";
import { StateType as FormStateType } from "./store";
import FooterToolbar from '@/layouts/IndexLayout/components/FooterToolbar.vue';
import TableForm from './components/TableForm/index.vue';

interface FormComplexPageSetupData {
    resetFields: (newValues?: Props) => void;
    validateInfos: validateInfos;
    modelRef: FormDataType;
    submitLoading: Ref<boolean>;
    handleSubmit: (e: MouseEvent) => void;
}

export default defineComponent({
    name: 'FormComplexPage',
    components: {
        TableForm,
        FooterToolbar
    },
    setup(): FormComplexPageSetupData {

        const store = useStore<{FormComplex: FormStateType}>();

        // 表单值
        const modelRef = reactive<FormDataType>({
            title: '',
            date: [],
            select: '',
            radio1: '',
            radio2: '',
            checkbox: [],
            remark: '',
            users: []
        });
        // 表单验证
        const rulesRef = reactive({
            title: [
                {
                    required: true,
                    message: '必填',
                },
            ],
            date: [
                {
                    required: true,
                    message: '必填',
                    trigger: 'change', 
                    type: 'array' 
                },
            ],  
            select: [
                {
                    required: true,
                    message: '请选择',
                },
            ],  
            radio1: [],  
            radio2: [
                {
                    required: true,
                    message: '请选择',
                },
            ],
            checkbox:[],
            remark: [],
            users: []      
        });
        // 获取表单内容
        const { resetFields, validate, validateInfos } = useForm(modelRef, rulesRef);
        // 重置 validateInfos 如果用到国际化需要此步骤
        //const validateInfosNew = useI18nAntdFormVaildateInfos(validateInfos);

        // 登录loading
        const submitLoading = ref<boolean>(false);
        // 登录
        const handleSubmit = async (e: MouseEvent) => {
            e.preventDefault();
            submitLoading.value = true;
            try {
                const fieldsValue = await validate<FormDataType>();
                const res: boolean = await store.dispatch('FormComplex/create',fieldsValue);                
                if (res === true) {
                    message.success('提交成功');
                    resetFields();                    
                }
            } catch (error) {
                // console.log('error', error);
            }
            submitLoading.value = false;
        };

        return {
            resetFields,
            validateInfos,
            modelRef,
            submitLoading,
            handleSubmit,
        }



    }
})
</script>