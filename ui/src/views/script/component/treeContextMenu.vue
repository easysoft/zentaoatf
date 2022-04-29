<template>
  <div class="dp-tree-context-menu">
    <a-menu @click="menuClick" mode="inline">
      <template v-if="displayBy === 'workspace'">
        <a-menu-item key="rename" class="menu-item" v-if="treeNode.type !== 'workspace'">
          <EditOutlined />
          <span>重命名</span>
        </a-menu-item>

        <a-menu-item key="add_brother_node" class="menu-item" v-if="treeNode.type !== 'workspace'">
          <PlusOutlined />
          <span>创建同级脚本</span>
        </a-menu-item>

        <a-menu-item key="add_child_node" class="menu-item" v-if="treeNode.type !== 'file'">
          <PlusOutlined />
          <span>创建子级脚本</span>
        </a-menu-item>

        <a-menu-item key="add_brother_dir" class="menu-item" v-if="treeNode.type !== 'workspace'">
          <PlusOutlined />
          <span>创建同级目录</span>
        </a-menu-item>

        <a-menu-item key="add_child_dir" class="menu-item" v-if="treeNode.type !== 'file'">
          <PlusOutlined />
          <span>创建子级目录</span>
        </a-menu-item>

        <a-menu-item key="remove" class="menu-item" v-if="treeNode.type !== 'workspace'">
          <CloseOutlined />
          <span v-if="treeNode.type === 'dir'">删除目录</span>
          <span v-if="treeNode.type === 'file'">删除脚本</span>
        </a-menu-item>

      </template>

      <template v-if="displayBy === 'module'">
        <a-menu-item key="sync_from_zentao" class="menu-item">
          <ArrowDownOutlined />
          <span>从禅道同步</span>
        </a-menu-item>

        <a-menu-item key="sync_to_zentao" class="menu-item">
          <ArrowUpOutlined />
          <span>同步到禅道</span>
        </a-menu-item>
      </template>
    </a-menu>
  </div>
</template>

<script lang="ts">
import {defineComponent, PropType, Ref} from "vue";
import {useI18n} from "vue-i18n";
import {Form, message} from 'ant-design-vue';
import {EditOutlined, CloseOutlined, PlusOutlined, ArrowDownOutlined, ArrowUpOutlined} from "@ant-design/icons-vue";

const useForm = Form.useForm;

export default defineComponent({
  name: 'TreeContextMenu',
  props: {
    treeNode: {
      type: Object,
      required: true
    },
    displayBy: {
      type: String,
      required: true
    },
    onSubmit: {
      type: Function as PropType<(selectedKey: string, targetId: number) => void>,
      required: true
    }
  },
  components: {
    EditOutlined, PlusOutlined, CloseOutlined,
    ArrowDownOutlined, ArrowUpOutlined,
  },
  setup(props) {
    const {t} = useI18n();

    const menuClick = (e) => {
      console.log('menuClick', e, props.treeNode)
      const key = e.key
      const targetId = props.treeNode.path

      props.onSubmit(key, targetId);
    };

    return {
      menuClick
    }
  }
})
</script>

<style lang="less">
.dp-tree-context-menu {
  z-index: 9;
  .ant-tree-node-content-wrapper {
    display: block !important;
  }
  .ant-menu {
    border: 1px solid #dedfe1;
    background: #f0f2f5;
    .ant-menu-item.menu-item {
      padding-left: 12px !important;
      height: 22px;
      line-height: 21px;
      text-align: left;
      .ant-menu-title-content {
        height: 22px;
        line-height: 21px;
      }
    }
  }
}
</style>