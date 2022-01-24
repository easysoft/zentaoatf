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
                  row-key="seq"
                  :columns="columns"
                  :data-source="list"
                  :loading="loading"
                  :pagination="false"
              >
                <template #seq="{ text }">
                  {{text}}
                </template>
                <template #startTime="{ record }">
                  {{ momentTime(record.startTime) }}
                </template>
                <template #duration="{ record }">
                  {{record.duration}}秒
                </template>
                <template #result="{ record }">
                  合计{{record.total}}：
                  <span class="t-pass">{{record.pass}}（{{percent(record.pass, record.total)}}）通过</span>，
                  <span class="t-fail">{{record.fail}}（{{percent(record.fail, record.total)}}）失败</span>，
                  <span class="t-skip">{{record.skip}}（{{percent(record.skip, record.total)}}）忽略</span>。
                </template>
                <template #action="{ record }">
                  <a-button @click="() => viewExec(record)" type="link" size="small">查看</a-button>
                  <a-button @click="() => deleteExec(record)" type="link" size="small"
                            :loading="deleteLoading.includes(record.seq)">删除</a-button>
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
import moment from 'moment';

interface ListExecSetupData {
  columns: any;
  list: ComputedRef<Execution[]>;
  loading: Ref<boolean>;
  getList:  () => Promise<void>;
  viewExec: (item) => void;

  deleteLoading: Ref<string[]>;
  deleteExec:  (item) => void;

  execCase:  () => void;
  execModule:  () => void;
  execSuite:  () => void;
  execTask:  () => void;
  momentTime: (tm) => string;
  percent: (numb, total) => string;
}

export default defineComponent({
    name: 'ExecListPage',
    components: {
    },
    setup(): ListExecSetupData {
      const momentTime = (tm) => {
        return moment.unix(tm).format("YYYY-MM-DD HH:mm:ss")
      }
      const percent = (numb, total) => {
        return Number(numb / total * 100).toFixed(2) + '%'
      }

      const columns =[
        {
          title: '序号',
          dataIndex: 'index',
          width: 150,
          customRender: ({text, index}: { text: any; index: number}) => index + 1,
        },
        {
          title: '名称',
          dataIndex: 'seq',
        },
        {
          title: '开始时间',
          dataIndex: 'startTime',
          slots: { customRender: 'startTime' },
        },
        {
          title: '耗时',
          dataIndex: 'duration',
          slots: { customRender: 'duration' },
        },
        {
          title: '结果',
          dataIndex: 'result',
          slots: { customRender: 'result' },
        },
        {
          title: '操作',
          key: 'action',
          width: 260,
          slots: { customRender: 'action' },
        },
      ];

      const router = useRouter();
      const store = useStore<{ ListExecution: ListStateType}>();

      const list = computed<any[]>(() => store.state.ListExecution.listResult);
      const loading = ref<boolean>(true);
      const getList = async (): Promise<void> => {
        loading.value = true;
        await store.dispatch('ListExecution/listExecution', {});
        loading.value = false;
      }

      // 查看
      const viewExec = (item) => {
        router.push(`/exec/history/${item.seq}`)
      }

      // 删除
      const deleteLoading = ref<string[]>([]);
      const deleteExec = (item) => {
        Modal.confirm({
          title: '删除脚本',
          content: '确定删除吗？',
          okText: '确认',
          cancelText: '取消',
          onOk: async () => {
            deleteLoading.value = [item.seq];
            const res: boolean = await store.dispatch('ListExecution/deleteExecution', item.seq);
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
        momentTime,
        percent,
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