<template>
  <div :style="colStyle" :class="colClass" class="z-col">
    <slot></slot>
  </div>
</template>

<script setup lang="ts">
import {computed, defineProps, inject} from "vue";
import {ButtonProps} from "@/layouts/MainLayout/components/Button.vue";

export interface ColumnProps {
  width?: string,
  span?: number,
  flex?: number,
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

  const width = typeof(props.width) === "undefined" ? '' : props.width
  const span = typeof(props.span) === "undefined" ? -1 : props.span
  const flex = typeof(props.flex) === "undefined" ? -1 : props.flex

  if (width && width !== '') {
    style.width = width;
  } else if (span === 0) {
    style.display = 'none';
  } else if (flex > 0) {
    style.flex = flex;
  }

  if (gutter.value > 0) {
    style.paddingLeft = gutter.value / 2 + 'px';
    style.paddingRight = gutter.value / 2 + 'px';
  }

  return style
})

</script>

<style lang="less" scoped>

</style>
