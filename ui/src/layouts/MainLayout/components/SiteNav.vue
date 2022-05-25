<template>
  <ButtonGroup class="space-left">
    <Button id="siteMenuToggle"
            :label="currSite.name"
            icon="zentao"
            class="rounded border lighten-16"
            iconColor="var(--color-blue)"
            iconClass="off-off"
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
      :items="[...sites, {checked: false,id: -1,name: t('create_site'),password: '',url: '',username: '', titleClass: 'top-line padding-top'}]"
      keyName="id"
      :checkedKey="currSite.id"
      @click="selectSite"
      :replaceFields="replaceFields"
  >
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

  <FormSite
      :show="showCreateSiteModal"
      @submit="createSite"
      @cancel="modalClose"
      ref="formSite"
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
import {computed, onMounted, watch, ref} from "vue";
import {getInitStatus} from "@/utils/cache";
import {notification} from "ant-design-vue";
import FormSite from "./FormSite.vue";

const { t } = useI18n();
const router = useRouter();

const store = useStore<{ Zentao: ZentaoData }>();

const sites = computed<any[]>(() => store.state.Zentao.sites);
const products = computed<any>(() => store.state.Zentao.products);
console.log(1111,sites.value);
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
  if(item.key == -1){
      showCreateSiteModal.value = true;
  }
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

const showCreateSiteModal = ref(false)
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
.top-line {border-top: 1px dashed var(--color-green)}
</style>
