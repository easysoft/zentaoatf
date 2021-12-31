<template>
    <div class="main">
        <h1 class="title">
            {{t('page.user.register.form.title')}}
        </h1>
        <a-form :wrapper-col="{span:24}">
            <a-form-item label="" v-bind="validateInfos.username">
                <a-input v-model:value="modelRef.username" :placeholder="t('page.user.register.form-item-username')" @keyup.enter="handleSubmit" />
            </a-form-item>
            <a-form-item label="" v-bind="validateInfos.password">
                <a-input-password v-model:value="modelRef.password" :placeholder="t('page.user.register.form-item-password')" @keyup.enter="handleSubmit" />
            </a-form-item>
            <a-form-item label="" v-bind="validateInfos.confirm">
                <a-input-password v-model:value="modelRef.confirm" :placeholder="t('page.user.register.form-item-confirmpassword')" @keyup.enter="handleSubmit" />
            </a-form-item>
            <a-form-item>
                <a-button type="primary" class="submit" @click="handleSubmit" :loading="submitLoading">
                    {{t('page.user.register.form.btn-submit')}}
                </a-button>  
                <div class="text-align-right">
                    <router-link to="/user/login">
                        {{t('page.user.register.form.btn-jump')}}
                    </router-link>
                </div>              
            </a-form-item>

            <a-alert v-if="errorMsg !== '' && typeof errorMsg !== 'undefined' &&  !submitLoading" :message="errorMsg" type="error" :show-icon="true" />
            
        </a-form>
    </div>
</template>
<script lang="ts">
import { computed, ComputedRef, defineComponent, reactive, Ref, ref } from "vue";
import { useRouter } from 'vue-router';
import { useStore } from 'vuex';
import { useI18n } from "vue-i18n";

import { Props, validateInfos } from 'ant-design-vue/lib/form/useForm';
import { message, Form } from 'ant-design-vue';
const useForm = Form.useForm;

import useI18nAntdFormVaildateInfos from '@/composables/useI18nAntdFormVaildateInfos';
import { RegisterParamsType } from "./data.d";
import { StateType as RegisterStateType } from "./store";

interface UserRegisterSetupData {
    t: (key: string | number) => string;
    validateInfos: ComputedRef<validateInfos>;
    modelRef: RegisterParamsType;
    submitLoading: Ref<boolean>;
    handleSubmit: (e: MouseEvent) => void;
    errorMsg: ComputedRef<string | undefined>; 
}

export default defineComponent({
    name: 'UserRegister',
    setup(): UserRegisterSetupData {
        const router = useRouter();
        const store = useStore<{userregister: RegisterStateType}>();
        const { t } = useI18n();

        // 表单值
        const modelRef = reactive<RegisterParamsType>({
            username: '',
            password: '',
            confirm: ''
        });
        // 表单验证
        const rulesRef = reactive({
            username: [
                {
                    required: true,
                    message: 'page.user.register.form-item-username.required',
                },
            ],
            password: [
                {
                    required: true,
                    message: 'page.user.register.form-item-password.required',
                },
            ],
            confirm: [
                {
                    validator: (rule: any, value: string, callback: any) => {
                        if (value === '') {
                            return Promise.reject('page.user.register.form-item-password.required');
                        } else if (value !== modelRef.password) {
                            return Promise.reject("page.user.register.form-item-confirmpassword.compare");
                        } else {
                            return Promise.resolve();
                        }
                    }
                }
            ],          
        });
        // 获取表单内容
        const { validate, validateInfos } = useForm(modelRef, rulesRef);
        // 注册loading
        const submitLoading = ref<boolean>(false);
        // 注册
        const handleSubmit = async (e: MouseEvent) => {
            e.preventDefault();
            submitLoading.value = true;
            try {
                const fieldsValue = await validate<RegisterParamsType>();
                const res: boolean = await store.dispatch('userregister/register',fieldsValue);                
                if (res === true) {
                    message.success(t('page.user.register.form.register-success'));
                    router.replace('/user/login');
                }
            } catch (error) {
                // console.log('error', error);
                message.warning(t('app.global.form.validatefields.catch'));
            }
            submitLoading.value = false;
        };        

        // 重置 validateInfos
        const validateInfosNew = useI18nAntdFormVaildateInfos(validateInfos);

         // 注册状态
        const errorMsg = computed<string | undefined>(()=> store.state.userregister.errorMsg);


        return {
            t,
            modelRef,
            validateInfos: validateInfosNew,
            submitLoading,
            handleSubmit,
            errorMsg
        }

    }
})
</script>
<style lang="less" scoped>
.main {
  flex: none;
  width: 320px;
  padding: 36px;
  margin: 0 auto;
  border-radius: 4px;
  background-color: rgba(255, 255, 255, 0.2);
  .title {
    font-size: 28px;
    margin-top: 0;
    margin-bottom: 30px;
    text-align: center;
    color: #ffffff;
    /* background-image:-webkit-linear-gradient(right,#FFFFFF,#009688, #FFFFFF); 
        -webkit-background-clip: text; 
        -webkit-text-fill-color:transparent; */
  }
  .submit {
    width: 100%;
  }
}
</style>