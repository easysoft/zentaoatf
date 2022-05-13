<template>
  <div class="tab-page canvas">
    <component :is="PageTabComponent[tab.type] || PageTabComponent.unknown" :tab="tab" />
  </div>
</template>

<script setup lang="ts">
import { defineProps } from "vue";
import { PageTab } from "@/store/tabs";
import TabPageResult from './TabPageResult.vue';
import TabPageScript from './TabPageScript.vue';
import TabPageSettings from './TabPageSettings.vue';
import TabPageSites from './TabPageSites.vue';
import TabPageUnknown from './TabPageUnknown.vue';
import {useI18n} from "vue-i18n";
const { t } = useI18n();

const PageTabComponent = {
    script: TabPageScript,
    settings: TabPageSettings,
    result: TabPageResult,
    sites: TabPageSites,
    unknown: TabPageUnknown,
};

defineProps<{
    tab: PageTab
}>();
</script>
