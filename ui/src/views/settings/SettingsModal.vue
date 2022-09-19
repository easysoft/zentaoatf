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
        <span class="space-right">{{ t("curr_server") }} : {{defaultServerName}}</span>
        <span class="space">{{ t("curr_proxy") }} : {{defaultProxy.name == '' ? t('local_proxy') : defaultProxy.name}}</span>
      </div>
    </div>
    <p class="divider setting-space-top"></p>
    <div class="t-card-toolbar">
      <div class="left strong">
        {{ t("curr_proxy_interpreter") }}
      </div>
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
    </Table>
    <p v-else class="empty-tip">
    {{ t("empty_data") }}
    </p>

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
      <template #remote_proxy="record">
        <Table
        v-if="record.value.proxies!= undefined && record.value.proxies.length > 0"
        :columns="remoteProxyColumns"
        :rows="record.value.proxies"
        :isHideHeader="true"
        :isHideBorder="true"
        :isHidePaging="true"
        :isSlotMode="true"
        :sortable="{}"
        >
      <template #name="proxyRecord">
        <span :title="proxyRecord.value.path">{{ proxyRecord.value.name }}</span>
      </template>
      <template #action="proxyRecord">
        <div class="flex-end">
          <Button v-if="proxyRecord.value.id!=0" @click="() => handleEditProxy(proxyRecord, record)" class="tab-setting-btn" size="sm">{{t("edit")}}</Button>
          <Button v-if="proxyRecord.value.id!=0" @click="() => handleRemoveProxy(proxyRecord, record)" class="tab-setting-btn" size="sm">{{ t("delete") }}</Button>
          <Button @click="() => handleInterpreterManger(proxyRecord)" class="tab-setting-btn" size="sm">{{ t("interpreter") }}</Button>
        </div>
      </template>
    </Table>
      </template>

      <template #action="record">
        <Button @click="() => createProxy(record)" class="tab-setting-btn text-align-left" size="sm">{{t("create_remote_proxy")}}</Button>
        <Button v-if="record.value.id" @click="() => handleEditServer(record)" class="tab-setting-btn" size="sm">{{t("edit")}}</Button>
        <Button v-if="record.value.id" @click="() => handleRemoveServer(record)" class="tab-setting-btn" size="sm">{{ t("delete") }}</Button>
        <Button v-if="!record.value.is_default" @click="() => handleSetDefault(record)" class="tab-setting-btn" size="sm">{{ t("set_default") }}</Button>
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

    <FormProxy
      v-if="showCreateProxyModal"
      :show="showCreateProxyModal"
      :info="editProxyInfo"
      @submit="submitProxy"
      @cancel="modalProxyClose"
      ref="formProxy"
    />

    <InterpreterModal
      v-if="showInterpreterModal"
      :proxyInfo="currentProxy.value" 
      @cancel="interpreterModalClose"
      :showOkBtn="false"
      :showCancelBtn="false"
       />

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

import { defineProps, defineEmits, computed, onMounted, ref, watch } from "vue";
import { useI18n } from "vue-i18n";
import { useStore } from "vuex";
import { StateType } from "@/views/site/store";
import { momentUtcDef } from "@/utils/datetime";
import Table from "@/components/Table.vue";
import Modal from "@/utils/modal";
import Button from "@/components/Button.vue";
import LanguageSettings from "./LanguageSettings.vue";
import {listProxy, saveProxy, removeProxy} from "@/views/proxy/service";
import {listServer, saveServer, removeServer} from "@/views/server/service";
import FormInterpreter from "@/views/interpreter/FormInterpreter.vue";
import InterpreterModal from "@/views/interpreter/interpreterModal.vue";
import FormProxy from "@/views/proxy/FormProxy.vue";
import FormServer from "@/views/server/FormServer.vue";
import {setServerURL} from "@/utils/cache";
import { StateType as GlobalData } from "@/store/global";
import { ProxyData } from "@/store/proxy";
import {saveInterpreter,listInterpreter, removeInterpreter, getLangSettings} from "@/views/interpreter/service";

const props = defineProps<{
  show: boolean;
}>();

const emit = defineEmits<{
    (type: 'cancel', event: {event: any}) : void,
}>();

const { t, locale } = useI18n();
const momentUtc = momentUtcDef;

