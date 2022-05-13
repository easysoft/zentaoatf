<template>
  <div class="tab-page-sites">
    <div>{{t('site_management')}}</div>
    <div>{{JSON.stringify(tab)}}</div>
  </div>
</template>

<script setup lang="ts">
import { defineProps } from "vue";
import { PageTab } from "@/store/tabs";
import {useI18n} from "vue-i18n";
const { t } = useI18n();

defineProps<{
    tab: PageTab
}>();
</script>
