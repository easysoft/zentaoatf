<template>
  <div class="z-row">
    <slot></slot>
  </div>
</template>

<script setup lang="ts">
import {provide, defineProps, ref} from "vue";
import {ColumnProps} from "@/layouts/MainLayout/components/Column.vue";

export interface RowProps {
  gutter?: number,
}

const props = defineProps<RowProps>();
const gutter = ref(props.gutter);

provide('gutter', gutter);

</script>

<style lang="less" scoped>

</style>
