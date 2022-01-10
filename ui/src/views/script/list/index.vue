<template>
    <div class="indexlayout-main-conent">
      <div class="main">
        <div class="left">
          <div class="toolbar">
            <a-button @click="expandOrNot" type="link">
              <span v-if="!isExpand">展开全部</span>
              <span v-if="isExpand">收缩全部</span>
            </a-button>
          </div>

          <a-tree
            ref="tree"
            :tree-data="treeData"
            :replace-fields="replaceFields"
            show-icon
            @select="selectNode"
            @expand="expandNode"
            v-model:expandedKeys="expandedKeys"
          >
            <template #icon="slotProps">
              <icon-svg v-if="slotProps.isDir" type="folder-outlined"></icon-svg>
              <icon-svg v-if="!slotProps.isDir" type="file-outlined"></icon-svg>
            </template>
          </a-tree>
        </div>

        <div class="content">
          CONTENT
        </div>
      </div>
    </div>
</template>

<script lang="ts">
import {ComputedRef, defineComponent, ref, Ref, reactive, computed, onMounted} from "vue";
import {useStore} from "vuex";
import IconSvg from "@/components/IconSvg";
import {ProjectData} from "@/store/project";

interface ListScriptPageSetupData {
  treeData: ComputedRef<any[]>;
  replaceFields: any,
  selectNode: (keys, e) => void;
  isExpand: Ref<boolean>;
  expandOrNot: (e) => void;
  expandNode: (e) => void;
  expandedKeys: Ref<string[]>
  tree: Ref<any>;
}

export default defineComponent({
    name: 'ScriptListPage',
    components: {
      IconSvg
    },
    setup(): ListScriptPageSetupData {
      const replaceFields = {
        key: 'path',
      };

      let isExpand = ref(false);
      const store = useStore<{ project: ProjectData }>();
      const treeData = computed<any>(() => store.state.project.scriptTree);
      let expandedKeys = computed<any>(() => store.state.project.scriptTreeOpenKeys);

      let tree = ref(null)

      onMounted(()=> {
        console.log('onMounted', tree)
      })

      const selectNode = (selectedKeys, e) => {
        if (e.selectedNodes.length === 0) return
        console.log('selectNode', e.selectedNodes[0].props.path)
      }

      const expandNode = (e) => {
        console.log('expandNode')
      }
      const expandOrNot = (e) => {
        console.log('expandOrNot')
        isExpand.value = !isExpand.value

        store.dispatch('project/genOpenKeys', isExpand.value);
      }

      return {
        treeData,
        replaceFields,
        selectNode,
        isExpand,
        expandOrNot,
        expandNode,
        tree,
        expandedKeys,
      }
    }

})
</script>

<style lang="less">
.indexlayout-main-conent {
  margin: 0px;
  height: 100%;

  .main {
    display: flex;
    height: 100%;
    .left {
      padding: 6px;
      width: 300px;
      height: 100%;
      border-right: 1px solid #D0D7DE;
      .toolbar {
        text-align: right;
        .ant-btn-link {
          padding: 0px 5px;
          color: #1890ff;
        }
      }
      .ant-tree {
        font-size: 16px;
      }
    }
    .content {
      padding: 15px;
      flex: 1;
      height: 100%;
    }
  }
}
</style>
