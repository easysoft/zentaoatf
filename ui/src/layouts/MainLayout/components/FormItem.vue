<template>
  <Row class="z-form-item" :class="[size]">
    <Col :width="labelWidth" :class="[labelCls]" class="z-form-item-label">
      {{label}}
    </Col>
    <Col :width="wrapperWidth" :class="wrapperCls" class="z-form-item-wrapper">
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
import Col from "./Col.vue";
import {ContrlSize} from "@/utils/enum";

export interface FormItemProps {
  name?: string,
  label?: string,
  size?: ContrlSize,

  labelCol?: string,
  wrapperCol?: string,
}

const props = defineProps<FormItemProps>();

console.log(props)

let errors = inject('errors');

const size = ref(props.size)

const labelWidth = computed(() => {
  return getWidth(props.labelCol);
})
const labelCls = computed(() => {
  return getCls(props.labelCol);
})

const wrapperWidth = computed(() => {
  return getWidth(props.wrapperCol);
})
const wrapperCls = computed(() => {
  return getCls(props.wrapperCol);
})

const getWidth = (val: string) => {
  if (!val) return undefined

  val += ''
  if (val.indexOf('px') > 0) {
    return val;
  }

  return undefined;
}
const getCls = (val: string) => {
  if (!val) return undefined

  val += ''
  if (val.indexOf('px') < 0) {
    return [`z-col-${val}`]
  }

  return undefined;
}

</script>

<style lang="less">
.z-form-item {
  .z-form-item-label {
    padding-right: 5px;
    text-align: right;
    line-height: 32px;
  }

  &.small {
    .z-form-item-label {
      line-height: 25px;
    }

    .z-form-item-wrapper {
      .z-form-item-control {
        input {
          height: 25px;
          padding: 2px 5px;
        }
      }
    }
  }
}
</style>
