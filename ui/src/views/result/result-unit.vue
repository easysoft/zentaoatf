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
          <a-col :span="6">{{ report.duration }}秒</a-col>

          <a-col :span="2" class="t-bord t-label-right">{{ t('fail') }}</a-col>
          <a-col :span="6" class="t-fail">{{ report.fail }}（{{ percent(report.fail, report.total) }}）</a-col>
        </a-row>

        <a-row>
          <a-col :span="16"></a-col>

          <a-col :span="2" class="t-bord t-label-right">{{ t('ignore') }}</a-col>
          <a-col :span="6" class="t-skip">{{ report.skip }}（{{ percent(report.ignore, report.total) }}）</a-col>
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
                <span :class="'t-'+record.status">
                  <span class="dot"><icon-svg type="dot" /></span>
                  <span>{{ resultStatus(record.status) }}</span>
                </span>
              </template>
              <template #info="{ record }">
                <template v-if="record.failure">
                  <a-button type="link" @click="showInfo(record.id)">{{t('view_error')}}</a-button>

                  <a-modal
                      v-model:visible="visibleMap[record.id]"
                      :title="t('error_detail')"
                      width="1000px">

                    <p>{{ jsonStr(record.failure) }}</p>

                    <template #footer>
                      <a-button @click="closeInfo(record.id)" type="primary">确定</a-button>
                    </template>
                  </a-modal>

                </template>
              </template>
            </a-table>
          </a-col>
        </a-row>
      </div>

      <result-form
          v-if="resultFormVisible"
          :resultData="resultFormData"
          :onSubmit="submitResultForm"
          :onCancel="cancelResultForm"
      />

    </a-card>
  </div>
</template>

<script lang="ts">
import {computed, defineComponent, onMounted, reactive, Ref, ref, watch} from "vue";
import {useStore} from 'vuex';
import {StateType} from "./store";
import {useRouter} from "vue-router";
import {momentUnixDef, percentDef} from "@/utils/datetime";
import {execByDef, resultStatusDef, testEnvDef, testTypeDef} from "@/utils/testing";
import {jsonStrDef} from "@/utils/dom";
import {notification} from "ant-design-vue";
import {submitResultToZentao} from "./service";
import ResultForm from './component/result.vue'
import {useI18n} from "vue-i18n";
import IconSvg from "@/components/IconSvg/index";

export default defineComponent({
  name: 'UnitTestResultPage',
  components: {
    ResultForm,
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

    const jsonStr = jsonStrDef
    const visibleMap = reactive<any>({})

    const router = useRouter();
    const store = useStore<{ result: StateType }>();

    watch(locale, () => {
      console.log('watch locale', locale)
      setColumns()
    }, {deep: true})

    const columns = ref([] as any[])
    const setColumns = () => {
      columns.value = [
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
    }
    setColumns()

    const report = computed<any>(() => store.state.result.detailResult);
    const loading = ref<boolean>(true);

    const seq = router.currentRoute.value.params.seq as string

    const get = async (): Promise<void> => {
      loading.value = true;
      await store.dispatch('result/get', seq);
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

    const showInfo = (id): void => {
      visibleMap[id] = true
    }
    const closeInfo = (id): void => {
      visibleMap[id] = false
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
.dot {
  margin-right: 5px;
  font-size: 8px;
  vertical-align: 2px
}
</style>
