<template>
    <div class="indexlayout-main-conent">
        <a-card :bordered="false">
          <template #title>
            执行用例
          </template>

          <template #extra>
            <div class="opt">
              <a-button @click="expandAll" type="link">
                <span v-if="!isExpand">展开全部</span>
                <span v-if="isExpand">收缩全部</span>
              </a-button>

              <a-button type="primary" @click="exec">执行</a-button>
              <a-button @click="back()">返回</a-button>
            </div>
          </template>

          <div>
<!--            <div class="toolbar">
              <div class="left"></div>
              <div class="right"></div>
            </div>-->

            <div class="tree-panel">
              <a-tree
                  ref="tree"
                  :tree-data="treeData"
                  :replace-fields="replaceFields"
                  checkable
                  show-icon
                  @expand="expandNode"
                  @select="selectNode"
                  @check="checkNode"

                  v-model:selectedKeys="selectedKeys"
                  v-model:checkedKeys="checkedKeys"
                  v-model:expandedKeys="expandedKeys"
              >
                <template #icon="slotProps">
                  <icon-svg v-if="slotProps.isDir" type="folder-outlined"></icon-svg>
                  <icon-svg v-if="!slotProps.isDir" type="file-outlined"></icon-svg>
                </template>
              </a-tree>
            </div>
          </div>

        </a-card>
    </div>
</template>

<script lang="ts">
import {computed, ComputedRef, defineComponent, onMounted, ref, Ref, watch} from "vue";
import {Form, notification} from "ant-design-vue";
import {useRouter} from "vue-router";
import {useStore} from "vuex";
import {ProjectData} from "@/store/project";
import IconSvg from "@/components/IconSvg";
import {execCase} from "@/views/execution/exec/service";

const useForm = Form.useForm;

interface ExecCasePageSetupData {
  model: any

  treeData: ComputedRef<any[]>;
  replaceFields: any,
  expandNode: (expandedKeys: string[], e: any) => void;
  selectNode: (keys, e) => void;
  checkNode: (keys, e) => void;
  isExpand: Ref<boolean>;
  expandAll: (e) => void;
  selectedKeys: Ref<string[]>
  checkedKeys: Ref<string[]>
  expandedKeys: Ref<string[]>
  tree: Ref;

  exec: (keys) => void;
  back: () => void;
}

export default defineComponent({
  name: 'ExecutionListPage',
  components: {
    IconSvg
  },
  setup(): ExecCasePageSetupData {
    const router = useRouter();
    const model = {}

    const replaceFields = {
      key: 'path',
    };

    let isExpand = ref(false);
    const store = useStore<{ project: ProjectData }>();
    const treeData = computed<any>(() => store.state.project.scriptTree);
    const expandedKeys = ref<string[]>([]);
    const selectedKeys = ref<string[]>([]);
    const checkedKeys = ref<string[]>([]);

    const getOpenKeys = (treeNode, isAll) => {
      if (!treeNode) return

      expandedKeys.value.push(treeNode.path)
      if (treeNode.children && isAll) {
        treeNode.children.forEach((item, index) => {
          getOpenKeys(item, isAll)
        })
      }
    }

    getOpenKeys(treeData.value[0], false)
    watch(treeData,(currConfig)=> {
      expandedKeys.value = []
      getOpenKeys(treeData.value[0], false)
    })

    let tree = ref(null)

    onMounted(() => {
      console.log('onMounted', tree)
    })

    const expandNode = (keys: string[], e: any) => {
      console.log('expandNode', keys[0], e)
    }
    const selectNode = (keys, e) => {
      console.log('selectNode', selectedKeys)
    }
    const checkNode = (keys, e) => {
      console.log('checkNode', checkedKeys)
    }

    const expandAll = (e) => {
      console.log('expandAll')
      isExpand.value = !isExpand.value

      expandedKeys.value = []
      if (isExpand.value) {
        getOpenKeys(treeData.value[0], true)
      }
    }

    const exec = ():void =>  {
      console.log("exec")
      if (checkedKeys.value.length == 0) return
      execCase(checkedKeys.value).then((json) => {
        console.log('json', json)
        if (json.code === 0) {
          notification.success({
            message: `开始执行`,
          });
        }
      })
    }
    const back = ():void =>  {
      router.push(`/execution/history`)
    }

    return {
      model,

      treeData,
      replaceFields,
      expandNode,
      selectNode,
      checkNode,
      isExpand,
      expandAll,
      tree,
      expandedKeys,
      selectedKeys,
      checkedKeys,

      exec,
      back,
    }
  }

})
</script>

<style lang="less" scoped>
.opt {
  .space {
    display: inline-block;
    width: 50px;
  }
  .ant-btn {
    margin-left: 12px;
  }
}

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

</style>