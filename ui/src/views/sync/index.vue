<template>
  <div class="main">

  <a-card title="从禅道同步">
    <a-form :label-col="labelCol" :wrapper-col="wrapperCol">
      <a-form-item label="产品" v-bind="validateInfos.productId">
        <a-select v-model:value="model.productId" @change="selectProduct">
          <a-select-option v-for="item in products" :key="item.id" :value="item.id">{{ item.name }}</a-select-option>
        </a-select>
      </a-form-item>
      <a-form-item label="模块" v-bind="validateInfos.moduleId">
        <a-select v-model:value="model.moduleId">
          <a-select-option v-for="item in modules" :key="item.id" :value="item.id">{{ item.name }}</a-select-option>
        </a-select>
      </a-form-item>
      <a-form-item label="套件" v-bind="validateInfos.suiteId">
        <a-select v-model:value="model.suiteId">
          <a-select-option v-for="item in suites" :key="item.id" :value="item.id">{{ item.name }}</a-select-option>
        </a-select>
      </a-form-item>
      <a-form-item label="任务" v-bind="validateInfos.taskId">
        <a-select v-model:value="model.taskId">
          <a-select-option v-for="item in tasks" :key="item.id" :value="item.id">{{ item.name }}</a-select-option>
        </a-select>
      </a-form-item>
      <a-form-item label="语言" v-bind="validateInfos.lang">
        <a-input v-model:value="model.lang" />
      </a-form-item>
      <a-form-item label="期待结果为独立文件">
        <a-input v-model:value="model.independentFile" />
      </a-form-item>

      <a-form-item :wrapper-col="{ span: 14, offset: 4 }">
        <a-button type="primary" @click.prevent="syncFromZentao">执行</a-button>
        <a-button style="margin-left: 10px" @click="resetFields">重置</a-button>
      </a-form-item>
    </a-form>
  </a-card>

  <a-card title="同步到禅道">
    <a-form :label-col="labelCol" :wrapper-col="wrapperCol">

      <a-form-item :wrapper-col="{ span: 14, offset: 4 }">
        <a-button type="primary" @click.prevent="syncToZentao">执行</a-button>
      </a-form-item>
    </a-form>
  </a-card>

  </div>

</template>
<script lang="ts">
import {defineComponent, ref, reactive, computed, watch, ComputedRef} from "vue";
import { useI18n } from "vue-i18n";

import { Props, validateInfos } from 'ant-design-vue/lib/form/useForm';
import {message, Form, notification} from 'ant-design-vue';
const useForm = Form.useForm;

import { SyncSettings } from './data.d';
import {useStore} from "vuex";
import {ZentaoData} from "./store";
import {ProjectData} from "@/store/project";

interface ConfigFormSetupData {
  formRef: any
  model: SyncSettings
  rules: any
  labelCol: any
  wrapperCol: any
  validate: any
  validateInfos: validateInfos
  syncFromZentao:  () => void;
  syncToZentao:  () => void;
  resetFields:  () => void;

  products: ComputedRef<any[]>;
  modules: ComputedRef<any[]>;
  suites: ComputedRef<any[]>;
  tasks: ComputedRef<any[]>;
  selectProduct:  (item) => void;
}

export default defineComponent({
  name: 'ConfigFormForm',
  components: {
  },
  setup(props): ConfigFormSetupData {
    const { t } = useI18n();

    const storeProject = useStore<{ project: ProjectData }>();
    const currConfig = computed<any>(() => storeProject.state.project.currConfig);

    const store = useStore<{zentao: ZentaoData}>();
    const products = computed<any[]>(() => store.state.zentao.products);
    const modules = computed<any[]>(() => store.state.zentao.modules);
    const suites = computed<any[]>(() => store.state.zentao.suites);
    const tasks = computed<any[]>(() => store.state.zentao.tasks);

    store.dispatch('zentao/fetchProducts')
    watch(currConfig, (currConfig)=> {
      store.dispatch('zentao/fetchProducts')
    })

    const formRef = ref();

    const model = reactive<SyncSettings>({} as SyncSettings);
    const rules = reactive({
      productId: [
        { required: true, type: 'number', message: '请选择产品', trigger: 'change' },
      ],
    });

    const { resetFields, validate, validateInfos } = useForm(model, rules);

    const selectProduct = (item) => {
      console.log('selectProduct', item)
      store.dispatch('zentao/fetchModules', item)
      store.dispatch('zentao/fetchSuites', item)
      store.dispatch('zentao/fetchTasks', item)
    };

    const syncFromZentao = () => {
      validate()
        .then(() => {
          console.log(model);
        })
        .catch(err => {
          console.log('error', err);
        });
    };

    const syncToZentao = () => {
      console.log('syncToZentao')
    };

    return {
      formRef,
      labelCol: { span: 6 },
      wrapperCol: { span: 12 },
      rules,
      resetFields,
      validate,
      validateInfos,
      syncFromZentao,
      syncToZentao,

      model,
      products,
      modules,
      suites,
      tasks,
      selectProduct,
    }

  }
})
</script>

<style lang="less" scoped>
.main {
  padding: 0 20%;
}
</style>