<template>
  <div class="list" :class="{compact, divider}">
    <template v-if="items">
      <ListItem
        v-for="({key, ...btnProps}) in itemList"
        :key="key"
        :data-key="key"
        v-bind="btnProps"
        @click="_handleItemClick"
      />

    </template>
    <slot />
  </div>
</template>

<script setup lang="ts">
import { defineProps, computed, defineEmits, ref } from 'vue';
import ListItem, { ListItemProps, ListItemKey } from './ListItem.vue';
import {useI18n} from "vue-i18n";
const { t } = useI18n();

const props = defineProps<{
    compact?: boolean,
    divider?: boolean,
    items?: ListItemProps[] | Record<string, any>[],
    keyName?: string,
    checkedKey?: ListItemKey,
    activeKey?: ListItemKey,
    replaceFields?: Record<string, string>,
}>();

const keyItemMap = ref<Record<NonNullable<ListItemKey>, ListItemProps | Record<string, any>>>({});

const itemList = computed(() => {
    if (!props.items) {
        return null;
    }
    return props.items.map((x, i) => {
        let item: (ListItemProps | Record<string, any>) & {key: NonNullable<ListItemKey>};
        if (props.replaceFields && ListItem.props) {
            item = Object.keys(ListItem.props).reduce((item2, propName) => {
                const replacePropName = props.replaceFields ? props.replaceFields[propName] : null;
                item2[propName] = x[typeof replacePropName === 'string' ? replacePropName : propName];
                return item2;
            }, {key: x.key !== undefined ? x.key : i});
        } else {
            item = {
                key: i,
                ...x
            };
        }
        if (props.keyName && props.keyName !== 'key') {
            item.key = x[props.keyName];
        }
        if (item.key === props.activeKey) {
            item.active = true;
        }

        if (item.key === props.checkedKey) {
            item.checked = true;
        }
        keyItemMap.value[item.key] = x;
        return item;
    });
});

const emit = defineEmits<{(type: 'click', event: {originalEvent: Event, key: ListItemKey, item: ListItemProps | Record<string, any>}) : void}>();

function _handleItemClick(event) {
    emit('click', {
        ...event,
        item: keyItemMap.value[event.key],
    });
}
</script>
