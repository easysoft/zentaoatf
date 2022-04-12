<template>
  <div class="main">
    <a-dropdown class="dropdown-list">

      <a class="t-white" @click.prevent>
        <span class="name">{{t('settings')}}</span>
        <span class="icon2"><icon-svg type="down"></icon-svg></span>
      </a>

      <template #overlay>
        <a-menu class="menu">
          <a-menu-item @click="setSite"><span class="t-link">{{t('zentao_site')}}</span></a-menu-item>
          <a-menu-item @click="setEnv"><span class="t-link">{{ t('interpreter') }}</span></a-menu-item>
          <a-menu-item @click="setLang"><span class="t-link">{{t('ui_lang')}}</span></a-menu-item>
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
import IconSvg from "@/components/IconSvg";
import {useRouter} from "vue-router";

import TopSelectLang from "./TopSelectLang.vue";

export default defineComponent({
  name: 'Settings',
  components: {
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

    return {
      t,
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