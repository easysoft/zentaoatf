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
      <a-button key="submit" type="primary" @click="onFinish">确定</a-button>
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
  modelRef: Ref<Interpreter>
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

    const modelRef = reactive<any>({productId: props.model.productId + '' || ''})

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
      if (modelRef.value.productId) return
      getDataForBugSubmition(modelRef.value.productId).then((jsn) => {
        products.value = jsn.data
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