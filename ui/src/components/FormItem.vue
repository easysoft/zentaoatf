<template>
  <div
    class="form-item"
    :class="{
      'form-item-single': !multiLine,
      'form-item-multi-lines': multiLine,
      'form-item-divider': divider && !outline,
      'form-item-disabled': disabled,
      'form-item-outline': outline,
    }"
    @click="$emit('click')"
  >
    <div
      :class="[
        'form-item-label',
        labelClass,
        {
          'form-item-label-single': !multiLine,
          'form-item-label-multi-lines': multiLine,
          'form-item-label-outline': outline,
          'form-item-label-require-mark': requiredMark,
        }
      ]"
      :style="multiLine ? undefined : {width: labelWidth ?? '60px'}"
    >
      <div v-if="label" class="form-item-label-text">
        <text>{{label}}</text>
      </div>
      <slot name="label" />
    </div>
    <div class="form-item-content">
      <div
        class="form-item-field"
        :class="{
          'form-item-field-single': !multiLine,
          'form-item-field-outline': outline,
        }"
      >
        <div
          v-if="leadingIcon"
          class="form-item-leading-icon"
          @click="$emit('clickLeadingIcon')"
        >
          <Icon :icon="leadingIcon" />
        </div>
        <slot />
        <div
          v-if="trailingIcon"
          class="form-item-trailing-icon"
          @click="$emit('clickTrailingIcon')"
        >
          <Icon :icon="trailingIcon" />
        </div>
      </div>
      <div v-if="info" class="form-item-error">
        <template v-for="(arr, key) in info" :key="key">
          <div v-for="(item, index) in arr" :key="index" class="text-red">
            {{item}}
          </div>
        </template>
      </div>
      <slot name="help-text">
        <div v-if="helpText" class="form-item-help-text" :class="helpTextClass ?? 'text-gray'">{{helpText}}</div>
      </slot>
    </div>
  </div>
</template>

<script setup lang="ts">
import {defineProps, computed, ref} from "vue";
import Icon from '@/components/Icon.vue';

export interface FormItemProps {
    name?: string,
    label?: string,
    labelClass?: string,
    requiredMark?: boolean,
    labelWidth?: string,
    leadingIcon?: string,
    trailingIcon?: string,
    helpText?: string,
    helpTextClass?: string,
    multiLine?: boolean,
    divider?: boolean,
    disabled?: boolean,
    outline?: boolean,
    info?: string | Record<string, string>,
}

defineProps<FormItemProps>();
</script>

<style lang="less" scoped>
.form-item {
  position:    relative;
  padding:    var(--space-sm) var(--space-base);
}
.form-item + .form-item {
  margin-top: var(--space-vert-xs);
}
.form-item-single {
  display:        flex;
  flex-direction: row;
  align-items:    center;
}
.form-item-label {
  display:         flex;
  flex-direction:  row;
  align-items:     center;
  padding-right:   var(--space-base);
  padding-top:     var(--space-sm);
  align-self:      start;
}
.form-item-label-multi-lines.form-item-label-outline {
  padding: var(--space-vert-sm) 0;
}
.form-item-label-single {
  flex:         none;
  margin-right: var(--space-horz-base);
}
.form-item-label-require-mark::after {
  content:     '*';
  color:       var(--color-red);
  font-size:   var(--font-size-xl);
  position:    relative;
  top:         3px;
  margin-left: 4px;
  line-height: var(--size-line-base);
}
.form-item-content {
  flex: auto
}
.form-item-field {
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: flex-start;
  position: relative;
  gap: var(--space-base);
}
.form-item-field:deep(input[type="text"]),
.form-item-field:deep(input[type="number"]),
.form-item-field:deep(input[type="password"]),
.form-item-field:deep(select),
.form-item-field:deep(.select),
.form-item-field:deep(textarea)  {
  flex: auto;
}
.form-item-field-single {
  flex: auto;
}
.form-item-divider:after {
  content:    ' ';
  display:    block;
  position:   absolute;
  left:       var(--space-horz-base);
  bottom:     0;
  right:      var(--space-horz-base);
  height:     var(--border-size);
  box-shadow: 0 1px 0 0 var(--color-alpha-10);
}
.form-item-disabled {
  opacity: var(--opacity-disabled);
  cursor:  not-allowed;
}
.form-item-field-outline {
  padding: var(--space-vert-xs) var(--space-horz-base);
  /* padding-right: var(--space-horz-xs); */
  box-shadow: inset 0 0 0 1px var(--color-alpha-10);
  text-align: left;
  justify-content: stretch;
  min-height: var(--size-line-lg);
  border-radius: 2px;
}
.form-item-field-outline.form-item-field-focused {
  box-shadow: inset 0 0 0 2rpx var(--color-primary);
}
.form-item-leading-icon {
  flex: none;
  margin-right: var(--space-horz-sm);
  display: flex;
  align-items: center;
}
.form-item-trailing-icon {
  position: absolute;
  top: 0;
  right: var(--space-horz-sm);
  bottom: 0;
  pointer-events: none;
  display: flex;
  justify-content: center;
  align-items: center;
  opacity: .7;
}
.form-item-error,
.form-item-help-text {
  width: 100%;
  margin-top: var(--space-sm)
}

.form-item-field:deep(input[type="text"]),
.form-item-field:deep(input[type="number"]),
.form-item-field:deep(input[type="password"]),
.form-item-field:deep(textarea),
.form-item-field:deep(select) {
  display: block;
  box-sizing: border-box;
  border: 1px solid var(--input-border-color);
  line-height: 1;
  padding: var(--space-sm) var(--space-base);
  width: 100%;
  background-color: var(--color-canvas);
  -webkit-appearance: none;
  min-height: var(--input-height);
  transition: border .1s, box-shadow .1s;
  box-shadow: none;
  color: inherit;
  border-radius: var(--input-border-radius);

  &:hover {
    border-color: var(--input-border-color-hover);
  }

  &:focus {
    border-color: var(--color-focus);
    outline: none;
    box-shadow: 0 0 0 2px var(--color-focus-pale);
  }

  &[readonly],
  &[disabled], .disabled {
    opacity: 1!important;
    background-color: var(--input-back-color-disabled);
  }
}

.form-item-field:deep(.select) {
  position: relative;
  width: 100%;

  > select {
    outline: none;

    &:not([multiple]) {padding-right: 25px;}

    &[multiple] {
      max-height: 75px;
      overflow-y: auto;
    }
  }

  &:not(.multiple):after {
    content: ' ';
    display: block;
    position: absolute;
    right: 10px;
    top: 11px;
    width: 0;
    height: 0;
    border-style: solid;
    border-width: (7px) (5px) 0 (5px);
    border-color: var(--input-border-color) transparent transparent transparent;
    pointer-events: none;
  }
}
</style>
