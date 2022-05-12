<template>
  <Row class="z-form-item">
    <Col :width="labelWidth" :flex="labelFlex" class="z-form-item-label">
      {{label}}
    </Col>
    <Col :width="wrapperWidth" :flex="wrapperFlex" class="z-form-item-content">
      <div class="z-form-item-control">
        <slot></slot>
      </div>
      <div class="z-form-item-error">
          {{errors[name]}}
      </div>
    </Col>
  </Row>
</template>

<script setup lang="ts">
import {computed, defineProps, inject, ref} from "vue";
import Row from "./Row.vue";
import Col from "./Column.vue";

export interface FormItemProps {
  name?: string,
  label?: string,

  labelCol?: string,
  wrapperCol?: string,
}

const props = defineProps<FormItemProps>();

console.log(props)

let errors = inject('errors');

const labelWidth = computed(() => {
  return getWidth(props.labelCol);
})
const labelFlex = computed(() => {
  return getFlex(props.labelCol);
})

const wrapperWidth = computed(() => {
  return getWidth(props.wrapperCol);
})
const wrapperFlex = computed(() => {
  return getFlex(props.wrapperCol);
})

const getWidth = (val: string) => {
  if (!val) return -1

  val += ''
  if (val.indexOf('px') > 0) {
    return val;
  }

  return -1;
}
const getFlex = (val: string) => {
  if (!val) return -1

  val += ''
  if (val.indexOf('px') < 0) {
    return val;
  }

  return -1;
}

</script>

<style lang="less" scoped>
.z-form-item {

}
</style>
