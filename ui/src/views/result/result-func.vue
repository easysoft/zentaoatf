<template>
  <div class="indexlayout-main-conent">
    <a-card :bordered="false">
      <template #title>
        {{t('test_result')}}
      </template>
      <template #extra>
        <div class="opt">
          <a-button @click="exec('all')" type="primary">{{ t('re_exec_all') }}</a-button>
          <a-button @click="exec('fail')" type="primary">{{ t('re_exec_failed') }}</a-button>

          <a-button @click="openResultForm()">{{ t('submit_result_to_zentao') }}</a-button>
          <a-button type="link" @click="() => back()">{{ t('back') }}</a-button>
        </div>
      </template>

      <div class="main">
        <div class="summary">
          <a-row>
            <a-col :span="2" class="t-bord t-label-right">{{ t('test_env') }}</a-col>
            <a-col :span="6">{{ testEnv(report.testEnv) }}</a-col>

            <a-col :span="2" class="t-bord t-label-right">{{ t('start_time') }}</a-col>
            <a-col :span="6">{{ momentTime(report.startTime) }}</a-col>

            <a-col :span="2" class="t-bord t-label-right">{{ t('case_num') }}</a-col>
            <a-col :span="6">{{ report.total }}</a-col>
          </a-row>
          <a-row>
            <a-col :span="2" class="t-bord t-label-right">{{ t('test_type') }}</a-col>
            <a-col :span="6">{{ testType(report.testType) }}</a-col>

            <a-col :span="2" class="t-bord t-label-right">{{ t('end_time') }}</a-col>
            <a-col :span="6">{{ momentTime(report.endTime) }}</a-col>

            <a-col :span="2" class="t-bord t-label-right">{{ t('pass') }}</a-col>
            <a-col :span="6" class="t-pass">{{ report.pass }}（{{ percent(report.pass, report.total) }}）</a-col>
          </a-row>

          <a-row>
            <a-col :span="2" class="t-bord t-label-right">{{ t('exec_type') }}</a-col>
            <a-col :span="6">{{ execBy(report) }}</a-col>

            <a-col :span="2" class="t-bord t-label-right">{{ t('duration') }}</a-col>
            <a-col :span="6">{{ report.duration }}{{ t('sec') }}</a-col>

            <a-col :span="2" class="t-bord t-label-right">{{ t('fail') }}</a-col>
            <a-col :span="6" class="t-fail">{{ report.fail }}（{{ percent(report.fail, report.total) }}）</a-col>
          </a-row>

          <a-row>
            <a-col :span="16"></a-col>

            <a-col :span="2" class="t-bord t-label-right">{{ t('ignore') }}</a-col>
            <a-col :span="6" class="t-skip">{{ report.skip }}（{{ percent(report.skip, report.total) }}）</a-col>
          </a-row>

          <div class="v-line v-line1"></div>
          <div class="v-line v-line2"></div>
        </div>

        <a-row>
          <a-col :span="2" class="t-bord t-label-right">{{ t('case_detail') }}</a-col>
        </a-row>
        <a-row>
          <a-col :span="2"></a-col>
          <a-col :span="22">
            <template v-for="cs in report.funcResult" :key="cs.id">
              <div class="case-info">
                <div class="info">
                  <span>{{ cs.id }}. {{ cs.path }}</span> &nbsp;
                  <span :class="'t-'+cs.status">
                    <icon-svg type="pass" v-if="cs.status==='pass'"></icon-svg>
                    <icon-svg type="fail" v-if="cs.status==='fail'"></icon-svg>
                    <icon-svg type="skip" v-if="cs.status==='skip'"></icon-svg>
                  </span>
                </div>
                <div class="buttons" v-if="cs.status==='fail'">
                  <a-button @click="openBugForm(cs)">{{ t('submit_bug_to_zentao') }}</a-button>
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
                  <span :class="'t-'+record.status">
                    <span class="dot"><icon-svg type="dot" /></span>
                    <span>{{ resultStatus(record.status) }}</span>
                  </span>
                </template>
                <template #checkPoints="{ record }">
                  <div v-for="checkPoint in record.checkPoints" :key="checkPoint.numb">
                    {{ checkPoint.numb }}.&nbsp;
                    <span :class="'t-'+checkPoint.status">
                      {{ resultStatus(checkPoint.status) }}
                    </span>
                    <span>"{{ checkPoint.expect }}"</span>
                    /
                    <span :class="'t-'+checkPoint.status">"{{ checkPoint.actual }}"</span>
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
        :resultData="resultFormData"
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
import {computed, defineComponent, onMounted, ref, Ref, watch} from "vue";
import {useStore} from 'vuex';
import {StateType} from "./store";
import {useRouter} from "vue-router";
import {momentUnixDef, percentDef} from "@/utils/datetime";
import {execByDef, resultStatusDef, testEnvDef, testTypeDef} from "@/utils/testing";
import {submitResultToZentao, submitBugToZentao} from "./service";
import {notification} from "ant-design-vue";
import ResultForm from "./component/result.vue";
import BugForm from "./component/bug.vue";
import {useI18n} from "vue-i18n";
import IconSvg from "@/components/IconSvg/index";
import {get, getCaseIdsFromReport} from "@/views/exec/service";

