<template>
<ZModal
    :showModal="props.show"
    :title="t('interpreter')"
    :contentStyle="{width: '90vw', height: '90vh'}"
    @onCancel="emit('cancel', {event: $event})"
  >
  <div class="site-main space-top space-left space-right">
    <div class="t-card-toolbar">
      <div class="left strong">
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
    <p v-else class="empty-tip">
    {{ t("empty_data") }}
    </p>

    <FormInterpreter
      v-if="showCreateInterpreterModal"
      :show="showCreateInterpreterModal"
      :info="editInfo"
      @submit="createInterpreter"
      @cancel="modalClose"
      ref="formInterpreter"
      :proxyPath="props.proxyInfo.path"
      :proxyId="props.proxyInfo.id"
    />
  </div>
</ZModal>
</template>

<script setup lang="ts">
import { defineProps, defineEmits } from "vue";
import { PageTab } from "@/store/tabs";
import { useI18n } from "vue-i18n";
import {
  computed,
  onMounted,
  ref,
  watch,
  withDefaults,
} from "vue";
import { useStore } from "vuex";
import { momentUtcDef } from "@/utils/datetime";
import Table from "@/components/Table.vue";
import Modal from "@/utils/modal";
import Button from "@/components/Button.vue";
import {listInterpreter, saveInterpreter, removeInterpreter} from "@/views/interpreter/service";
import FormInterpreter from "@/views/interpreter/FormInterpreter.vue";
import { getLangSettings } from "@/views/interpreter/service";
import { StateType as GlobalData } from "@/store/global";

const props = withDefaults(defineProps<{
  show?: boolean;
  proxyInfo: any;
}>(), {
  show: true,
  proxyInfo: ref({
    id: 0,
    path: "",
  }),
});

const emit = defineEmits<{
    (type: 'cancel', event: {event: any}) : void,
}>();

const { t, locale } = useI18n();
const momentUtc = momentUtcDef;

const interpreters = ref<any>([]);

const editInfo = ref({});

onMounted(() => {
  console.log("onMounted");
});

const store = useStore<{ global: GlobalData }>();
const serverUrl = computed<any>(() => store.state.global.serverUrl);

watch(serverUrl, () => {
  console.log('watch serverUrl', serverUrl.value)
  list()
}, { deep: true })
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
      field: "tableIndex",
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
  const data = await getLangSettings(props.proxyInfo.path);
  languageMap.value = data.languageMap;
};
getInterpretersA();

onMounted(() => {
  console.log("onMounted");
});

const list = () => {
  listInterpreter(props.proxyInfo.path).then((json) => {
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
      await removeInterpreter(item.value.id, props.proxyInfo.path);
      list();
      store.dispatch('proxy/fetchInterpreters', {})
    },
  });
};

const modalClose = () => {
  showCreateInterpreterModal.value = false;
};
const formInterpreter = ref({} as any);
const createInterpreter = (formData) => {
    if (formData.path == '') {
        listInterpreter(formData.proxyPath).then((json) => {
            if (json.code === 0) {
                if (json.data.length == 0) {
                    return;
                }
                json.data.forEach(item => {
                    saveInterpreter({ name: item.name, path: item.path, lang: item.lang }, props.proxyInfo.path).then((json) => {
                        if (json.code === 0) {
                            formInterpreter.value.clearFormData();
                            showCreateInterpreterModal.value = false;
                            list();
                            store.dispatch('proxy/fetchInterpreters', {})
                        }
                    }, (json) => { console.log(json) })
                });
            }
        });
    } else {
        saveInterpreter(formData, props.proxyInfo.path).then((json) => {
            if (json.code === 0) {
                formInterpreter.value.clearFormData();
                showCreateInterpreterModal.value = false;
                list();
                store.dispatch('proxy/fetchInterpreters', {})
            }
        }, (json) => { console.log(json) })
    }
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
