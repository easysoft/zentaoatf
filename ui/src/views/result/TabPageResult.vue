<template>
  <div v-if="report" class="page-result dock scrollbar-y">
    <header class="single row align-center padding canvas sticky shadow-border-bottom">
      <Icon v-if="isTestPassed" icon="checkmark-circle" class="text-green" size="1.8em" />
      <Icon v-else icon="close-circle" class="text-red" size="2em" />
      <strong class="space-left">{{ t('case_num_format', {count: report.total}) }}</strong>
      <small class="rounded inline-block space-left" :class="isTestPassed ? 'green-pale padding-sm-h' : 'red-pale padding-sm-h'">{{t('pass')}} {{percent(report.pass, report.total)}}</small>
      <div class="flex-auto row justify-end result-action">
        <Button
          v-if="currProduct.id"
          @click="openResultForm()"
          :label="t('submit_result_to_zentao')"
          class="space-left border blue-pale rounded"
        />
        <Button
          v-if="report.testType != 'unit'"
          class="space-left border primary-pale rounded"
          icon="refresh"
          :label="t('re_exec_all')"
          @click="exec('all')"
        />
        <Button
          v-if="report.testType != 'unit' && report.fail > 0"
          class="space-left primary rounded"
          icon="bug-arrow-counterclockwise"
          @click="exec('fail')"
        >
          {{t('re_exec_failed') }}
        </Button>
        <Button v-if="report.testType === 'unit'" @click="exec('')" class="space-left primary rounded">
          {{t('re_exec_unit') }}
        </Button>
      </div>
    </header>
    <div class="padding-lg darken-1 divider">
      <div class="row justify-center space-bottom result-summary">
        <div v-if="report.pass"><i style="background: #16a34a;"></i> {{t("pass")}} {{report.pass}}<small class="muted">（{{percent(report.pass, report.total)}}）</small></div>
        <div v-if="report.fail"><i style="background: #dc2626;"></i> {{t("fail")}} {{report.fail}}<small class="muted">（{{percent(report.fail, report.total)}}）</small></div>
      <div v-if="report.skip"><i style="background: #94a3b8;"></i> {{t("skip")}} {{report.skip}}<small class="muted">（{{percent(report.skip, report.total)}}）</small></div>
      </div>
      <ProgressBar
        v-if="report.total"
        class="result-progress shadow space-xl-h"
        :progress="[100 * report.pass / report.total, 100 * report.fail / report.total, 100 * report.skip / report.total]"
        :height="20"
        :radius="4"
        colors="#16a34a,#dc2626,#94a3b8"
      />
      <div class="result-infos space-xl-top row justify-between gap-lg">
        <div class="row single gap-sm align-center">
          <span class="muted">{{t("test_env")}}</span>
          <span>{{testEnv(report.testEnv)}}</span>
        </div>
        <div class="row single gap-sm align-center">
          <span class="muted">{{t("test_type")}}</span>
          <span>{{report.testType}}</span>
        </div>
        <div class="row single gap-sm align-center">
          <span class="muted">{{t("exec_type")}}</span>
          <span>{{execBy(report) && te(execBy(report)) ? t(execBy(report)) : execBy(report)}}</span>
        </div>
        <div class="row single gap-sm align-center">
          <Icon icon="timer" class="muted" />
          <span>{{t("duration")}}</span>
          <span>{{report.duration}}{{ t("sec") }}</span>
          <small class="muted">(<span :title="t('start_time')">{{ momentTime(report.startTime, "YYYY-MM-DD HH:mm:ss") }}</span> ~ <span :title="t('end_time')">{{ momentTime(report.endTime, "YYYY-MM-DD HH:mm:ss") }}</span>)</small>
        </div>
      </div>
    </div>
    <div class="result-cases" v-if="report.unitResult?.length || report.funcResult?.length">
      <div class="padding single row gap align-center">
        <Icon icon="task-list-square" class="muted" />
        <strong>{{t('case_list') }}</strong>
      </div>
      <div v-if="report.testType === 'unit'" class="unit-result-list divider">
        <div v-for="result in report.unitResult" :key="`${result.testSuite}.${result.title}`" class="unit-result">
          <ListItem @click="toggleItemCollapsed(result)" no-state class="divider-top">
            <template #leading>
              <Icon v-if="result.status === 'pass'" icon="checkmark-circle-filled" class="text-green" />
              <Icon v-else-if="result.status === 'skip'" icon="subtract-circle-filled" class="muted" />
              <Icon v-else icon="close-circle-filled" class="text-red" />
            </template>
            <div clas="unit-result-title">
              <div>{{result.title}}</div>
            </div>
            <template #trailing>
              <Button
                v-if="result.status === 'fail' && currProduct.id"
                icon="bug"
                @click="showInfo(result)"
                class="space-left rounded pure text-red"
                :label="t('view_error')"
              />
              <Button
                v-if="result.status === 'fail' && currProduct.id"
                icon="bug"
                @click="openBugForm(result)"
                class="space-left rounded pure text-blue"
                :label="t('submit_bug_to_zentao')"
              />
            </template>
          </ListItem>
          <div class="unit-result-info row gap-xl padding-bottom small">
            <div class="muted">{{t('suite')}}: <code>{{result.testSuite}}</code></div>
            <div class="row single gap-sm align-center muted">
              <Icon icon="timer" class="muted" />
              <span>{{t("duration")}}</span>
              <span>{{result.duration}}{{ t("sec") }}</span>
            </div>
          </div>
        </div>
      </div>
      <div v-else class="func-result-list divider">
        <div v-for="result in report.funcResult" :key="result.path || result.key" class="func-result" :class="isCollapsed(result) ? 'collapsed' : 'expaned'">
          <ListItem @click="toggleItemCollapsed(result)">
            <template #leading>
              <Icon v-if="result.status === 'pass'" icon="checkmark-circle-filled" class="text-green" />
              <Icon v-else-if="result.status === 'skip'" icon="subtract-circle-filled" class="muted" />
              <Icon v-else icon="close-circle-filled" class="text-red" />
            </template>
            <div class="func-result-title"><code class="small">{{result.path || result.relativePath}}</code></div>
            <template #trailing>
              <Button
                v-if="result.status === 'fail' && currProduct.id"
                icon="bug"
                @click="openBugForm(result)"
                class="space-left rounded pure text-blue"
                :label="t('submit_bug_to_zentao')"
              />
              <Icon :icon="isCollapsed(result) ? 'chevron-right' : 'chevron-down'" class="rounded pure" />
            </template>
          </ListItem>
          <template v-if="!isCollapsed(result)">
            <div v-if="result.steps.length" class="result-step-list padding-xl-left">
              <div v-for="step in result.steps" :key="step.id" class="result-step-item padding">
                <div class="row single align-center gap" :class="step.status === 'fail' ? 'red-pale' : 'green-pale'">
                  <div class="padding-sm-h small" :class="step.status === 'fail' ? 'red' : 'green'">{{resultStatus(step.status)}}</div>
                  <span>{{step.name}}</span>
                </div>
                <div class="result-step-checkpoints padding-sm-top">
                  <div v-for="checkpoint in step.checkPoints" class="result-step-checkpoint padding-sm-v row single align-stretch" :key="checkpoint.numb">
                    <div class="flex-none padding-lg-right">
                      <div class="row single align-center gap ">
                        <Icon :icon="checkpoint.status === 'fail' ? 'close-circle' : 'checkmark-circle'" :class="checkpoint.status === 'fail' ? 'text-red' : 'text-green'" />
                        <span>{{t("checkpoint")}} {{checkpoint.numb}}</span>
                      </div>
                    </div>
                    <pre class="flex-1 small darken-1 space-0">
                      <div class="text-gray darken-1 padding-sm-h">{{t('expect')}}</div>
                      <code class="padding-sm scrollbar-y">{{expectDesc(checkpoint.expect)}}&nbsp;</code>
                    </pre>
                    <pre class="flex-1 small darken-1 space-0 divider-left-dark" :class="checkpoint.status === 'fail' ? 'red-pale' : ''">
                      <div class="text-gray darken-1 padding-sm-h">{{t('actual')}}</div>
                      <code class="padding-sm scrollbar-y">{{actualDesc(checkpoint.actual)}}&nbsp;</code>
                    </pre>
                  </div>
                </div>
              </div>
            </div>
            <div v-else class="padding center muted">{{t('empty_data')}}</div>
          </template>
        </div>
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
  <Loading class="dock" v-else />
