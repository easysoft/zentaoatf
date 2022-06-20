<template>
  <div class="popper layer rounded" :style="menuFinalStyle" ref="menuRef">
    <template v-if="state.show">
      <div class="tab-group">
        <div
          ref="barRef"
          class="tab-bar"
          :style="{ width: widthRef + 'px' }"
        ></div>
        <div ref="titsRef" class="tab-header" layout="row" layout-wrap>
          <div
            ref="tabRef"
            v-for="(item, index) in tabs"
            :class="[{ active: activeKey == item.key }, 'tab-nav']"
            :key="item.key"
            @click="onTabClick($event, item, index)"
          >
            {{ item.title }}
          </div>
        </div>
        <div class="tab-panel">
          <List
            :items="list"
            :replaceFields="replaceFields"
            :checkedKey="checkedKey"
            :keyName="keyName"
            @click="onListClick($event)"
          />
          <slot></slot>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import {
  withDefaults,
  defineProps,
  ref,
  reactive,
  onMounted,
  onUnmounted,
  computed,
  defineEmits,
} from "vue";
import { useWindowSize, onClickOutside } from "@vueuse/core";
import { ListItemKey, ListItemProps } from "./ListItem.vue";
import List from "@/components/List.vue";

export type PopperPosition =
  | "bottom"
  | "left"
  | "right"
  | "top"
  | "bottom-left"
  | "bottom-right"
  | "top-left"
  | "top-right"
  | "left-top"
  | "left-bottom"
  | "right-top"
  | "right-bottom"
  | { left: number; top: number }
  | ((toggleElement: HTMLElement) => { left: number; top: number });

const props = withDefaults(
  defineProps<{
    toggle?: object | string;
    triggerEvent?: string;
    position?: PopperPosition;
    tabs?: ListItemProps[] | Record<string, any>[];
    list?: ListItemProps[] | Record<string, any>[];
    keyName?: string;
    checkedKey?: ListItemKey;
    checkedTab?: ListItemKey;
    replaceFields?: Record<string, string>; // {title: 'name'}
    defaultShow?: boolean;
  }>(),
  {}
);

let widthRef = ref();
const menuRef = ref<HTMLElement>();
const cleanUpRef = ref();
const state = reactive({ show: !!props.defaultShow, showed: false });
const showTimerRef = ref<number>(0);
const tabs = computed(() => props.tabs);
const list = computed(() => props.list);

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

const emit = defineEmits<{
  (
    type: "tabChanged",
    event: {
      key: ListItemKey;
      item: ListItemProps | Record<string, any>;
    }
  ): void;
  (
    type: "click",
    event: {
      originalEvent: Event;
      key: ListItemKey;
      item: ListItemProps | Record<string, any>;
    }
  ): void;
}>();

onClickOutside(menuRef, (event) => {
  if (state.showed) {
    _toggle(false);
  }
});

function getToggleElement() {
  const { toggle } = props;
  if (!toggle) {
    if (menuRef.value) {
      return (
        (menuRef.value
          .closest(".dropdown")
          ?.querySelector(".dropdown-toggle") as HTMLElement) || undefined
      );
    }
    return undefined;
  }
  if (toggle instanceof HTMLElement) {
    return toggle;
  }
  if (typeof toggle === "string") {
    return (document.querySelector(toggle) as HTMLElement) || undefined;
  }
  return undefined;
}

const windowSize = useWindowSize();

const menuFinalStyle = computed(() => {
  if (!state.show) {
    return { display: "none" };
  }
  const style: Record<string, any> = {
    display: "block",
    opacity: 1,
    top: "0px",
    left: "0px",
  };
  const element = getToggleElement();
  if (!element || !menuRef.value) {
    return style;
  }
  if (!state.showed) {
    style.opacity = 0;
    return style;
  }
  const { position = "bottom-left" } = props;
  const bounding = element.getBoundingClientRect();
  const menuBounding = menuRef.value.getBoundingClientRect();
  if (typeof position === "function") {
    Object.assign(style, position(element));
  } else if (typeof position === "object") {
    Object.assign(style, position);
  } else if (position === "bottom") {
    style.top = bounding.bottom;
    style.left = Math.round(
      bounding.left + bounding.width / 2 - menuBounding.width / 2
    );
  } else if (position === "bottom-left") {
    style.top = bounding.bottom;
    style.left = bounding.left;
  } else if (position === "bottom-right") {
    style.top = bounding.bottom;
    style.left = bounding.right - menuBounding.width;
  } else if (position === "top") {
    style.top = bounding.top - menuBounding.height;
    style.left = bounding.left + bounding.width / 2 - menuBounding.width / 2;
  } else if (position === "top-left") {
    style.top = bounding.top - menuBounding.height;
    style.left = bounding.left;
  } else if (position === "top-right") {
    style.top = bounding.top - menuBounding.height;
    style.left = bounding.right - menuBounding.width;
  } else if (position === "left") {
    style.left = bounding.left - menuBounding.width;
    style.top = bounding.top + bounding.height / 2 - menuBounding.height / 2;
  } else if (position === "left-top") {
    style.left = bounding.left - menuBounding.width;
    style.top = bounding.top;
  } else if (position === "left-bottom") {
    style.left = bounding.left - menuBounding.width;
    style.top = bounding.bottom - menuBounding.height;
  } else if (position === "right") {
    style.left = bounding.left - menuBounding.width;
  } else if (position === "right-top") {
    style.left = bounding.right;
    style.top = bounding.top;
  } else if (position === "right-bottom") {
    style.left = bounding.right;
    style.top = bounding.bottom - menuBounding.height;
  }

  style.maxHeight = `${Math.round(windowSize.height.value - style.top - 10)}px`;
  style.maxWidth = `${Math.round(windowSize.width.value - style.left - 10)}px`;
  style.left = `${Math.round(style.left)}px`;
  style.top = `${Math.round(style.top)}px`;
  style.overflow = "auto";

  return style;
});

onMounted(() => {
  const toggleElement = getToggleElement();
  if (!toggleElement) {
    return;
  }

  const triggerEvent = props.triggerEvent ?? "click";
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
  onMounted(() => {
    // 设置状态线初始化宽度
    widthRef.value = tabRef.value.clientWidth;
  });
});

const activeKey = computed(() => props.checkedTab);
const tabRef = ref(null);
const onTabClick = (event, tab, index) => {
  emit("tabChanged", tab);
};

const onListClick = (e) => {
  if (state.showed) {
    _toggle(false);
  }
  emit("click", e);
};

onUnmounted(() => {
  if (cleanUpRef.value) {
    cleanUpRef.value();
  }
});
</script>

<style  lang="less" scoped>
.popper {
  position: fixed;
  min-width: 100px;
  z-index: 100;
  padding: 10px;
}

.tab-group {
  .tab-header {
    display: flex;
    padding-bottom: 10px;
    .tab-nav {
      flex: 1;
      text-align: center;
      color: #5c5c5c;
      line-height: 40px;
      cursor: pointer;
      border-bottom: 2px solid transparent;
      &.active {
        color: #2f5cd5;
        border-bottom-color: #2f5cd5;
      }
    }
  }
  .tab-panel {
    padding: 0 10px;
  }
}
</style>
