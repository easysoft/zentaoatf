<template>
  <div class="indexlayout-main-content space-top">
    <div>
      <div class="opt">
        <Button v-if="reportRef.testType != 'unit'"
            @click="exec('all')"
            class="space-left state primary"
            :label="t('re_exec_all')"
            type="button"
        />
        <Button v-if="reportRef.testType != 'unit'" @click="exec('fail')" class="space-left state primary">
          {{t('re_exec_failed') }}
        </Button>
        <Button v-else @click="exec('')" class="space-left state primary">{{t('re_exec_unit') }}
        </Button>

        <Button
            v-if="currProduct.id"
            @click="openResultForm()"
            :label="t('submit_result_to_zentao')"
            type="button"
            class="space-left"/>
      </div>

      <div class="main">
        <div class="summary">
          <Row>
            <Col :span="3" class="t-bord t-label-right">{{t("test_env")}}
            </Col>
            <Col :span="5">{{ testEnv(reportRef.testEnv) }}</Col>

            <Col :span="3" class="t-bord t-label-right">{{t("start_time")}}
            </Col>
            <Col :span="5">{{ momentTime(reportRef.startTime) }}</Col>

            <Col :span="3" class="t-bord t-label-right">{{t("case_num")}}
            </Col>
            <Col :span="5">{{ reportRef.total }}</Col>
          </Row>
          <Row>
            <Col :span="3" class="t-bord t-label-right">{{t("test_type")}}
            </Col>
            <Col :span="5">{{ testType(reportRef.testType) }}</Col>

            <Col :span="3" class="t-bord t-label-right">{{t("end_time")}}
            </Col>
            <Col :span="5">{{ momentTime(reportRef.endTime) }}</Col>

            <Col :span="3" class="t-bord t-label-right">{{ t("pass") }}</Col>
            <Col :span="5" class="t-pass"
            >{{ reportRef.pass }}（{{percent(reportRef.pass, reportRef.total)}}）
            </Col
            >
          </Row>

          <Row>
            <Col :span="3" class="t-bord t-label-right">{{t("exec_type")}}
            </Col>
            <Col :span="5">{{ execBy(reportRef) && te(execBy(reportRef)) ? t(execBy(reportRef)) : execBy(reportRef) }}</Col>

            <Col :span="3" class="t-bord t-label-right">{{t("duration")}}
            </Col>
            <Col :span="5">{{ reportRef.duration }}{{ t("sec") }}</Col>

            <Col :span="3" class="t-bord t-label-right">{{ t("fail") }}</Col>
            <Col :span="5" class="t-fail">{{ reportRef.fail }}（{{percent(reportRef.fail, reportRef.total)}}）
            </Col>
          </Row>

          <Row>
            <Col :span="16"></Col>

            <Col :span="3" class="t-bord t-label-right">{{ t("ignore") }}</Col>
            <Col :span="5" class="t-skip">{{ reportRef.skip }}（{{percent(reportRef.skip, reportRef.total)}}）
            </Col>
          </Row>

          <div class="v-line v-line1"></div>
          <div class="v-line v-line2"></div>
        </div>

        <Row>
          <Col :span="3" class="t-bord t-label-right">{{t("case_detail") }}</Col>
        </Row>
        <Row class="case-result-item">
          <Col :span="24" v-if="reportRef.testType != 'unit'">
            <template v-for="cs in reportRef.funcResult" :key="cs.id">
              <div class="case-info">
                <div class="info">
                  <span>{{ cs.id }}. {{ cs.path }}</span> &nbsp;
                  <span :class="'t-' + cs.status">
                    <icon-svg type="pass" v-if="cs.status === 'pass'"></icon-svg>
                    <icon-svg type="fail" v-if="cs.status === 'fail'"></icon-svg>
                    <icon-svg type="skip" v-if="cs.status === 'skip'"></icon-svg>
                  </span>
                </div>
                <div class="buttons" v-if="cs.status === 'fail'">
                  <Button
                      v-if="currProduct.id"
                      @click="openBugForm(cs)"
                      class="space-left"
                      :label="t('submit_bug_to_zentao')"/>
                </div>
              </div>

              <Table
                  v-if="cs.steps.length"
                  :columns="columns"
                  :rows="cs.steps"
                  :isHidePaging="true"
                  :isSlotMode="true">
                <template #no="record">
                  {{ record.value.id }}
                </template>

                <template #name="record">
                  {{ record.value.name }}
                </template>

                <template #status="record">
                  <span :class="'t-' + record.value.status">
                    <span class="dot"><icon-svg type="dot"/></span>
                    <span>{{ resultStatus(record.value.status) }}</span>
                  </span>
                </template>

                <template #checkpoint="record">
                  <div
                      v-for="checkPointItem in record.value.checkPoints"
                      :key="checkPointItem.numb">
                    <span class="checkpoint-num">
                      {{ checkPointItem.numb }}.
                    </span>
                    <span :class="'t-' + checkPointItem.status">
                      {{ resultStatus(checkPointItem.status) }}
                    </span>
                    &nbsp; (
                    <span>{{ expectDesc(checkPointItem.expect) }}</span>
                    /
                    <span :class="'t-' + checkPointItem.status">
                      {{ actualDesc(checkPointItem.actual) }}
                    </span>
                    )
                  </div>
                </template>
              </Table>
              <br/>
            </template>
          </Col>
          <Col :span="24" v-else>
            <Table
                v-if="reportRef.unitResult.length"
                :columns="columns"
                :rows="reportRef.unitResult"
                :isHidePaging="true"
                :isSlotMode="true"
            >
              <template #status="record">
                  <span :class="'t-' + record.value.status">
                    <span class="dot"><icon-svg type="dot"/></span>
                    <span>{{ resultStatus(record.value.status) }}</span>
                  </span>
              </template>
              <template #duration="record">
                {{ record.value.duration }}
              </template>
              <template #opt="record">
                <template v-if="record.value.failure">
                  <span @click="showInfo(record.value)" class="t-link t-primary">{{ t('view_error') }}</span>
                </template>
              </template>
            </Table>
            <br/>
          </Col>
        </Row>
      </div>
    </div>

    <!-- use v-if, each time will call init func in popup page  -->
    <FormResult
        v-if="showSubmitResultModal"
        :show="showSubmitResultModal"
        :finish="closeResultForm"
        :data="report"
        ref="formSite"
    />
    <FormBug
        v-if="showSubmitBugModal"
        :show="showSubmitBugModal"
        :finish="closeBugForm"
        :data="bugFormData"
        ref="formSite"
    />
  </div>
