<template>
  <a-modal
      :destroy-on-close="true"
      :mask-closable="false"
      title="新建项目"
      :visible="visible"
      :onCancel="onCancel"
  >
    <template #footer>
      <a-button key="back" @click="() => onCancel()">取消</a-button>
      <a-button key="submit" type="primary" @click="onFinish">确定</a-button>
    </template>

    <div>
      <a-form :labelCol="{ span: 4 }" :wrapper-col="{span:20}">
        <a-form-item label="项目路径" v-bind="validateInfos.path">
          <a-input v-model:value="modelRef.path" placeholder="" />
        </a-form-item>

        <a-form-item label="项目类型" v-bind="validateInfos.type">
          <a-select v-model:value="modelRef.type">
            <a-select-option key="func" value="func">ZTF自动化</a-select-option>
            <a-select-option key="unit" value="unit">单元测试</a-select-option>
          </a-select>
        </a-form-item>

      </a-form>

    </div>

  </a-modal>
</template>

<script lang="ts">
import {defineComponent, onMounted, PropType, reactive, ref, Ref} from "vue";
import {Interpreter} from "@/views/config/data";
import {useI18n} from "vue-i18n";
import {Form} from "ant-design-vue";
import { validateInfos } from 'ant-design-vue/lib/form/useForm';

interface ProjectCreateFormSetupData {
  modelRef: Ref<Interpreter>
  validateInfos: validateInfos
  onFinish: () => Promise<void>;
}

export default defineComponent({
  name: 'ProjectCreateForm',
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

  setup(props): ProjectCreateFormSetupData {
    const { t } = useI18n();

    const modelRef = reactive<any>({path: '', type: 'func'})
    const rulesRef = reactive({
      path: [ { required: true, message: '请输入项目完整路径' } ],
      type: [ { required: true, message: '请选择项目类型' } ],
    });

    const { validate, validateInfos } = Form.useForm(modelRef, rulesRef);

    const onFinish = async () => {
      console.log('onFinish')

      validate().then(() => {
        props.onSubmit(modelRef);
      }).catch(err => { console.log('') })
    }

    onMounted(()=> {
      console.log('onMounted')
    })

    return {
      modelRef,
      validateInfos,
      onFinish
    }
  }
})
</script>