<template>
  <div class="tree">
    <TreeNode
      v-for="data in dataList"
      :key="data.id"
      class="tree-node-root"
      v-bind="_convertNodeData(data, undefined, 0)"
      :childrenConverter="_convertNodeData"
      @click="_handleClick"
      @toggle="_handleToggle"
      @clickToolbar="emit('clickToolbar', $event)"
      @check="_handleCheck"
    />
    <slot />
  </div>
</template>

<script setup lang="ts">
import { defineProps, withDefaults, defineExpose, reactive, ref, computed, watch, defineEmits } from 'vue';
import TreeNode, {TreeNodeData} from './TreeNode.vue';

export interface TreeProps {
    data: TreeNodeData | TreeNodeData[],
    nodeConverter?: (item: TreeNodeData) => TreeNodeData,
    rootIcon?: string,
    rootCollapsedIcon?: string,
    nodeIcon?: string,
    nodeCollapsedIcon?: string,
    leafIcon?: string,
    inheritStyleFromParent?: boolean,
    checkable?: boolean,
    multiSelect?: boolean,
    defaultCollapsedMap?: Record<string, boolean>,
    defaultActiveID?: string,
    defaultCheckedMap?: Record<string, boolean | 'indeterminate'>,
    toggleOnClick?: boolean,
}

const props = withDefaults(defineProps<TreeProps>(), {
    toggleOnClick: true,
    rootIcon: 'folder-open-filled',
    rootCollapsedIcon: 'folder-filled',
    nodeIcon: 'folder-open-filled',
    nodeCollapsedIcon: 'folder-filled',
    leafIcon: 'file-text',
    inheritStyleFromParent: true,
    checkable: undefined
});

const collapsedMap = reactive(props.defaultCollapsedMap ?? {});
const activeID = ref(props.defaultActiveID ?? '');
const checkedMap = reactive(props.defaultCheckedMap ?? {});

const dataList = computed(() => {
    if (Array.isArray(props.data)) {
        return props.data;
    }
    return [props.data];
});

const treeMap = computed(() => {
    const map: {[id: string]: {id: string, children?: string[], parent?: string, checked?: boolean | 'indeterminate'}} = {};
    function generateMap(data: TreeNodeData, parent?: string) {
        const {id} = data;
        const item: {id: string, children?: string[], parent?: string, checked?: boolean | 'indeterminate'} = {id: id, parent, checked: data.checked};
        if (data.children) {
            item.children = [];
            data.children.forEach(child => {
                if (item.children) {
                    item.children.push(child.id);
                }
                generateMap(child, id);
            });
        }
        map[id] = item;
        return item;
    }
    dataList.value.forEach(data => generateMap(data));
    return map;
});

const emit = defineEmits<{
    (type: 'collapse', event: {collapsed: Record<string, boolean>}) : void,
    (type: 'check', event: {checked: Record<string, boolean | 'indeterminate'>}) : void,
    (type: 'active', event: {activeID: string}) : void,
    (type: 'clickToolbar', event: {node: TreeNodeData, parent?: TreeNodeData, event: any}) : void,
}>();

function getNodeCheckState(id: string): boolean | 'indeterminate' {
    let checked: boolean | 'indeterminate' | undefined = checkedMap[id];
    if (checked === undefined) {
        checked = treeMap.value[id]?.checked;
    }
    if (checked === undefined) {
        checked = false;
    }
    return checked;
}

function _convertNodeData(node: TreeNodeData, parent: TreeNodeData | undefined, index: number) {
    node = {
      ...node,
      level: parent ? ((parent.level ?? 0) + 1) : 0
    };
    if (collapsedMap[node.id] !== undefined) {
        node.collapsed = collapsedMap[node.id];
    }
    if (checkedMap[node.id] !== undefined) {
        node.checked = checkedMap[node.id];
    }
    node.selected = node.id === activeID.value;
    if (node.icon === undefined) {
        if (node.children) {
            node.icon = parent ? props.nodeIcon : props.rootIcon;
            if (node.collapsed) {
                if(parent) {
                    if (props.nodeCollapsedIcon) {
                        node.icon = props.nodeCollapsedIcon;
                    }
                } else {
                    if (props.rootCollapsedIcon) {
                        node.icon = props.rootCollapsedIcon;
                    }
                }
            }
        } else {
            node.icon = props.leafIcon;
        }
    }

    if (parent && props.inheritStyleFromParent) {
        ['iconClass', 'iconStyle', 'titleClass', 'titleStyle', 'collapsedIconClass', 'collapsedIconStyle', 'expandedIconClass', 'expandedIconStyle', 'toolbarClass', 'toolbarStyle', 'collapsedIcon', 'expandedIcon', 'toolbarItems', 'toolbarShowOnHover', 'indent'].forEach(propName => {
            if (node[propName] === undefined && parent[propName]) {
                node[propName] = parent[propName];
            }
        });
    }

    if (props.checkable !== undefined) {
        node.checkable = props.checkable;
    }
    return node;
}