</template>

<script setup lang="ts">
import {computed, defineProps, onMounted, ref, toRefs, watch} from "vue";
import {PageTab} from "@/store/tabs";
import {useI18n} from "vue-i18n";
import {useStore} from "vuex";
import {momentUnixDef, percentDef} from "@/utils/datetime";
import Button from "./Button.vue";
import Row from "./Row.vue";
import Col from "./Col.vue";
import Table from "./Table.vue";
import FormResult from "./FormResult.vue";
import FormBug from "./FormBug.vue";
import Modal from "@/utils/modal"
import {jsonStrDef} from "@/utils/dom";
import {actualDesc, execByDef, expectDesc, resultStatusDef, testEnvDef, testTypeDef,} from "@/utils/testing";

import {ZentaoData} from "@/store/zentao";
import {StateType} from "@/views/result/store";
import IconSvg from "@/components/IconSvg/index";
import bus from "@/utils/eventBus";
import settings from "@/config/settings";


const {t, te, locale} = useI18n();

const store = useStore<{ Result: StateType }>();
const report = computed<any>(() => store.state.Result.detailResult);

const zentaoStore = useStore<{ Zentao: ZentaoData }>();
const currProduct = computed<any>(() => zentaoStore.state.Zentao.currProduct);

const jsonStr = jsonStrDef
const execBy = execByDef;
const momentTime = momentUnixDef;
const percent = percentDef;
const testEnv = testEnvDef;
const testType = testTypeDef;
const resultStatus = (code) => {
  const s = resultStatusDef(code);
  return t(s);
};
const props = defineProps<{
  tab: PageTab;
}>();

const {tab} = toRefs(props);
let {seq, workspaceId} = tab.value.data;

watch(
    locale,
    () => {
      console.log("watch locale", locale);
      setColumns();
    },
    {deep: true}
);

const reportRef = ref({});

