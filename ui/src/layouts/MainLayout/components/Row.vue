<template>
  <div class="z-row" :style="rowStyle">
    <slot></slot>
  </div>
</template>

<script setup lang="ts">
import {provide, defineProps, ref, computed} from "vue";
import {useI18n} from "vue-i18n";
const { t } = useI18n();

export interface RowProps {
  gutter?: number,
}

const props = defineProps<RowProps>();
const gutter = ref(props.gutter);

provide('gutter', gutter);

const rowStyle = computed(() => {
  const style: Record<string, any> = {};

  if (gutter.value > 0) {
    style.marginLeft = gutter.value / -2 + 'px';
    style.marginRight = gutter.value / -2 + 'px';
  }

  return style
})

</script>

<style lang="less" scoped>

</style>