function _handleClick(event) {
    const {node} = event;
    activeID.value = node.id;
    if (props.toggleOnClick) {
        collapsedMap[node.id] = !(collapsedMap[node.id] === undefined ? node.collapsed : collapsedMap[node.id]);
    }
    emit('active', {activeID: node.id});
}

function _handleToggle(event) {
    const {node} = event;
    collapsedMap[node.id] = !(collapsedMap[node.id] === undefined ? node.collapsed : collapsedMap[node.id]);
}

function _updateParentCheckState(id) {
    const parentID = treeMap.value[id]?.parent;
    if (parentID === undefined) {
        return;
    }
    let parentInfo = treeMap.value[parentID];
    if (!parentInfo || !parentInfo.children) {
        return;
    }

    const currentState = checkedMap[id];
    if (currentState === 'indeterminate') {
        checkedMap[parentID] = currentState;
    } else if (currentState) {
        const isAllChildChecked = parentInfo.children.every(x => getNodeCheckState(x) === true);
        checkedMap[parentID] = isAllChildChecked ? true : 'indeterminate';
    } else {
        const isAllChildUnchecked = parentInfo.children.every(x => getNodeCheckState(x) === false);
        checkedMap[parentID] = isAllChildUnchecked ? false : 'indeterminate';
    }
    _updateParentCheckState(parentID);
}

function _updateChildrenCheckState(id) {
    const node = treeMap.value[id];
    if (!node || !node.children) {
        return;
    }
    const currentState = checkedMap[id];
    node.children.forEach(childId => {
        checkedMap[childId] = currentState;
        _updateChildrenCheckState(childId);
    });
}

function _handleCheck(event) {
    const {node} = event;
    const {id} = node;
    const oldCheck = getNodeCheckState(id);
    checkedMap[id] = oldCheck === 'indeterminate' ? true : !oldCheck;
    _updateParentCheckState(id);
    _updateChildrenCheckState(id);
}

function collapse(id: string, includeChildren = false) {
    const item = treeMap.value[id];
    if (!item) {
        return;
    }
    collapsedMap[id] = true;
    if (includeChildren && item.children) {
        item.children.forEach(child => {
            collapse(child, includeChildren);
        });
    }
}

function expand(id: string, includeChildren = false) {
    const item = treeMap.value[id];
    if (!item) {
        return;
    }
    collapsedMap[id] = false;
    if (includeChildren && item.children) {
        item.children.forEach(child => {
            expand(child, includeChildren);
        });
    }
}

function collaseAll() {
    Object.keys(treeMap.value).forEach(id => {collapsedMap[id] = true;});
    console.log('collaseAll', {...collapsedMap});
}

function expandAll() {
    Object.keys(collapsedMap).forEach(id => {
        collapsedMap[id] = false;
    });
}

function isAllCollapsed() {
    return Object.keys(treeMap.value).every(id => collapsedMap[id]);
}

function toggleAllCollapsed(toggle?: boolean) {
    if (toggle === undefined) {
        toggle = !isAllCollapsed();
    }
    if (toggle) {
        collaseAll();
    } else {
        expandAll();
    }
}

defineExpose({
    collapse,
    expand,
    collaseAll,
    expandAll,
    collapsedMap,
    activeID,
    checkedMap,
    isAllCollapsed,
    toggleAllCollapsed
});

watch(collapsedMap, () => {
    emit('collapse', {collapsed: collapsedMap});
});
watch(checkedMap, () => {
    emit('check', {checked: checkedMap});
});
</script>

<style scoped>
.tree :deep(.tree-node.has-children > .tree-node-item > .tree-node-icon) {
  color: var(--tree-toggle-icon-color);
  opacity: 1;
}
.tree-node-root > :deep(.tree-node-item > .tree-node-icon) {
  position: relative;
}
.tree-node-root > :deep(.tree-node-item > .tree-node-title) {
  font-weight: bold;
}
.tree-node-root > :deep(.tree-node-item > .tree-node-icon:before) {
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
