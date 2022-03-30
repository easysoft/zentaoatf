<template>
  <div v-if="!currWorkspace.path">
    <a-empty :image="simpleImage" :description="t('pls_create_workspace')"/>
  </div>

  <div v-if="currWorkspace.path">
    <div v-if="currWorkspace.type === 'unit'" class="panel">
      {{ t('no_sync_for_unittest') }}
    </div>

    <div class="main" v-if="currWorkspace.type === 'func'">
      <a-tabs>
        <a-tab-pane key="1" :tab="t('sync_from_zentao')">

        </a-tab-pane>

        <a-tab-pane key="2" :tab="t('sync_to_zentao')" force-render>
            <a-form :label-col="labelCol" :wrapper-col="wrapperCol">
              <a-form-item :label="t('product')" v-bind="validateInfosTo.productId">
                <a-select v-model:value="modelTo.productId">
                  <a-select-option key="" value="">&nbsp;</a-select-option>
                  <a-select-option v-for="item in products" :key="item.id" :value="item.id">{{ item.name }}</a-select-option>
                </a-select>
              </a-form-item>

              <a-form-item :wrapper-col="{ span: 14, offset: 4 }">
                <a-button type="primary" @click.prevent="syncToZentaoSubmit">{{ t('submit') }}</a-button>
                <a-button style="margin-left: 10px" @click="resetFieldsTo">{{ t('reset') }}</a-button>
              </a-form-item>
            </a-form>
        </a-tab-pane>
      </a-tabs>

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
import {WorkspaceData} from "@/store/workspace";
import {ZentaoData} from "@/store/zentao";
import {syncFromZentao, syncToZentao} from "@/views/sync/service";
import throttle from "lodash.throttle";
import {useRouter} from "vue-router";

const useForm = Form.useForm;

export default defineComponent({
  name: 'ConfigFormForm',
  components: {},
  setup(props) {
    const {t} = useI18n();
    const router = useRouter();

    const storeWorkspace = useStore<{ workspace: WorkspaceData }>();
    const currConfig = computed<any>(() => storeWorkspace.state.workspace.currConfig);
    const currWorkspace = computed<any>(() => storeWorkspace.state.workspace.currWorkspace);

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

    const modelFrom = reactive<SyncSettings>({
      productId: '',
      lang: 'python',
      independentFile: false
    } as SyncSettings);
    const rulesFrom = reactive({
      productId: [
        {required: true, message: t('pls_product')},
      ],
      lang: [
        {required: true, message: t('pls_lang'), trigger: 'change'}
      ],
    });

    const syncFromForm = useForm(modelFrom, rulesFrom);
    const resetFieldsFrom = syncFromForm.resetFields
    const validateFrom = syncFromForm.validate
    const validateInfosFrom = syncFromForm.validateInfos

    const modelTo = reactive<SyncSettings>({productId: ''} as SyncSettings);
    const rulesTo = reactive({
      productId: [
        {required: true, message: t('pls_product')},
      ]
    })
    const syncToForm = useForm(modelTo, rulesTo);
    const resetFieldsTo = syncToForm.resetFields
    const validateTo = syncToForm.validate
    const validateInfosTo = syncToForm.validateInfos

    const selectProduct = (item) => {
      console.log('selectProduct', item)
      if (!item) return

      store.dispatch('zentao/fetchModules', item)
      store.dispatch('zentao/fetchSuites', item)
      store.dispatch('zentao/fetchTasks', item)
    };

    const syncFromZentaoSubmit = () => {
      console.log('syncFromZentaoSubmit')

      validateFrom()
          .then(() => {
            syncFromZentao(modelFrom).then((json) => {
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

      validateTo()
          .then(() => {
            console.log('then', modelTo);
            syncToZentao(modelTo.productId).then((json) => {
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
      currWorkspace,

      formRef,
      labelCol: {span: 6},
      wrapperCol: {span: 12},
      rulesFrom,
      validateFrom,
      validateInfosFrom,
      resetFieldsFrom,
      syncFromZentaoSubmit,

      modelFrom,
      langs,
      products,
      modules,
      suites,
      tasks,
      selectProduct,

      modelTo,
      rulesTo,
      validateTo,
      validateInfosTo,
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