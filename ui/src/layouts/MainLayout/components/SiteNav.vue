<template>
  <ButtonGroup class="space-left">
    <Button id="siteMenuToggle"
            :label="currSite.name"
            :icon="currSite.username ? 'zentao' : 'hard-drive-filled'"
            class="rounded border lighten-16"
            :iconClass="currSite.username ? 'text-blue' : 'text-secondary'"
            suffix-icon="caret-down"/>
    <Button v-if="products.length > 0"
            id="productMenuToggle"
            :label="currProduct.name"
            icon="cube"
            class="rounded border lighten-16"
            suffix-icon="caret-down"/>
  </ButtonGroup>

  <DropdownMenu
      toggle="#siteMenuToggle"
      class="padding-0-bottom"
      :items="sites"
      keyName="id"
      :checkedKey="currSite.id"
      @click="selectSite"
      :replaceFields="replaceFields"
  >
    <div class="divider space-sm-top"></div>
    <ListItem class="darken-1" icon="globe" :title="t('site_management')" @click="openSiteManagementTab()" />
  </DropdownMenu>

  <DropdownMenu
      v-if="products.length > 0"
      toggle="#productMenuToggle"
      :items="products"
      keyName="id"
      :checkedKey="currProduct.id"
      @click="selectProduct"
      :replaceFields="replaceFields"
  />
</template>

<script setup lang="ts">
import Button from './Button.vue';
import ButtonGroup from './ButtonGroup.vue';
import DropdownMenu from './DropdownMenu.vue';
import ListItem from './ListItem.vue';
import {useI18n} from "vue-i18n";
import {useRouter} from "vue-router";
import {useStore} from "vuex";
import {ZentaoData} from "@/store/zentao";
import {computed, onMounted, watch, ref} from "vue";
import {getInitStatus} from "@/utils/cache";
import {notification} from "ant-design-vue";

const { t } = useI18n();
const router = useRouter();
const store = useStore<{ Zentao: ZentaoData }>();

const products = computed<any>(() => store.state.Zentao.products);
const currSite = computed<any>(() => store.state.Zentao.currSite);
const sites = computed<any[]>(() => store.state.Zentao.sites.map(site => ({
    icon: site.username ? 'zentao' : 'hard-drive-filled',
    iconClass: 'muted',
    ...site
})));
const currProduct = computed<any>(() => store.state.Zentao.currProduct);

store.dispatch('Zentao/fetchSitesAndProduct', {}).then((payload) => {
  showZentaoMsg(payload)
})

watch(currSite, ()=> {
  console.log(`watch currSite id = ${currSite.value.id}`)

  if (currSite.value.id <= 0) {
    getInitStatus().then((initStatus) => {
      console.log('initStatus', initStatus)
      if (!initStatus) {
        router.push(`/site/list`)
      }
    })
  }
})

onMounted(() => {
  console.log('onMounted')
})

const showZentaoMsg = (payload): void => {
  if (payload.zentaoErr) {
    notification.error({
      message: t('zentao_request_failed_title'),
      description: t('zentao_request_failed_desc'),
      duration: null,
    });
  }
}

const openSiteManagementTab = (showCreateSiteModal?: boolean) => {
    console.log('openSiteManagementTab');
    store.dispatch('tabs/open', {
        id: 'sites',
        title: t('site_management'),
        type: 'sites',
        data: {showCreateSiteModal}
    });
};

const selectSite = (item): void => {
  console.log('selectSite', item.key)
  store.dispatch('Zentao/fetchSitesAndProduct', {currSiteId: item.key}).then((payload) => {
    showZentaoMsg(payload)
  })
}

const selectProduct = (item): void => {
  console.log('selectProduct', item.key)
  store.dispatch('Zentao/fetchSitesAndProduct', {currProductId: item.key}).then((payload) => {
    showZentaoMsg(payload)
  })
}

const replaceFields = {
  key: 'id',
  title: 'name',
}

</script>

<style>
.top-line {border-top: 1px dashed var(--color-green)}
</style>
