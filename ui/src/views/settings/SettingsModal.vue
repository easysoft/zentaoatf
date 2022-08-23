<template>
<ZModal
    id="settingModal"
    :showModal="props.show"
    :title="t('settings')"
    :contentStyle="{width: '90vw', height: '90vh'}"
    @onCancel="emit('cancel', {event: $event})"
  >
  <div class="site-main space-top space-left space-right">
    <LanguageSettings></LanguageSettings>
    <p class="divider setting-space-top"></p>
    <div class="t-card-toolbar">
      <div class="left strong">
        {{ t("interpreter") }}
      </div>
      <Button class="state primary" size="sm" @click="create()">
        {{ t("create_interpreter") }}
      </Button>
    </div>
    <Table
      v-if="interpreters.length > 0"
      :columns="columns"
      :rows="interpreters"
      :isHidePaging="true"
      :isSlotMode="true"
      :sortable="{}"
    >
      <template #lang="record">
        {{ languageMap[record.value.lang]?.name }}
      </template>

      <template #createdAt="record">
        <span v-if="record.value.createdAt">{{ momentUtc(record.value.createdAt) }}</span>
      </template>

      <template #action="record">
        <Button @click="() => edit(record)" class="tab-setting-btn" size="sm">{{
          t("edit")
        }}</Button>
        <Button @click="() => remove(record)" class="tab-setting-btn" size="sm"
          >{{ t("delete") }}
        </Button>
      </template>
    </Table>
    <p v-else class="empty-tip">
    {{ t("empty_data") }}
    </p>

    <FormInterpreter
      :show="showCreateInterpreterModal"
      :info="editInfo"
      @submit="createInterpreter"
      @cancel="modalClose"
      ref="formInterpreter"
    />
  </div>
</ZModal>
</template>

<script setup lang="ts">

import { defineProps, defineEmits, computed, onMounted, ref, watch } from "vue";
import { useI18n } from "vue-i18n";
import { momentUtcDef } from "@/utils/datetime";
import Table from "@/components/Table.vue";
import Modal from "@/utils/modal";
import Button from "@/components/Button.vue";
import LanguageSettings from "./LanguageSettings.vue";
import {saveInterpreter,listInterpreter, removeInterpreter, getLangSettings} from "@/views/interpreter/service";
import FormInterpreter from "@/views/interpreter/FormInterpreter.vue";

const props = defineProps<{
  show: boolean;
}>();

const emit = defineEmits<{
    (type: 'cancel', event: {event: any}) : void,
}>();

const { t, locale } = useI18n();
const momentUtc = momentUtcDef;

let interpreters = ref<any>([]);

const editInfo = ref(0);

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
      label: t("no"),
      field: "id",
      width: "60px",
    },
    {
      label: t("lang"),
      field: "lang",
      width: "60px",
    },
    {
      label: t("interpreter_path"),
      field: "path",
    },
    {
      label: t("create_time"),
      field: "createdAt",
      width: "160px",
    },
    {
      label: t("opt"),
      field: "action",
      width: "160px",
    },
  ];
};
setColumns();

const showCreateInterpreterModal = ref(false);

let languageMap = ref<any>({});
const getInterpretersA = async () => {
  const data = await getLangSettings();
  languageMap.value = data.languageMap;
};
getInterpretersA();

onMounted(() => {
  console.log("onMounted");
});

const list = () => {
  listInterpreter().then((json) => {
    console.log("---", json);

    if (json.code === 0) {
      interpreters.value = json.data;
    }
  });
};
list();

const create = () => {
  editInfo.value = {};
  showCreateInterpreterModal.value = true;
};
const edit = (item) => {
  editInfo.value = item;
  showCreateInterpreterModal.value = true;
};

const remove = (item) => {
  Modal.confirm({
    title: t("confirm_delete", {
      name: languageMap.value[item.value.lang].name,
      typ: t("script_lang"),
    }),
    content: '',
    okText: t("confirm"),
    cancelText: t("cancel"),
    onOk: async () => {
      await removeInterpreter(item.value.id);
      list();
    },
  });
};

const modalClose = () => {
  showCreateInterpreterModal.value = false;
};
const formInterpreter = ref({} as any);
const createInterpreter = (formData) => {
    saveInterpreter(formData).then((json) => {
        if (json.code === 0) {
          formInterpreter.value.clearFormData();
          showCreateInterpreterModal.value = false;
          list();
        }
  }, (json) => {console.log(json)})
};
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
