<template>
  <div class="tree-context-menu">
    <div class="menu">
      <div @click="menuClick('sync-from-zentao')" class="menu-item">
        <span>{{t('sync-from-zentao')}}</span>
      </div>
      <div @click="menuClick('sync-to-zentao')" class="menu-item">
        <span>{{t('sync-to-zentao')}}</span>
      </div>
      <div @click="menuClick('exec')" class="menu-item">
        <span>{{t('exec')}}</span>
      </div>

      <div @click="menuClick('copy')" class="menu-item">
        <span>{{t('copy')}}</span>
      </div>
      <div @click="menuClick('cut')" class="menu-item">
        <span>{{t('cut')}}</span>
      </div>
      <div @click="menuClick('paste')" class="menu-item">
        <span>{{t('paste')}}</span>
      </div>

      <div @click="menuClick('open-in-explore')" class="menu-item">
        <span>{{t('open-in-explore')}}</span>
      </div>
      <div @click="menuClick('open-in-terminal')" class="menu-item">
        <span>{{t('open-in-terminal')}}</span>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import {defineComponent, PropType, Ref} from "vue";
import {useI18n} from "vue-i18n";

export default defineComponent({
  name: 'TreeContextMenu',
  props: {
    treeNode: {
      type: Object,
      required: true
    },
    onMenuClick: {
      type: Function as PropType<(menuKey: string, targetId: number) => void>,
      required: true
    }
  },
  components: {
  },
  setup(props) {
    const { t } = useI18n();

    const menuClick = (menuKey) => {
      props.onMenuClick(menuKey, props.treeNode.id);
    };

    return {
      t,
      menuClick
    }
  }
})
</script>

<style lang="less">
.tree-context-menu {
  .menu {
    padding: 0;
    border: 1px solid #dedfe1;
    background-color: #fff;
    .menu-item {
      margin: 0;
      padding: 5px 10px;
      height: 22px;
      line-height: 22px;
      cursor: pointer;
      &:hover {
        background-color: #f5f5f5;
      }
    }
  }
}
</style>