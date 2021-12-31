<template>
    <div class="indexlayout-main-conent">
        <a-card :bordered="false" title="友情链接">
            <template #extra>
                    <a-radio-group defaultValue="all">
                        <a-radio-button value="all">全部</a-radio-button>
                        <a-radio-button value="header">头部</a-radio-button>
                        <a-radio-button value="footer">底部</a-radio-button>
                    </a-radio-group>
                    <a-input-search placeholder="请输入"  style="width:270px;margin-left: 16px;" />
            </template>

            <a-button type="dashed" style="width:100%; margin-bottom: 8px;" @click="() => setCreateFormVisible(true)">
                <PlusOutlined /> 新增
            </a-button>

            <a-list
                item-layout="horizontal"
                row-key="id"
                :loading="loading"
                :pagination="{
                    ...pagination,
                    onChange: (page) => {
                        getList(page);
                    }
                }"
                :data-source="list"
            >
                <template #renderItem="{ item }">
                    <a-list-item>
                        <template #actions>
                            <a-button type="link" @click="() => detailUpdateData(item.id)" :loading="detailUpdateLoading.includes(item.id)">编辑</a-button>
                            <a-button type="link" @click="() => deleteTableData(item.id)" :loading="deleteLoading.includes(item.id)">删除</a-button>
                        </template>
                        <a-list-item-meta>
                            <template #title>
                                 <a :href="item.href" target="_blank">{{item.name}}</a>
                            </template>
                            <template #description>
                                {{item.desc}}
                            </template>
                        </a-list-item-meta>
                        <div>
                            <a-tag v-if="item.type === 'header'" color="green">头部</a-tag>
                            <a-tag v-else color="cyan">底部</a-tag>
                        </div>
                    </a-list-item>

                </template>
            </a-list>

            <create-form 
                :visible="createFormVisible" 
                :onCancel="() => setCreateFormVisible(false)" 
                :onSubmitLoading="createSubmitLoading" 
                :onSubmit="createSubmit"
            />

            <update-form
                v-if="updateFormVisible===true"
                :visible="updateFormVisible"
                :values="updateData"
                :onCancel="() => updataFormCancel()"
                :onSubmitLoading="updateSubmitLoading"
                :onSubmit="updateSubmit"
            />


        </a-card>
    </div>
</template>
<script lang="ts">
import { computed, ComputedRef, defineComponent, onMounted, Ref, ref } from "vue";
import { useStore } from "vuex";
import { message, Modal } from "ant-design-vue";
import { PlusOutlined } from '@ant-design/icons-vue';
import { Props } from 'ant-design-vue/lib/form/useForm';
import CreateForm from './components/CreateForm.vue';
import UpdateForm from './components/UpdateForm.vue';
import { StateType as ListStateType } from "./store";
import { PaginationConfig, TableListItem } from './data.d';

interface ListBasicPageSetupData {
    list:  ComputedRef<TableListItem[]>;
    pagination: ComputedRef<PaginationConfig>;
    loading: Ref<boolean>;
    getList:  (current: number) => Promise<void>;
    createFormVisible: Ref<boolean>;
    setCreateFormVisible:  (val: boolean) => void;
    createSubmitLoading: Ref<boolean>;
    createSubmit: (values: Omit<TableListItem, 'id'>, resetFields: (newValues?: Props | undefined) => void) => Promise<void>;
    detailUpdateLoading: Ref<number[]>;
    detailUpdateData: (id: number) => Promise<void>;
    updateData: ComputedRef<Partial<TableListItem>>;
    updateFormVisible: Ref<boolean>;
    updataFormCancel:  () => void;
    updateSubmitLoading: Ref<boolean>;
    updateSubmit:  (values: TableListItem, resetFields: (newValues?: Props | undefined) => void) => Promise<void>;
    deleteLoading: Ref<number[]>;
    deleteTableData:  (id: number) => void;
}

