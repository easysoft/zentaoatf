<template>
  <div class="tab-page canvas">
    <component :is="PageTabComponent[tab.type] || PageTabComponent.unknown" :tab="tab" ref="pageRef"/>
  </div>
</template>

<script setup lang="ts">
import {defineProps, defineExpose, ref} from "vue";
import {PageTab} from "@/store/tabs";
import TabPageResult from './TabPageResult.vue';
import TabPageScript from './TabPageScript.vue';
import TabPageSettings from './TabPageSettings.vue';
import TabPageSites from './TabPageSites.vue';
import TabPageExecUnit from './TabPageExecUnit.vue';
import TabPageUnknown from './TabPageUnknown.vue';
import {useI18n} from "vue-i18n";
import { SaveFilledIconType } from "@ant-design/icons-vue/lib/icons/SaveFilled";

const {t} = useI18n();

const PageTabComponent = {
  script: TabPageScript,
  settings: TabPageSettings,
  result: TabPageResult,
  sites: TabPageSites,
  execUnit: TabPageExecUnit,
  unknown: TabPageUnknown,
};

const pageRef = ref<
  InstanceType< typeof TabPageScript> | 
  InstanceType< typeof TabPageScript> | null
  >(null)
const save = () => {
  if (typeof TabPageScript == typeof pageRef.value) {
    pageRef.value?.save();
  }
}

const props = defineProps<{
  tab: PageTab,
}>();

defineExpose({
  save
});
</script>
