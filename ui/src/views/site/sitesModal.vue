<template>
<ZModal
    :showModal="props.show"
    id="siteModal"
    :title="t('site_management')"
    :contentStyle="{width: '90vw', height: '90vh'}"
    @onCancel="emit('cancel', {event: $event})"
  >

  <div class="dock scrollbar-y">
    <header class="single row align-center padding canvas sticky shadow-border-bottom">
      <strong>{{t('site_num', {count: sites.length})}}</strong>
      <div class="flex-auto row justify-end">
        <Button
          class="state primary rounded"
          icon="add"
          size="sm"
          @click="create()"
        >
          {{t('create_site')}}
        </Button>
      </div>
    </header>

    <List v-if="sites.length">
      <ListItem
        v-for="site in sites"
        :key="site.id"
        icon="zentao"
        iconClass="text-blue"
        iconSize="2em"
        divider
        no-state
      >
        <div class="row single align-center gap space-xs-v">
          <span class="text-primary">{{site.name}}</span>
          <div class="row single align-center gap-sm small muted"><Icon icon="link" size="1em" /> {{site.url}}</div>
          <div class="row single align-center gap muted">
            <Icon icon="person" size="1em" /> {{site.username}}
          </div>
        </div>
        <template #trailing>
          <Button
            class="pure rounded text-primary"
            icon="edit"
            @click="() => edit(site.id)"
          >
            {{t('edit')}}
          </Button>
          <Button
            class="pure rounded text-primary"
            icon="delete"
            @click="() => remove(site)"
          >
            {{t('delete')}}
          </Button>
        </template>
      </ListItem>
    </List>
    <div v-else class="center padding-xl empty-tip">{{ t("empty_data") }}</div>

    <FormSite
      :show="showCreateSiteModal"
      :id="editId"
      @submit="createSite"
      @cancel="modalClose"
      ref="formSite"
     />
  </div>
</ZModal>
</template>

<script setup lang="ts">
import { defineProps, defineEmits, ref, watch, computed } from "vue";
import { useI18n } from "vue-i18n";
import { useStore } from "vuex";
import { StateType } from "@/views/site/store";
import { momentUtcDef } from "@/utils/datetime";
import { StateType as GlobalData } from "@/store/global";

import List from "@/components/List.vue";
import ListItem from "@/components/ListItem.vue";
import Icon from "@/components/Icon.vue";
import Modal from "@/utils/modal";
import Button from "@/components/Button.vue";
import FormSite from "@/views/site/FormSite.vue";
import useSites from '@/hooks/use-sites';

const { t } = useI18n();

const props = defineProps<{
  show: boolean;
}>();

const emit = defineEmits<{
    (type: 'cancel', event: {event: any}) : void,
}>();

const editId = ref(0);
const store = useStore<{ Site: StateType, global: GlobalData }>();
const showCreateSiteModal = ref(false);

const {fetchSites, sites} = useSites();

const serverUrl = computed<any>(() => store.state.global.serverUrl);
watch(serverUrl, () => {
  console.log('watch serverUrl', serverUrl.value)
  fetchSites()
}, { deep: true })

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
      name: item.name,
      typ: t("zentao_site"),
    }),
    okText: t("confirm"),
    cancelText: t("cancel"),
    onOk: async () => {
      store.dispatch("Site/delete", item.id).then((_success) => {
        store.dispatch("Zentao/fetchSitesAndProduct").then((_success2) => {
          fetchSites();
        });
      });
    },
  });
};

const modalClose = () => {
  showCreateSiteModal.value = false;
}
const formSite = ref({} as any)
const createSite = (formData) => {
  let formDataNew = {
    ...formData,
  };
  if(formDataNew.url.indexOf("http") !== 0) {
    formDataNew.url = "http://" + formDataNew.url;
  }
  store.dispatch('Site/save', formDataNew).then((success) => {
    if (success) {
      formSite.value.clearFormData()
      showCreateSiteModal.value = false;
      store.dispatch('Zentao/fetchSitesAndProduct').then((_success2) => {
        fetchSites();
      });
    }
  })
};
</script>

<style>
</style>
