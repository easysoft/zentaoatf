<template>
  <div v-if="!currProject.path">
    <a-empty :image="simpleImage" :description="t('pls_create_project')"/>
  </div>

  <div v-if="currProject.path">
    <div v-if="currProject.type === 'unit'" class="panel">
      {{ t('no_sync_for_unittest') }}
    </div>
    <div class="main" v-if="currProject.type === 'func'">

      <a-card :title="t('sync_from_zentao')">
        <a-form :label-col="labelCol" :wrapper-col="wrapperCol">
          <a-form-item :label="t('product')" v-bind="validateInfos.productId">
            <a-select v-model:value="model.productId" @change="selectProduct">
              <a-select-option key="" value="">&nbsp;</a-select-option>
              <a-select-option v-for="item in products" :key="item.id" :value="item.id">{{ item.name }}</a-select-option>
            </a-select>
          </a-form-item>
          <a-form-item :label="t('module')" v-bind="validateInfos.moduleId">
            <a-select v-model:value="model.moduleId">
              <a-select-option key="" value="">&nbsp;</a-select-option>
              <a-select-option v-for="item in modules" :key="item.id" :value="item.id"><span v-html="item.name"></span>
              </a-select-option>
            </a-select>
          </a-form-item>
          <a-form-item :label="t('suite')" v-bind="validateInfos.suiteId">
            <a-select v-model:value="model.suiteId">
              <a-select-option key="" value="">&nbsp;</a-select-option>
              <a-select-option v-for="item in suites" :key="item.id" :value="item.id">{{ item.name }}</a-select-option>
            </a-select>
          </a-form-item>
          <a-form-item :label="t('task')" v-bind="validateInfos.taskId">
            <a-select v-model:value="model.taskId">
              <a-select-option key="" value="">&nbsp;</a-select-option>
              <a-select-option v-for="item in tasks" :key="item.id" :value="item.id">{{ item.name }}</a-select-option>
            </a-select>
          </a-form-item>
          <a-form-item :label="t('lang')" v-bind="validateInfos.lang">
            <a-select v-model:value="model.lang">
              <a-select-option v-for="item in langs" :key="item.code" :value="item.code">{{ item.name }}</a-select-option>
            </a-select>
          </a-form-item>
          <a-form-item :label="t('independent_expect')">
            <a-switch v-model:checked="model.independentFile"/>
          </a-form-item>
          <a-form-item :wrapper-col="{ span: 14, offset: 4 }">
            <a-button type="primary" @click.prevent="syncFromZentaoSubmit">{{ t('submit') }}</a-button>
            <a-button style="margin-left: 10px" @click="resetFieldsFrom">{{ t('reset') }}</a-button>
          </a-form-item>
        </a-form>
      </a-card>

      <a-card :title="t('sync_to_zentao')">
        <a-form :label-col="labelCol" :wrapper-col="wrapperCol">
          <a-form-item :label="t('product')" v-bind="validateInfosCommit.productId">
            <a-select v-model:value="modelCommit.productId">
              <a-select-option key="" value="">&nbsp;</a-select-option>
              <a-select-option v-for="item in products" :key="item.id" :value="item.id">{{ item.name }}</a-select-option>
            </a-select>
          </a-form-item>

          <a-form-item :wrapper-col="{ span: 14, offset: 4 }">
            <a-button type="primary" @click.prevent="syncToZentaoSubmit">{{ t('submit') }}</a-button>
            <a-button style="margin-left: 10px" @click="resetFieldsTo">{{ t('reset') }}</a-button>
          </a-form-item>
        </a-form>
      </a-card>

    </div>
  </div>
</template>
<script lang="ts">
import {computed, ComputedRef, defineComponent, reactive, ref, watch} from "vue";
import {useI18n} from "vue-i18n";

import {validateInfos} from 'ant-design-vue/lib/form/useForm';
import {Empty, Form, notification} from 'ant-design-vue';
import {SyncSettings} from './data.d';
import {useStore} from "vuex";
import {ProjectData} from "@/store/project";
import {ZentaoData} from "@/store/zentao";
import {syncFromZentao, syncToZentao} from "@/views/sync/service";
import throttle from "lodash.debounce";
import {useRouter} from "vue-router";

const useForm = Form.useForm;

interface ConfigFormSetupData {
  t: (key: string | number) => string;
  currProject: ComputedRef;

  formRef: any
  model: SyncSettings
  rules: any

  labelCol: any
  wrapperCol: any
  validate: any
  validateInfos: validateInfos
  syncFromZentaoSubmit: () => void;
  resetFieldsFrom: () => void;

