<template>
  <div class="main">
    <a-dropdown class="dropdown-list">

      <a class="t-link-btn" @click.prevent>
        <span class="name">设置</span>
        <span class="icon2"><icon-svg type="down"></icon-svg></span>
      </a>

      <template #overlay>
        <a-menu class="menu">
          <a-menu-item @click="setSite"><span class="t-link">禅道站点</span></a-menu-item>
          <a-menu-item @click="setEnv"><span class="t-link">运行环境</span></a-menu-item>
          <a-menu-item @click="setLang"><span class="t-link">界面语言</span></a-menu-item>
        </a-menu>
      </template>

    </a-dropdown>

    <a-modal v-model:visible="selectLangVisible" title="请选择语言">
      <div>
        <TopSelectLang></TopSelectLang>
      </div>
      <template #footer>
        <a-button key="back" @click="selectLangVisible=false">关闭</a-button>
      </template>
    </a-modal>

  </div>
</template>
<script lang="ts">
import {defineComponent, ref, Ref} from "vue";
import {useI18n} from "vue-i18n";
import IconSvg from "@/components/IconSvg";
import {useRouter} from "vue-router";

import TopSelectLang from "./TopSelectLang.vue";

interface SettingsSetupData {
  setSite: () => void;
  setEnv: () => void;
  setLang: () => void;

  selectLangVisible: Ref<boolean>
}

export default defineComponent({
  name: 'Settings',
  components: {
    IconSvg, TopSelectLang
  },
  setup(): SettingsSetupData {
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

    return {
      setSite,
      setEnv,
      setLang,
      selectLangVisible,
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