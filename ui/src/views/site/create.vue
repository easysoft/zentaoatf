<template>
  <div v-if="!currProject.path">
    <a-empty :image="simpleImage" :description="t('pls_create_project')"/>
  </div>

  <a-card v-if="currProject.path" :title="t('edit_config')">
    <a-form :label-col="labelCol" :wrapper-col="wrapperCol">
      <a-form-item :label="t('zentao_url')" v-bind="validateInfos.url">
        <a-input v-model:value="model.url"
                 @blur="validate('url', { trigger: 'blur' }).catch(() => {})" placeholder="https://zentao.site.com" />
      </a-form-item>
      <a-form-item :label="t('username')" v-bind="validateInfos.username">
        <a-input v-model:value="model.username"
                 @blur="validate('username', { trigger: 'blur' }).catch(() => {})" placeholder="" />
      </a-form-item>
      <a-form-item :label="t('password')" v-bind="validateInfos.password">
        <a-input v-model:value="model.password"
                 @blur="validate('password', { trigger: 'blur' }).catch(() => {})" placeholder="" />
      </a-form-item>

      <a-form-item :wrapper-col="{ span: 14, offset: 4 }">
        <a-button type="primary" @click.prevent="submitForm">{{t('save')}}</a-button> &nbsp;
        <a-button style="margin-left: 10px" @click="resetFields">{{t('reset')}}</a-button>
      </a-form-item>
    </a-form>

  </a-card>

</template>
<script lang="ts">
import {defineComponent, ref, reactive, computed, watch, ComputedRef, Ref, toRaw, toRef} from "vue";
import { useI18n } from "vue-i18n";

import { Props, validateInfos } from 'ant-design-vue/lib/form/useForm';
import {message, Form, notification, Empty} from 'ant-design-vue';
import _ from "lodash";
const useForm = Form.useForm;

import {useStore} from "vuex";
import {ProjectData} from "@/store/project";
import CreateInterpreterForm from './interpreter/create.vue';
import UpdateInterpreterForm from './interpreter/update.vue';
import IconSvg from "@/components/IconSvg/index";

interface SiteFormSetupData {
  t: (key: string | number) => string;
  currProject: ComputedRef;

  currSiteRef: Ref
  model: Ref
  rules: any
  labelCol: any
  wrapperCol: any
  validate: any
  validateInfos: validateInfos
  submitForm:  () => void;
  resetFields:  () => void;

  languages: Ref
  languageMap: Ref
  interpreters: Ref<[]>
  interpreter: Ref

  createSubmitLoading: Ref<boolean>;
  createFormVisible: Ref<boolean>;
  setCreateFormVisible:  (val: boolean) => void;
  addInterpreter: () => void;
  createSubmit: (values: any, resetFields: (newValues?: Props | undefined) => void) => Promise<void>;

  updateSubmitLoading: Ref<boolean>;
  updateFormVisible: Ref<boolean>;
  editInterpreter: (item) => void;
  deleteInterpreter: (item) => void;
  updateFormCancel:  () => void;
  updateSubmit:  (values: any, resetFields: (newValues?: Props | undefined) => void) => Promise<void>;
  simpleImage: any
}

export default defineComponent({
  name: 'SiteForm',
  components: {
  },
  setup(props): SiteFormSetupData {
    const { t } = useI18n();

    const store = useStore<{ project: ProjectData }>();
    const currSiteRef = computed<any>(() => store.state.project.currSite);
    const currProject = computed<any>(() => store.state.project.currProject);

    let model = reactive<any>(currSiteRef.value);

    const rules = reactive({
      url: [
        {
          required: true,
          trigger: 'blur',
          validator: async (rule: any, value: string) => {
            if (value === '' || !value) {
              throw new Error(t('pls_zentao_url'));
            }

            const regx = /(http?|https):\/\/[-A-Za-z0-9+&@#/%?=~_|!:,.;]+[-A-Za-z0-9+&@#/%=~_|]/i
            if (!regx.test(value)) {
              throw new Error(t('wrong_url'));
            }
          }
        },
      ],
      username: [
        { required: true, message: t('pls_username'), trigger: 'blur' },
      ],
      password: [
        { required: true, message: t('pls_password'), trigger: 'blur' },
      ],
    });

    const { resetFields, validate, validateInfos } = useForm(model, rules);
    const submitForm = () => {
      validate()
          .then(() => {
            console.log(model);
            store.dispatch('project/saveSite', model).then((json) => {
              if (json.code === 0) {
                notification.success({
                  message: t('save_success'),
                });
              } else {
                notification.error({
                  message: t('save_fail'),
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
      labelCol: { span: 4 },
      wrapperCol: { span: 14 },
      currSiteRef,
      currProject,
      model,

      rules,
      resetFields,
      validate,
      validateInfos,
      submitForm,
      simpleImage: Empty.PRESENTED_IMAGE_SIMPLE,
    }

  }
})
</script>

<style lang="less" scoped>

.interpreter-header {
  margin: 5px 30px;
  padding-bottom: 6px;
  border-bottom: 1px solid #f0f0f0;
}

.interpreter-item {
  margin: 5px 30px;

}

</style>
