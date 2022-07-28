<template>
<ZModal
    :showModal="props.show"
    title="BUG"
    :contentStyle="{width: '90vw', height: '90vh'}"
    @onCancel="emit('cancel', {event: $event})"
    :showOkBtn="false"
  >
  <div class="site-main space-top space-left space-right">
    <Table
      v-if="bugs.length > 0"
      :columns="columns"
      :rows="bugs"
      :isHidePaging="true"
      :isSlotMode="true"
      :sortable="{}"
    >
    <template #steps="record">
        <span v-html="record.value.steps"></span>
      </template>
      <template #openedDate="record">
        <span v-if="record.value.openedDate">{{ momentUtc(record.value.openedDate) }}</span>
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
import { momentUtcDef } from "@/utils/datetime";
import Table from "@/components/Table.vue";
import notification from "@/utils/notification";
import Modal from "@/utils/modal";
import Button from "@/components/Button.vue";
import { ZentaoData } from "@/store/zentao";
import {useStore} from 'vuex';

const props = withDefaults(defineProps<{
  show?: boolean;
  caseIds: number[];
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

const { t, locale } = useI18n();
const momentUtc = momentUtcDef;

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
      isKey: true,
      label: 'ID',
      field: "id",
      width: "60px",
    },
    {
      label: t("severity"),
      field: "severity",
      width: "60px",
    },
    {
      label: t("title"),
      field: "title",
      width: "160px",
    },
    {
      label: t("status"),
      field: "statusName",
      width: "60px",
    },
    {
      label: t("step"),
      field: "steps",
    },
    {
      label: t("created_by"),
      field: "openedBy",
      width: "60px",
    },
    {
      label: t("create_time"),
      field: "openedDate",
      width: "160px",
    },
  ];
};
setColumns();

const bugs = ref([] as any[]);
const list = () => {
  props.caseIds.forEach((id) => {
    const caseBugs = bugMap.value[id] || [];
    if (caseBugs.length > 0) {
      bugs.value = bugs.value.concat(caseBugs);
    }
  });
};
list();

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
