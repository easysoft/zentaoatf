<template>
  <div>
    <div class="title strong space-bottom">{{ t("ui_lang") }}</div>
    <Form labelCol="16" wrapperCol="1" class="settting-form">
        <FormItem name="taskId" v-for="item in locales" :key="item" :label="languageLabels[item]">
            <input type="radio" class="language-input" :value="item" v-model="locale" @change="changeLang(item)" />
        </FormItem>
    </Form>
  </div>
</template>

<script setup lang="ts">
import { setI18nLanguage } from "@/config/i18n";
import { useI18n } from "vue-i18n";
import Form from "./Form.vue";
import FormItem from "./FormItem.vue";

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

<style>
.language-input {
  min-width: 60px;
}
.settting-form{
    display: flex;
}
</style>