const columns = ref([] as any[]);
const setColumns = () => {
  if (reportRef.value.testType === 'unit') {
    columns.value = [
      {
        label: t('no'),
        field: "id",
        width: "10%",
        isKey: true,
      },
      {
        label: t('case'),
        width: "20%",
        field: 'title',
      },
      {
        label: t('suite'),
        width: "40%",
        field: 'testSuite',
      },
      {
        label: t('duration_sec'),
        width: "10%",
        field: 'duration',
      },
      {
        label: t('status'),
        width: "10%",
        field: 'status',
      },
      {
        label: t('opt'),
        width: "10%",
        field: 'opt',
      },
    ];
    return;
  }
  columns.value = [
    {
      label: t('no'),
      width: "10%",
      field: "id",
      isKey: true,
    },
    {
      label: t("step"),
      field: "name",
      width: "20%",
    },
    {
      label: t("status"),
      field: "status",
      width: "10%",
    },
    {
      label: t("checkpoint"),
      field: "checkpoint",
      width: "50%",
    },
  ];
};
setColumns();

const loading = ref<boolean>(true);

watch(
    report,
    () => {
      if (seq !== report.value.seq || workspaceId !== report.value.workspaceId) {
        return;
      }
      console.log("watch report", report.value);
      reportRef.value = report.value;
      setColumns();
    },
    {deep: true}
);

const get = async (): Promise<void> => {
  loading.value = true;
  await store.dispatch("Result/get", {workspaceId: workspaceId, seq: seq});
  loading.value = false;
};
get();

const exec = (scope): void => {
  console.log('exec', report.value);

  const testType = report.value.testType;

  if (testType === "func") {
    const caseMap = getCaseIdsInReport(report.value)
    const cases = caseMap[scope]
    bus.emit(settings.eventExec, {execType: 'ztf', scripts: cases});

  } else if (testType === "unit") {
    const data = {
      execType: 'unit',
      cmd: report.value.testCommand,
      id: report.value.workspaceId,
      type: report.value.workspaceType,
      submitResult: report.value.submitResult,
    }
    console.log(data)
    bus.emit(settings.eventExec, data);
  }
};

// 提交结果
const showSubmitResultModal = ref(false)
const openResultForm = () => {
  console.log("openResultForm");
  showSubmitResultModal.value = true
};
const closeResultForm = () => {
  console.log("closeResultForm");
  showSubmitResultModal.value = false
};

// 提交缺陷
const showSubmitBugModal = ref(false)
const bugFormData = ref({})
const openBugForm = (cs) => {
  console.log("openBugForm");

  if (cs.product === 0) cs.product = ''
  cs.workspaceId = report.value.workspaceId
  cs.seq = report.value.seq
  bugFormData.value = cs

  showSubmitBugModal.value = true
};
const closeBugForm = () => {
  console.log("closeBugForm");
  showSubmitBugModal.value = false
};

const showInfo = (item): void => {
  Modal.confirm({
    title: t("error_detail"),
    showOkBtn: false,
    content: jsonStr(item.failure),
    cancelTitle: t("close"),
    contentStyle: {width: '800px'}
  })
}

onMounted(() => {
  console.log("onMounted");
});

const getCaseIdsInReport = (reportVal) => {
  const allCases = [] as string[]
  const failedCases = [] as string[]

  reportVal.funcResult.forEach(cs => {
    const item = {path: cs.path, workspaceId: reportVal.workspaceId}
    allCases.push(item)
    if (cs.status === 'fail') failedCases.push(item)
  })

  return {all: allCases, fail: failedCases}
}

</script>


<style lang="less" scoped>
.main {
  padding: 20px var(--space-base);
}

.dot {
  margin-right: 5px;
  font-size: 8px;
  vertical-align: 2px;
}

.summary {
  position: relative;

  .v-line {
    position: absolute;
    top: 10px;
    width: 1px;
    height: 90px;
    background: #e4e4e4;
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

.checkpoint-num {
  display: inline-block;
  width: 18px;
}

.tab-result-link {
  border: none;
  background: none;
  color: #1890ff;
  border-style: hidden !important;
}
.opt{
    display: flex;
    padding-right: 20px;
    flex-direction: row-reverse;
}
.case-result-item{
    padding-left: 1rem;
}
</style>
