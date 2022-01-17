<template>
  <div class="indexlayout-main-conent">
    <div id="main">
      <div id="left">
        <div class="toolbar">
          <div class="left"></div>
          <div class="right">
            <a-button @click="expandAll" type="link">
              <span v-if="!isExpand">展开全部</span>
              <span v-if="isExpand">收缩全部</span>
            </a-button>
          </div>
        </div>

        <div class="tree-panel">
          <a-tree
              ref="tree"
              :tree-data="treeData"
              :replace-fields="replaceFields"
              show-icon
              @expand="expandNode"
              @select="selectNode"
              v-model:expandedKeys="expandedKeys"
          >
            <template #icon="slotProps">
              <icon-svg v-if="slotProps.isDir" type="folder-outlined"></icon-svg>
              <icon-svg v-if="!slotProps.isDir" type="file-outlined"></icon-svg>
            </template>
          </a-tree>
        </div>
      </div>

      <div id="resize"></div>

      <div id="content">
        <div class="toolbar">
          <a-button @click="expandAll" type="primary">执行</a-button>
        </div>
        <div class="panel">
          <pre><code>{{ script.code }}</code></pre>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import {computed, ComputedRef, defineComponent, onMounted, ref, Ref, watch} from "vue";
import {useStore} from "vuex";
import IconSvg from "@/components/IconSvg";
import {ProjectData} from "@/store/project";
import {ScriptData} from "../store";
import {resizeWidth} from "@/utils/dom";

interface ListScriptPageSetupData {
  treeData: ComputedRef<any[]>;
  replaceFields: any,
  expandNode: (expandedKeys: string[], e: any) => void;
  selectNode: (keys, e) => void;
  isExpand: Ref<boolean>;
  expandAll: (e) => void;
  expandedKeys: Ref<string[]>
  tree: Ref;
  script: any
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
    const expandedKeys = ref<string[]>([]);

    const getOpenKeys = (treeNode, isAll) => {
      if (!treeNode) return

      expandedKeys.value.push(treeNode.path)
      if (treeNode.children && isAll) {
        treeNode.children.forEach((item, index) => {
          getOpenKeys(item, isAll)
        })
      }
    }

    expandedKeys.value = []
    getOpenKeys(treeData.value[0], false)
    watch(treeData,(currConfig)=> {
      expandedKeys.value = []
      getOpenKeys(treeData.value[0], false)
    })

    let tree = ref(null)

    const storeScript = useStore<{ script: ScriptData }>();
    const script = computed<any>(() => storeScript.state.script.detail);

    onMounted(() => {
      console.log('onMounted', tree)
      resizeWidth('main', 'left', 'resize', 'content', 280, 800)
    })

    const expandNode = (keys: string[], e: any) => {
      console.log('expandNode', keys[0], e)
    }

    const selectNode = (selectedKeys, e) => {
      console.log('selectNode', e.selectedNodes)
      if (e.selectedNodes.length === 0) return

      storeScript.dispatch('script/getScript', e.selectedNodes[0].props);
    }

    const expandAll = (e) => {
      console.log('expandAll')
      isExpand.value = !isExpand.value

      expandedKeys.value = []
      if (isExpand.value) {
        getOpenKeys(treeData.value[0], true)
      }
    }

    return {
      treeData,
      replaceFields,
      expandNode,
      selectNode,
      isExpand,
      expandAll,
      tree,
      expandedKeys,
      script,
    }
  }

})
</script>

<style lang="less">
.indexlayout-main-conent {
  margin: 0px;
  height: 100%;

  #main {
    display: flex;
    height: 100%;
    width: 100%;

    #left {
      width: 380px;
      height: 100%;
      padding: 3px;

      .toolbar {
        display: flex;
        padding: 0 3px;
        border-bottom: 1px solid #D0D7DE;

        .left {
          flex: 1;
        }
        .right {
          width: 70px;
          text-align: right;
        }
        .ant-btn-link {
          padding: 0px 3px;
          color: #1890ff;
        }
      }

      .tree-panel {
        height: calc(100% - 35px);
        overflow: auto;

        .ant-tree {
          font-size: 16px;
        }
      }

    }

    #resize {
      width: 2px;
      height: 100%;
      background-color: #D0D7DE;
      cursor: ew-resize;

      &.active {
        background-color: #a9aeb4;
      }
    }

    #content {
      width: 80%;
      height: 100%;
      flex: 1;

      .toolbar {
        padding: 5px 10px;
        height: 46px;
        text-align: right;
      }

      .panel {
        padding: 16px;
        height: calc(100% - 46px);
        overflow: auto;
      }
    }
  }
}
</style>
