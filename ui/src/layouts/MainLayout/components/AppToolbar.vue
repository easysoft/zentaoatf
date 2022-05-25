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
const { t, locale } = useI18n();

let  language = locale.value
const switchLanguage = () => {
  console.info("switchLanguage", getLocale(), language, defaultLang)
  language = (language === defaultLang ? 'en-US' : defaultLang)
  setI18nLanguage(language, true); 
}

</script>

<style lang="less" scoped>
</style>
