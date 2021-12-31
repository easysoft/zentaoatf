<template>
    <slot v-if="isPermission"></slot>
    <slot v-else name="otherwise">
        <a-result  status="403" title="403" sub-title="Sorry, you are not authorized to access this page.">
            <template #extra>
                <router-link to="/">
                    <a-button type="primary">
                        Back Home
                    </a-button>
                </router-link>
            </template>
        </a-result>
    </slot>
    
</template>
<script lang="ts">
import { computed, ComputedRef, defineComponent, PropType } from "vue";
import { useStore } from "vuex";
import { StateType as UserStateType } from "@/store/user";
import { hasPermissionRouteRoles } from "@/utils/routes";


interface PermissionSetupData {
    isPermission: ComputedRef<boolean>;
}

export default defineComponent({
    name: 'Permission',
    props: {
        roles: {
            type: [String , Array] as PropType<string[] | string>,
        }
    },
    setup(props): PermissionSetupData {
        const store = useStore<{user: UserStateType}>(); 

        // 是否有权限
        const isPermission = computed(()=> hasPermissionRouteRoles(store.state.user.currentUser.roles, props.roles));

        return {
            isPermission
        }

    }
})
</script>