<template>

  <a-card title="修改配置">
    <a-form :label-col="labelCol" :wrapper-col="wrapperCol">
      <a-form-item label="禅道地址" v-bind="validateInfos.url">
        <a-input v-model:value="model.url"
                 @blur="validate('url', { trigger: 'blur' }).catch(() => {})" placeholder="https://zentao.site.com" />
      </a-form-item>
      <a-form-item label="用户名" v-bind="validateInfos.username">
        <a-input v-model:value="model.username"
                 @blur="validate('username', { trigger: 'blur' }).catch(() => {})" placeholder="" />
      </a-form-item>
      <a-form-item label="密码" v-bind="validateInfos.password">
        <a-input v-model:value="model.password"
                 @blur="validate('password', { trigger: 'blur' }).catch(() => {})" placeholder="" />
      </a-form-item>

      <a-form-item label="执行器">
        <div>
          <a-row class="interpreter-header">
            <a-col :span="10" class="t-center t-bord">语言</a-col>
            <a-col :span="10" class="t-center t-bord">解析器</a-col>
            <a-col :span="4" class="t-center t-bord">
              <a-button @click="setCreateFormVisible(true)" type="link" size="small">新建</a-button>
            </a-col>
          </a-row>
          <a-row  v-for="item in interpreters" :key="item.id" class="interpreter-item">
            <a-col :span="10" class="t-center t-bord">
              {{item.lang}}
            </a-col>
            <a-col :span="10" class="t-center t-bord">
              {{item.value}}
            </a-col>
            <a-col :span="4" class="t-center t-bord">
              <a-button type="link" size="small">编辑</a-button>
            </a-col>
          </a-row>
        </div>
      </a-form-item>

      <a-form-item :wrapper-col="{ span: 14, offset: 4 }">
        <a-button type="primary" @click.prevent="submitForm">保存</a-button>
        <a-button style="margin-left: 10px" @click="resetFields">重置</a-button>
      </a-form-item>
    </a-form>

    <create-form
        v-if="createFormVisible===true"
        :visible="createFormVisible"
        :onCancel="() => setCreateFormVisible(false)"
        :onSubmitLoading="createSubmitLoading"
        :onSubmit="createSubmit"
    />

    <update-form
        v-if="updateFormVisible===true"
        :visible="updateFormVisible"
        :values="interpreter"
        :onCancel="() => updateFormCancel()"
        :onSubmitLoading="updateSubmitLoading"
        :onSubmit="updateSubmit"
    />

  </a-card>

</template>
<script lang="ts">
import {defineComponent, ref, reactive, computed, watch, ComputedRef, Ref} from "vue";
import { useI18n } from "vue-i18n";

import { Props, validateInfos } from 'ant-design-vue/lib/form/useForm';
import {message, Form, notification} from 'ant-design-vue';
import _ from "lodash";
const useForm = Form.useForm;

import {Config} from './data.d';
import {saveConfig} from "./service";
import {useStore} from "vuex";
import {ProjectData} from "@/store/project";
import CreateForm from './create.vue';
import UpdateForm from './update.vue';

interface ConfigFormSetupData {
  currConfig: any
  model: Partial<Config>
  rules: any
  labelCol: any
  wrapperCol: any
  validate: any
  validateInfos: validateInfos
  submitForm:  () => void;
  resetFields:  () => void;

  interpreters: Ref
  interpreter: Ref

  createFormVisible: Ref<boolean>;
  setCreateFormVisible:  (val: boolean) => void;
  createSubmitLoading: Ref<boolean>;
  createSubmit: (values: any, resetFields: (newValues?: Props | undefined) => void) => Promise<void>;

  updateFormVisible: Ref<boolean>;
  updateFormCancel:  () => void;
  updateSubmitLoading: Ref<boolean>;
  updateSubmit:  (values: any, resetFields: (newValues?: Props | undefined) => void) => Promise<void>;
}

export default defineComponent({
  name: 'ConfigForm',
  components: {
    CreateForm, UpdateForm,
  },
  setup(props): ConfigFormSetupData {
    const { t } = useI18n();

    const store = useStore<{ project: ProjectData }>();
    const currConfig = computed<any>(() => store.state.project.currConfig);

    let model = reactive<Partial<Config>>(currConfig.value);
    watch(currConfig,(currConfig)=> {
      _.extend(model, currConfig)
    })

    const rules = reactive({
      url: [
        {
          required: true,
          trigger: 'blur',
          validator: async (rule: any, value: string) => {
            if (value === '' || !value) {
              throw new Error('请输入禅道地址');
            }

            const regx = /(http?|https):\/\/[-A-Za-z0-9+&@#/%?=~_|!:,.;]+[-A-Za-z0-9+&@#/%=~_|]/i
            if (!regx.test(value)) {
              throw new Error('请输入http或https开头的禅道地址');
            }
          }
        },
      ],
      username: [
        { required: true, message: '请输入用户名', trigger: 'blur' },
      ],
      password: [
        { required: true, message: '请输入密码', trigger: 'blur' },
      ],
    });

    const { resetFields, validate, validateInfos } = useForm(model, rules);

    const submitForm = () => {
      validate()
        .then(() => {
          console.log(model);
          saveConfig(model).then((json) => {
            console.log('json', json)
            if (json.code === 0) {
              notification.success({
                message: `保存配置成功`,
              });
            } else {
              notification.error({
                message: `保存配置失败`,
                description: json.msg,
              });
            }
          })
        })
        .catch(err => {
          console.log('error', err);
        });
    };

    // 解析器
    const interpreters = reactive<any>({})
    const interpreter = reactive<any>({})
    // 新建
    const createFormVisible = ref<boolean>(false);
    const setCreateFormVisible = (val: boolean) => {
      createFormVisible.value = val;
    };
    const createSubmitLoading = ref<boolean>(false);
    const createSubmit = async (values: any, resetFields: (newValues?: Props | undefined) => void) => {
      createSubmitLoading.value = true;

      console.log(values)

      createSubmitLoading.value = false;
      setCreateFormVisible(false)
    }

    // 更新
    const updateFormVisible = ref<boolean>(false);
    const setUpdateFormVisible = (val: boolean) => {
      updateFormVisible.value = val;
    }
    const updateFormCancel = () => {
      setUpdateFormVisible(false);
    }
    const updateSubmitLoading = ref<boolean>(false);
    const updateSubmit = async (values: any, resetFields: (newValues?: Props | undefined) => void) => {
      updateSubmitLoading.value = true;

      console.log(values)

      updateSubmitLoading.value = false;
      setUpdateFormVisible(false)
    }

    return {
      labelCol: { span: 4 },
      wrapperCol: { span: 14 },
      currConfig,
      model,

      rules,
      resetFields,
      validate,
      validateInfos,
      submitForm,

      interpreters,
      interpreter,

      createFormVisible,
      setCreateFormVisible,
      createSubmitLoading,
      createSubmit,

      updateFormVisible,
      updateFormCancel,
      updateSubmitLoading,
      updateSubmit,
    }

  }
})
</script>

<style lang="less" scoped>

.interpreter-header {
  margin: 5px 50px;
  padding-bottom: 6px;
  border-bottom: 1px solid #f0f0f0;
}

.interpreter-item {
  margin: 5px 50px;

}

</style>
