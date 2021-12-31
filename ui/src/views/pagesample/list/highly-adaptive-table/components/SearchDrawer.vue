<template>
    <a-drawer
      placement="right"
      :width="360"
      :title="title"
      :closable="false"
      :onClose="onClose"
      :visible="visible"
    >
        
        <a-form layout="vertical"  :wrapper-col="{span:24}">
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

        <div :style="{
            position: 'absolute',
            bottom: 0,
            width: '100%',
            borderTop: '1px solid #e8e8e8',
            padding: '10px 16px',
            textAlign: 'right',
            left: 0,
            background: '#fff',
            borderRadius: '0 0 4px 4px',
        }">
            <div class="text-align-right">
                <a-button style="margin-right: 8px" @click="onClose">
                    取消
                </a-button>
                <a-button type="primary" @click="onSearch">
                    搜索
                </a-button>
            </div>
        </div>

      
    </a-drawer>
</template>
<script lang="ts">
import { defineComponent, PropType, reactive } from "vue";

import { Form } from 'ant-design-vue';
const useForm = Form.useForm;

import TypeSelect from './TypeSelect.vue';
import { TableListItem } from "../data.d";

export default defineComponent({
    name: 'SearchDrawer',
    props: {
        visible: {
            type: Boolean,
            required: true
        },
        title: {
            type: String,
            default: '高级搜索'
        },
        onClose: {
            type: Function,
            required: true
        },
        onSubmit: {
            type: Function as PropType<(values: Omit<TableListItem, 'id'>) => void>,
            required: true
        }
    },
    components: {
        TypeSelect
    },
    setup(props) {

        // 表单值
        const modelRef = reactive<Omit<TableListItem, 'id'>>({
            name: '',
            desc: '',
            href: '',
            type: ''
        });
        // 表单验证
        const rulesRef = reactive({
            name: [],
            desc: [], 
            href: [],
            type: []         
        });
        // 获取表单内容
        const { resetFields, validate, validateInfos } = useForm(modelRef, rulesRef);

        const onSearch = async () => {
             try {
                const fieldsValue = await validate<Omit<TableListItem, 'id'>>();
                props.onSubmit(fieldsValue);
            } catch (error) {
                // console.log('error', error);
            }
        }

        return {
            modelRef,
            validateInfos,
            resetFields,
            onSearch
        }


    }

})
</script>