<template>
  <a-card>
    <template #title>
      {{t('edit_site')}}
    </template>

    <template #extra>
      <div class="opt">
        <a-button @click="back" type="link">{{ t('back') }}</a-button>
      </div>
    </template>

    <a-form :label-col="labelCol" :wrapper-col="wrapperCol">
      <a-form-item :label="t('name')" v-bind="validateInfos.name">
        <a-input v-model:value="modelRef.name"
                 @blur="validate('name', { trigger: 'blur' }).catch(() => {})" />
      </a-form-item>

      <a-form-item :label="t('zentao_url')" v-bind="validateInfos.url">
        <a-input v-model:value="modelRef.url"
                 @blur="validate('url', { trigger: 'blur' }).catch(() => {})" placeholder="https://zentao.site.com" />
      </a-form-item>
      <a-form-item :label="t('username')" v-bind="validateInfos.username">
        <a-input v-model:value="modelRef.username"
                 @blur="validate('username', { trigger: 'blur' }).catch(() => {})" placeholder="" />
      </a-form-item>
      <a-form-item :label="t('password')" v-bind="validateInfos.password">
        <a-input-password v-model:value="modelRef.password"
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
import {defineComponent, ref, reactive, computed, Ref,watch, ComputedRef} from "vue";
import {useRouter} from "vue-router";
import { useI18n } from "vue-i18n";

import { validateInfos } from 'ant-design-vue/lib/form/useForm';
import {Form, notification} from 'ant-design-vue';
const useForm = Form.useForm;

import {useStore} from "vuex";
import {StateType} from "@/views/site/store";

interface SiteFormSetupData {
  t: (key: string | number) => string;
  modelRef: ComputedRef;
  labelCol: any
  wrapperCol: any
  validate: any
  validateInfos: validateInfos
  submitForm:  () => void;
  resetFields:  () => void;
  back:  () => void;
}

export default defineComponent({
  name: 'SiteForm',
  components: {
  },
  setup(props): SiteFormSetupData {
    const { t } = useI18n();
    const router = useRouter();

    const store = useStore<{ Site: StateType }>();
    const modelRef = computed(() => store.state.Site.detailResult);

    const get = async (id: number): Promise<void> => {
      await store.dispatch('Site/get', id);
    }
    const id = +router.currentRoute.value.params.id
    get(id)

    const rules = reactive({
      name: [
        { required: true, message: t('pls_name'), trigger: 'blur' },
      ],
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

    const { resetFields, validate, validateInfos } = useForm(modelRef, rules);
    const submitForm = () => {
      validate()
        .then(() => {
          console.log(modelRef.value);
          store.dispatch('Site/save', modelRef.value).then((success) => {
            if (success) {
              notification.success({ message: t('save_success') });
              router.push(`/site/list`)
            }
          })
        })
      .catch((e) => {console.log('')})
    };

    const back = () => {
      console.log('back')
      router.push(`/site/list`)
    }

    return {
      t,
      labelCol: { span: 4 },
      wrapperCol: { span: 14 },
      modelRef,

      resetFields,
      validate,
      validateInfos,
      submitForm,
      back,
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