</template>

<script setup lang="ts">
import {computed, defineProps, onMounted, reactive, ref, toRefs, watch} from "vue";
import {PageTab} from "@/store/tabs";
import {useI18n} from "vue-i18n";
import {useStore} from "vuex";
import {momentUnixFormat, percentDef} from "@/utils/datetime";
import Modal from "@/utils/modal"
import {jsonStrDef} from "@/utils/dom";
import {actualDesc, execByDef, expectDesc, resultStatusDef, testEnvDef, testTypeDef,} from "@/utils/testing";
import {ZentaoData} from "@/store/zentao";
import bus from "@/utils/eventBus";
import settings from "@/config/settings";

import Button from "@/components/Button.vue";
import Icon from "@/components/Icon.vue";
import FormResult from "@/views/result/FormResult.vue";
import FormBug from "@/views/result/FormBug.vue";

import ProgressBar from '@/components/ProgressBar.vue';
import Loading from '@/components/Loading.vue';
import ListItem from '@/components/ListItem.vue';
import useTestReport from '@/hooks/use-test-report';

const {t, te, locale} = useI18n();

const store = useStore<{ Zentao: ZentaoData }>();
const currProduct = computed<any>(() => store.state.Zentao.currProduct);

const jsonStr = jsonStrDef
const execBy = execByDef;
const momentTime = momentUnixFormat;
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

