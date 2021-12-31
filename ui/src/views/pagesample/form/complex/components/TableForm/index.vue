<template>
        <a-table
            :pagination="false"
            :loading="TableLoading"
            :dataSource="TableData"
            :columns="columns"
        >
            <template #name="{ text, record }">
                <a-input v-if="record.edit" :value="text" @change="e => handleChange(e.target.value, record.key, 'name')"  placeholder="姓名" />
                <span v-else>{{text}}</span>
            </template>
            <template #workid="{ text, record }">
                <a-input v-if="record.edit" :value="text" @change="e => handleChange(e.target.value, record.key, 'workId')" placeholder="工号" />
                <span v-else>{{text}}</span>
            </template>
            <template #action="{ record }">
                <template v-if="record.edit">
                    <span v-if="record.isNew">
                        <a @click="saveRow(record)">添加</a>
                        <a-divider type="vertical" />
                        <a-popconfirm title="是否要删除此行？" @confirm="remove(record.key)">
                            <a>删除</a>
                        </a-popconfirm>
                    </span>
                    <span v-else>
                        <a @click="saveRow(record)">保存</a>
                        <a-divider type="vertical" />
                        <a @click="cancel(record.key)">取消</a>
                    </span>
                </template>
                <span v-else>
                    <a  @click="toggle(record.key)">编辑</a>
                    <a-divider type="vertical" />
                    <a-popconfirm title="是否要删除此行？" @confirm="remove(record.key)">
                        <a>删除</a>
                    </a-popconfirm>
                </span>
            </template>
        </a-table>
        <a-button
            style="width: 100%;margin-top: 16px;margin-bottom: 8px; "
            type="dashed"
            @click="newTableData"
        >
            <PlusOutlined />
            新增内容
      </a-button>
</template>
<script lang="ts">
import { defineComponent, PropType, ref, toRefs, watch } from "vue";
import { message } from "ant-design-vue";
import { PlusOutlined } from '@ant-design/icons-vue';
import { TableFormDataType } from "./data.d";

interface TableFormSetupData {
    columns: any;
    TableData: TableFormDataType;
    TableLoading: boolean;
    newTableData:  () => void;
    saveRow: (record: TableFormDataType) => void;
    remove: (key: string) => void;
    cancel: (key: string) => void;
    toggle: (key: string) => void;
    handleChange: (value: string, key: string, column: 'name' | 'workId') => void;
}

export default defineComponent({
    name: 'TableForm',
    props: {
        value: {
            type: Array as PropType<TableFormDataType[]>,
            required: true
        }
    },
    components: {
        PlusOutlined
    },
    setup(props, { emit }): TableFormSetupData {

        const { value } = toRefs(props);        

        const columns = [
            {
                title: '姓名',
                dataIndex: 'name',
                key: 'name',
                width: '35%',
                slots: { customRender: 'name' },
            },
            {
                title: '工号',
                dataIndex: 'workId',
                key: 'workId',
                width: '35%',
                slots: { customRender: 'workid' },
            },
            {
                title: '操作',
                key: 'action',
                slots: { customRender: 'action' },
            },
        ];
        const TableData = ref<TableFormDataType[]>(props.value);
        const TableLoading = ref<boolean>(false);

        // 新增内容
        const newIndex = ref<number>(0);
        const newTableData = () => {

            const newData = TableData.value.map(item => ({ ...item }));

            newData.push({
                key: `NEW_TEMP_ID_${newIndex.value}`,
                workId: '',
                name: '',
                edit: true,
                isNew: true,
            });

            newIndex.value ++;
            TableData.value = newData;
        }

        // 添加、保存
        const saveRow = (record: TableFormDataType) => {
            TableLoading.value = true;
            const { key, name, workId } = record
            if (!name || !workId) {
                TableLoading.value = false;
                message.error('请填写完整成员信息。')
                return
            }

            const target: any = TableData.value.find(item => item.key === key);
            if (target) {
                target.edit = false;
                target.isNew = false;
                target._originalData = undefined;
            }
            TableLoading.value = false;

            const newData = TableData.value.map(item => ({ ...item }));

            emit('update:value', newData);
        }

        // 删除
        const remove = (key: string) => {
            const newData = TableData.value.filter(item => item.key !== key);
            TableData.value = newData;
            emit('update:value', newData);
        }

        // 取消编辑
        const cancel = (key: string) => {
            const target: any = TableData.value.find(item => item.key === key);
            if(target) {
                Object.keys(target).forEach(key => { target[key] = target._originalData[key] });
                target._originalData = undefined;
            }
        }

        // 编辑显示
        const toggle = (key: string) => {
            const target: any = TableData.value.find(item => item.key === key);
            target._originalData = { ...target };
            target.edit = !target.edit;
        }

        // 输入框修改内容
        const handleChange = (value: string, key: string, column: 'name' | 'workId') => {
            const newData = [...TableData.value];
            const target = newData.find(item => key === item.key)
            if (target) {
                target[column] = value;
                TableData.value = newData;
            }
        }

        watch(value,()=> {
            const newData = value.value.map(item => ({ ...item }));
            TableData.value = newData;
        })

        return {
            columns,
            TableData: TableData as unknown as TableFormDataType,
            TableLoading: TableLoading as unknown as boolean,
            newTableData,
            saveRow,
            remove,
            cancel,
            toggle,
            handleChange
        }

    }

})
</script>