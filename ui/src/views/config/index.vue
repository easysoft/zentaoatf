<template>

  <a-card :title="t('edit_config')">
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

      <a-form-item v-if="currProject.type === 'func' && currConfigRef.isWin" label="执行器">
        <div>
          <a-row class="interpreter-header">
            <a-col :span="4" class="t-center t-bord">{{t('script_lang')}}</a-col>
            <a-col :span="18" class="t-center t-bord">{{t('interpreter')}}</a-col>
            <a-col :span="2" class="t-center t-bord">
              <a-button @click="addInterpreter" type="link" size="small">{{ t('create') }}</a-button>
            </a-col>
          </a-row>

          <a-row  v-for="item in interpreters" :key="item.lang" class="interpreter-item">
            <a-col :span="4" class="t-bord">
              {{languageMap[item.lang]}}
            </a-col>
            <a-col :span="18" class="t-bord">
              {{item.val}}
            </a-col>
            <a-col :span="2" class="t-center t-bord">
              <icon-svg @click="editInterpreter(item)" type="edit" class="t-icon"></icon-svg> &nbsp;
              <icon-svg @click="deleteInterpreter(item)" type="close" class="t-icon"></icon-svg>
            </a-col>
          </a-row>

        </div>
      </a-form-item>

      <a-form-item :wrapper-col="{ span: 14, offset: 4 }">
        <a-button type="primary" @click.prevent="submitForm">{{t('save')}}</a-button> &nbsp;
        <a-button style="margin-left: 10px" @click="resetFields">{{t('reset')}}</a-button>
      </a-form-item>
    </a-form>

    <create-interpreter-form
        v-if="createFormVisible===true"
        :visible="createFormVisible"
        :onCancel="() => setCreateFormVisible(false)"
        :onSubmitLoading="createSubmitLoading"
        :onSubmit="createSubmit"
        :languages="languages"
        :languageMap="languageMap"
    />

    <update-interpreter-form
        v-if="updateFormVisible===true"
        :visible="updateFormVisible"
        :model="interpreter"
        :onCancel="() => updateFormCancel()"
        :onSubmitLoading="updateSubmitLoading"
        :onSubmit="updateSubmit"
        :languageMap="languageMap"
    />

  </a-card>

</template>
<script lang="ts">
import {defineComponent, ref, reactive, computed, watch, ComputedRef, Ref, toRaw, toRef} from "vue";
import { useI18n } from "vue-i18n";

import { Props, validateInfos } from 'ant-design-vue/lib/form/useForm';
import {message, Form, notification} from 'ant-design-vue';
import _ from "lodash";
const useForm = Form.useForm;

import {Config, Interpreter} from './data.d';
import {useStore} from "vuex";
import {ProjectData} from "@/store/project";
import CreateInterpreterForm from './interpreter/create.vue';
import UpdateInterpreterForm from './interpreter/update.vue';
import {
  getInterpretersFromConfig,
  setInterpreter,
} from "@/utils/testing";
import IconSvg from "@/components/IconSvg/index";
import throttle from "lodash.debounce";

interface ConfigFormSetupData {
  t: (key: string | number) => string;
  currProject: ComputedRef;

  currConfigRef: Ref
  model: Partial<Config>
  rules: any
  labelCol: any
  wrapperCol: any
  validate: any
  validateInfos: validateInfos
  submitForm:  () => void;
  resetFields:  () => void;

  languages: Ref
  languageMap: Ref
  interpreters: Ref<Interpreter[]>
  interpreter: Ref<Interpreter>

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
}

export default defineComponent({
  name: 'ConfigForm',
  components: {
    IconSvg, CreateInterpreterForm, UpdateInterpreterForm,
  },
  setup(props): ConfigFormSetupData {
    const { t } = useI18n();

    const store = useStore<{ project: ProjectData }>();
    const currConfigRef = computed<any>(() => store.state.project.currConfig);
    const currProject = computed<any>(() => store.state.project.currProject);

    let interpreter = reactive<any>({} as Interpreter)
    let interpreters = ref<any>([] as Interpreter[])
    let languages = ref<any>({})
    let languageMap = ref<any>({})

    const getInterpreters = throttle((currConfig) => {
      const data = getInterpretersFromConfig(currConfig)
      interpreters.value = data.interpreters
      languages.value = data.languages
      languageMap.value = data.languageMap
    }, 600)
    getInterpreters(currConfigRef.value)

    let model = reactive<any>(currConfigRef.value);
    watch(currConfigRef,()=> {
      console.log('watch currConfigRef', currConfigRef)
      _.extend(model, currConfigRef.value)
      getInterpreters(currConfigRef.value)
    }, {deep: true})

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
          setInterpreter(model, interpreters)
          console.log(model);
          store.dispatch('project/saveConfig', model).then((json) => {
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

    // 新建解析器
    const createFormVisible = ref<boolean>(false);
    const setCreateFormVisible = (val: boolean) => {
      createFormVisible.value = val;
    };
    const createSubmitLoading = ref<boolean>(false);
    const addInterpreter = () => {
      interpreter.value = {}
      setCreateFormVisible(true)
    }
    const createSubmit = async (values: any, resetFields: (newValues?: Props | undefined) => void) => {
      // createSubmitLoading.value = true;
      currConfigRef.value[values.lang] = values.val
      // createSubmitLoading.value = false;
      setCreateFormVisible(false)
    }

    // 更新解析器
    const updateFormVisible = ref<boolean>(false);
    const setUpdateFormVisible = (val: boolean) => {
      updateFormVisible.value = val;
    }
    const updateFormCancel = () => {
      setUpdateFormVisible(false);
    }
    const updateSubmitLoading = ref<boolean>(false);
    const editInterpreter = (item) => {
      interpreter.value = item
      setUpdateFormVisible(true)
    }
    const updateSubmit = async (values: any, resetFields: (newValues?: Props | undefined) => void) => {
      currConfigRef.value[values.lang] = values.val
      setUpdateFormVisible(false)
    }
    const deleteInterpreter = (item) => {
      delete currConfigRef.value[item.lang]
    }

    return {
      t,
      labelCol: { span: 4 },
      wrapperCol: { span: 14 },
      currConfigRef,
      currProject,
      model,

      rules,
      resetFields,
      validate,
      validateInfos,
      submitForm,

      languages,
      languageMap,
      interpreters,
      interpreter,

      setCreateFormVisible,
      createSubmitLoading,
      createFormVisible,
      addInterpreter,
      createSubmit,

      updateSubmitLoading,
      updateFormVisible,
      updateFormCancel,
      editInterpreter,
      deleteInterpreter,
      updateSubmit,
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
