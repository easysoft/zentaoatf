<template>
    <div class="indexlayout-main-conent">
      <div class="main">
        <div class="left">
          <a-tree
            :tree-data="scriptTree"
            :replace-fields="replaceFields"
            show-icon
            @select="selectNode"
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
  scriptTree: ComputedRef<any[]>;
  replaceFields: any,
  selectNode: (keys, e) => void;
}

export default defineComponent({
    name: 'ScriptListPage',
    components: {
      IconSvg
    },
    setup(): ListScriptPageSetupData {
      const replaceFields = {
        key: 'title',
      };

      const store = useStore<{ project: ProjectData }>();
      const scriptTree = computed<any>(() => store.state.project.scriptTree);

      onMounted(()=> {
        console.log('onMounted')
      })


      const selectNode = (selectedKeys, e) => {
        if (e.selectedNodes.length === 0) return
        console.log('selectNode', e.selectedNodes[0].props.path)
      }

      return {
        scriptTree,
        replaceFields,
        selectNode,
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
