<template>
  <div class="site-main space-top space-left space-right">
    <div class="t-card-toolbar">
      <div class="left strong">
        {{ t("interpreter") }}
      </div>
      <Button class="state primary" size="sm" @click="create()">
        {{ t("create_interpreter") }}
      </Button>
    </div>
    <Table
      :columns="columns"
      :rows="interpreters"
      :isHidePaging="true"
      :isSlotMode="true"
    >
      <template #lang="record">
        {{ languageMap[record.value.lang].name }}
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

    <FormInterpreter
      :show="showCreateInterpreterModal"
      :info="editInfo"
      @submit="createInterpreter"
      @cancel="modalClose"
      ref="formInterpreter"
    />
  </div>
  <hr>

  <LanguageSettings></LanguageSettings>
</template>

<script setup lang="ts">
import { defineProps } from "vue";
import { PageTab } from "@/store/tabs";
import { useI18n } from "vue-i18n";
import {
  computed,
  ComputedRef,
  defineComponent,
  onMounted,
  ref,
  Ref,
  watch,
  reactive,
} from "vue";
import { useStore } from "vuex";
import { StateType } from "@/views/site/store";
import { momentUtcDef } from "@/utils/datetime";
import Table from "./Table.vue";
import notification from "@/utils/notification";
import Modal from "@/utils/modal";
import Button from "./Button.vue";
import LanguageSettings from "./LanguageSettings.vue";
import {getLangInterpreter, saveInterpreter} from "@/views/interpreter/service";
import {
  listInterpreter,
  removeInterpreter,
} from "@/views/interpreter/service";
import FormInterpreter from "./FormInterpreter.vue";
import { getLangSettings } from "@/views/interpreter/service";

const props = defineProps<{
  tab: PageTab;
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
      width: "15%",
    },
    {
      label: t("lang"),
      field: "lang",
      width: "15%",
    },
    {
      label: t("interpreter_path"),
      field: "path",
      width: "30%",
    },
    {
      label: t("create_time"),
      field: "createdAt",
      width: "20%",
    },
    {
      label: t("opt"),
      field: "action",
      width: "20%",
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
const formInterpreter = ref(null);
const createInterpreter = (formData) => {
    saveInterpreter(formData).then((json) => {
        if (json.code === 0) {
        formInterpreter.value.clearFormData();
        notification.success({ message: t("save_success") });
        showCreateInterpreterModal.value = false;
        list();
        }
  })
};
</script>
