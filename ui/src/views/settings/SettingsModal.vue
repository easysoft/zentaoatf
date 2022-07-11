<template>
<ZModal
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
    />

    <p class="divider setting-space-top"></p>
    <div class="t-card-toolbar">
      <div class="left strong">
        {{ t("remote_proxy") }}
      </div>
      <Button class="state primary" size="sm" @click="createProxy()">
        {{ t("create_remote_proxy") }}
      </Button>
    </div>
    <Table
      v-if="remoteProxies.length > 0"
      :columns="remoteProxyColumns"
      :rows="remoteProxies"
      :isHidePaging="true"
      :isSlotMode="true"
      :sortable="{}"
    >
      <template #createdAt="record">
        <span v-if="record.value.createdAt">{{ momentUtc(record.value.createdAt) }}</span>
      </template>

      <template #action="record">
        <Button @click="() => handleEditProxy(record)" class="tab-setting-btn" size="sm">{{
          t("edit")
        }}</Button>
        <Button @click="() => handleRemoveProxy(record)" class="tab-setting-btn" size="sm"
          >{{ t("delete") }}
        </Button>
      </template>
    </Table>
    <p v-else class="empty-tip">
    {{ t("empty_data") }}
    </p>

    <FormProxy
      v-if="showCreateProxyModal"
      :show="showCreateProxyModal"
      :info="editProxyInfo"
      @submit="submitProxy"
      @cancel="modalProxyClose"
      ref="formProxy"
    />

    <p class="divider setting-space-top"></p>
    <div class="t-card-toolbar">
      <div class="left strong">
        {{ t("remote_server") }}
      </div>
      <Button class="state primary" size="sm" @click="createServer()">
        {{ t("create_remote_server") }}
      </Button>
    </div>
    <Table
      v-if="remoteServers.length > 0"
      :columns="remoteServerColumns"
      :rows="remoteServers"
      :isHidePaging="true"
      :isSlotMode="true"
      :sortable="{}"
    >
      <template #createdAt="record">
        <span v-if="record.value.createdAt">{{ momentUtc(record.value.createdAt) }}</span>
      </template>

      <template #action="record">
        <Button v-if="record.value.id" @click="() => handleEditServer(record)" class="tab-setting-btn" size="sm">{{
          t("edit")
        }}</Button>
        <Button v-if="record.value.id" @click="() => handleRemoveServer(record)" class="tab-setting-btn" size="sm"
          >{{ t("delete") }}
        </Button>
        <Button v-if="!record.value.is_default" @click="() => handleSetDefault(record)" class="tab-setting-btn" size="sm"
          >{{ t("set_default") }}
        </Button>
      </template>
    </Table>
    <p v-else class="empty-tip">
    {{ t("empty_data") }}
    </p>

    <FormServer
      v-if="showCreateServerModal"
      :show="showCreateServerModal"
      :info="editServerInfo"
      @submit="submitServer"
      @cancel="modalServerClose"
      ref="formServer"
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
import Table from "@/components/Table.vue";
import notification from "@/utils/notification";
import Modal from "@/utils/modal";
import Button from "@/components/Button.vue";
import LanguageSettings from "./LanguageSettings.vue";
import {listInterpreter, saveInterpreter, removeInterpreter} from "@/views/interpreter/service";
import {listProxy, saveProxy, removeProxy} from "@/views/proxy/service";
import {listServer, saveServer, removeServer} from "@/views/server/service";
import FormInterpreter from "@/views/interpreter/FormInterpreter.vue";
import { getLangSettings } from "@/views/interpreter/service";
import FormProxy from "@/views/proxy/FormProxy.vue";
import FormServer from "@/views/server/FormServer.vue";
import {setServerURL} from "@/utils/cache";
import { StateType as GlobalData } from "@/store/global";
import { ProxyData } from "@/store/proxy";

const props = defineProps<{
  show: boolean;
}>();

const emit = defineEmits<{
    (type: 'cancel', event: {event: any}) : void,
}>();

const { t, locale } = useI18n();
const momentUtc = momentUtcDef;

const interpreters = ref<any>([]);
const remoteServers = ref<any>([]);

const editInfo = ref({});
const editProxyInfo = ref({});
const editServerInfo = ref({});

