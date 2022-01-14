<template>
    <div class="indexlayout-main-conent">
        <a-card :bordered="false">
            <template #title>
                测试执行
            </template>
            <template #extra>
              <div class="opt">
                <a-select @change="onSearch" v-model:value="queryParams.enabled" :options="statusArr">
                </a-select>
                <a-input-search @change="onSearch" @search="onSearch" v-model:value="queryParams.keywords"
                                placeholder="输入关键字搜索" style="width:270px;margin-left: 16px;" />

                <span class="space"></span>

                <a-button @click="execCase" type="primary">执行用例</a-button>
                <a-button @click="execSuite" type="primary">执行套件</a-button>
                <a-button @click="execTask" type="primary">执行任务</a-button>
              </div>
            </template>

            <div>
              <a-table
                  row-key="id"
                  :columns="columns"
                  :data-source="list"
                  :loading="loading"
                  :pagination="{
                    ...pagination,
                    onChange: (page) => {
                        getList(page);
                    },
                    onShowSizeChange: (page, size) => {
                        pagination.pageSize = size
                        getList(page);
                    },
                }"
              >
                <template #name="{ text, record  }">
                  <a :href="record.href" target="_blank">{{text}}</a>
                </template>
                <template #status="{ record }">
                  <a-tag v-if="record.disabled == 0" color="green">启用</a-tag>
                  <a-tag v-else color="cyan">禁用</a-tag>
                </template>
                <template #action="{ record }">
                  <a-button type="link"
                            @click="() => designExecution(record.id)"
                            :loading="getLoading.includes(record.id)">设计</a-button>
                  <a-button type="link"
                            @click="() => editExecution(record.id)"
                            :loading="getLoading.includes(record.id)">编辑</a-button>
                  <a-button type="link"
                            @click="() => deleteExecution(record.id)"
                            :loading="deleteLoading.includes(record.id)">删除</a-button>
                </template>

              </a-table>

              <create-form
                  :visible="createFormVisible"
                  :onCancel="() => setCreateFormVisible(false)"
                  :onSubmitLoading="createSubmitLoading"
                  :onSubmit="createSubmit"
              />

              <update-form
                  v-if="updateFormVisible===true"
                  :visible="updateFormVisible"
                  :values="item"
                  :onCancel="() => updateFormCancel()"
                  :onSubmitLoading="updateSubmitLoading"
                  :onSubmit="updateSubmit"
              />
            </div>
        </a-card>
    </div>
</template>

<script lang="ts">
import {ComputedRef, defineComponent, ref, Ref, reactive, computed, onMounted} from "vue";
import { SelectTypes } from 'ant-design-vue/es/select';
import {Execution, QueryParams, PaginationConfig} from '../data.d';
import {useStore} from "vuex";

import { Props } from 'ant-design-vue/lib/form/useForm';
import { message, Modal, Form } from "ant-design-vue";
const useForm = Form.useForm;

import CreateForm from './components/CreateForm.vue';
import UpdateForm from './components/UpdateForm.vue';

import {StateType as ListStateType} from "../store";
import debounce from "lodash.debounce";
import {useRoute, useRouter} from "vue-router";

interface ListExecutionPageSetupData {
  statusArr,
  queryParams,
  columns: any;
  list: ComputedRef<Execution[]>;
  pagination: ComputedRef<PaginationConfig>;
  loading: Ref<boolean>;
  getList:  (current: number) => Promise<void>;
  createFormVisible: Ref<boolean>;
  setCreateFormVisible:  (val: boolean) => void;
  createSubmitLoading: Ref<boolean>;
  createSubmit: (values: Omit<Execution, 'id'>, resetFields: (newValues?: Props | undefined) => void) => Promise<void>;

  getLoading: Ref<number[]>;
  designExecution: (id: number) => void;
  editExecution: (id: number) => Promise<void>;
  item: ComputedRef<Partial<Execution>>;
  updateFormVisible: Ref<boolean>;
  updateFormCancel:  () => void;
  updateSubmitLoading: Ref<boolean>;
  updateSubmit:  (values: Execution, resetFields: (newValues?: Props | undefined) => void) => Promise<void>;

  deleteLoading: Ref<number[]>;
  deleteExecution:  (id: number) => void;

  onSearch:  () => void;
  execCase:  () => void;
  execSuite:  () => void;
  execTask:  () => void;
}

