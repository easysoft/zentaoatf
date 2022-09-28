<template>
    <div class="space-left">
    <div class="title space">{{t('remote_proxy')}}：</div>
    <Button id="proxyMenuToggle"
        :label="currProxy.id == 0 || currProxy.id == undefined ? t('local_proxy') : currProxy.name"
        class="rounded border lighten-16"
        suffix-icon="caret-down"/>
    </div>
  <DropdownMenu
      v-if="proxies.length>0"
      toggle="#proxyMenuToggle"
      class="padding-0-bottom"
      :items="proxies"
      keyName="id"
      :checkedKey="currProxy.id == undefined ? 0 : currProxy.id"
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
    const proxiesList = [...store.state.proxy.proxies]
    if(proxiesList.length > 0){
        proxiesList.push({id:0, name: t('local_proxy'), path: 'local'})
    }
    return proxiesList
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
