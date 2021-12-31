<template>
    <div class="main">
        <h1 class="title">
            {{t('page.user.login.form.title')}}
        </h1>
        <a-form :wrapper-col="{span:24}">
            <a-form-item label="" v-bind="validateInfos.username">
                <a-input v-model:value="modelRef.username" :placeholder="t('page.user.login.form-item-username')" @keyup.enter="handleSubmit">
                    <template #prefix><user-outlined /></template>
                </a-input>
            </a-form-item>
            <a-form-item label="" v-bind="validateInfos.password">
                <a-input-password v-model:value="modelRef.password" :placeholder="t('page.user.login.form-item-password')" @keyup.enter="handleSubmit">
                    <template #prefix><unlock-outlined /></template>
                </a-input-password>
            </a-form-item>
            <a-form-item>
                <a-button type="primary" class="submit" @click="handleSubmit" :loading="submitLoading">
                    {{t('page.user.login.form.btn-submit')}}
                </a-button>  
                <div class="text-align-right">
                    <router-link to="/user/register">
                        {{t('page.user.login.form.btn-jump')}}
                    </router-link>
                </div>              
            </a-form-item>

            <a-alert v-if="loginStatus === 'error' && !submitLoading" :message="t('page.user.login.form.login-error')" type="error" :show-icon="true" />

        </a-form>
    </div>
</template>
<script lang="ts">
import { computed, ComputedRef, defineComponent, reactive, Ref, ref, watch } from "vue";
import { useRouter } from 'vue-router';
import { useStore } from 'vuex';
import { useI18n } from "vue-i18n";

import { Props, validateInfos } from 'ant-design-vue/lib/form/useForm';
import { message, Form } from 'ant-design-vue';
const useForm = Form.useForm;

import { UserOutlined, UnlockOutlined } from '@ant-design/icons-vue';
import useI18nAntdFormVaildateInfos from '@/composables/useI18nAntdFormVaildateInfos';
import { LoginParamsType } from './data.d';
import { StateType as UserLoginStateType } from './store';

interface UserLoginSetupData {
    t: (key: string | number) => string;    
    resetFields: (newValues?: Props) => void;
    validateInfos: ComputedRef<validateInfos>;
    modelRef: LoginParamsType;
    submitLoading: Ref<boolean>;
    handleSubmit: (e: MouseEvent) => void;
    loginStatus: ComputedRef<"error" | "ok" | undefined>;
}

export default defineComponent({
    name: 'UserLogin',
    components: {
        UserOutlined,
        UnlockOutlined,
    },
    setup(): UserLoginSetupData {
        const router = useRouter();
        const { currentRoute } = router;
        const store = useStore<{userlogin: UserLoginStateType}>();
        const { t } = useI18n();

        // 表单值
        const modelRef = reactive<LoginParamsType>({
            username: '',
            password: ''
        });
        // 表单验证
        const rulesRef = reactive({
            username: [
                {
                    required: true,
                    message: 'page.user.login.form-item-username.required',
                },
            ],
            password: [
                {
                    required: true,
                    message: 'page.user.login.form-item-password.required',
                },
            ],            
        });
        // 获取表单内容
        const { resetFields, validate, validateInfos } = useForm(modelRef, rulesRef);
        // 登录loading
        const submitLoading = ref<boolean>(false);
        // 登录
        const handleSubmit = async (e: MouseEvent) => {
            e.preventDefault();
            submitLoading.value = true;
            try {
                const fieldsValue = await validate<LoginParamsType>();
                const res: boolean = await store.dispatch('userlogin/login',fieldsValue);        
                if (res === true) {
                    message.success(t('page.user.login.form.login-success'));
                    const { redirect, ...query } = currentRoute.value.query;
                    router.replace({
                        path: redirect as string || '/',
                        query: {
                            ...query
                        }
                    });
                }
            } catch (error) {
                // console.log('error', error);
                message.warning(t('app.global.form.validatefields.catch'));
            }
            submitLoading.value = false;
        };

        // 登录状态
        const loginStatus = computed<"ok" | "error" | undefined>(()=> store.state.userlogin.loginStatus);

        // 重置 validateInfos
        const validateInfosNew = useI18nAntdFormVaildateInfos(validateInfos);

        return {
            t,
            resetFields,
            validateInfos: validateInfosNew,
            modelRef,
            submitLoading,
            handleSubmit,
            loginStatus
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