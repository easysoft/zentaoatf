<template>
  <div class="indexlayout-main-conent">
    <a-card :bordered="false">
      <template #title>
        {{t('test_result')}}
      </template>
      <template #extra>
        <div class="opt">
          <a-button @click="openResultForm()" type="primary">{{ t('submit_result_to_zentao') }}</a-button>

          <a-button type="link" @click="() => back()">{{ t('back') }}</a-button>
        </div>
      </template>

      <div class="main">
        <a-row>
          <a-col :span="2" class="t-bord t-label-right">{{ t('test_env') }}</a-col>
          <a-col :span="4">{{ testEnv(report.testEnv) }}</a-col>
          <a-col :span="2" class="t-bord t-label-right">{{ t('test_type') }}</a-col>
          <a-col :span="4">{{ testType(report.testType) }}</a-col>
          <a-col :span="2" class="t-bord t-label-right">{{ t('exec_type') }}</a-col>
          <a-col :span="4">{{ execBy(report) }}</a-col>
        </a-row>
        <a-row>
          <a-col :span="2" class="t-bord t-label-right">{{ t('start_time') }}</a-col>
          <a-col :span="4">{{ momentTime(report.startTime) }}</a-col>
          <a-col :span="2" class="t-bord t-label-right">{{ t('end_time') }}</a-col>
          <a-col :span="4">{{ momentTime(report.endTime) }}</a-col>
          <a-col :span="2" class="t-bord t-label-right">{{ t('duration') }}</a-col>
          <a-col :span="4">{{ report.duration }}秒</a-col>
        </a-row>
        <a-row>
          <a-col :span="2" class="t-bord t-label-right">{{ t('case_num') }}</a-col>
          <a-col :span="4">{{ report.total }}</a-col>
          <a-col :span="2" class="t-bord t-label-right">{{ t('pass') }}</a-col>
          <a-col :span="4">{{ report.pass }}（{{ percent(report.pass, report.total) }}）</a-col>
          <a-col :span="2" class="t-bord t-label-right">{{ t('fail') }}</a-col>
          <a-col :span="4">{{ report.fail }}（{{ percent(report.fail, report.total) }}）</a-col>
          <a-col :span="2" class="t-bord t-label-right">{{ t('skip') }}</a-col>
          <a-col :span="4">{{ report.skip }}（{{ percent(report.skip, report.total) }}）</a-col>
        </a-row>

        <a-row>
          <a-col :span="2" class="t-bord t-label-right">{{ t('case_detail') }}</a-col>
        </a-row>
        <a-row>
          <a-col :span="2"></a-col>
          <a-col :span="22">
            <a-table
                :columns="columns"
                :data-source="report.unitResult"
                row-key="id"
                :pagination="false">
              <template #seq="{ record }">
                {{ record.id }}
              </template>
              <template #duration="{ record }">
                {{ record.duration }}
              </template>
              <template #status="{ record }">
                <span :class="'t-'+record.status">{{ resultStatus(record.status) }}</span>
              </template>
              <template #info="{ record }">
                <template v-if="record.failure">
                  <a-button type="link" @click="showInfo(record.id)">{{t('view_error')}}</a-button>
                  <a-modal v-model:visible="visibleMap[record.id]"
                           :title="t('error_detail')"
                           @ok="closeInfo(record.id)"
                           width="1000px">
                    <p>{{ jsonStr(record.failure) }}</p>
                  </a-modal>
                </template>
              </template>
            </a-table>
          </a-col>
        </a-row>
      </div>

      <result-form
          v-if="resultFormVisible"
          :data="resultFormData"
          :onSubmit="submitResultForm"
          :onCancel="cancelResultForm"
      />

    </a-card>
  </div>
</template>

<script lang="ts">
import {computed, defineComponent, onMounted, reactive, Ref, ref} from "vue";
import {useStore} from 'vuex';
import {StateType as ListStateType} from "@/views/exec/store";
import {useRouter} from "vue-router";
import {momentTimeDef, percentDef} from "@/utils/datetime";
import {execByDef, resultStatusDef, testEnvDef, testTypeDef} from "@/utils/testing";
import {jsonStrDef} from "@/utils/dom";
import {notification} from "ant-design-vue";
import {submitResultToZentao} from "@/views/exec/service";
import ResultForm from './component/result.vue'
import {useI18n} from "vue-i18n";

interface UnitTestResultPageSetupData {
  t: (key: string | number) => string;
  report: Ref;
  columns: any[]

  loading: Ref<boolean>;
  exec: (scope) => void;
  back: () => void;

  visibleMap: Ref
  showInfo: (id) => void;
  closeInfo: (id) => void;

  resultFormData: Ref
  resultFormVisible: Ref<boolean>;
  setResultFormVisible: (val: boolean) => void;
  openResultForm: () => void
  submitResultForm: (model) => void
  cancelResultForm: () => void;

  execBy: (record) => string;
  momentTime: (tm) => string;
  percent: (numb, total) => string;
  testEnv: (code) => string;
  testType: (code) => string;
  resultStatus: (code) => string;
  jsonStr: (record) => string;
}

export default defineComponent({
  name: 'UnitTestResultPage',
  components: {
    ResultForm
  },

  setup(): UnitTestResultPageSetupData {
    const { t } = useI18n();

    const execBy = execByDef
    const momentTime = momentTimeDef
    const percent = percentDef
    const testEnv = testEnvDef
    const testType = testTypeDef
    const resultStatus = resultStatusDef
    const jsonStr = jsonStrDef
    const visibleMap = reactive<any>({})

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
        title: t('index'),
        dataIndex: 'title',
        slots: {customRender: 'title'},
      },
      {
        title: t('suite'),
        dataIndex: 'testSuite',
      },
      {
        title: t('duration_sec'),
        dataIndex: 'duration',
        slots: {customRender: 'duration'},
      },
      {
        title: t('status'),
        dataIndex: 'status',
        slots: {customRender: 'status'},
      },
      {
        title: t('info'),
        dataIndex: 'info',
        slots: {customRender: 'info'},
      },
    ]

    const report = computed<any>(() => store.state.History.item);
    const loading = ref<boolean>(true);

    const seq = router.currentRoute.value.params.seq as string

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

    const resultFormData = ref({})
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

    const showInfo = (id): void => {
      visibleMap[id] = true
    }
    const closeInfo = (id): void => {
      visibleMap[id] = false
    }

    const back = (): void => {
      router.push(`/exec/history`)
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

      visibleMap,
      showInfo,
      closeInfo,

      resultFormData,
      resultFormVisible,
      setResultFormVisible,
      openResultForm,
      submitResultForm,
      cancelResultForm,

      execBy,
      momentTime,
      percent,
      testEnv,
      testType,
      resultStatus,
      jsonStr,
    }
  }
})
</script>

<style lang="less" scoped>
.main {
  padding: 20px;
}
</style>
