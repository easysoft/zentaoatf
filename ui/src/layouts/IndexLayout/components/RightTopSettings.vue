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
  </div>
</template>
<script lang="ts">
import {defineComponent, WritableComputedRef} from "vue";
import {useI18n} from "vue-i18n";
import IconSvg from "@/components/IconSvg";
import {useRouter} from "vue-router";

interface SelectLangSetupData {
  setSite: () => void;
  setEnv: () => void;
  setLang: () => void;
}

export default defineComponent({
  name: 'SelectLang',
  components: {
    IconSvg
  },
  setup(): SelectLangSetupData {
    const router = useRouter();

    const setSite = (): void => {
      console.log('setSite')
      router.push(`/site/list`)
    }
    const setEnv = (): void => {
      console.log('setEnv')
    }
    const setLang = (): void => {
      console.log('setLang')
    }

    return {
      setSite,
      setEnv,
      setLang,
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