<template>
    <div class="indexlayout-main-conent">
        <a-card :bordered="false">
            <template #title>
                <a-button type="primary" @click="() => setCreateFormVisible(true)">新建脚本</a-button>
            </template>
            <template #extra>
              <a-select @change="onSearch" v-model:value="queryParams.enabled" :options="statusArr" placeholder="状态">
              </a-select>
              <a-input-search @change="onSearch" @search="onSearch" v-model:value="queryParams.keywords"
                              placeholder="请输入" style="width:270px;margin-left: 16px;" />
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
                            @click="() => designScript(record.id)"
                            :loading="getLoading.includes(record.id)">设计</a-button>
                  <a-button type="link"
                            @click="() => editScript(record.id)"
                            :loading="getLoading.includes(record.id)">编辑</a-button>
                  <a-button type="link"
                            @click="() => deleteScript(record.id)"
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
import {Script, QueryParams, PaginationConfig} from '../data.d';
import {useStore} from "vuex";

import { Props } from 'ant-design-vue/lib/form/useForm';
import { message, Modal, Form } from "ant-design-vue";
const useForm = Form.useForm;

import CreateForm from './components/CreateForm.vue';
import UpdateForm from './components/UpdateForm.vue';

import {StateType as ListStateType} from "../store";
import debounce from "lodash.debounce";
import {useRoute, useRouter} from "vue-router";

interface ListScriptPageSetupData {
  statusArr,
  queryParams,
  columns: any;
  list: ComputedRef<Script[]>;
  pagination: ComputedRef<PaginationConfig>;
  loading: Ref<boolean>;
  getList:  (current: number) => Promise<void>;
  createFormVisible: Ref<boolean>;
  setCreateFormVisible:  (val: boolean) => void;
  createSubmitLoading: Ref<boolean>;
  createSubmit: (values: Omit<Script, 'id'>, resetFields: (newValues?: Props | undefined) => void) => Promise<void>;

  getLoading: Ref<number[]>;
  designScript: (id: number) => void;
  editScript: (id: number) => Promise<void>;
  item: ComputedRef<Partial<Script>>;
  updateFormVisible: Ref<boolean>;
  updateFormCancel:  () => void;
  updateSubmitLoading: Ref<boolean>;
  updateSubmit:  (values: Script, resetFields: (newValues?: Props | undefined) => void) => Promise<void>;

  deleteLoading: Ref<number[]>;
  deleteScript:  (id: number) => void;

  onSearch:  () => void;
}

export default defineComponent({
    name: 'ScriptListPage',
    components: {
      CreateForm,
      UpdateForm
    },
    setup(): ListScriptPageSetupData {
      const statusArr = ref<SelectTypes['options']>([
          {
            label: '所有',
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
      const store = useStore<{ ListScript: ListStateType}>();

      const list = computed<Script[]>(() => store.state.ListScript.queryResult.list);
      let pagination = computed<PaginationConfig>(() => store.state.ListScript.queryResult.pagination);
      let queryParams = reactive<QueryParams>({keywords: '', enabled: '1',
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

        await store.dispatch('ListScript/queryScript', {
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
      const createSubmit = async (values: Omit<Script, 'id'>, resetFields: (newValues?: Props | undefined) => void) => {
        createSubmitLoading.value = true;
        const res: boolean = await store.dispatch('ListScript/createScript',values);
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
        store.commit('ListScript/setItem',{});
      }
      // 更新弹框 - 提交 loading
      const updateSubmitLoading = ref<boolean>(false);
      // 更新弹框 - 提交
      const updateSubmit = async (values: Script, resetFields: (newValues?: Props | undefined) => void) => {
        updateSubmitLoading.value = true;
        const res: boolean = await store.dispatch('ListScript/updateScript',values);
        if(res === true) {
          updateFormCancel();
          message.success('编辑成功！');
          getList(pagination.value.current);
        }
        updateSubmitLoading.value = false;
      }

      const item = computed<Partial<Script>>(() => store.state.ListScript.detailResult);
      // 编辑
      const getLoading = ref<number[]>([]);
      const editScript = async (id: number) => {
        getLoading.value = [id];
        const res: boolean = await store.dispatch('ListScript/getScript',id);
        if(res===true) {
          setUpdateFormVisible(true);
        }
        getLoading.value = [];
      }

      // 设计
      const designScript = (id: number) => {
        router.push(`/~/script/design/${id}`)
      }

      // 删除
      const deleteLoading = ref<number[]>([]);
      const deleteScript = (id: number) => {
        Modal.confirm({
          title: '删除脚本',
          content: '确定删除吗？',
          okText: '确认',
          cancelText: '取消',
          onOk: async () => {
            deleteLoading.value = [id];
            const res: boolean = await store.dispatch('ListScript/deleteScript',id);
            if (res === true) {
              message.success('删除成功！');
              await getList(pagination.value.current);
            }
            deleteLoading.value = [];
          }
        });
      }

      // 搜索
      const onSearch = debounce(() =>  {
        getList(1)
      }, 500);

      onMounted(()=> {
        getList(1);
      })

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
        designScript,
        editScript,
        item,
        updateFormVisible,
        updateFormCancel,
        updateSubmitLoading,
        updateSubmit,
        deleteLoading,
        deleteScript,

        onSearch,
      }
    }

})
</script>