<template>
  <div :style="colStyle" :class="colClass" class="z-col">
    <slot></slot>
  </div>
</template>

<script setup lang="ts">
import {computed, defineProps, inject} from "vue";
import {ButtonProps} from "@/layouts/MainLayout/components/Button.vue";

export interface ColumnProps {
  width?: number,
  span?: number,
  offset?: number,
}

const props = defineProps<ColumnProps>();
let gutter = inject('gutter');

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

  console.log('123', gutter.value)

  if (gutter.value > 0) {
    style.marginLeft = gutter.value / 2 + 'px';

    style.marginRight = gutter.value / 2 + 'px';
  }

  return style
})

</script>

<style lang="less" scoped>

</style>
