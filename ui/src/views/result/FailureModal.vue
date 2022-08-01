<template>
<ZModal
    :showModal="props.show"
    :title="t('fail_result')"
    :contentStyle="{width: '90vw', height: '90vh'}"
    @onCancel="emit('cancel', {event: $event})"
    :showOkBtn="false"
  >
  <div class="site-main space-top space-left space-right">
    <Table
      v-if="failureList.length > 0"
      :columns="columns"
      :rows="failureList"
      :isHidePaging="true"
      :isSlotMode="true"
      :sortable="{}"
    >
    <template #testScriptName="record">
        <span>{{record.value.testType === "unit" || record.value.total != 1 ? record.value.workspaceName + '(' + record.value.total + ')' : record.value.testScriptName}}</span>
      </template>
      <template #execType="record">
        <span>{{execBy(record.value) && te(execBy(record.value)) ? t(execBy(record.value)) : execBy(record.value)}}</span>
      </template>
      <template #startTime="record">
        <span v-if="record.value.startTime">{{ momentUnixFormat(record.value.startTime, 'YYYY-MM-DD HH:mm') }}</span>
      </template>
      <template #action="record">
        <Button @click="() => showDetail(record.value)" class="tab-setting-btn" size="sm"
          >{{ t("view") }}
        </Button>
      </template>
    </Table>
    <p v-else class="empty-tip">
    {{ t("empty_data") }}
    </p>

  </div>
</ZModal>
</template>

<script setup lang="ts">
import { useI18n } from "vue-i18n";
import {
  computed,
  onMounted,
  ref,
  watch,
  defineProps,
  defineEmits,
  withDefaults,
} from "vue";
import {momentUnixFormat} from "@/utils/datetime";
import Table from "@/components/Table.vue";
import Modal from "@/utils/modal";
import Button from "@/components/Button.vue";
import { ZentaoData } from "@/store/zentao";
import {useStore} from 'vuex';
import {getFailure} from './service'
import { execByDef as execBy } from "@/utils/testing";

const props = withDefaults(defineProps<{
  show?: boolean;
  path: string;
}>(), {
    show: true,
});

const showModalRef = computed(() => {
  return props.show;
});

const store = useStore<{Zentao: ZentaoData }>();
const bugMap = computed<any>(() => store.state.Zentao.bugMap);

const emit = defineEmits<{
    (type: 'cancel', event: {event: any}) : void,
}>();

const { t, te, locale } = useI18n();

onMounted(() => {
  console.log("onMounted");
});

watch(
  locale,
  () => {
    console.log("watch locale", locale);
    setColumns();
  },
  { deep: true }
);

const columns = ref([] as any[]);
const setColumns = () => {
  columns.value = [
    {
      label: t("title"),
      field: "testScriptName",
      width: "60px",
    },
    {
      label: t("test_env"),
      field: "testEnv",
      width: "60px",
    },
    {
      label: t("test_type"),
      field: "testType",
      width: "60px",
    },
    {
      label: t("exec_type"),
      field: "execType",
      width: "60px",
    },
    {
      label: t("duration"),
      field: "duration",
      width: "60px",
    },
    {
      label: t("pass"),
      field: "pass",
    },
    {
      label: t("fail"),
      field: "fail",
    },
    {
      label: t("create_time"),
      field: "startTime",
      width: "160px",
    },
    {
      label: t("opt"),
      field: "action",
      width: "60px",
    },
  ];
};
setColumns();

const failureList = ref([] as any[]);
const list = () => {
  getFailure(props.path).then((res) => {
    failureList.value = res.data;
  })
};
list();
const showDetail = (item) => {
    const displayName = item.testType === "unit" || item.total != 1 ? item.workspaceName + '(' + item.total + ')' : item.testScriptName

    store.dispatch('tabs/open', {
      id: 'result-' + item.no,
      title: displayName,
      type: 'result',
      data: {seq: item.seq, workspaceId: item.workspaceId}
    });
    emit('cancel', {event: {}})
}
</script>

<style>
.empty-tip {
  text-align: center;
  padding: 20px 0;
}
.setting-space-top{
    margin-top: 1rem;
}
.t-card-toolbar{
    display: flex;
    justify-content: space-between;
    align-items: flex-end;
    margin-bottom: 1rem;
}
.tab-setting-btn {
  border: none;
  background: none;
  color: #1890ff;
  border-style: hidden !important;
}
</style>
