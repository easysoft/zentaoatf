<template>
  <div class="lang-settings">
    <div class="title strong space-bottom">{{ t("ui_lang") }}</div>

    <Row>
      <template v-for="(item,index) in locales" :key="index">
        <Col width="60px" class="lang-label">
          {{ languageLabels[item] }}
        </Col>
        <Col width="30px" class="lang-input">
          <input type="radio" :value="item" v-model="locale"  @change="changeLang(item)"/>
        </Col>
      </template>
    </Row>
  </div>
</template>

<script setup lang="ts">
import {setI18nLanguage} from "@/config/i18n";
import {useI18n} from "vue-i18n";
import Row from "@/components/Row.vue";
import Col from "@/components/Col.vue";
import {setLang} from "@/services/settings";

const {t, locale} = useI18n();

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

  setLang(key)
};
</script>

<style lang="less" scoped>
.lang-settings {
  .lang-label {
    display: inline-block;
    margin-right: 5px;
    text-align: right;
  }

  .lang-input {
    line-height: 24px;
  }
}
</style>