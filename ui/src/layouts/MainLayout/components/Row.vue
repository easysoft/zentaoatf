<template>
  <div class="z-row">
    <slot></slot>
  </div>
</template>

<script setup lang="ts">
import {computed, defineProps} from "vue";
import {ColumnProps} from "@/layouts/MainLayout/components/Column.vue";

export interface RowProps {
  type?: string,
}

const props = defineProps<RowProps>();

</script>

<style lang="less" scoped>

</style>
