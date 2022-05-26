<template>
  <div>{{ t("ui_lang") }}</div>
  <div>
    <a-radio-group name="radioGroup" v-model:value="locale">
      <a-radio
        v-for="item in locales"
        @change="changeLang(item)"
        :key="item"
        :value="item"
        :style="radioStyle"
      >
        <span style="margin-left: 5px">{{ languageLabels[item] }}</span>
      </a-radio>
    </a-radio-group>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { setI18nLanguage } from "@/config/i18n";
import { useI18n } from "vue-i18n";

const { t, locale } = useI18n();

const radioStyle = ref({
  display: "block",
  height: "30px",
  lineHeight: "30px",
});

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
  setI18nLanguage(key);
};
</script>