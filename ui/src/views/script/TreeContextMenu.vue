<template>
  <div class="tree-context-menu">
    <div class="menu">
      <div v-if="siteId != 1 && (treeNode.type !== 'file' || treeDataMap[treeNode.id]?.caseId > 0)" @click="menuClick('sync-from-zentao')" class="menu-item">
        <span>{{t('sync-from-zentao')}}</span>
      </div>
      <div v-if="siteId != 1 && (treeNode.type !== 'file' || treeDataMap[treeNode.id]?.caseId > 0)" @click="menuClick('sync-to-zentao')" class="menu-item">
        <span>{{t('sync-to-zentao')}}</span>
      </div>
      <div @click="menuClick('exec')" class="menu-item">
        <span>{{t('exec')}}</span>
      </div>

      <div v-if="treeNode.type != 'workspace'" @click="menuClick('copy')" class="menu-item">
        <span>{{t('copy')}}</span>
      </div>
      <div v-if="treeNode.type != 'workspace'" @click="menuClick('cut')" class="menu-item">
        <span>{{t('cut')}}</span>
      </div>
      <div v-if="clipboardData?.id && (treeNode?.type == 'workspace' || treeNode?.type == 'dir')" @click="menuClick('paste')" class="menu-item">
        <span>{{t('paste')}}</span>
      </div>
      <div v-if="treeNode?.type != 'workspace'" @click="menuClick('delete')" class="menu-item">
        <span>{{t('delete')}}</span>
      </div>

      <div v-if="isElectron" @click="menuClick('open-in-explore')" class="menu-item">
        <span>{{t('open-in-explore')}}</span>
      </div>
      <div v-if="isElectron" @click="menuClick('open-in-terminal')" class="menu-item">
        <span>{{t('open-in-terminal')}}</span>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import {computed, defineComponent, PropType, ref} from "vue";
import {useI18n} from "vue-i18n";
import {useStore} from "vuex";
import {ScriptData} from "@/views/script/store";

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
    },
    clipboardData: {
      type: Object,
      required: false
    },
    siteId: {
        type: Number,
        required: true
    }
  },
  components: {
  },
  setup(props) {
    const { t } = useI18n();

    const isElectron = ref(!!window.require)

    const store = useStore<{ Script: ScriptData }>();
    const treeDataMap = computed<any>(() => store.state.Script.treeDataMap);

    const menuClick = (menuKey) => {
      props.onMenuClick(menuKey, props.treeNode.id);
    };

    return {
      t,
      isElectron,
      treeDataMap,
      menuClick
    }
  }
})
</script>

<style lang="less">
.tree-context-menu {
  .menu {
    padding: 0;
    min-width: 90px;
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