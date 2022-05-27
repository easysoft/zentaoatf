<template>
  <div class="site-main space-top space-left space-right">
      <div class="t-card-toolbar">
        <div class="left">
          {{ t("site_management") }}
        </div>
        <Button 
        class="state primary"
          size="sm"
         @click="create()"
         >
          {{t('create_site')}}
        </Button>
      </div>
    <Table
      :columns="columns"
      :rows="models"
      :isHidePaging="true"
      :isSlotMode="true"
    >
      <template #status="record">
        {{ disableStatus(record.value.disabled) }}
      </template>
      <template #createdAt="record">
        <span v-if="record.value.createdAt">{{
          momentUtc(record.value.createdAt)
        }}</span>
      </template>

      <template #action="record">
        <Button
          class="tab-setting-btn"
          v-if="record.value.url"
          @click="() => edit(record.value.id)"
          size="sm"
          >{{ t("edit") }}</Button
        >
        <Button
          class="tab-setting-btn"
          v-if="record.value.url"
          @click="() => remove(record)"
          size="sm"
          >{{ t("delete") }}
        </Button>
      </template>
    </Table>

    <FormSite
      :show="showCreateSiteModal"
      :id="editId"
      @submit="createSite"
      @cancel="modalClose"
      ref="formSite"
     />
  </div>
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
} from "vue";
import { useStore } from "vuex";
import { StateType } from "@/views/site/store";
import { momentUtcDef } from "@/utils/datetime";
import { PaginationConfig, QueryParams } from "@/types/data";
import { ZentaoData } from "@/store/zentao";
import { disableStatusDef } from "@/utils/decorator";
import Table from "./Table.vue";
import notification from "@/utils/notification";
import Modal from "@/utils/modal";
import Button from "./Button.vue";
import FormSite from "./FormSite.vue";

const { t, locale } = useI18n();
const momentUtc = momentUtcDef;
const disableStatus = disableStatusDef;

const props = defineProps<{
  tab: PageTab;
}>();

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
      width: "8%",
    },
    {
      label: t("name"),
      field: "name",
      width: "15%",
    },
    {
      label: t("zentao_url"),
      field: "url",
      width: "15%",
    },
    {
      label: t("username"),
      field: "username",
      width: "15%",
    },
    {
      label: t("status"),
      field: "status",
      width: "15%",
    },
    {
      label: t("create_time"),
      field: "createdAt",
      width: "15%",
    },
    {
      label: t("opt"),
      field: "action",
      width: "10%",
    },
  ];
};
setColumns();

const zentaoStore = useStore<{ zentao: ZentaoData }>();
const store = useStore<{ Site: StateType }>();
const showCreateSiteModal = ref(false)
const models = computed<any[]>(() => store.state.Site.queryResult.result);

const queryParams = ref<QueryParams>({
  keywords: "",
  enabled: "1",
  page: 1,
  pageSize: 100,
});

const model = ref({} as any);

const loading = ref<boolean>(true);
const list = (page: number) => {
  loading.value = true;
  store.dispatch("Site/list", {
    keywords: queryParams.value.keywords,
    enabled: queryParams.value.enabled,
    pageSize: queryParams.value.pageSize,
    page: page,
  });
  loading.value = false;
};
list(1);

onMounted(() => {
  console.log("onMounted");
});

const create = () => {
  console.log("create");
  editId.value = 0;
  showCreateSiteModal.value = true;
};
const edit = (id) => {
  console.log("edit", id);
  editId.value = id;
  showCreateSiteModal.value = true;
};

const remove = (item) => {
  Modal.confirm({
    title: t("confirm_delete", {
      name: item.value.name,
      typ: t("zentao_site"),
    }),
    okText: t("confirm"),
    cancelText: t("cancel"),
    onOk: async () => {
      store.dispatch("Site/delete", item.value.id).then((success) => {
        zentaoStore.dispatch("Zentao/fetchSitesAndProduct").then((success) => {
          notification.success(t("delete_success"));
          list(1);
        });
      });
    },
  });
};

const modalClose = () => {
  showCreateSiteModal.value = false;
}
const formSite = ref(null)
const createSite = (formData) => {
    store.dispatch('Site/save', formData).then((response) => {
        if (response) {
            formSite.value.clearFormData()
            notification.success({message: t('save_success')});
            showCreateSiteModal.value = false;
        }
    })
};

</script>

<style>
.tab-setting-btn {
  border: none;
  background: none;
  color: #1890ff;
  border-style: hidden !important;
}
.t-card-toolbar{
    display: flex;
    justify-content: space-between;
}
</style>