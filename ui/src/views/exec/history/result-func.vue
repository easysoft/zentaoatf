<template>
  <div class="indexlayout-main-conent">
    <a-card :bordered="false">
      <template #title>
        执行结果详情
      </template>
      <template #extra>
        <div class="opt">
          <a-button @click="exec('all')" type="primary">重新执行所有用例</a-button>
          <a-button @click="exec('fail')" type="primary">重新执行失败用例</a-button>

          <a-button @click="openResultForm()">提交结果到禅道</a-button>
          <a-button type="link" @click="() => back()">返回</a-button>
        </div>
      </template>

      <div class="main">
        <a-row>
          <a-col :span="2" class="t-bord t-label-right">测试环境</a-col>
          <a-col :span="4">{{ testEnv(report.testEnv) }}</a-col>
          <a-col :span="2" class="t-bord t-label-right">测试类型</a-col>
          <a-col :span="4">{{ testType(report.testType) }}</a-col>
          <a-col :span="2" class="t-bord t-label-right">执行类型</a-col>
          <a-col :span="4">{{ execBy(report) }}</a-col>
        </a-row>
        <a-row>
          <a-col :span="2" class="t-bord t-label-right">开始时间</a-col>
          <a-col :span="4">{{ momentTime(report.startTime) }}</a-col>
          <a-col :span="2" class="t-bord t-label-right">结束时间</a-col>
          <a-col :span="4">{{ momentTime(report.endTime) }}</a-col>
          <a-col :span="2" class="t-bord t-label-right">耗时</a-col>
          <a-col :span="4">{{ report.duration }}秒</a-col>
        </a-row>
        <a-row>
          <a-col :span="2" class="t-bord t-label-right">用例数</a-col>
          <a-col :span="4">{{ momentTime(report.startTime) }}</a-col>
          <a-col :span="2" class="t-bord t-label-right">通过</a-col>
          <a-col :span="4">{{ report.pass }}（{{ percent(report.pass, report.total) }}）</a-col>
          <a-col :span="2" class="t-bord t-label-right">失败</a-col>
          <a-col :span="4">{{ report.fail }}（{{ percent(report.fail, report.total) }}）</a-col>
          <a-col :span="2" class="t-bord t-label-right">跳过</a-col>
          <a-col :span="4">{{ report.skip }}（{{ percent(report.skip, report.total) }}）</a-col>
        </a-row>

        <a-row>
          <a-col :span="2" class="t-bord t-label-right">用例详情</a-col>
        </a-row>
        <a-row>
          <a-col :span="2"></a-col>
          <a-col :span="22">
            <template v-for="cs in report.funcResult" :key="cs.id">

              <div class="case-info">
                <div class="info">
                  <span>{{ cs.id }}. {{ cs.path }}</span> &nbsp;
                  <span :class="'t-'+cs.status">{{ resultStatus(cs.status) }}</span>
                </div>
                <div class="buttons" v-if="cs.status==='fail'">
                  <a-button @click="openBugForm(cs)">提交缺陷到禅道</a-button>
                </div>
              </div>

              <a-table
                  :columns="columns"
                  :data-source="cs.steps"
                  row-key="id"
                  :pagination="false">
                <template #seq="{ record }">
                  {{ record.id }}
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
                    <span :class="'t-'+item.status">{{ resultStatus(item.status) }}</span> &nbsp;
                    <span>"{{ item.expect }}"</span> / <span>"{{ item.actual }}"</span>
                  </div>
                </template>
              </a-table>

              <br/>
            </template>
          </a-col>
        </a-row>
      </div>

      <result-form
          v-if="resultFormVisible"
          :onSubmit="submitResultForm"
          :onCancel="cancelResultForm"
      />
      <bug-form
          v-if="bugFormVisible"
          :model="bugFormData"
          :onSubmit="submitBugForm"
          :onCancel="cancelBugForm"
      />

    </a-card>
  </div>
</template>

<script lang="ts">
import {computed, defineComponent, onMounted, ref, Ref} from "vue";
import {useStore} from 'vuex';
import {StateType as ListStateType} from "@/views/exec/store";
import {useRouter} from "vue-router";
import {momentTimeDef, percentDef} from "@/utils/datetime";
import {execByDef, resultStatusDef, testEnvDef, testTypeDef} from "@/utils/testing";
import {submitResultToZentao, submitBugToZentao} from "@/views/exec/service";
import {notification} from "ant-design-vue";
import ResultForm from "@/views/exec/history/component/result.vue";
import BugForm from "@/views/exec/history/component/bug.vue";

interface UnitTestResultPageSetupData {
  report: Ref;
  columns: any[]

