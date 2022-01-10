<template>
  <a-modal
      :destroy-on-close="true"
      :mask-closable="false"
      title="输入目录"
      :visible="visible"
      :onCancel="onCancel"
  >
    <template #footer>
      <a-button key="back" @click="() => onCancel()">取消</a-button>
      <a-button key="submit" type="primary" @click="onFinish">确定</a-button>
    </template>

    <div>
      <a-input v-model:value="parentDir" placeholder="" />
    </div>

  </a-modal>
</template>

<script lang="ts">
import {defineComponent, onMounted, PropType, ref, Ref} from "vue";

interface DirSelectionSetupData {
  parentDir: Ref<string>;
  onFinish: () => Promise<void>;
}

export default defineComponent({
  name: 'DirSelection',
  props: {
    visible: {
      type: Boolean,
      required: true
    },
    onCancel: {
      type: Function,
      required: true
    },
    onSubmit: {
      type: Function as PropType<(values: any) => void>,
      required: true
    }
  },
  components: {},
  setup(props): DirSelectionSetupData {

    const onFinish = async () => {
      console.log('finish')
      props.onSubmit(parentDir.value);
    };

    let parentDir = ref("");

    onMounted(()=> {
      console.log('onMounted')
    })

    return {
      parentDir,
      onFinish
    }
  }
})
</script>