onMounted(() => {
  console.log("onMounted");
});

const store = useStore<{ global: GlobalData, proxy: ProxyData }>();
store.dispatch("proxy/fetchProxies");
const remoteProxies = computed<any[]>(() => store.state.proxy.proxies);
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
const remoteProxyColumns = ref([] as any[]);
const remoteServerColumns = ref([] as any[]);
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
  remoteProxyColumns.value = [
    {
      isKey: true,
      label: t("no"),
      field: "id",
      width: "60px",
    },
    // {
    //   label: t("lang"),
    //   field: "lang",
    //   width: "60px",
    // },
    {
      label: t("name"),
      field: "name",
    },
    {
      label: t("proxy_link"),
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
  remoteServerColumns.value = [
    {
      isKey: true,
      label: t("no"),
      field: "id",
      width: "60px",
    },
    {
      label: t("name"),
      field: "name",
    },
    {
      label: t("server_link"),
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
const showCreateProxyModal = ref(false);
const showCreateServerModal = ref(false);

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
  listServer().then((json) => {
    if (json.code === 0) {
      let defaultServerId = 0;
      json.data.forEach(server => {
        if(server.is_default) {
          defaultServerId = server.id;
        }
      });
      json.data.splice(0, 0, {id: 0, path: t("local"), is_default: defaultServerId > 0 ? false : true});
      remoteServers.value = json.data;
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

const formProxy = ref({} as any);
const createProxy = () => {
  editProxyInfo.value = {};
  showCreateProxyModal.value = true;
};
const handleEditProxy = (item) => {
  editProxyInfo.value = item;
  showCreateProxyModal.value = true;
};
const submitProxy = (formData) => {
    // if(formData.type !== 'ztf'){
    //   formData.lang = '';
    // }
    saveProxy(formData).then((json) => {
        if (json.code === 0) {
          formProxy.value.clearFormData();
          showCreateProxyModal.value = false;
          list();
          store.dispatch("proxy/fetchProxies");
        }
  }, (json) => {console.log(json)})
};


const handleRemoveProxy = (item) => {
  Modal.confirm({
    title: t("confirm_delete", {
      name: item.value.path,
      typ: t("proxy_link"),
    }),
    content: '',
    okText: t("confirm"),
    cancelText: t("cancel"),
    onOk: async () => {
      await removeProxy(item.value.id);
      list();
    },
  });
};

const modalProxyClose = () => {
  showCreateProxyModal.value = false;
};

const formServer = ref({} as any);
const createServer = () => {
  editServerInfo.value = {};
  showCreateServerModal.value = true;
};
const handleEditServer = (item) => {
  editServerInfo.value = item;
  showCreateServerModal.value = true;
};
const handleSetDefault = (item) => {
    remoteServers.value.forEach(server => {
        if(server.id){
            let serverClone = JSON.parse(JSON.stringify(server));
            serverClone.is_default = item.value.id == server.id;
            saveServer(serverClone).then((json) => {
                if (json.code === 0) {
                    if(serverClone.is_default){
                        setServerURL(server.path);
                        store.commit('global/setServerUrl', server.path);
                        store.dispatch('proxy/fetchProxies');
                        store.dispatch('Zentao/fetchSitesAndProduct', {})
                        list();
                    }
                }
            }, (json) => {console.log(json)})
        }
    });
    if(!item.value.id){
        setTimeout(() => {
            setServerURL('local');
            store.commit('global/setServerUrl', 'local');
            store.dispatch('proxy/fetchProxies');
            store.dispatch('Zentao/fetchSitesAndProduct', {})
        }, 200);
    }
};
const submitServer = (formData) => {
    saveServer(formData).then((json) => {
        if (json.code === 0) {
          formServer.value.clearFormData();
          showCreateServerModal.value = false;
          list();
        }
  }, (json) => {console.log(json)})
};


const handleRemoveServer = (item) => {
  Modal.confirm({
    title: t("confirm_delete", {
      name: item.value.path,
      typ: t("server_link"),
    }),
    content: '',
    okText: t("confirm"),
    cancelText: t("cancel"),
    onOk: async () => {
      await removeServer(item.value.id);
      list();
    },
  });
};

const modalServerClose = () => {
  showCreateServerModal.value = false;
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
