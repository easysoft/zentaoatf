<template>
  <a-modal
      title="提交结果到禅道"
      :destroy-on-close="true"
      :mask-closable="false"
      :visible="true"
      :onCancel="onCancel"
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

        <a-form-item label="任务" v-bind="validateInfos.taskId">
          <a-select v-model:value="modelRef.taskId">
            <a-select-option key="" value="">&nbsp;</a-select-option>
            <a-select-option v-for="item in tasks" :key="item.id" :value="item.id+''">{{ item.name }}</a-select-option>
          </a-select>
        </a-form-item>
      </a-form>

    </div>

  </a-modal>
</template>

<script lang="ts">
import {ComputedRef, defineComponent, onMounted, PropType, reactive, ref, Ref} from "vue";
import {Interpreter} from "@/views/config/data";
import { validateInfos } from 'ant-design-vue/lib/form/useForm';
import {message, Form} from 'ant-design-vue';
import {queryProduct, queryTask} from "@/services/zentao";
const useForm = Form.useForm;

interface ResultFormSetupData {
  modelRef: Ref<Interpreter>
  onFinish: () => Promise<void>;

  labelCol: any
  wrapperCol: any
  rules: any
  validate: any
  validateInfos: validateInfos,
  resetFields:  () => void;
  products: Ref<any[]>;
  tasks: Ref<any[]>;
  selectProduct:  (item) => void;
}

export default defineComponent({
  name: 'ResultForm',
  props: {
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

  setup(props): ResultFormSetupData {
    const rules = reactive({
      productId: [
        { required: true, message: '请选择产品' },
      ],
      taskId: [
        { required: true, message: '请选择任务' },
      ],
    });

    const modelRef = reactive<any>({productId: '', taskId: ''})
    let products = ref([])
    let tasks = ref([])

    const { resetFields, validate, validateInfos } = useForm(modelRef, rules);

    const listProduct = () => {
      queryProduct().then((jsn) => {
        products.value = jsn.data
      })
    }
    listProduct()

    const selectProduct = (item) => {
      if (!item) return

      queryTask(item).then((jsn) => {
        tasks.value = jsn.data
      })
    }

    const onFinish = async () => {
      console.log('onFinish', modelRef)

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
      tasks,
      selectProduct,

      modelRef,
      onFinish,
    }
  }
})
</script>