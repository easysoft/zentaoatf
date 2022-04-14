<template>
  <div class="main">
    <a-dropdown class="dropdown-list">

      <a class="t-white" @click.prevent>
        <span class="name"><SettingOutlined /></span>
        <span class="icon2"><icon-svg type="down"></icon-svg></span>
      </a>

      <template #overlay>
        <a-menu class="menu">
          <a-menu-item @click="setSite"><span class="t-link">{{t('zentao_site')}}</span></a-menu-item>
          <a-menu-item @click="setEnv"><span class="t-link">{{ t('interpreter') }}</span></a-menu-item>
          <a-menu-item @click="setLang"><span class="t-link">{{t('ui_lang')}}</span></a-menu-item>

          <template v-if="isElectron">
            <a-menu-divider />

            <a-menu-item @click="fullScreen">
              <span class="t-link">
                {{t('fullScreen')}}
                <FullscreenOutlined v-if="!fullScreenDef" />
                <FullscreenExitOutlined v-if="fullScreenDef" />
              </span>
            </a-menu-item>
            <a-menu-item @click="help">
              <span class="t-link">
                {{ t('help') }}
                <QuestionCircleOutlined />
              </span>
            </a-menu-item>
            <a-menu-item @click="exit">
              <span class="t-link">
                {{t('exit')}}
                <LogoutOutlined />
              </span>
            </a-menu-item>
          </template>
        </a-menu>
      </template>

    </a-dropdown>

    <a-modal v-model:visible="selectLangVisible" :title="t('select_ui_lang')">
      <div>
        <TopSelectLang></TopSelectLang>
      </div>
      <template #footer>
        <a-button @click="selectLangVisible=false" type="primary">{{t('close')}}</a-button>
      </template>
    </a-modal>

  </div>
</template>
<script lang="ts">
import {defineComponent, ref, Ref} from "vue";
import {useI18n} from "vue-i18n";
import {SettingOutlined, FullscreenOutlined, FullscreenExitOutlined, QuestionCircleOutlined, LogoutOutlined} from '@ant-design/icons-vue';
import IconSvg from "@/components/IconSvg";
import {useRouter} from "vue-router";

import TopSelectLang from "./TopSelectLang.vue";
import {getElectron} from "@/utils/comm";

export default defineComponent({
  name: 'Settings',
  components: {
    SettingOutlined, FullscreenOutlined, FullscreenExitOutlined, QuestionCircleOutlined, LogoutOutlined,
    IconSvg, TopSelectLang
  },
  setup() {
    const { t } = useI18n();
    const router = useRouter();

    const selectLangVisible = ref(false)

    const setSite = (): void => {
      console.log('setSite')
      router.push(`/site/list`)
    }
    const setEnv = (): void => {
      console.log('setEnv')
      router.push(`/interpreter/list`)
    }
    const setLang = (): void => {
      console.log('setLang')
      selectLangVisible.value = true
    }

    const isElectron = ref(getElectron())
    console.log(`isElectron ${isElectron.value}`)

    const fullScreenDef = ref(false)
    const fullScreen = (): void => {
      console.log('fullScreen')
      fullScreenDef.value = !fullScreenDef.value

      const remote = window.require('@electron/remote')
      const mainWin = remote.getCurrentWindow();
      mainWin.setFullScreen(!mainWin.isFullScreen());
    }
    const help = (): void => {
      console.log('help')
      const shell = window.require('@electron/remote').shell
      shell.openExternal('https://ztf.im');
    }
    const exit = (): void => {
      console.log('exit')

      const app = window.require('@electron/remote').app
      app.quit()
    }

    return {
      t,
      setSite,
      setEnv,
      setLang,
      selectLangVisible,

      isElectron,
      fullScreenDef,
      fullScreen,
      help,
      exit,
    }
  }
})
</script>
<style lang="less" scoped>
.main {
  .dropdown-list {
    display: inline-block;
    margin-right: 16px;
    padding-top: 13px;
    font-size: 15px !important;

    .name {
      margin-right: 5px;
    }

    .icon2 {
      .svg-icon {
        vertical-align: -3px !important;
      }
    }
  }
}
</style>