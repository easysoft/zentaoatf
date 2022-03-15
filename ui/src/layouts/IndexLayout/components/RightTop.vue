<template>
    <div id="right-top">
        <div class="right-top-top">

            <div @click="gotoSite" class="logo-wrapper">
              <span class="logo"><img src="../../../assets/images/logo.png"></span>
            </div>

            <div class="top-project-wrapper">
              <RightTopProject class="top-select-project"></RightTopProject>
            </div>

            <div class="menu-wrapper">
              <RightTopMenu :menuData="menuData" :belongTopMenu="belongTopMenu"></RightTopMenu>
            </div>

            <div class="settings-wrapper">
              <RightTopSettings />
            </div>
        </div>
    </div>
</template>

<script lang="ts">
import { defineComponent, PropType, Ref, toRefs } from "vue";
import { useI18n } from "vue-i18n";

import { BreadcrumbType, RoutesDataItem } from '@/utils/routes';
import RightTopSettings from './RightTopSettings.vue';
import useTopMenuWidth from "../composables/useTopMenuWidth";
import RightTopProject from './RightTopProject.vue';
import RightTopMenu from './RightTopMenu.vue';

interface RightTopSetupData {
  t: (key: string | number) => string;
  topMenuCon: Ref;
  pathToId: (val) => void
  gotoSite: (val) => void
}

export default defineComponent({
    name: 'RightTop',
    components: {
      RightTopSettings,
      RightTopProject, RightTopMenu,
    },
    props: {
      collapsed: {
        type: Boolean,
        default: false
      },
      topNavEnable: {
        type: Boolean,
        default: true
      },
      belongTopMenu: {
        type: String,
        default: ''
      },
      toggleCollapsed: {
        type: Function as PropType<() => void>
      },
      breadCrumbs: {
        type: Array as PropType<BreadcrumbType[]>,
        default: () => {
          return [];
        }
      },
      menuData: {
        type: Array as PropType<RoutesDataItem[]>,
        default: () => {
          return [];
        }
      },
      routeItem: {
        type: Object as PropType<RoutesDataItem>,
        required: true
      }
    },
    setup(props): RightTopSetupData {
      const { t } = useI18n();
      const { topNavEnable } = toRefs(props);

      const { topMenuCon } = useTopMenuWidth(topNavEnable);
      const pathToId = (path) => {
        return path.replaceAll('/', 'menu-')
      }

      const gotoSite =() => {
        window.open('https://ztf.im','_blank')
      }

      return {
        t,
        topMenuCon,
        pathToId,
        gotoSite,
      }
    }
})
</script>

<style lang="less" scoped>
@import '../../../assets/css/global.less';

#right-top {
  width: 100%;
  height: @headerHeight;
  box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);
  z-index: 9;

  .right-top-top {
    display: flex;
    width: 100%;
    height: @headerHeight;
    background-color: @menu-dark-bg;
    color: #FFFFFF;

    .logo-wrapper {
      margin: 0 30px 0 30px;
      width: auto;
      height: @headerHeight;
      line-height: @headerHeight;
      text-align: center;
      cursor: pointer;
      &:hover {
        background-color: @menu-dark-bg;
        color: @menu-dark-highlight-color;
      }
      .logo {
        img {
          height: 30px;
          vertical-align: -10px
        }
      }
      .title {
        display: inline-block;
        margin-left: 8px;
        padding-top: 3px;
        line-height: 28px;
      }
    }

    .top-project-wrapper {
      margin-right: 16px;
      width: 320px;
    }

    .menu-wrapper {
      flex: 1;
    }

    .settings-wrapper {
      display: flex;
      width: 90px;
    }
  }
}
</style>