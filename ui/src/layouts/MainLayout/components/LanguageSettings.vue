<template>
  <div>{{ t("ui_lang") }}</div>
  <div>
    <div  v-for="item in locales" :key="item" >
      <input type="radio" :id="'ui-lang-setting-radio-' + item" :value="item" v-model="locale" @change="changeLang(item)" />
      <label :for="'ui-lang-setting-radio-' + item" > {{languageLabels[item] }}</label>
    </div>
  </div>
</template>

<script setup lang="ts">
import { setI18nLanguage } from "@/config/i18n";
import { useI18n } from "vue-i18n";

const { t, locale } = useI18n();

const locales: string[] = ["zh-CN", "en-US"];
const languageLabels: { [key: string]: string } = {
  "zh-CN": "简体中文",
  "en-US": "English",
};
const languageIcons: { [key: string]: string } = {
  "zh-CN": "🇨🇳",
  "en-US": "🇺🇸",
};

const changeLang = (key): void => {
console.info(key)
  setI18nLanguage(key);
};
</script>
