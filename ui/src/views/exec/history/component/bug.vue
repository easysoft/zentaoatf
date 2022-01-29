<template>
  <a-modal
      title="提交缺陷到禅道"
      :destroy-on-close="true"
      :mask-closable="false"
      :visible="true"
      :onCancel="onCancel"
      width="800px"
  >
    <template #footer>
      <a-button key="back" @click="() => onCancel()">取消</a-button>
      <a-button key="submit" type="primary" @click="onFinish">提交</a-button>
    </template>

    <div>
      <a-form :label-col="labelCol" :wrapper-col="wrapperCol">
        <a-form-item label="产品" v-bind="validateInfos.productId">
          <a-select v-model:value="modelRef.productId" @change="selectProduct">
            <a-select-option key="" value="">&nbsp;</a-select-option>
            <a-select-option v-for="item in products" :key="item.id" :value="item.id+''">{{item.name}}</a-select-option>
          </a-select>
        </a-form-item>

        <a-form-item label="模块" v-bind="validateInfos.moduleId">
          <a-select v-model:value="modelRef.moduleId">
            <a-select-option key="" value="">&nbsp;</a-select-option>
            <a-select-option v-for="item in modules" :key="item.id" :value="item.id+''">{{ item.name }}</a-select-option>
          </a-select>
        </a-form-item>

        <a-form-item label="分类">
          <a-select v-model:value="modelRef.categoryId">
            <a-select-option key="" value="">&nbsp;</a-select-option>
            <a-select-option v-for="item in categories" :key="item.id" :value="item.id+''">{{ item.name }}</a-select-option>
          </a-select>
        </a-form-item>

        <a-form-item label="版本">
          <a-select v-model:value="modelRef.versionId">
            <a-select-option key="" value="">&nbsp;</a-select-option>
            <a-select-option v-for="item in versions" :key="item.id" :value="item.id+''">{{ item.name }}</a-select-option>
          </a-select>
        </a-form-item>

        <a-form-item label="严重程度">
          <a-select v-model:value="modelRef.severityId">
            <a-select-option key="" value="">&nbsp;</a-select-option>
            <a-select-option v-for="item in severities" :key="item.id" :value="item.id+''">{{ item.name }}</a-select-option>
          </a-select>
        </a-form-item>

        <a-form-item label="优先级">
          <a-select v-model:value="modelRef.priorityId">
            <a-select-option key="" value="">&nbsp;</a-select-option>
            <a-select-option v-for="item in priorities" :key="item.id" :value="item.id+''">{{ item.name }}</a-select-option>
          </a-select>
        </a-form-item>

        <a-form-item label="步骤">
          <a-textarea v-model:value="modelRef.steps" :auto-size="{ minRows: 5, maxRows: 8 }" />
        </a-form-item>

      </a-form>

    </div>

  </a-modal>
</template>

<script lang="ts">
import {defineComponent, onMounted, PropType, reactive, ref, Ref} from "vue";
import {Interpreter} from "@/views/config/data";
import { validateInfos } from 'ant-design-vue/lib/form/useForm';
import {Form} from 'ant-design-vue';
import {
  getDataForBugSubmition, queryProduct,
} from "@/services/zentao";
const useForm = Form.useForm;

interface BugFormSetupData {
  modelRef: Ref
  onFinish: () => Promise<void>;

  labelCol: any
  wrapperCol: any
  rules: any
  validate: any
  validateInfos: validateInfos,
  resetFields:  () => void;
  products: Ref<any[]>;
  modules: Ref<any[]>;
  categories: Ref<any[]>;
  versions: Ref<any[]>;
  severities: Ref<any[]>;
  priorities: Ref<any[]>;
  selectProduct:  (item) => void;
}

export default defineComponent({
  name: 'BugForm',
  props: {
    model: {
      type: Object as PropType<any>,
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

  setup(props): BugFormSetupData {
    const rules = reactive({
      productId: [
        { required: true, message: '请选择产品' },
      ],
      taskId: [
        { required: true, message: '请选择任务' },
      ],
    });

    const modelRef = reactive<any>({})

    let products = ref([])
    let modules = ref([])
    let categories = ref([])
    let versions = ref([])
    let severities = ref([])
    let priorities = ref([])

    const { resetFields, validate, validateInfos } = useForm(modelRef, rules);

    const getProductData = () => {
      queryProduct().then((jsn) => {
        products.value = jsn.data
      })
    }
    getProductData()

    const getBugData = () => {
      if (!modelRef.productId) return

      getDataForBugSubmition(props.model).then((jsn) => {
        modelRef.value.steps = jsn.data.steps.join('\n')

        modules.value = jsn.data.fields.modules
        categories.value = jsn.data.fields.categories
        versions.value = jsn.data.fields.versions
        severities.value = jsn.data.fields.severities
        priorities.value = jsn.data.fields.priorities
      })
    }
    getBugData()

    const selectProduct = (item) => {
      if (!item) return
      getBugData()
    }

    const onFinish = async () => {
      console.log('onFinish', modelRef.value)

      validate().then(() => {
        props.onSubmit(modelRef);
      }).catch(err => { console.log('') })
    }

    onMounted(()=> {
      console.log('onMounted')
    })

    return {
      labelCol: { span: 6 },
      wrapperCol: { span: 16 },
      rules,
      validate,
      validateInfos,
      resetFields,

      products,
      modules,
      categories,
      versions,
      severities,
      priorities,
      selectProduct,

      modelRef,
      onFinish,
    }
  }
})
</script>