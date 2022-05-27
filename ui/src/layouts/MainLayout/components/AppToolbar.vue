<template>
  <Toolbar>
<!--    <Button class="rounded pure" icon="search" :hint="t('search')" />-->
    <Button class="rounded pure" icon="settings" :hint="t('settings')" @click="openSettings" />
  </Toolbar>
</template>

<script setup lang="ts">
import Button from './Button.vue';
import Toolbar from './Toolbar.vue';
import {useI18n} from "vue-i18n";
import {useStore} from "vuex";
import {ZentaoData} from "@/store/zentao";

const { t, locale } = useI18n();
const store = useStore<{ Zentao: ZentaoData }>();

const openSettings = () => {
    store.dispatch('tabs/open', {
        id: 'settings',
        title: t('settings'),
        type: 'settings',
    });
}

</script>
