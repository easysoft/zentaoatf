<template>
  <form class="z-form">
    <slot></slot>
  </form>
</template>

<script setup lang="ts">
import {defineProps, provide} from "vue";
import {FormItemProps} from "@/layouts/MainLayout/components/FormItem.vue";

export interface FormProps {
  msgs?: any[],
  labelCol?: string,
  wrapperCol?: string,
}

const props = defineProps<FormProps>();
console.log(props)

provide('errors', {});
provide('labelCol', props.labelCol);
provide('wrapperCol', props.wrapperCol);



</script>

<style lang="less" scoped>
  .z-form {

  }
</style>