  loading: Ref<boolean>;
  exec: (scope) => void;
  back: () => void;

  resultFormVisible: Ref<boolean>;
  setResultFormVisible: (val: boolean) => void;
  openResultForm: () => void
  submitResultForm: (model) => void
  cancelResultForm: () => void;

  bugFormData: Ref
  bugFormVisible: Ref
  setBugFormVisible: (id, val) => void;
  openBugForm: (cs) => void
  submitBugForm: (model) => void
  cancelBugForm: () => void;

  execBy: (item) => string;
  momentTime: (tm) => string;
  percent: (numb, total) => string;
  testEnv: (code) => string;
  testType: (code) => string;
  resultStatus: (code) => string;
}

export default defineComponent({
  name: 'UnitTestResultPage',
  components: {
    ResultForm, BugForm,
  },

  setup(): UnitTestResultPageSetupData {
    const execBy = execByDef
    const momentTime = momentTimeDef
    const percent = percentDef
    const testEnv = testEnvDef
    const testType = testTypeDef
    const resultStatus = resultStatusDef

    const router = useRouter();
    const store = useStore<{ History: ListStateType }>();

    const columns = [
      {
        title: '序号',
        dataIndex: 'seq',
        width: 150,
        customRender: ({text, index}: { text: any; index: number }) => index + 1,
      },
      {
        title: '步骤',
        dataIndex: 'name',
        slots: {customRender: 'name'},
      },
      {
        title: '状态',
        dataIndex: 'status',
        slots: {customRender: 'status'},
      },
      {
        title: '检查点（编号 状态 期待结果 实际结果）',
        dataIndex: 'checkPoints',
        slots: {customRender: 'checkPoints'},
      },
    ]

    const report = computed<any>(() => store.state.History.item);
    const loading = ref<boolean>(true);

    const seq = router.currentRoute.value.params.seq

    const get = async (): Promise<void> => {
      loading.value = true;
      await store.dispatch('History/get', seq);
      loading.value = false;
    }
    get()

    const exec = (scope): void => {
      console.log(report)

      const productId = report.value.productId
      const execBy = report.value.execBy
      const execById = report.value.execById

      if (execBy === 'case') router.push(`/exec/run/${execBy}/${seq}/${scope}`)
      else router.push(`/exec/run/${execBy}/${productId}/${execById}/${seq}/${scope}`)
    }

    // 提交结果
    const resultFormVisible = ref<boolean>(false);
    const setResultFormVisible = (val: boolean) => {
      resultFormVisible.value = val;
    }
    const openResultForm = () => {
      console.log('openResultForm')
      setResultFormVisible(true)
    }
    const submitResultForm = (formData) => {
      console.log('submitResultForm', formData)

      const data = Object.assign({seq: seq}, formData)
      console.log('data', data)
      submitResultToZentao(data).then((json) => {
        console.log('json', json)
        if (json.code === 0) {
          notification.success({
            message: '提交成功',
          });
          setResultFormVisible(false)
        } else {
          notification.error({
            message: '提交失败',
            description: json.msg,
          });
        }
      })
    }
    const cancelResultForm = () => {
      setResultFormVisible(false);
    }

    // 提交缺陷
    const bugFormData = ref({})
    const bugFormVisible = ref<boolean>(false);
    const setBugFormVisible = (val: boolean) => {
      bugFormVisible.value = val;
    }
    const openBugForm = (cs) => {
      console.log('openBugForm', cs)
      bugFormData.value = cs
      setBugFormVisible(true)
    }
    const submitBugForm = (formData) => {
      console.log('submitBugForm', formData)

      const data = Object.assign({seq: seq}, formData)
      console.log('data', data)
      submitBugToZentao(data).then((json) => {
        console.log('json', json)
        if (json.code === 0) {
          notification.success({
            message: '提交成功',
          });
          setResultFormVisible(false)
        } else {
          notification.error({
            message: '提交失败',
            description: json.msg,
          });
        }
      })
    }
    const cancelBugForm = () => {
      setBugFormVisible(false);
    }

    const back = (): void => {
      router.push(`/exec/history`)
    }

    onMounted(() => {
      console.log('onMounted')
    })

    return {
      report,
      columns,
      loading,
      exec,
      back,

      resultFormVisible,
      setResultFormVisible,
      openResultForm,
      submitResultForm,
      cancelResultForm,

      bugFormData,
      bugFormVisible,
      setBugFormVisible,
      openBugForm,
      submitBugForm,
      cancelBugForm,

      execBy,
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

.case-info {
  display: flex;
  .info {
    flex: 1;
  }
  .buttons {
    width: 200px;
    text-align: right;
  }
}
</style>
