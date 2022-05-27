<template>
  <Row class="z-form-item" :class="[size]">
    <Col :span="labelCol" class="z-form-item-label">
      {{label}}
    </Col>
    <Col :span="wrapperCol" class="z-form-item-wrapper">
      <div class="z-form-item-control">
        <slot></slot>
        <span v-if="errorMap.required" class="z-err tips">*</span>
      </div>
      <div class="z-form-item-error">
        <template v-for="(arr, key) in errorMap" :key="key">
          <div v-for="(item, index) in arr" :key="index" class="z-err">
            {{item}}
          </div>
        </template>
      </div>
    </Col>
  </Row>
</template>

<script setup lang="ts">
import {computed, defineProps, inject, ref} from "vue";
import Row from "./Row.vue";
import Col from "./Col.vue";
import {ContrlSize} from "@/utils/enum";
import {useI18n} from "vue-i18n";
const { t } = useI18n();

export interface FormItemProps {
  name?: string,
  label?: string,
  info?: any,
  size?: ContrlSize,
}

const props = defineProps<FormItemProps>();

let labelColStr = inject('labelCol') as string
let wrapperColStr = inject('wrapperCol') as string

let labelCol = parseInt(labelColStr)
let wrapperCol = parseInt(wrapperColStr)

const size = ref(props.size)
const errorMap = computed(() => {
  return props.info ? props.info : {};
})

</script>

<style lang="less">
.z-form-item {
  margin-bottom: 5px;

  button {
    margin-right: 5px;
  }

  .z-form-item-label {
    padding-right: 10px;
    text-align: right;
    line-height: 28px;
  }
  .z-form-item-wrapper {
    .z-form-item-control {
      height: 100%;
      input, select {
        width: 90%;
        height: 28px;
        vertical-align: middle;
      }
      input[type="checkbox"] {
        width: auto;
        height: 14px;
        margin-top: 6px;
      }
      label {
        display: inline-block;
        padding: 0 3px;
        line-height: 28px;
        vertical-align: middle;
      }
    }
  }
  .tips {
    display: inline-block;
    margin-left: 5px;
    line-height: 28px;
    vertical-align: middle;
  }

  &.small {
    .z-form-item-label {
      line-height: 25px;
    }
    .z-form-item-wrapper {
      .z-form-item-control {
        input, select {
          height: 25px;
        }
        label {
          line-height: 25px;
        }
      }
    }
    .tips {
      line-height: 25px;
    }
  }
}
</style>
