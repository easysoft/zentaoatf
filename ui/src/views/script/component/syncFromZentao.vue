<template>
  <a-form :label-col="labelCol" :wrapper-col="wrapperCol">

    <a-form-item :label="t('workspace')" v-bind="validateInfosFrom.workspaceId">
      <a-select v-model:value="model.workspaceId">
        <a-select-option key="" value="">&nbsp;</a-select-option>
        <a-select-option v-for="item in workspaces" :key="item.id" :value="item.id"><span v-html="item.name"></span>
        </a-select-option>
      </a-select>
    </a-form-item>

    <a-form-item :label="t('module')" v-bind="validateInfosFrom.moduleId">
      <a-select v-model:value="model.moduleId">
        <a-select-option key="" value="">&nbsp;</a-select-option>
        <a-select-option v-for="item in modules" :key="item.id" :value="item.id"><span v-html="item.name"></span>
        </a-select-option>
      </a-select>
    </a-form-item>

    <a-form-item :label="t('suite')" v-bind="validateInfosFrom.suiteId">
      <a-select v-model:value="model.suiteId">
        <a-select-option key="" value="">&nbsp;</a-select-option>
        <a-select-option v-for="item in suites" :key="item.id" :value="item.id">{{ item.name }}</a-select-option>
      </a-select>
    </a-form-item>

    <a-form-item :label="t('task')" v-bind="validateInfosFrom.taskId">
      <a-select v-model:value="model.taskId">
        <a-select-option key="" value="">&nbsp;</a-select-option>
        <a-select-option v-for="item in tasks" :key="item.id" :value="item.id">{{ item.name }}</a-select-option>
      </a-select>
    </a-form-item>

    <a-form-item :label="t('lang')" v-bind="validateInfosFrom.lang">
      <a-select v-model:value="model.lang">
        <a-select-option v-for="item in langs" :key="item.code" :value="item.code">{{ item.name }}</a-select-option>
      </a-select>
    </a-form-item>

    <a-form-item :label="t('independent_by_module')">
      <a-switch v-model:checked="model.byModule"/>
    </a-form-item>

    <a-form-item :label="t('independent_expect')">
      <a-switch v-model:checked="model.independentFile"/>
    </a-form-item>

    <a-form-item :wrapper-col="{ span: 14, offset: 4 }">
      <a-button type="primary" @click.prevent="syncFromZentaoSubmit">{{ t('submit') }}</a-button>
      <a-button style="margin-left: 10px" @click="resetFieldsFrom">{{ t('reset') }}</a-button>
    </a-form-item>

  </a-form>
</template>

<script lang="ts">
import {defineComponent, reactive, ref, Ref, PropType, computed, ComputedRef, watch} from "vue";
import {useI18n} from "vue-i18n";

import {Empty, Form, notification} from 'ant-design-vue';
import {useStore} from "vuex";
import {ZentaoData} from "@/store/zentao";
import {SyncSettings} from "@/views/sync/data";
import {useRouter} from "vue-router";
import {WorkspaceData} from "../../workspace/store";
import {syncFromZentao} from "@/views/script/service";

const useForm = Form.useForm;

export default defineComponent({
  name: 'SyncFromZentao',
  props: {
    onClose: {
      type: Function,
      required: true
    },
  },
  components: {},
  setup(props) {
    const {t} = useI18n();
    const router = useRouter();

    const workspaceStore = useStore<{ Workspace: WorkspaceData }>();
    const zentaoStore = useStore<{ zentao: ZentaoData }>();

    const currSite = computed<any>(() => zentaoStore.state.zentao.currSite);
    const currProduct = computed<any>(() => zentaoStore.state.zentao.currProduct);

    const workspaces = computed<any[]>(() => workspaceStore.state.Workspace.listResult);
    const langs = computed<any[]>(() => zentaoStore.state.zentao.langs);
    const modules = computed<any[]>(() => zentaoStore.state.zentao.modules);
    const suites = computed<any[]>(() => zentaoStore.state.zentao.suites);
    const tasks = computed<any[]>(() => zentaoStore.state.zentao.tasks);

    const fetchData = () => {
      workspaceStore.dispatch('Workspace/list', currProduct.value.id)
      zentaoStore.dispatch('zentao/fetchModules', {siteId: currSite.value.id, productId: currProduct.value.id})
      zentaoStore.dispatch('zentao/fetchSuites', {siteId: currSite.value.id, productId: currProduct.value.id})
      zentaoStore.dispatch('zentao/fetchTasks', {siteId: currSite.value.id, productId: currProduct.value.id})
    }
    fetchData()

    const formRef = ref();

    const model = reactive<SyncSettings>({
      workspaceId: '',
      lang: 'python',
      independentFile: false
    } as SyncSettings);

    const rulesFrom = reactive({
      workspaceId: [
        {required: true, message: '请选择导出到的工作目录'},
      ],
      lang: [
        {required: true, message: t('pls_lang'), trigger: 'change'}
      ],
    });

    const syncFromForm = useForm(model, rulesFrom);
    const resetFieldsFrom = syncFromForm.resetFields
    const validateFrom = syncFromForm.validate
    const validateInfosFrom = syncFromForm.validateInfos

    const syncFromZentaoSubmit = () => {
      console.log('syncFromZentaoSubmit')

      validateFrom()
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

    return {
      t,

      formRef,
      labelCol: {span: 6},
      wrapperCol: {span: 12},
      rulesFrom,
      validateFrom,
      validateInfosFrom,
      resetFieldsFrom,

      model,
      workspaces,
      langs,
      modules,
      suites,
      tasks,

      syncFromZentaoSubmit,

      simpleImage: Empty.PRESENTED_IMAGE_SIMPLE,
    }

  }
})
</script>