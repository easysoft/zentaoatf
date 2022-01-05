<template>
    <div class="indexlayout-main-conent">
      <div class="main">
        <div class="left">
          <a-tree
            checkable
            :tree-data="scriptTree"
            :replace-fields="replaceFields"
          >
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
import {ProjectData} from "@/store/project";

interface ListScriptPageSetupData {
  scriptTree: ComputedRef<any[]>;
  replaceFields: any,
}

export default defineComponent({
    name: 'ScriptListPage',
    components: {
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

      return {
        scriptTree,
        replaceFields,
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
