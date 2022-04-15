<template>
    <div id="right-top">
        <div class="right-top-top">

            <div @click="gotoSite" class="avatar-wrapper">
              <div v-if="profile.avatar" class="avatar avatar-img" :style="{'background-image': 'url('+profile.avatar+')'}">
              </div>

              <div v-if="!profile.avatar" class="avatar avatar-text">
                  {{nameFirstCap(profile)}}
              </div>
            </div>

            <div class="top-site-product-wrapper">
              <TopSiteProduct></TopSiteProduct>
            </div>

            <div class="menu-wrapper">
              <TopMenu :menuData="menuData" :belongTopMenu="belongTopMenu"></TopMenu>
            </div>

            <div class="settings-wrapper">
              <TopSettings />
            </div>
        </div>

      <TopNotify />
      <TopWebsocket />
    </div>
</template>

<script lang="ts">
import {computed, defineComponent, PropType, Ref, toRefs, ComputedRef} from "vue";
import { useI18n } from "vue-i18n";

import { BreadcrumbType, RoutesDataItem } from '@/utils/routes';

import useTopMenuWidth from "../composables/useTopMenuWidth";

import {useStore} from "vuex";
import {ZentaoData} from "@/store/zentao";
import {nameFirstCapDef} from "@/utils/string";

import TopMenu from './TopMenu.vue';
import TopSiteProduct from './TopSiteProduct.vue';
import TopSettings from './TopSettings.vue';
import TopNotify from './TopNotify.vue';
import TopWebsocket from './TopWebsocket.vue';

export default defineComponent({
    name: 'Top',
    components: {
      TopMenu, TopSiteProduct, TopSettings,
      TopNotify, TopWebsocket,
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
    setup(props) {
      const { t } = useI18n();

      const store = useStore<{ Zentao: ZentaoData }>();
      const profile = computed<any[]>(() => store.state.Zentao.profile);
      store.dispatch('Zentao/getProfile');

      const { topNavEnable } = toRefs(props);

      const { topMenuCon } = useTopMenuWidth(topNavEnable);
      const pathToId = (path) => {
        return path.replaceAll('/', 'menu-')
      }

      const gotoSite =() => {
        window.open('https://ztf.im','_blank')
      }
      const nameFirstCap = nameFirstCapDef

      return {
        t,
        profile,
        topMenuCon,
        pathToId,
        gotoSite,
        nameFirstCap,
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

    .avatar-wrapper {
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
      .avatar {
        height: 30px;
        width: 30px;
        background-size: cover;
        border-radius: 50%;
        margin-top: 8px;
        &.avatar-img {

        }
        &.avatar-text {
          line-height: 30px;
          background-color: hsl(55, 40%, 60%);
        }
      }

      .title {
        margin-left: 8px;
        padding-top: 3px;
        line-height: 28px;
      }
    }

    .top-site-product-wrapper {
      margin-right: 16px;
      width: 320px;
    }

    .menu-wrapper {
      flex: 1;
    }

    .settings-wrapper {
      width: 100px;
      text-align: right;
    }
  }
}
</style>