  langs: ComputedRef<any[]>;
  products: ComputedRef<any[]>;
  modules: ComputedRef<any[]>;
  suites: ComputedRef<any[]>;
  tasks: ComputedRef<any[]>;
  selectProduct: (item) => void;

  modelCommit: SyncSettings
  rulesCommit: any
  syncToZentaoSubmit: () => void;
  validateCommit: any
  validateInfosCommit: validateInfos
  resetFieldsTo: () => void;
  simpleImage: any
}

export default defineComponent({
  name: 'ConfigFormForm',
  components: {},
  setup(props): ConfigFormSetupData {
    const {t} = useI18n();
    const router = useRouter();

    const storeProject = useStore<{ project: ProjectData }>();
    const currConfig = computed<any>(() => storeProject.state.project.currConfig);
    const currProject = computed<any>(() => storeProject.state.project.currProject);

    const store = useStore<{ zentao: ZentaoData }>();
    const langs = computed<any[]>(() => store.state.zentao.langs);
    const products = computed<any[]>(() => store.state.zentao.products);
    const modules = computed<any[]>(() => store.state.zentao.modules);
    const suites = computed<any[]>(() => store.state.zentao.suites);
    const tasks = computed<any[]>(() => store.state.zentao.tasks);

    const fetchProducts = throttle((): void => {
      store.dispatch('zentao/fetchLangs')
      store.dispatch('zentao/fetchProducts').catch((error) => {
        if (error.response.data.code === 2000) router.push(`/config`)
      })
    }, 600)
    fetchProducts()
    watch(currConfig, () => {
      fetchProducts()
    })

    const formRef = ref();

    const model = reactive<SyncSettings>({
      productId: '',
      lang: 'python',
      independentFile: false
    } as SyncSettings);

    const modelCommit = reactive<SyncSettings>({
      productId: '',
    } as SyncSettings);

    const rules = reactive({
      productId: [
        {required: true, message: t('pls_product')},
      ],
      lang: [
        {required: true, message: t('pls_lang'), trigger: 'change'}
      ],
    });
    const rulesCommit = reactive({
      productId: [
        {required: true, message: t('pls_product')},
      ]
    })

    const syncFromForm = useForm(model, rules);
    const resetFieldsFrom = syncFromForm.resetFields
    const  validate = syncFromForm.validate
    const  validateInfos = syncFromForm.validateInfos

    const commitForm = useForm(modelCommit, rulesCommit);
    const resetFieldsTo = commitForm.resetFields
    const validateCommit = commitForm.validate
    const validateInfosCommit = commitForm.validateInfos

    const selectProduct = (item) => {
      console.log('selectProduct', item)
      if (!item) return

      store.dispatch('zentao/fetchModules', item)
      store.dispatch('zentao/fetchSuites', item)
      store.dispatch('zentao/fetchTasks', item)
    };

    const syncFromZentaoSubmit = () => {
      console.log('syncFromZentaoSubmit')

      validate()
          .then(() => {
            syncFromZentao(model).then((json) => {
              console.log('json', json)
              if (json.code === 0) {
                notification.success({
                  message: `同步成功`,
                });
              } else {
                notification.error({
                  message: `同步失败`,
                  description: json.msg,
                });
              }
            })
          })
          .catch(err => {
            console.log('validate fail', err)
          });
    };

    const syncToZentaoSubmit = () => {
      console.log('syncToZentaoSubmit')

      validateCommit()
          .then(() => {
            console.log('then', modelCommit);
            syncToZentao(modelCommit.productId).then((json) => {
              console.log('json', json)
              if (json.code === 0) {
                notification.success({
                  message: t('sync_success'),
                });
              } else {
                notification.error({
                  message: t('sync_fail'),
                  description: json.msg,
                });
              }
            })
          })
          .catch(err => {
            console.log('error', err);
          });
    };

    return {
      t,
      currProject,

      formRef,
      labelCol: {span: 6},
      wrapperCol: {span: 12},
      rules,
      validate,
      validateInfos,
      resetFieldsFrom,
      syncFromZentaoSubmit,

      model,
      langs,
      products,
      modules,
      suites,
      tasks,
      selectProduct,

      modelCommit,
      rulesCommit,
      validateCommit,
      validateInfosCommit,
      resetFieldsTo,
      syncToZentaoSubmit,
      simpleImage: Empty.PRESENTED_IMAGE_SIMPLE,
    }

  }
})
</script>

<style lang="less" scoped>
.panel {
  padding: 20px;
}

.main {
  padding: 50px 20% 0 20%;
}
</style>