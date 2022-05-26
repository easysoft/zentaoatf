<template>
  <div class="indexlayout-main-content space-top">
    <div>
      <div class="opt">
        <Button
          v-if="reportRef.testType != 'unit'"
          @click="exec('all')"
          class="space-left state primary"
          :label="t('re_exec_all')"
          type="button"
        />
        <Button v-if="reportRef.testType != 'unit'" @click="exec('fail')" class="space-left state primary">{{
          t("re_exec_failed")
        }}</Button>
        <Button v-if="reportRef.testType == 'unit'" @click="exec('')" class="space-left state primary">{{
          t("re_exec_unit")
        }}</Button>

        <Button
          v-if="currProduct.id"
          @click="openResultForm()"
          :label="t('submit_result_to_zentao')"
          type="button"
          class="space-left"
        />
      </div>

      <div class="main">
        <div class="summary">
          <Row :gutter="10">
            <Col :span="2" class="t-bord t-label-right">{{
              t("test_env")
            }}</Col>
            <Col :span="6">{{ testEnv(reportRef.testEnv) }}</Col>

            <Col :span="2" class="t-bord t-label-right">{{
              t("start_time")
            }}</Col>
            <Col :span="6">{{ momentTime(reportRef.startTime) }}</Col>

            <Col :span="2" class="t-bord t-label-right">{{
              t("case_num")
            }}</Col>
            <Col :span="6">{{ reportRef.total }}</Col>
          </Row>
          <Row>
            <Col :span="2" class="t-bord t-label-right">{{
              t("test_type")
            }}</Col>
            <Col :span="6">{{ testType(reportRef.testType) }}</Col>

            <Col :span="2" class="t-bord t-label-right">{{
              t("end_time")
            }}</Col>
            <Col :span="6">{{ momentTime(reportRef.endTime) }}</Col>

            <Col :span="2" class="t-bord t-label-right">{{ t("pass") }}</Col>
            <Col :span="6" class="t-pass"
              >{{ reportRef.pass }}（{{
                percent(reportRef.pass, reportRef.total)
              }}）</Col
            >
          </Row>

          <Row>
            <Col :span="2" class="t-bord t-label-right">{{
              t("exec_type")
            }}</Col>
            <Col :span="6">{{ execBy(reportRef) }}</Col>

            <Col :span="2" class="t-bord t-label-right">{{
              t("duration")
            }}</Col>
            <Col :span="6">{{ reportRef.duration }}{{ t("sec") }}</Col>

            <Col :span="2" class="t-bord t-label-right">{{ t("fail") }}</Col>
            <Col :span="6" class="t-fail"
              >{{ reportRef.fail }}（{{
                percent(reportRef.fail, reportRef.total)
              }}）</Col
            >
          </Row>

          <Row>
            <Col :span="16"></Col>

            <Col :span="2" class="t-bord t-label-right">{{ t("ignore") }}</Col>
            <Col :span="6" class="t-skip"
              >{{ reportRef.skip }}（{{
                percent(reportRef.skip, reportRef.total)
              }}）</Col
            >
          </Row>

          <div class="v-line v-line1"></div>
          <div class="v-line v-line2"></div>
        </div>

        <Row>
          <Col :span="2" class="t-bord t-label-right">{{
            t("case_detail")
          }}</Col>
        </Row>
        <Row>
          <Col :width="'2'"></Col>
          <Col :width="'98'" v-if="reportRef.testType != 'unit'">
            <template v-for="cs in reportRef.funcResult" :key="cs.id">
              <div class="case-info">
                <div class="info">
                  <span>{{ cs.id }}. {{ cs.path }}</span> &nbsp;
                  <span :class="'t-' + cs.status">
                    <icon-svg
                      type="pass"
                      v-if="cs.status === 'pass'"
                    ></icon-svg>
                    <icon-svg
                      type="fail"
                      v-if="cs.status === 'fail'"
                    ></icon-svg>
                    <icon-svg
                      type="skip"
                      v-if="cs.status === 'skip'"
                    ></icon-svg>
                  </span>
                </div>
                <div class="buttons" v-if="cs.status === 'fail'">
                  <Button
                    v-if="currProduct.id"
                    @click="openBugForm(cs)"
                    class="space-left"
                    :label="t('submit_bug_to_zentao')"
                  />
                </div>
              </div>

              <Table
                :columns="columns"
                :rows="cs.steps"
                :isHidePaging="true"
                :isSlotMode="true"
              >
                <template #no="record">
                  {{ record.value.id }}
                </template>

                <template #name="record">
                  {{ record.value.name }}
                </template>

                <template #status="record">
                  <span :class="'t-' + record.value.status">
                    <span class="dot"><icon-svg type="dot" /></span>
                    <span>{{ resultStatus(record.value.status) }}</span>
                  </span>
                </template>

                <template #checkpoint="record">
                  <div
                    v-for="checkPointItem in record.value.checkPoints"
                    :key="checkPointItem.numb"
                  >
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
              <br />
            </template>
          </Col>
          <Col :width="'98'" v-else>
              <Table
                :columns="columns"
                :rows="reportRef.unitResult"
                :isHidePaging="true"
                :isSlotMode="true"
              >
                <template #status="record">
                  <span :class="'t-' + record.value.status">
                    <span class="dot"><icon-svg type="dot" /></span>
                    <span>{{ resultStatus(record.value.status) }}</span>
                  </span>
                </template>
                <template #duration="record">
                    {{ record.value.duration }}
                </template>
                <template #opt="record">
                    <template v-if="record.value.failure">
                    <span @click="showInfo(record.value)" class="t-link t-primary">{{t('view_error')}}</span>
                </template>
                </template>
              </Table>
              <br />
          </Col>
        </Row>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { defineProps, toRefs, reactive } from "vue";
