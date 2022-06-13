<template>
  <div class="tab-page canvas">
    <component :is="PageTabComponent[tab.type] || PageTabComponent.unknown" :tab="tab" ref="pageRef"/>
  </div>
</template>

<script setup lang="ts">
import {defineProps, defineExpose, ref} from "vue";
import {PageTab} from "@/store/tabs";
import TabPageResult from '@/views/result/TabPageResult.vue';
import TabPageScript from '@/views/script/TabPageScript.vue';
import TabPageExecUnit from '@/views/exec/TabPageExecUnit.vue';
import TabPageUnknown from './TabPageUnknown.vue';
import {useI18n} from "vue-i18n";
import { SaveFilledIconType } from "@ant-design/icons-vue/lib/icons/SaveFilled";

const {t} = useI18n();

const PageTabComponent = {
  script: TabPageScript,
  result: TabPageResult,
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

<style lang="less" scoped>
.tab-page {
  height: calc(100% - 30px);
  position: relative;
}
</style>
