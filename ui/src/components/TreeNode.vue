<template>
  <div class="tree-node" :class="{collapsed: isCollapsed, checked, selected, checkable, 'has-children': children, 'is-leaf': !children, 'tree-node-root': root}" :data-id="id">
    <div
      class="tree-node-item"
      :style="`padding-left: ${indent * (level || 0)}px`"
      @mouseenter="_handleMouseEnter"
      @mouseleave="_handleMouseLeave"
      @click.stop="emit('click', {node: props, event: $event})"
    >
      <template v-if="children">
        <Button v-if="isCollapsed" size="sm" :icon="collapsedIcon" :class="collapsedIconClass" :style="collapsedIconStyle" class="tree-node-toggle" @click="emit('toggle', {node: props, event: $event})" />
        <Button v-else size="sm" :icon="expandedIcon" :class="expandedIconClass" :style="expandedIconStyle" class="tree-node-toggle" @click="emit('toggle', {node: props, event: $event})" />
      </template>
      <Button
        v-if="checkable"
        class="tree-node-check rounded pure"
        :class="{indeterminate: checked === 'indeterminate'}"
        :icon="checked === true ? 'checkbox-checked' : 'checkbox-unchecked'"
        size="sm"
        @click="emit('check', {node: props, event: $event})"
      />
      <div class="tree-node-icon">
        <Icon v-if="icon" :icon="icon" :class="iconClass" :style="iconStyle" />
      </div>
      <div class="tree-node-title" :class="titleClass" :style="titleClass">{{title}}</div>
      <Toolbar
        v-if="toolbarItems && showToolbar"
        class="tree-node-toolbar"
        :buttons="typeof toolbarItems === 'function' ? toolbarItems(props) : toolbarItems"
        @click="emit('clickToolbar', {node: props, event: $event})"
      />
    </div>
    <div v-if="children && !isCollapsed" class="tree-node-children">
      <TreeNode
        v-for="(child, index) in children"
        :key="child.id"
        v-bind="childrenConverter ? childrenConverter(child, props, index) : child"
        :childrenConverter="childrenConverter"
        @click="emit('click', {parent: props, ...$event})"
        @toggle="emit('toggle', {parent: props, ...$event})"
        @clickToolbar="emit('clickToolbar', {parent: props, ...$event})"
        @check="emit('check', {parent: props, ...$event})"
      />
    </div>
    <slot />
  </div>
</template>

<script lang="ts">
export default {
    name: 'TreeNode',
    inheritAttrs: false
}
</script>

<script setup lang="ts">
import { defineProps, defineEmits, withDefaults, computed, ref } from 'vue';
import Toolbar, { ToolbarItemProps } from './Toolbar.vue';
import Icon from './Icon.vue';
import Button from './Button.vue';

export interface TreeNodeData {
    id: string,
    title: string,
    level?: number,
    indent?: number,
    root?: boolean,
    titleClass?: string,
    titleStyle?: string | object,
    collapsed?: boolean | ((item: TreeNodeData) => boolean),
    checkable?: boolean,
    checked?: boolean | 'indeterminate',
    selected?: boolean,
    icon?: string,
    iconClass?: string,
    iconStyle?: string | object,
    collapsedIcon?: string,
    collapsedIconClass?: string,
    collapsedIconStyle?: string | object,
    expandedIcon?: string,
    expandedIconClass?: string,
    expandedIconStyle?: string,
    toolbarItems?: ToolbarItemProps[] | ((item: TreeNodeData) => ToolbarItemProps[]),
    toolbarClass?: string,
    toolbarShowOnHover?: boolean,
    children?: TreeNodeData[],
    childrenConverter?: (item: TreeNodeData, parent: TreeNodeData | undefined, index: number) => TreeNodeData,
}

const props = withDefaults(defineProps<TreeNodeData>(), {
    toolbarShowOnHover: true,
    collapsedIcon: 'chevron-right',
    expandedIcon: 'chevron-down',
    collapsedIconClass: 'rounded pure',
    expandedIconClass: 'rounded pure',
    indent: 15
});

const showToolbar = ref(!props.toolbarShowOnHover);

const emit = defineEmits<{
    (type: 'click', event: {node: TreeNodeData, parent?: TreeNodeData, event: any}) : void,
    (type: 'check', event: {node: TreeNodeData, parent?: TreeNodeData, event: any}) : void,
    (type: 'toggle', event: {node: TreeNodeData, parent?: TreeNodeData, event: any}) : void,
    (type: 'clickToolbar', event: {node: TreeNodeData, parent?: TreeNodeData, event: any}) : void,
}>();

const isCollapsed = computed(() => {
    if (typeof props.collapsed === 'function') {
        return props.collapsed(props);
    }
    return !!props.collapsed;
});

function _handleMouseEnter() {
    showToolbar.value = true;
}
function _handleMouseLeave() {
    showToolbar.value = false;
}
</script>

<style scoped>
.tree-node-item {
  display: flex;
  align-items: center;
  flex-direction: row;
  flex-wrap: nowrap;
  height: calc(2em + 2px);
  cursor: pointer;
  position: relative;
}
.tree-node.selected > .tree-node-item,
.tree-node-item:hover {
  background-color: var(--color-darken-1);
}
.tree-node.selected > .tree-node-item:hover {
  background-color: var(--color-darken-2);
}
.tree-node-icon {
  margin-right: var(--space-sm);
  opacity: .7;
  display: flex;
  align-items: center;
  justify-content: center;
}
.tree-node-item:hover > .tree-node-icon {
  opacity: 1;
}
.tree-node-check {
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: -1px;
  margin-left: -3px;
  color: var(--color-primary);
}
.tree-node-check.indeterminate:before {
  content: ' ';
  display: block;
  height: 1px;
  background-color: var(--color-primary);
  left: 34%;
  right: 34%;
  position: absolute;
}
.tree-node-title {
  flex: auto;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.tree-node.is-leaf > .tree-node-item > :first-child {
  margin-left: calc(2em - 2px);
}

.tree-node-root > .tree-node-item > .tree-node-icon {
  position: relative;
}
.tree-node-root > .tree-node-item > .tree-node-title {
  font-weight: bold;
}
.tree-node-root > .tree-node-item > .tree-node-icon:before {
  content: ' ';
  display: block;
  position: absolute;
  width: 6px;
  height: 6px;
  border-radius: 50%;
  left: 0;
  bottom: 1px;
  border: 1px solid var(--color-gray);
  background: var(--color-gray);
  box-shadow: inset var(--color-canvas) 0 0 0 1px, 0 0 0 1px var(--color-canvas);
}
</style>
