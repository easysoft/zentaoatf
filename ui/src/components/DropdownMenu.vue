<template>
  <div
    class="dropdown-menu"
    :class="menuClass ?? 'layer rounded'"
    :style="menuFinalStyle"
    @click="_handleClickMenu"
    ref="menuRef"
  >
    <template v-if="state.show">
      <List
        :items="items"
        :replaceFields="replaceFields"
        :checkedKey="checkedKey"
        :activeKey="activeKey"
        :keyName="keyName"
        :class="listClass"
        :compact="listCompact"
        :divider="listDivider"
        @click="emit('click', $event)"
      />
      <slot />
    </template>
  </div>
</template>

<script setup lang="ts">
import { withDefaults, defineProps, ref, reactive, onMounted, onUnmounted, computed, defineEmits } from 'vue';
import { useWindowSize, onClickOutside } from '@vueuse/core'
import List from './List.vue';
import {ListItemKey, ListItemProps} from './ListItem.vue';

export type DropdownMenuPosition = 'bottom' | 'left' | 'right' | 'top' | 'bottom-left' | 'bottom-right' | 'top-left' | 'top-right' | 'left-top' | 'left-bottom' | 'right-top' | 'right-bottom' |{left: number, top: number} | ((toggleElement: HTMLElement) => {left: number, top: number});

const props = withDefaults(defineProps<{
    toggle?: object | string,
    triggerEvent?: string,
    position?: DropdownMenuPosition,
    items?: ListItemProps[] | Record<string, any>[],
    keyName?: string,
    checkedKey?: ListItemKey,
    activeKey?: ListItemKey,
    replaceFields?: Record<string, string>, // {title: 'name'}
    listClass?: string,
    listCompact?: boolean,
    listDivider?: boolean,
    defaultShow?: boolean,
    showOnHover?: boolean,
    hideOnClickAway?: boolean,
    hideOnClickMenu?: boolean,
    menuClass?: string,
    limitInWindow?: boolean,
    menuStyle?: object
}>(), {hideOnClickAway: true, hideOnClickMenu: true});

const menuRef = ref<HTMLElement>();
const cleanUpRef = ref();
const state = reactive({show: !!props.defaultShow, showed: false});
const showTimerRef = ref<number>(0);

function _toggle(show?: boolean) {
    if (show === undefined) {
        show = !state.show;
    }
    if (state.show === show) {
        return;
    }
    if (showTimerRef.value) {
        clearTimeout(showTimerRef.value);
        showTimerRef.value = 0;
    }
    state.show = show;
    state.showed = false;
    if (show) {
        showTimerRef.value = setTimeout(() => {
            state.showed = true;
            showTimerRef.value = 0;
        }, 100);
    }
}

function _handleClickMenu(event: MouseEvent) {
    if (!props.hideOnClickMenu || event.target && event.target instanceof HTMLElement && event.target.closest('.not-hide-menu')) {
        return;
    }
    if (state.showed) {
      _toggle(false);
    }
}

const emit = defineEmits<{(type: 'click', event: {originalEvent: Event, key: ListItemKey, item: ListItemProps | Record<string, any>}) : void}>();

onClickOutside(menuRef, event => {
    if (props.hideOnClickAway && state.showed) {
        _toggle(false);
    }
});

function getToggleElement() {
    const {toggle} = props;
    if (!toggle) {
        if (menuRef.value) {
            return menuRef.value.closest('.dropdown')?.querySelector('.dropdown-toggle') as HTMLElement || undefined;
        }
        return undefined;
    }
    if (toggle instanceof HTMLElement) {
        return toggle;
    }
    if (typeof toggle === 'string') {
        return document.querySelector(toggle) as HTMLElement || undefined;
    }
    return undefined;
}

const windowSize = useWindowSize();