export default defineComponent({
  name: 'UnitTestResultPage',
  components: {
    ResultForm, BugForm,
    IconSvg,
  },

  setup() {
    const { t, locale } = useI18n();

    const execBy = execByDef
    const momentTime = momentUnixDef
    const percent = percentDef
    const testEnv = testEnvDef
    const testType = testTypeDef
    const resultStatus = (code) => {
      const s = resultStatusDef(code)
      return t(s)
    }

    const router = useRouter();

    watch(locale, () => {
      console.log('watch locale', locale)
      setColumns()
    }, {deep: true})

    const columns = ref([] as any[])
    const setColumns = () => {
      columns.value = [
        {
          title: t('no'),
          dataIndex: 'seq',
          width: 150,
          customRender: ({text, index}: { text: any; index: number }) => index + 1,
        },
        {
          title: t('step'),
          dataIndex: 'name',
          slots: {customRender: 'name'},
        },
        {
          title: t('status'),
          dataIndex: 'status',
          slots: {customRender: 'status'},
        },
        {
          title: t('checkpoint'),
          dataIndex: 'checkPoints',
          slots: {customRender: 'checkPoints'},
        },
      ]
    }
    setColumns()

    const store = useStore<{ result: StateType }>();
    const report = computed<any>(() => store.state.result.detailResult);
    const loading = ref<boolean>(true);

    const workspaceId = router.currentRoute.value.params.workspaceId
    const seq = router.currentRoute.value.params.seq

    const get = async (): Promise<void> => {
      loading.value = true;
      await store.dispatch('result/get', {workspaceId: workspaceId, seq: seq});
      console.log('===', report)
      loading.value = false;
    }
    get()

    const exec = (scope): void => {
      console.log(report)

      const productId = report.value.productId
      const execBy = report.value.execBy
      const execById = report.value.execById

      if (execBy === 'case')
        router.push(`/script/index/${workspaceId}/${seq}/${scope}`)
      else
        router.push(`/script/index/${execBy}/${productId}/${execById}/${seq}/${scope}`)
    }

    // 提交结果
    const resultFormData = ref({})
    const resultFormVisible = ref<boolean>(false);
    const setResultFormVisible = (val: boolean) => {
      resultFormVisible.value = val;
    }
    const openResultForm = () => {
      console.log('openResultForm')
      resultFormData.value = report.value
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
            message: t('submit_success'),
          });
          setResultFormVisible(false)
        } else {
          notification.error({
            message: t('submit_failed'),
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
      if (cs.product === 0) cs.product = ''
      cs.workspaceId = report.value.workspaceId
      cs.seq = report.value.seq
      bugFormData.value = cs
      setBugFormVisible(true)
    }
    const submitBugForm = (formData) => {
      console.log('submitBugForm', formData)

      const data = Object.assign({
        workspaceId: report.value.workspaceId,
        seq: report.value.seq
      }, formData)

      data.module = parseInt(data.module)
      data.severity = parseInt(data.severity)
      data.pri = parseInt(data.pri)

      submitBugToZentao(data).then((json) => {
        console.log('json', json)
        if (json.code === 0) {
          notification.success({
            message: t('submit_success'),
          });
          setBugFormVisible(false)
        } else {
          notification.error({
            message: t('submit_failed'),
            description: json.msg,
          });
        }
      })
    }
    const cancelBugForm = () => {
      setBugFormVisible(false);
    }

    const back = (): void => {
      router.push(`/result/list`)
    }

    onMounted(() => {
      console.log('onMounted')
    })

    return {
      t,
      report,
      columns,
      loading,
      exec,
      back,

      resultFormData,
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
.dot {
  margin-right: 5px;
  font-size: 8px;
  vertical-align: 2px
}

.summary {
  position: relative;
  .v-line {
    position: absolute;
    top: 10px;
    width: 1px;
    height: 90px;
    background: #E4E4E4;
  }
  .v-line1 {
    left: 30%;
  }
  .v-line2 {
    left: 65%;
  }
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
