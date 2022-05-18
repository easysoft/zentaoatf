<template>
  <Row class="z-form-item" :class="[size]">
    <Col :width="labelWidth" :class="[labelCls]" class="z-form-item-label">
      {{label}}
    </Col>
    <Col :width="wrapperWidth" :class="wrapperCls" class="z-form-item-wrapper">
      <div class="z-form-item-control">
        <slot></slot>
        <span v-if="errorMap.required" class="z-err tips">*</span>
      </div>
      <div class="z-form-item-error">
        <template v-for="(val, key) in errorMap" :key="key">
          <div v-for="(item, index) in val" :key="index" class="z-err">
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
console.log(props)

let labelCol = inject('labelCol') + '';
let wrapperCol = inject('wrapperCol') + '';

const size = ref(props.size)
const errorMap = computed(() => {
  return props.info ? props.info : [];
})

const labelWidth = computed(() => {
  return getWidth(labelCol);
})
const labelCls = computed(() => {
  return getCls(labelCol);
})

const wrapperWidth = computed(() => {
  return getWidth(wrapperCol);
})
const wrapperCls = computed(() => {
  return getCls(wrapperCol);
})

const getWidth = (val: string) => {
  if (!val) return undefined

  val += ''
  if (val.indexOf('px') > 0) {
    return val;
  }

  return val;
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
  margin-bottom: 5px;

  button {
    margin-right: 5px;
  }

  .z-form-item-label {
    padding-right: 5px;
    text-align: right;
    line-height: 28px;
  }
  .z-form-item-wrapper {
    .z-form-item-control {
      input, select {
        width: 90%;
        height: 28px;
        vertical-align: middle;
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