import { PageTab } from "@/store/tabs";
import { useI18n } from "vue-i18n";
import { computed, defineComponent, onMounted, ref, Ref, watch } from "vue";
import { useStore } from "vuex";
import { momentUnixDef, percentDef } from "@/utils/datetime";
import Button from "./Button.vue";
import Row from "./Row.vue";
import Col from "./Col.vue";
import Table from "./Table.vue";
import Modal from "@/utils/modal"
import {jsonStrDef} from "@/utils/dom";
import {
  execByDef,
  resultStatusDef,
  testEnvDef,
  testTypeDef,
  expectDesc,
  actualDesc,
} from "@/utils/testing";

import notification from "@/utils/notification";

import { submitResultToZentao } from "@/views/result/service";
import { submitBugToZentao } from "@/services/bug";
import { ZentaoData } from "@/store/zentao";
import { StateType } from "@/views/result/store";
import IconSvg from "@/components/IconSvg/index";
import bus from "@/utils/eventBus";
import settings from "@/config/settings";
import modal from "@/utils/modal";

const { t, locale } = useI18n();

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
const { tab } = toRefs(props);
let { seq, workspaceId } = tab.value.data;

watch(
  locale,
  () => {
    console.log("watch locale", locale);
    setColumns();
  },
  { deep: true }
);

const reportRef = ref({});

const columns = ref([] as any[]);
const setColumns = () => {
  if(reportRef.value.testType === 'unit') {
    columns.value = [
      {
        label: t('no'),
        field: "id",
        width: "3%",
        isKey: true,
      },
      {
        label: t('case'),
        field: 'title',
      },
      {
        label: t('suite'),
        field: 'testSuite',
      },
      {
        label: t('duration_sec'),
        field: 'duration',
      },
      {
        label: t('status'),
        field: 'status',
      },
      {
        label: t('opt'),
        field: 'opt',
      },
    ];
    return;
  }
  columns.value = [
    {
      label: "ID",
      field: "id",
      width: "3%",
      isKey: true,
    },
    {
      label: t("step"),
      field: "name",
      width: "10%",
    },
    {
      label: t("status"),
      field: "status",
      width: "5%",
    },
    {
      label: t("checkpoint"),
      field: "checkpoint",
      width: "20%",
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
  { deep: true }
);

const get = async (): Promise<void> => {
  loading.value = true;
  await store.dispatch("Result/get", { workspaceId: workspaceId, seq: seq });
  loading.value = false;
};
get();

const exec = (scope): void => {
  console.log('exec', report.value);

  const testType = report.value.testType;
  const productId = report.value.productId;
  const workspaceId = report.value.workspaceId;
  const execById = report.value.execById;

  if (testType === "func") {
    const caseMap = getCaseIdsInReport(report.value)
    const cases = caseMap[scope]
    bus.emit(settings.eventExec, { execType: 'ztf', scripts: cases });

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
const resultFormData = ref({});
const resultFormVisible = ref<boolean>(false);
const setResultFormVisible = (val: boolean) => {
  resultFormVisible.value = val;
};
const openResultForm = () => {
  console.log("openResultForm");
  resultFormData.value = report.value;
  setResultFormVisible(true);
};
const submitResultForm = (formData) => {
  console.log("submitResultForm", formData);

  const data = Object.assign(
    {
      workspaceId: report.value.workspaceId,
      seq: report.value.seq,
    },
    formData
  );

  data.taskId = parseInt(data.taskId);

  submitResultToZentao(data).then((json) => {
    console.log("json", json);
    if (json.code === 0) {
      notification.success({
        message: t("submit_success"),
      });
      setResultFormVisible(false);
    } else {
      notification.error({
        message: t("submit_failed"),
        description: json.msg,
      });
    }
  });
};
const cancelResultForm = () => {
  setResultFormVisible(false);
};

// 提交缺陷
const bugFormData = ref({});
const bugFormVisible = ref<boolean>(false);
const setBugFormVisible = (val: boolean) => {
  bugFormVisible.value = val;
};
const openBugForm = (cs) => {
  console.log("openBugForm", cs);
  if (cs.product === 0) cs.product = "";
  cs.workspaceId = report.value.workspaceId;
  cs.seq = report.value.seq;
  bugFormData.value = cs;
  setBugFormVisible(true);
};
const submitBugForm = (formData) => {
  console.log("submitBugForm", formData);

  const data = Object.assign(
    {
      workspaceId: report.value.workspaceId,
      seq: report.value.seq,
    },
    formData
  );

  data.module = parseInt(data.module);
  data.severity = parseInt(data.severity);
  data.pri = parseInt(data.pri);

  submitBugToZentao(data).then((json) => {
    console.log("json", json);
    if (json.code === 0) {
      notification.success({
        message: t("submit_success"),
      });
      setBugFormVisible(false);
    } else {
      notification.error({
        message: t("submit_failed"),
        description: json.msg,
      });
    }
  });
};
const cancelBugForm = () => {
  setBugFormVisible(false);
};

const showInfo = (item): void => {
    Modal.confirm({
        title: t("error_detail"),
        showOkBtn: false,
        content: jsonStr(item.failure),
        cancelTitle: t("close"),
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
  padding: 20px;
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
</style>
