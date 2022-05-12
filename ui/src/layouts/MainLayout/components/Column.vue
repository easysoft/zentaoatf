<template>
  <div :style="colStyle" :class="colClass" class="z-col">
    <slot></slot>
  </div>
</template>

<script setup lang="ts">
import {computed, defineProps} from "vue";
import {ButtonProps} from "@/layouts/MainLayout/components/Button.vue";

export interface ColumnProps {
  width?: number,
  span?: number,
  offset?: number,
}

const props = defineProps<ColumnProps>();

const colClass = computed(() => {
  const classes: string[] = [];

  const span = typeof(props.span) === "undefined" ? -1 : props.span

  if (span > 0) {
    classes.push(`z-col-${span}`);
  }

  if (props.offset > 0) {
    classes.push(`z-col-offset-${props.offset}`);
  }

  return classes;
})

const colStyle = computed(() => {
  const style: Record<string, any> = {};

  const width = typeof(props.width) === "undefined" ? -1 : props.width
  const span = typeof(props.span) === "undefined" ? -1 : props.span

  if (width > 0) {
    style.width = `${width}px`;
  } else if (span === 0) {
    style.display = 'none';
  }

  return style
})

</script>

<style lang="less" scoped>

</style>
