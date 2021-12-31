<template>
    <router-link to="/" class="indexlayout-top-message">
      <BellOutlined :style="{ fontSize: '16px' }" />
      <a-badge
        class="indexlayout-top-message-badge"
        :count="message"
        :numberStyle="{ boxShadow: 'none', height: '15px', 'line-height': '15px' }"
      />
    </router-link>
</template>
<script lang="ts">
import { computed, ComputedRef, defineComponent, onMounted } from "vue";
import { useStore } from "vuex";
import { BellOutlined } from '@ant-design/icons-vue';
import { StateType as UserStateType } from "@/store/user";

interface RightTopMessageSetupData {
    message: ComputedRef<number>;
}

export default defineComponent({
    name: 'RightTopMessage',
    components: {
        BellOutlined
    },
    setup(): RightTopMessageSetupData {

        const store = useStore<{user: UserStateType}>();
        
        const message = computed<number>(()=> store.state.user.message);


        onMounted(()=>{
            store.dispatch("user/fetchMessage");
        })


        return {
            message
        }
    }
})
</script>