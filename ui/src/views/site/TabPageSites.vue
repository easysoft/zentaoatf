<template>
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
</template>

<script setup lang="ts">
import { defineProps } from "vue";
import { PageTab } from "@/store/tabs";
import { useI18n } from "vue-i18n";
import { ref } from "vue";
import { useStore } from "vuex";
import { StateType } from "@/views/site/store";
import { momentUtcDef } from "@/utils/datetime";

import List from "@/components/List.vue";
import ListItem from "@/components/ListItem.vue";
import Icon from "@/components/Icon.vue";
import notification from "@/utils/notification";
import Modal from "@/utils/modal";
import Button from "@/components/Button.vue";
import FormSite from "@/views/site/FormSite.vue";
import useSites from '@/hooks/use-sites';

const { t } = useI18n();
const momentUtc = momentUtcDef;

const props = defineProps<{
  tab: PageTab;
}>();

const editId = ref(0);
const store = useStore<{ Site: StateType }>();
const showCreateSiteModal = ref(!!props.tab?.data?.showCreateSiteModal);

const {fetchSites, sites} = useSites();

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
      store.dispatch("Site/delete", item.id).then((success) => {
        store.dispatch("Zentao/fetchSitesAndProduct").then((success) => {
        //   notification.success(t("delete_success"));
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
  store.dispatch('Site/save', formData).then((response) => {
    if (response) {
      formSite.value.clearFormData()
      showCreateSiteModal.value = false;
      store.dispatch('Zentao/fetchSitesAndProduct').then((success) => {
        // notification.success({message: t('save_success')});
        fetchSites();
      });
    }
  })
};
</script>

<style>
</style>
