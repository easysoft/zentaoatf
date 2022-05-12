<template>
  <ButtonGroup class="space-left">
    <Button id="siteMenuToggle"
            :label="currSite.name"
            icon="zentao"
            class="rounded border lighten-16"
            iconColor="var(--color-blue)"
            iconClass="off-off"
            suffix-icon="caret-down"/>
    <Button id="productMenuToggle"
            :label="currProduct.name"
            icon="cube"
            class="rounded border lighten-16"
            suffix-icon="caret-down"/>
  </ButtonGroup>
{{currSite.id}}
  <DropdownMenu
      toggle="#siteMenuToggle"
      :items="sites"
      @click="selectSite"
      :checkedKey="currSite.id"
      keyName="id"
      :replaceFields="replaceFields"
  >
  </DropdownMenu>

  <DropdownMenu
      toggle="#productMenuToggle"
      :items="products"
      :checkedKey="currProduct.id"
      keyName="id"
      :replaceFields="replaceFields"
  />

</template>

<script setup lang="ts">
import Button from './Button.vue';
import ButtonGroup from './ButtonGroup.vue';
import DropdownMenu from './DropdownMenu.vue';
import {useI18n} from "vue-i18n";
import {useRouter} from "vue-router";
import {useStore} from "vuex";
import {ZentaoData} from "@/store/zentao";
import {computed, onMounted, watch} from "vue";
import {getInitStatus} from "@/utils/cache";
import {notification} from "ant-design-vue";

const { t } = useI18n();
const router = useRouter();

const store = useStore<{ Zentao: ZentaoData }>();

const sites = computed<any[]>(() => store.state.Zentao.sites);
const products = computed<any>(() => store.state.Zentao.products);

const currSite = computed<any>(() => store.state.Zentao.currSite);
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

const selectSite = (item): void => {
  console.log('selectSite', item.key)
  store.dispatch('Zentao/fetchSitesAndProduct', {currSiteId: item.key}).then((payload) => {
    showZentaoMsg(payload)
  })
}
const selectProduct = (product): void => {
  console.log('selectProduct', product)
  store.dispatch('Zentao/fetchSitesAndProduct', {currProductId: product.id}).then((payload) => {
    showZentaoMsg(payload)
  })
}

const replaceFields = {
  key: 'id',
  title: 'name',
}

</script>
