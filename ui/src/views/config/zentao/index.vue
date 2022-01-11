<template>

  <a-card title="修改配置">
    <a-form :label-col="labelCol" :wrapper-col="wrapperCol">
      <a-form-item label="禅道地址" v-bind="validateInfos.url">
        <a-input v-model:value="model.url"
                 @blur="validate('url', { trigger: 'blur' }).catch(() => {})" placeholder="https://zentao.site.com" />
      </a-form-item>
      <a-form-item label="用户名" v-bind="validateInfos.username">
        <a-input v-model:value="model.username" placeholder="" />
      </a-form-item>
      <a-form-item label="密码" v-bind="validateInfos.password">
        <a-input v-model:value="model.password" placeholder="" />
      </a-form-item>

      <a-form-item :wrapper-col="{ span: 14, offset: 4 }">
        <a-button type="primary" @click.prevent="submitForm">保存</a-button>
        <a-button style="margin-left: 10px" @click="resetFields">重置</a-button>
      </a-form-item>
    </a-form>
  </a-card>

</template>
<script lang="ts">
import { defineComponent, ref, reactive } from "vue";
import { useI18n } from "vue-i18n";

import { Props, validateInfos } from 'ant-design-vue/lib/form/useForm';
import {message, Form, notification} from 'ant-design-vue';
const useForm = Form.useForm;

import { Config } from '../data.d';
import {saveConfig} from "@/services/project";

interface ConfigFormSetupData {
  formRef: any
  model: Config
  rules: any
  labelCol: any
  wrapperCol: any
  validate: any
  validateInfos: validateInfos
  submitForm:  () => void;
  resetFields:  () => void;
}

export default defineComponent({
  name: 'ConfigFormForm',
  components: {
  },
  setup(props): ConfigFormSetupData {
    const { t } = useI18n();

    const formRef = ref();

    const model = reactive<Config>({
      url: '',
      username: '',
      password: '',
    });

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
                message: `保存配置错误`,
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
      formRef,
      labelCol: { span: 4 },
      wrapperCol: { span: 14 },
      model,
      rules,
      resetFields,
      validate,
      validateInfos,
      submitForm,
    }

  }
})
</script>