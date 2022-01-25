<template>
  <div class="indexlayout-main-conent">
    <a-card :bordered="false">
      <template #title>
       执行结果详情
      </template>
      <template #extra>
        <a-button type="link" @click="() => back()">返回</a-button>
      </template>

      <div class="main">
        <a-row>
          <a-col :span="2" class="t-bord t-label-right">测试环境</a-col>
          <a-col :span="4">{{ testEnv(report.testEnv) }}</a-col>
          <a-col :span="2" class="t-bord t-label-right">测试类型</a-col>
          <a-col :span="4">{{ testType(report.testType) }}</a-col>
        </a-row>
        <a-row>
          <a-col :span="2" class="t-bord t-label-right">开始时间</a-col>
          <a-col :span="4">{{ momentTime(report.startTime) }}</a-col>
          <a-col :span="2" class="t-bord t-label-right">结束时间</a-col>
          <a-col :span="4">{{ momentTime(report.endTime) }}</a-col>
          <a-col :span="2" class="t-bord t-label-right">耗时</a-col>
          <a-col :span="4">{{report.duration}}秒</a-col>
        </a-row>
        <a-row>
          <a-col :span="2" class="t-bord t-label-right">用例数</a-col>
          <a-col :span="4">{{ momentTime(report.startTime) }}</a-col>
          <a-col :span="2" class="t-bord t-label-right">通过</a-col>
          <a-col :span="4">{{report.pass}}（{{percent(report.pass, report.total)}}）</a-col>
          <a-col :span="2" class="t-bord t-label-right">失败</a-col>
          <a-col :span="4">{{report.fail}}（{{percent(report.fail, report.total)}}）</a-col>
          <a-col :span="2" class="t-bord t-label-right">跳过</a-col>
          <a-col :span="4">{{report.skip}}（{{percent(report.skip, report.total)}}）</a-col>
        </a-row>

        <a-row><a-col :span="2" class="t-bord t-label-right">用例详情</a-col></a-row>
        <a-row>
          <a-col :span="2"></a-col>
          <a-col :span="18">
          <template v-for="result in report.funcResult" :key="result.id">
            <div>{{result.id}}. {{result.path}}
              <span :class="'t-'+result.status">{{resultStatus(result.status)}}</span>
            </div>

            <a-table
                :columns="columns"
                :data-source="result.steps"
                row-key="id"
                :pagination="false">
              <template #seq="{ record }">
                {{record.id}}
              </template>
              <template #name="{ record }">
                {{ record.name }}
              </template>
              <template #status="{ record }">
                <span :class="'t-'+record.status">{{ resultStatus(record.status) }}</span>
              </template>
              <template #checkPoints="{ record }">
                <div v-for="item in record.checkPoints" :key="item.numb">
                  {{ item.numb }}.&nbsp;
                  <span :class="'t-'+result.status">{{ resultStatus(item.status) }}</span>
                  &nbsp;&nbsp;&nbsp;
                  "{{ item.expect }}"
                  "{{ item.actual }}"
                </div>
              </template>
            </a-table>
          </template>
          </a-col>
          <a-col :span="2"></a-col>
        </a-row>
      </div>

    </a-card>
  </div>
</template>

<script lang="ts">
import {
  defineComponent,
  onMounted,
  Ref,
  ref,
  computed
} from "vue";
import { useStore } from 'vuex';
import {StateType as ListStateType} from "@/views/exec/store";
import {useRouter} from "vue-router";
import {momentTimeDef, percentDef} from "@/utils/datetime";
import { resultStatusDef, testEnvDef, testTypeDef} from "@/utils/testing";

interface DesignExecutionPageSetupData {
  report: Ref;
  columns: any[]

  loading: Ref<boolean>;
  back: () => void;

  momentTime: (tm) => string;
  percent: (numb, total) => string;
  testEnv: (code) => string;
  testType: (code) => string;
  resultStatus: (code) => string;
}

export default defineComponent({
    name: 'ExecutionResultPage',
    setup(): DesignExecutionPageSetupData {
      const momentTime = momentTimeDef
      const percent = percentDef
      const testEnv = testEnvDef
      const testType = testTypeDef
      const resultStatus = resultStatusDef

      const router = useRouter();
      const store = useStore<{ History: ListStateType}>();

      const columns = [
        {
          title: '序号',
          dataIndex: 'seq',
          width: 150,
          customRender: ({text, index}: { text: any; index: number}) => index + 1,
        },
        {
          title: '步骤',
          dataIndex: 'name',
          slots: { customRender: 'name' },
        },
        {
          title: '状态',
          dataIndex: 'status',
          slots: { customRender: 'status' },
        },
        {
          title: '检查点（编号 状态 期待结果 实际结果）',
          dataIndex: 'checkPoints',
          slots: { customRender: 'checkPoints' },
        },
      ]

      const report = computed<any[]>(() => store.state.History.item);
      const loading = ref<boolean>(true);

      const seq = router.currentRoute.value.params.seq

      const get = async (): Promise<void> => {
        loading.value = true;
        await store.dispatch('History/get', seq);
        loading.value = false;
      }
      get(seq)

      const back = ():void =>  {
        router.push(`/exec/history`)
      }

      onMounted(() => {
        console.log('onMounted')
      })

      return {
        report,
        columns,
        loading,
        back,
        momentTime,
        percent,
        testEnv,
        testType,
        resultStatus,
      }
    }
})
</script>

<style lang="less" scoped>
  .main {
    padding: 20px;
  }
</style>