const report = useTestReport(seq as string, workspaceId as number);

const isTestPassed = computed(() => {
    const {value} = report;
    return value && !value.fail;
});

const collapsedMap = reactive({});

function isCollapsed(item) {
    if (report.value.testType === 'unit') {
        return false;
    } else {
        const collapsed = collapsedMap[item.path];
        if (typeof collapsed === 'boolean') {
            return collapsed;
        }
        return item.status === 'pass' || item.status === 'skip';
    }
}

function toggleItemCollapsed(item) {
    collapsedMap[item.path] = !isCollapsed(item);
}

const exec = (scope): void => {
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
const getCaseIdsInReport = (reportVal) => {
  const allCases: object[] = [];
  const failedCases: object[] = [];

  reportVal.funcResult.forEach(cs => {
    const item = {path: cs.path, workspaceId: reportVal.workspaceId}
    allCases.push(item)
    if (cs.status === 'fail') failedCases.push(item)
  })

  return {all: allCases, fail: failedCases}
}

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
</script>

<style scoped>
.result-progress {
  border: 1px solid rgba(0,0,0,.5);
}
.result-progress:deep(.progress-bar-percent) {
  box-shadow: inset 0 0 1px 1px rgba(255,255,255,.7);
  opacity: .9;
}
.result-progress:deep(.progress-bar-percent:hover) {
  opacity: 1;
}
.result-progress:deep(.progress-bar-percent + .progress-bar-percent) {
  border-left: 1px solid rgba(0,0,0,.1);
}
.result-summary {
  gap: 10px;
}
.result-summary > div > i {
  display: inline-block;
  width: 10px;
  height: 10px;
  border-radius: 2px;
  margin-right: 4px;
}
.func-result {
  border-top: 1px solid var(--color-darken-2);
}
.func-result.expaned {
  outline: 1px solid transparent;
  outline-offset: -1px;
}
.func-result.expaned:hover {
  outline-color: var(--color-focus);
}
.func-result.expaned :deep(.func-result-title) {
  font-weight: bold;
}
.unit-result :deep(.list-item),
.func-result :deep(.list-item) {
  min-height: 28px;
}
.result-step-checkpoint pre {
  white-space: normal;
}
.result-step-checkpoint pre > code {
  white-space: pre-wrap;
  display: block;
  max-height: 150px;
  overflow-y: auto;
}
.unit-result:hover {
  background-color: var(--color-darken-1);
}
.unit-result-info {
  padding-left: 30px;
}
.result-action .btn {
  height: auto;
  padding: 3px .7em;
  min-height: calc(2em + 2px);
}
</style>
