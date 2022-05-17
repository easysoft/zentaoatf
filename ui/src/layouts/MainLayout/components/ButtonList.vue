<template>
  <div class="btn-list" :style="`gap:${gap ?? 8}px`">
    <template v-if="buttonPropsList">
      <Button
        v-for="({key, ...btnProps}) in buttonPropsList"
        :key="key"
        :data-key="key"
        v-bind="btnProps"
        @click="_handleButtonClick"
      />
    </template>
    <slot />
  </div>
</template>

<script setup lang="ts">
import {useI18n} from "vue-i18n";
const { t,te } = useI18n();

import {defineProps, computed, defineEmits} from 'vue';
import Button, {ButtonProps} from './Button.vue';

export type ButtonListItemProps = ButtonProps;

const props = defineProps<{
    buttons?: ButtonProps[] | Record<string, any>[],
    replaceFields?: Record<string, string>,
    defaultBtnClass?: string,
    defaultIconClass?: string,
    defaultIconSize?: string | number,
    defaultSuffixIconClass?: string,
    defaultSuffixIconSize?: string | number,
    defaultLabelIconClass?: string,
    gap?: number
}>();

const buttonPropsList = computed(() => {
    if (!props.buttons) {
        return null;
    }
    return props.buttons.map((x, i) => {
        let item: (ButtonProps | Record<string, any>) & {key: string | number | symbol};
        
        if (props.replaceFields && Button.props) {
            item = Object.keys(Button.props).reduce((item, propName) => {
                const replacePropName = props.replaceFields ? props.replaceFields[propName] : null;
                item[propName] = x[typeof replacePropName === 'string' ? replacePropName : propName];
                return item;
            }, {key: x.key !== undefined ? x.key : i});
        } else {
            item = {
                key: i,
                'class': props.defaultBtnClass,
                iconClass: props.defaultIconClass,
                suffixIconClass: props.defaultSuffixIconClass,
                labelClass: props.defaultLabelIconClass,
                iconSize: props.defaultIconSize,
                suffixIconSize: props.defaultSuffixIconSize,
                ...x
            };
        }

        if(te(item.hint)){
            item.hint = t(item.hint)
        }
        if(item.hintI18n!=undefined){
            item.hint += t(item.hintI18n)
        }
        return item;
    });
});

const emit = defineEmits<{(type: 'click', event: {originalEvent: Event, key: string | number | symbol}) : void}>();

function _handleButtonClick(event) {
    emit('click', event);
}
</script>

<style scoped>
.btn-list {
  display: flex;
  flex-direction: row;
  align-items: center;
}
</style>