const menuFinalStyle = computed(() => {
    if (!state.show) {
        return {display: 'none'};
    }
    const style: Record<string, any> = {
        display: 'block',
        opacity: 1,
        top: '0px',
        left: '0px',
    };
    const element = getToggleElement();
    if (!element || !menuRef.value) {
        return style;
    }
    if (!state.showed) {
        style.opacity = 0;
        return style;
    }
    const {position = 'bottom-left'} = props;
    const bounding = element.getBoundingClientRect();
    const menuBounding = menuRef.value.getBoundingClientRect();
    if (typeof position === 'function') {
        Object.assign(style, position(element));
    } else if (typeof position === 'object') {
        Object.assign(style, position);
    } else if (position === 'bottom') {
        style.top = bounding.bottom;
        style.left = Math.round(bounding.left + (bounding.width / 2) - (menuBounding.width / 2));
    } else if (position === 'bottom-left') {
        style.top = bounding.bottom;
        style.left = bounding.left;
    } else if (position === 'bottom-right') {
        style.top = bounding.bottom;
        style.left = bounding.right - menuBounding.width;
    } else if (position === 'top') {
        style.top = bounding.top - menuBounding.height;
        style.left = bounding.left + (bounding.width / 2) - (menuBounding.width / 2);
    } else if (position === 'top-left') {
        style.top = bounding.top - menuBounding.height;
        style.left = bounding.left;
    } else if (position === 'top-right') {
        style.top = bounding.top - menuBounding.height;
        style.left = bounding.right - menuBounding.width;
    } else if (position === 'left') {
        style.left = bounding.left - menuBounding.width;
        style.top = bounding.top + (bounding.height / 2) - (menuBounding.height / 2);
    } else if (position === 'left-top') {
        style.left = bounding.left - menuBounding.width;
        style.top = bounding.top;
    } else if (position === 'left-bottom') {
        style.left = bounding.left - menuBounding.width;
        style.top = bounding.bottom - menuBounding.height;
    } else if (position === 'right') {
        style.left = bounding.left - menuBounding.width;
    } else if (position === 'right-top') {
        style.left = bounding.right;
        style.top = bounding.top;
    } else if (position === 'right-bottom') {
        style.left = bounding.right;
        style.top = bounding.bottom - menuBounding.height;
    }

    if (props.limitInWindow !== false) {
        style.left = Math.min(windowSize.width.value - menuBounding.width, style.left);
        style.top = Math.min(windowSize.height.value - menuBounding.height, style.top);
    }
    style.maxHeight = `${Math.round(windowSize.height.value - style.top - 10)}px`;
    style.maxWidth = `${Math.round(windowSize.width.value - style.left - 10)}px`;
    style.left = `${Math.round(style.left)}px`;
    style.top = `${Math.round(style.top)}px`;
    style.overflow = 'auto';

    if (props.menuStyle) {
        Object.assign(style, props.menuStyle);
    }

    return style;
});

onMounted(() => {
    const toggleElement = getToggleElement();
    if (!toggleElement) {
        return;
    }

    const triggerEvent = props.triggerEvent ?? 'click';
    const handler = (event: Event) => {
        if (state.showed) {
            _toggle(false);
        } else if (!state.show) {
            _toggle(true);
        }
    };
    toggleElement.addEventListener(triggerEvent, handler, false);
    cleanUpRef.value = () => {
        toggleElement.removeEventListener(triggerEvent, handler, false);
    };
});

onUnmounted(() => {
    if (cleanUpRef.value) {
        cleanUpRef.value();
    }
})
</script>

<style scoped>
.dropdown-menu {
  position: fixed;
  min-width: 100px;
  z-index: 100;
  padding-top: var(--space-sm);
  padding-bottom: var(--space-sm);
}
.dropdown-menu :deep(.list-item) {
  padding-left: var(--space-lg);
  padding-right: var(--space-lg);
}
.dropdown-menu :deep(.list-item.has-checkmark) {
  padding-right: var(--space-sm);
}
</style>