const remoteServers = ref<any>([]);

const editInfo = ref({});
const editProxyInfo = ref({});
const editServerInfo = ref({});

onMounted(() => {
  console.log("onMounted");
});

const store = useStore<{ global: GlobalData, proxy: ProxyData }>();
store.dispatch("proxy/fetchProxies");
store.dispatch('proxy/fetchInterpreters')

const interpreters = computed<any>(() => store.state.proxy.interpreters);
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
      width: "500px",
    },
    {
      label: t("create_time"),
      field: "createdAt",
      width: "160px",
    },
  ];
  remoteProxyColumns.value = [
    {
      label: t("name"),
      field: "name",
      width: "400px",
    },
    {
      label: t("opt"),
      field: "action",
      width: "250px",
    },
  ];
  remoteServerColumns.value = [
    {
      isKey: true,
      label: t("no"),
      field: "tableIndex",
      width: "60px",
    },
    {
      label: t("name"),
      field: "name",
      width: "160px",
    },
    {
      label: t("server_link"),
      field: "path",
      width: "260px",
    },
    {
      label: t("remote_proxy"),
      field: "remote_proxy",
      width: "700px",
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
const showInterpreterModal = ref(false);
const defaultProxy = computed<any>(() => store.state.proxy.currProxy);
const defaultServerName = ref('')

let languageMap = ref<any>({});
const getInterpretersA = async () => {
  const data = await getLangSettings();
  languageMap.value = data.languageMap;
};
getInterpretersA();

onMounted(() => {
  console.log("onMounted");
  list();
});

watch(defaultProxy, () => {
    console.log("watch default proxy")
    list()
})

const list = () => {
  listServer().then((json) => {
    if (json.code === 0) {
      let defaultServerId = 0;
      json.data.forEach(server => {
        if(server.is_default) {
          defaultServerName.value = server.name
          defaultServerId = server.id;
        }
      });
      json.data.splice(0, 0, {id: 0, path: 'local', name: t('local'), is_default: defaultServerId > 0 ? false : true});
      defaultServerName.value = defaultServerId > 0 ? defaultServerName.value : t('local');
      remoteServers.value = json.data;
      json.data.forEach((server, index) => {
        listProxy({proxyPath: server.path}).then((proxies) => {
            proxies.data.push({id:0, name: t('local_proxy'), path: 'local'})
            remoteServers.value[index].proxies = proxies.data;
        });
      });
    }
  });
};


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
const createProxy = (server) => {
  editProxyInfo.value = {proxyPath: server.value.path};
  showCreateProxyModal.value = true;
};
const handleEditProxy = (item, server) => {
  editProxyInfo.value = item;
  editProxyInfo.value.proxyPath = server.value.path
  showCreateProxyModal.value = true;
};
const submitProxy = (formData) => {
    saveProxy(formData).then((json) => {
        if (json.code === 0) {
          formProxy.value.clearFormData();
          showCreateProxyModal.value = false;
          list();
          store.dispatch("proxy/fetchProxies");
        }
  }, (json) => {console.log(json)})
};


const handleRemoveProxy = (item, server) => {
  Modal.confirm({
    title: t("confirm_delete", {
      name: item.value.name,
      typ: t("proxy_link"),
    }),
    content: '',
    okText: t("confirm"),
    cancelText: t("cancel"),
    onOk: async () => {
      await removeProxy(item.value.id, server.value.path);
      store.dispatch("proxy/fetchProxies");
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
    let serverClone = {} as any
    if(editServerInfo.value.value != undefined){
        serverClone = JSON.parse(JSON.stringify(editServerInfo.value.value));
    }
    serverClone.name = formData.name;
    serverClone.path = formData.path;
    saveServer(serverClone).then((json) => {
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
      name: item.value.name,
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

const currentProxy = ref({} as any);
const handleInterpreterManger = (proxy) => {
    currentProxy.value = proxy;
    showInterpreterModal.value = true;
}
const interpreterModalClose = () => {
    showInterpreterModal.value = false;
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
.flex-end{
  display: flex;
  justify-content: flex-end;
  flex-wrap: wrap;
}
.text-align-left{
    text-align: left;
    margin-bottom: 1rem;
}
</style>