export default defineComponent({
    name: 'ExecutionListPage',
    components: {
      CreateForm,
      UpdateForm
    },
    setup(): ListExecutionPageSetupData {
      const statusArr = ref<SelectTypes['options']>([
          {
            label: '所有状态',
            value: '',
          },
          {
            label: '启用',
            value: '1',
          },
          {
            label: '禁用',
            value: '0',
          },
        ]);

      const router = useRouter();
      const store = useStore<{ ListExecution: ListStateType}>();

      const list = computed<Execution[]>(() => store.state.ListExecution.queryResult.list);
      let pagination = computed<PaginationConfig>(() => store.state.ListExecution.queryResult.pagination);
      let queryParams = reactive<QueryParams>({keywords: '', enabled: '',
        page: pagination.value.current, pageSize: pagination.value.pageSize});

      const columns =[
        {
          title: '序号',
          dataIndex: 'index',
          width: 80,
          customRender: ({text, index}: { text: any; index: number}) => (pagination.value.current - 1) * pagination.value.pageSize + index + 1,
        },
        {
          title: '名称',
          dataIndex: 'name',
          slots: { customRender: 'name' },
        },
        {
          title: '描述',
          dataIndex: 'desc',
        },
        {
          title: '状态',
          dataIndex: 'status',
          slots: { customRender: 'status' },
        },
        {
          title: '操作',
          key: 'action',
          width: 260,
          slots: { customRender: 'action' },
        },
      ];

      const loading = ref<boolean>(true);
      const getList = async (current: number): Promise<void> => {
        loading.value = true;

        await store.dispatch('ListExecution/queryExecution', {
          keywords: queryParams.keywords,
          enabled: queryParams.enabled,
          pageSize: pagination.value.pageSize,
          page: current,
        });
        loading.value = false;
      }

      // 创建弹框 - visible
      const createFormVisible = ref<boolean>(false);
      const setCreateFormVisible = (val: boolean) => {
        createFormVisible.value = val;
      };
      // 创建弹框 - 提交 loading
      const createSubmitLoading = ref<boolean>(false);
      // 创建弹框 - 提交
      const createSubmit = async (values: Omit<Execution, 'id'>, resetFields: (newValues?: Props | undefined) => void) => {
        createSubmitLoading.value = true;
        const res: boolean = await store.dispatch('ListExecution/createExecution',values);
        if(res === true) {
          resetFields();
          setCreateFormVisible(false);
          message.success('新增成功！');
          getList(1);
        }
        createSubmitLoading.value = false;
      }

      // 更新弹框 - visible
      const updateFormVisible = ref<boolean>(false);
      const setUpdateFormVisible = (val: boolean) => {
        updateFormVisible.value = val;
      }
      const updateFormCancel = () => {
        setUpdateFormVisible(false);
        store.commit('ListExecution/setItem',{});
      }
      // 更新弹框 - 提交 loading
      const updateSubmitLoading = ref<boolean>(false);
      // 更新弹框 - 提交
      const updateSubmit = async (values: Execution, resetFields: (newValues?: Props | undefined) => void) => {
        updateSubmitLoading.value = true;
        const res: boolean = await store.dispatch('ListExecution/updateExecution',values);
        if(res === true) {
          updateFormCancel();
          message.success('编辑成功！');
          getList(pagination.value.current);
        }
        updateSubmitLoading.value = false;
      }

      const item = computed<Partial<Execution>>(() => store.state.ListExecution.detailResult);
      // 编辑
      const getLoading = ref<number[]>([]);
      const editExecution = async (id: number) => {
        getLoading.value = [id];
        const res: boolean = await store.dispatch('ListExecution/getExecution',id);
        if(res===true) {
          setUpdateFormVisible(true);
        }
        getLoading.value = [];
      }

      // 设计
      const designExecution = (id: number) => {
        router.push(`/~/execution/design/${id}`)
      }

      // 删除
      const deleteLoading = ref<number[]>([]);
      const deleteExecution = (id: number) => {
        Modal.confirm({
          title: '删除脚本',
          content: '确定删除吗？',
          okText: '确认',
          cancelText: '取消',
          onOk: async () => {
            deleteLoading.value = [id];
            const res: boolean = await store.dispatch('ListExecution/deleteExecution',id);
            if (res === true) {
              message.success('删除成功！');
              await getList(pagination.value.current);
            }
            deleteLoading.value = [];
          }
        });
      }

      onMounted(()=> {
        getList(1);
      })

      const onSearch = debounce(() =>  {
        getList(1)
      }, 500);

      const execCase = () =>  {
        console.log("execCase")
        router.push(`/execution/exec/case`)
      }
      const execSuite = () =>  {
        console.log("execSuite")
      }
      const execTask = () =>  {
        console.log("execSuite")
      }

      return {
        statusArr,
        queryParams,
        columns,
        list,
        pagination,
        loading,
        getList,

        createFormVisible,
        setCreateFormVisible,
        createSubmitLoading,
        createSubmit,
        getLoading,
        designExecution,
        editExecution,
        item,
        updateFormVisible,
        updateFormCancel,
        updateSubmitLoading,
        updateSubmit,
        deleteLoading,
        deleteExecution,

        onSearch,
        execCase,
        execSuite,
        execTask,
      }
    }

})
</script>

<style lang="less" scoped>
  .opt {
    .space {
      display: inline-block;
      width: 50px;
    }
    .ant-btn {
      margin-left: 12px;
    }
  }
</style>