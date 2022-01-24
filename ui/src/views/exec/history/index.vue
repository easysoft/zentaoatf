<template>
    <div class="indexlayout-main-conent">
        <a-card :bordered="false">
            <template #title>
                测试执行
            </template>
            <template #extra>
              <div class="opt">
                <a-button @click="execCase" type="primary">执行用例</a-button>
                <a-button @click="execModule" type="primary">执行模块</a-button>
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
                            @click="() => viewExec(record.id)"
                            :loading="getLoading.includes(record.id)">查看</a-button>
                  <a-button type="link"
                            @click="() => editExec(record.id)"
                            :loading="getLoading.includes(record.id)">编辑</a-button>
                  <a-button type="link"
                            @click="() => deleteExec(record.id)"
                            :loading="deleteLoading.includes(record.id)">删除</a-button>
                </template>

              </a-table>
            </div>
        </a-card>
    </div>
</template>

<script lang="ts">
import {ComputedRef, defineComponent, ref, Ref, reactive, computed, onMounted} from "vue";
import { SelectTypes } from 'ant-design-vue/es/select';
import {Execution} from '../data.d';
import {QueryParams, PaginationConfig} from '@/types/data.d';
import {useStore} from "vuex";

import { Props } from 'ant-design-vue/lib/form/useForm';
import { message, Modal, Form } from "ant-design-vue";
const useForm = Form.useForm;

import {StateType as ListStateType} from "../store";
import debounce from "lodash.debounce";
import {useRouter} from "vue-router";

interface ListExecPageSetupData {
  columns: any;
  list: ComputedRef<Execution[]>;
  loading: Ref<boolean>;
  getList:  () => Promise<void>;
  viewExec: (id: number) => void;

  deleteLoading: Ref<number[]>;
  deleteExec:  (id: number) => void;

  execCase:  () => void;
  execModule:  () => void;
  execSuite:  () => void;
  execTask:  () => void;
}

export default defineComponent({
    name: 'ExecListPage',
    components: {
    },
    setup(): ListExecPageSetupData {
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

      const list = computed<Execution[]>(() => store.state.ListExecution.queryResult.data);

      const columns =[
        {
          title: '序号',
          dataIndex: 'index',
          width: 80,
          customRender: ({text, index}: { text: any; index: number}) => index + 1,
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
      const getList = async (): Promise<void> => {
        loading.value = true;

        await store.dispatch('ListExecution/queryExecution', {});
        loading.value = false;
      }

      // 查看
      const viewExec = (id: number) => {
        router.push(`/~/execution/result/${id}`)
      }

      // 删除
      const deleteLoading = ref<number[]>([]);
      const deleteExec = (id: number) => {
        Modal.confirm({
          title: '删除脚本',
          content: '确定删除吗？',
          okText: '确认',
          cancelText: '取消',
          onOk: async () => {
            deleteLoading.value = [id];
            const res: boolean = await store.dispatch('ListExecution/deleteExec',id);
            if (res === true) {
              message.success('删除成功！');
              await getList();
            }
            deleteLoading.value = [];
          }
        });
      }

      onMounted(()=> {
        getList();
      })

      const execCase = () =>  {
        console.log("execCase")
        router.push(`/exec/exec/case`)
      }
      const execModule = () =>  {
        console.log("execModule")
        router.push(`/exec/exec/module`)
      }
      const execSuite = () =>  {
        console.log("execSuite")
        router.push(`/exec/exec/suite`)
      }
      const execTask = () =>  {
        console.log("execSuite")
        router.push(`/exec/exec/task`)
      }

      return {
        columns,
        list,
        loading,
        getList,

        viewExec,
        deleteLoading,
        deleteExec,

        execCase,
        execModule,
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