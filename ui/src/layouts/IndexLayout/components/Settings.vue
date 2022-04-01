<template>
  <div
      :style="{
          position: 'fixed',
          display: 'block',
          width: '45px',
          height: '45px',
          lineHeight: '48px',
          right: '0',
          top: '30%',
          backgroundColor: '#222834',
          textAlign: 'center',
          cursor: 'pointer',
          borderRadius: '5px 0 0 5px',
        }"
      @click="show"
  >
    <SettingOutlined :style="{ fontSize: '20px', color: '#fcfcfc' }"/>
  </div>
</template>
<script lang="ts">
import {computed, ComputedRef, defineComponent, Ref, ref} from "vue";
import {useStore} from 'vuex';
import {StateType as GlobalStateType} from '@/store/global';
import {SettingOutlined} from '@ant-design/icons-vue';
import {useI18n} from "vue-i18n";

interface SettingsSetupData {
  t: (key: string | number) => string;
  visible: Ref<boolean>;
  close: () => void;
  show: () => void;
  topNavEnable: ComputedRef<boolean>;
  onChangeTopNavEnable: (v: boolean) => void;
  headFixed: Ref<boolean>;
  onChangeHeadFixed: (v: boolean) => void;
  disabledHeadFixed: Ref<boolean>;
}

export default defineComponent({
  name: 'Settings',
  components: {
    SettingOutlined
  },
  setup(): SettingsSetupData {
    const {t} = useI18n();

    const store = useStore<{ Global: GlobalStateType }>();

    const visible = ref<boolean>(false);
    // 关闭
    const close = (): void => {
      visible.value = false;
    }
    // 显示
    const show = (): void => {
      visible.value = true;
    }

    // 固定右侧头部
    const disabledHeadFixed = ref<boolean>(true);
    const headFixed = computed<boolean>(() => store.state.Global.headFixed);
    const onChangeHeadFixed = (v: boolean): void => {
      store.commit('Global/setHeadFixed', v);
    }

    // 启用顶部导航
    const topNavEnable = computed<boolean>(() => store.state.Global.topNavEnable);
    const onChangeTopNavEnable = (v: boolean): void => {
      store.commit('Global/setTopNavEnable', v);

      if (v) {
        disabledHeadFixed.value = true;
        onChangeHeadFixed(true);
      } else {
        disabledHeadFixed.value = false;
      }

    }

    return {
      t,
      visible,
      close,
      show,
      topNavEnable,
      onChangeTopNavEnable,
      headFixed,
      onChangeHeadFixed,
      disabledHeadFixed
    }

  }
})
</script>