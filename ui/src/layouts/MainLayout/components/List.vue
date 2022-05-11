<template>
  <div class="list" :class="{compact, divider}">
    <template v-if="items">
      <ListItem
        v-for="({key, ...btnProps}) in itemList"
        :key="key"
        v-bind="btnProps"
      />
    </template>
    <slot />
  </div>
</template>

<script setup lang="ts">
import { defineProps, computed } from 'vue';
import ListItem, { ListItemProps } from './ListItem.vue';

const props = defineProps<{
    compact?: boolean,
    divider?: boolean,
    items?: ListItemProps[] | Record<string, any>[],
    replaceFields?: Record<string, string>, // {title: 'name'}
}>();

const itemList = computed(() => {
    if (!props.items) {
        return null;
    }
    return props.items.map((x, i) => {
        if (props.replaceFields && ListItem.props) {
            return Object.keys(ListItem.props).reduce((item, propName) => {
                const replacePropName = props.replaceFields ? props.replaceFields[propName] : null;
                item[propName] = x[typeof replacePropName === 'string' ? replacePropName : propName];
                return item;
            }, {key: i})
        }
        return {
            key: i,
            ...x
        };
    });
});

</script>
