<template>
  <div class="site-main space-top">
    <div>
      <div class="t-card-toolbar">
        <div class="left">
          {{ t("zentao_site") }}
        </div>
        <div class="right">
          <Form labelCol="6" wrapperCol="16" class="site-search">
            <FormItem
              name="type"
              :label="t('type')"
              :info="validateInfos.enabled"
            >
              <select
                name="type"
                v-model="modelRef.enabled"
                class="form-control"
              >
                <option
                  @change="onSearch"
                  v-for="item in statusArr"
                  :key="item.value"
                  :value="item.value"
                >
                  {{ item.label }}
                </option>
              </select>
            </FormItem>
            <FormItem
              name="keywords"
              :label="t('input_keyword_to_search')"
              :info="validateInfos.keywords"
            >
              <input
                v-model="modelRef.keywords"
                @change="onSearch"
                @search="onSearch"
                class="form-control"
              />
            </FormItem>
          </Form>
        </div>
      </div>
    </div>

    <div>
      <Table
        :is-loading="false"
        :columns="columns"
        :rows="models"
        :isHidePaging="true"
        :isSlotMode="true"
      >
        <template #status="record">
          {{ disableStatus(record.value.disabled) }}
        </template>
        <template #createdAt="record">
          <span v-if="record.value.createdAt">{{ momentUtc(record.value.createdAt) }}</span>
        </template>

        <template #action="record">
          <a-button
            v-if="record.value.url"
            @click="() => edit(record.value.id)"
            type="link"
            size="small"
            >{{ t("edit") }}</a-button
          >
          <a-button
            v-if="record.value.url"
            @click="() => remove(record)"
            type="link"
            size="small"
            :loading="removeLoading.includes(record.value.seq)"
            >{{ t("delete") }}
          </a-button>
        </template>
      </Table>
    </div>

    <a-modal
      :title="t('confirm_to_delete_site')"
      v-if="confirmVisible"
      :visible="true"
      :destroy-on-close="true"
    >
      <template #footer>
        <div :class="{ 't-dir-right': !isWin }" class="t-right">
          <a-button
            @click="removeConfirmed()"
            type="primary"
            class="t-btn-gap"
            >{{ t("confirm") }}</a-button
          >
          <a-button @click="confirmVisible = false" class="t-btn-gap">{{
            t("cancel")
          }}</a-button>
        </div>
      </template>
    </a-modal>
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

import { PlusCircleOutlined } from "@ant-design/icons-vue";

import { StateType } from "@/views/site/store";
import { useRouter } from "vue-router";
import { momentUtcDef } from "@/utils/datetime";
import { PaginationConfig, QueryParams } from "@/types/data";
import debounce from "lodash.debounce";
import { ZentaoData } from "@/store/zentao";
import { disableStatusDef } from "@/utils/decorator";
import { disableStatusMap } from "@/utils/const";
import { getInitStatus, setInitStatus } from "@/utils/cache";
import Table from "./Table.vue";
import { isWindows } from "@/utils/comm";
import Form from "./Form.vue";
import FormItem from "./FormItem.vue";
import { useForm } from "@/utils/form";
import notification from "@/utils/notification";
import Modal from "@/utils/modal";

const { t, locale } = useI18n();
const isWin = isWindows();
const momentUtc = momentUtcDef;
const disableStatus = disableStatusDef;

const props = defineProps<{
  tab: PageTab;
}>();

console.log(111111, props);

const modelRef = ref({});
const rulesRef = ref({
  //   name: [{ required: true, msg: t("pls_name") }],
});
const { validate, reset, validateInfos } = useForm(modelRef, rulesRef);

const submit = () => {
  if (validate()) {
    console.log("submit");
  }
};

const localeConf = {} as any;
getInitStatus().then((initStatus) => {
  console.log("initStatus", initStatus);
  if (!initStatus) {
    localeConf.emptyText = t("pls_add_zentao_site");

    setTimeout(() => {
      setInitStatus();
    }, 1000);
  }
});

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
      field: "index",
    },
    {
      label: t("name"),
      field: "name",
    },
    {
      label: t("zentao_url"),
      field: "url",
    },
    {
      label: t("username"),
      field: "username",
    },
    {
      label: t("status"),
      field: "status",
    },
    {
      label: t("create_time"),
      field: "createdAt",
    },
    {
      label: t("opt"),
      field: "action",
      width: 260,
    },
  ];
};
setColumns();

const statusArr = ref(disableStatusMap);

const router = useRouter();
const zentaoStore = useStore<{ zentao: ZentaoData }>();
const store = useStore<{ Site: StateType }>();
const models = computed<any[]>(() => store.state.Site.queryResult.result);
const pagination = computed<PaginationConfig>(
  () => store.state.Site.queryResult.pagination
);
const queryParams = ref<QueryParams>({
  keywords: "",
  enabled: "1",
  page: pagination.value.page,
  pageSize: pagination.value.pageSize,
});

const confirmVisible = ref(false);
const model = ref({} as any);

const loading = ref<boolean>(true);
const list = (page: number) => {
  loading.value = true;
  store.dispatch("Site/list", {
    keywords: queryParams.value.keywords,
    enabled: queryParams.value.enabled,
    pageSize: pagination.value.pageSize,
    page: page,
  });
  loading.value = false;
};
list(1);

const onSearch = debounce(() => {
  list(1);
}, 500);

onMounted(() => {
  console.log("onMounted");
});

const create = () => {
  console.log("create");
  router.push(`/site/edit/0`);
};
const edit = (id) => {
  console.log("edit");
  router.push(`/site/edit/${id}`);
};

const removeLoading = ref<string[]>([]);
const remove = (item) => {
  model.value = item;
  confirmVisible.value = true;
};
const removeConfirmed = async () => {
    Modal.confirm({
          title: '删除项目',
          content: t('confirm_delete'),
          okText: t('confirm'),
          cancelText: t('cancel'),
          onOk: async () => {
              console.log('ok')
          }
        });
//   removeLoading.value = [model.value.id];

//   store.dispatch("Site/delete", model.value.id).then((success) => {
//     zentaoStore.dispatch("Zentao/fetchSitesAndProduct").then((success) => {
//       notification.success(t("delete_success"));
//       list(pagination.value.page);

//       removeLoading.value = [];
//       confirmVisible.value = false;
//     });
//   });
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
  flex: 0 0 16.666667%;
  max-width: 16.666667%;
  padding-top: calc(0.375rem + 1px);
  padding-bottom: calc(0.375rem + 1px);
  margin-bottom: 0;
  line-height: 1.5;
}
.z-form-item {
  display: flex;
  align-items: center;
  width: 100%;
}
.form-control:focus {
  color: #495057;
  background-color: #fff;
  border-color: #80bdff;
  outline: 0;
  box-shadow: 0 0 0 0.2rem rgb(0 123 255 / 25%);
}
</style>