<template>
  <div class="site-main space-top space-left space-right">
    <div class="t-card-toolbar">
      <div class="left">
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

    <FormSite
      :show="showCreateInterpreterModal"
      :id="editId"
      @submit="createSite"
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
import { ZentaoData } from "@/store/zentao";
import Table from "./Table.vue";
import notification from "@/utils/notification";
import Modal from "@/utils/modal";
import Button from "./Button.vue";
import FormSite from "./FormSite.vue";
import LanguageSettings from "./LanguageSettings.vue";
import { getLangSettings } from "@/views/interpreter/service";
import {
  listInterpreter,
  removeInterpreter,
} from "@/views/interpreter/service";

const props = defineProps<{
  tab: PageTab;
}>();

const { t, locale } = useI18n();
const momentUtc = momentUtcDef;

let interpreters = ref<any>([]);
let interpreter = reactive<any>({});

const editId = ref(0);

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
      width: "10%",
    },
    {
      label: t("lang"),
      field: "lang",
      width: "25%",
    },
    {
      label: t("interpreter_path"),
      field: "path",
      width: "30%",
    },
    {
      label: t("create_time"),
      field: "createdAt",
      width: "25%",
    },
    {
      label: t("opt"),
      field: "action",
      width: "25%",
    },
  ];
};
setColumns();

const zentaoStore = useStore<{ zentao: ZentaoData }>();
const store = useStore<{ Site: StateType }>();
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
  console.log("create");
  editId.value = 0;
  showCreateInterpreterModal.value = true;
};
const edit = (id) => {
  console.log("edit", id);
  editId.value = id;
  showCreateInterpreterModal.value = true;
};

const remove = (item) => {
  Modal.confirm({
    title: "",
    content: t("confirm_delete", {
      name: item.value.name,
      typ: t("zentao_site"),
    }),
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
const createSite = (formData) => {
  store.dispatch("Site/save", formData).then((response) => {
    if (response) {
      formInterpreter.value.clearFormData();
      notification.success({ message: t("save_success") });
      showCreateInterpreterModal.value = false;
    }
  });
};
</script>

<style>
.site-search {
  display: flex;
  justify-content: flex-end;
}
.form-control {
  width: 100%;
  color: #495057;
  background-color: #fff;
  border: 1px solid #ced4da;
  border-radius: 0.25rem;
  transition: border-color 0.15s ease-in-out, box-shadow 0.15s ease-in-out;
}
.z-form-item-label {
  font-weight: 400;
  color: #212529;
  text-align: left;
  box-sizing: border-box;
  display: inline-block;
  position: relative;
  width: 100%;
  padding-right: 15px;
  padding-left: 15px;
  padding-top: calc(0.375rem + 1px);
  padding-bottom: calc(0.375rem + 1px);
  margin-bottom: 0;
  line-height: 1.5;
}
.z-form-item {
  display: flex;
  align-items: center;
  word-break: keep-all;
}
.form-control:focus {
  color: #495057;
  background-color: #fff;
  border-color: #80bdff;
  outline: 0;
  box-shadow: 0 0 0 0.2rem rgb(0 123 255 / 23%);
}
.tab-setting-btn {
  border: none;
  background: none;
  color: #1890ff;
  border-style: hidden !important;
}
.t-card-toolbar {
  display: flex;
  justify-content: space-between;
}
</style>