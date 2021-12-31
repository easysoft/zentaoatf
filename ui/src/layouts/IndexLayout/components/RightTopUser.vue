<template>
    <a-dropdown>
        <a class="indexlayout-top-usermenu ant-dropdown-link" @click="e => e.preventDefault()">
            {{currentUser.name}} <DownOutlined />
        </a>
        <template #overlay>
            <a-menu @click="onMenuClick">
                <a-menu-item key="userinfo">
                    {{t('index-layout.topmenu.userinfo')}}
                </a-menu-item>
                <a-menu-item key="logout">
                   {{t('index-layout.topmenu.logout')}}
                </a-menu-item>
            </a-menu>
        </template>
    </a-dropdown>
</template>
<script lang="ts">
import { DownOutlined } from '@ant-design/icons-vue';
import { computed, ComputedRef, defineComponent } from "vue";
import { useStore } from "vuex";
import { useRouter } from "vue-router";
import { useI18n } from "vue-i18n";
import { StateType as UserStateType, CurrentUser } from "@/store/user";
interface RightTopUserSetupData {
    t: (key: string | number) => string;
    currentUser: ComputedRef<CurrentUser>;
    onMenuClick: (event: any) => Promise<void>;
}
export default defineComponent({
    name: 'RightTopUser',
    components: {
        DownOutlined
    },
    setup(): RightTopUserSetupData {
        const store = useStore<{user: UserStateType}>();
        const router = useRouter();
        const { t } = useI18n();


        // 获取当前登录用户信息
        const currentUser = computed<CurrentUser>(()=> store.state.user.currentUser);

        // 点击菜单
        const onMenuClick = async (event: any) => {
            const { key } = event;

            if (key === 'logout') {
                const res: boolean = await store.dispatch('user/logout');
                if(res === true) {
                    router.replace({
                        path: '/user/login',
                        query: {
                            redirect: router.currentRoute.value.path,
                            ...router.currentRoute.value.query
                        }
                    })
                }
            }

        }



        return {
            t,
            currentUser,
            onMenuClick
        }
    }
})
</script>