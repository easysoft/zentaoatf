<template>
    <Button id="proxyMenuToggle"
            :label="currProxy.id == 0 ? t('local_proxy') : currProxy.name"
            class="rounded border lighten-16"
            suffix-icon="caret-down"/>

  <DropdownMenu
      toggle="#proxyMenuToggle"
      class="padding-0-bottom"
      :items="proxies"
      keyName="id"
      :checkedKey="currProxy.id"
      @click="selectProxy"
      :replaceFields="replaceFields"
  >
  </DropdownMenu>
</template>

<script setup lang="ts">
import Button from '@/components/Button.vue';
import DropdownMenu from '@/components/DropdownMenu.vue';
import {useI18n} from "vue-i18n";
import {useStore} from "vuex";
import {computed, onMounted, watch, ref} from "vue";
import { ProxyData } from "@/store/proxy";

const { t } = useI18n();
const store = useStore<{ proxy: ProxyData }>();

store.dispatch("proxy/fetchProxies");
const proxies = computed<any[]>(() => {
    const proxies = store.state.proxy.proxies
    proxies.push({id:0, name: t('local'), path: ''})
    return proxies
});
const currProxy = computed<any>(() => store.state.proxy.currProxy);

onMounted(() => {
  console.log('onMounted')
})


const selectProxy = (item): void => {
  console.log('selectProxy', item.item.id)
  store.dispatch('proxy/selectProxy', {currProxyId: item.item.id})
}

const replaceFields = {
  key: 'id',
  title: 'name',
}

</script>

<style>
.top-line {border-top: 1px dashed var(--color-green)}
</style>
