<template>
  <Toolbar>
<!--    <Button class="rounded pure" icon="search" :hint="t('search')" />-->
    <Button class="rounded pure" icon="local-language" :hint="t('ui_lang')" @click="switchLanguage" />
    <Button class="rounded pure" icon="settings" :hint="t('settings')" @click="switchLanguage" />
  </Toolbar>
</template>

<script setup lang="ts">
import Button from './Button.vue';
import Toolbar from './Toolbar.vue';
import {setI18nLanguage} from "@/config/i18n";
import { getLocale, defaultLang } from "@/utils/i18n";
import {useI18n} from "vue-i18n";
import {useStore} from "vuex";
import {ZentaoData} from "@/store/zentao";

const { t, locale } = useI18n();
const store = useStore<{ Zentao: ZentaoData }>();

let  language = locale.value
const switchLanguage = () => {
<<<<<<< HEAD
//   console.info("switchLanguage", getLocale(), language, defaultLang)
//   language = (language === defaultLang ? 'en-US' : defaultLang)
//   setI18nLanguage(language);
    store.dispatch('tabs/open', {
        id: 'settings',
        title: t('settings'),
        type: 'settings',
    });
=======
  console.info("switchLanguage", getLocale(), language, defaultLang)
  language = (language === defaultLang ? 'en-US' : defaultLang)
  setI18nLanguage(language, true); 
>>>>>>> d2fef79bec49f3055e18cc844f66d303785a70ab
}

</script>

<style lang="less" scoped>
</style>