export default defineComponent({
    name: 'ListBasicPage',
    components: {
        PlusOutlined,
        CreateForm,
        UpdateForm
    },
    setup(): ListBasicPageSetupData {

        const store = useStore<{ ListBasic: ListStateType}>();

        // 列表数据
        const list = computed<TableListItem[]>(() => store.state.ListBasic.tableData.list);

        // 列表分页
        const pagination = computed<PaginationConfig>(() => store.state.ListBasic.tableData.pagination);

        // 获取数据
        const loading = ref<boolean>(true);
        const getList = async (current: number): Promise<void> => {
            loading.value = true;
            await store.dispatch('ListBasic/queryTableData', {
                per: pagination.value.pageSize,
                page: current,
            });
            loading.value = false;
        }


        // 新增弹框 - visible
        const createFormVisible = ref<boolean>(false);
        const setCreateFormVisible = (val: boolean) => {
            createFormVisible.value = val;
        };
        // 新增弹框 - 提交 loading
        const createSubmitLoading = ref<boolean>(false);
        // 新增弹框 - 提交
        const createSubmit = async (values: Omit<TableListItem, 'id'>, resetFields: (newValues?: Props | undefined) => void) => {
            createSubmitLoading.value = true;
            const res: boolean = await store.dispatch('ListBasic/createTableData',values);
            if(res === true) {
                resetFields();
                setCreateFormVisible(false);
                message.success('新增成功！');
                getList(1);
            }
            createSubmitLoading.value = false;
        }


        // 编辑弹框 - visible
        const updateFormVisible = ref<boolean>(false);
        const setUpdateFormVisible = (val: boolean) => {
            updateFormVisible.value = val;
        }
        const updataFormCancel = () => {
            setUpdateFormVisible(false);
            store.commit('ListBasic/setUpdateData',{});
        }
        // 编辑弹框 - 提交 loading
        const updateSubmitLoading = ref<boolean>(false);
        // 编辑弹框 - 提交
        const updateSubmit = async (values: TableListItem, resetFields: (newValues?: Props | undefined) => void) => {
            updateSubmitLoading.value = true;
            const res: boolean = await store.dispatch('ListBasic/updateTableData',values);
            if(res === true) {
                updataFormCancel();                
                message.success('编辑成功！');
                getList(pagination.value.current);
            }
            updateSubmitLoading.value = false;
        }

        // 编辑弹框 data
        const updateData = computed<Partial<TableListItem>>(() => store.state.ListBasic.updateData);
        const detailUpdateLoading = ref<number[]>([]);
        const detailUpdateData = async (id: number) => {
            detailUpdateLoading.value = [id];
            const res: boolean = await store.dispatch('ListBasic/queryUpdateData',id);
            if(res===true) {
                setUpdateFormVisible(true);
            }
            detailUpdateLoading.value = [];
        }

        // 删除 loading
        const deleteLoading = ref<number[]>([]);
        // 删除
        const deleteTableData = (id: number) => {

            Modal.confirm({
                title: '删除',
                content: '确定删除吗？',
                okText: '确认',
                cancelText: '取消',
                onOk: async () => {
                    deleteLoading.value = [id];
                    const res: boolean = await store.dispatch('ListBasic/deleteTableData',id);
                    if (res === true) {
                        message.success('删除成功！');
                        getList(pagination.value.current);
                    }
                   deleteLoading.value = [];
                }
            });

        }




        onMounted(()=> {
           getList(1);
        })

        return {
            list,
            pagination,
            loading,
            getList,
            createFormVisible,
            setCreateFormVisible,
            createSubmitLoading,
            createSubmit,
            detailUpdateLoading,
            detailUpdateData,
            updateData,
            updateFormVisible,
            updataFormCancel,
            updateSubmitLoading,
            updateSubmit,
            deleteLoading,
            deleteTableData
        }

    }
    
})
</script>