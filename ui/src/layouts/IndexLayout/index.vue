<template>
    <div id="indexlayout">
        <div
          id="indexlayout-right"
          :class="{'fiexd-header': headFixed}"
        >
            <right-top
              :collapsed="collapsed"
              :topNavEnable="topNavEnable"
              :belongTopMenu="belongTopMenu"
              :toggleCollapsed="toggleCollapsed"
              :breadCrumbs="breadCrumbs"
              :menuData="permissionMenuData"
              :routeItem="routeItem"
            >              
            </right-top>
            <div class="indexlayout-right-main">
              <router-view></router-view>
            </div>
        </div>

    </div>
</template>
<script lang="ts">
import { defineComponent, computed, ComputedRef, watch, ref, Ref, nextTick } from "vue";
import { useStore } from 'vuex';
import { useRoute } from 'vue-router';
import { StateType as GlobalStateType } from '@/store/global';
import { 
  vueRoutes, RoutesDataItem, getRouteItem, getSelectLeftMenuPath, 
  formatRoutePathTheParents, getRouteBelongTopMenu, getBreadcrumbRoutes, 
  BreadcrumbType
} from '@/utils/routes';
import { mergeUnique as ArrayMergeUnique } from '@/utils/array';
import useTitle from '@/composables/useTitle';
import IndexLayoutRoutes from './routes';
import RightTop from '@/layouts/IndexLayout/components/RightTop.vue';

interface IndexLayoutSetupData {
  collapsed: ComputedRef<boolean>;
  toggleCollapsed: () => void;
  topNavEnable: ComputedRef<boolean>;
  belongTopMenu: ComputedRef<string>;
  headFixed: ComputedRef<boolean>;  
  selectedKeys: ComputedRef<string[]>;
  leftOpenKeys: Ref<string[]>;
  breadCrumbs: ComputedRef<BreadcrumbType[]>;
  permissionMenuData: ComputedRef<RoutesDataItem[]>;
  onOpenChange: (key: any) => void;
  routeItem: ComputedRef<RoutesDataItem>;
}

export default defineComponent({
    name: 'IndexLayout',
    components: {
        RightTop,
    },
    setup(): IndexLayoutSetupData {
      const store = useStore<{
        global: GlobalStateType;
      }>(); 
      const route = useRoute();

      // 所有菜单路由
      const menuData: RoutesDataItem[] = vueRoutes(IndexLayoutRoutes);

      // 当前路由 item
      const routeItem = computed<RoutesDataItem>(()=> getRouteItem(route.path, menuData));

      // 有权限的菜单
      const permissionMenuData = computed<RoutesDataItem[]>(()=> { return menuData });

      // 当前路由的顶部菜单path
      const belongTopMenu = computed<string>(()=>getRouteBelongTopMenu(routeItem.value))

      // 当前路由的父路由path[]
      const routeParentPaths = computed<string[]>(()=>formatRoutePathTheParents(routeItem.value.path));

      // 收缩左侧
      const collapsed = computed<boolean>(()=> store.state.global.collapsed);
      const toggleCollapsed = (): void => {
        store.commit('global/changeLayoutCollapsed', !collapsed.value);
      }

      // 右侧顶部导航是否开启
      const topNavEnable = computed<boolean>(()=> store.state.global.topNavEnable);

      // 右侧顶部是否固定
      const headFixed = computed<boolean>(()=> store.state.global.headFixed);

      // 左侧选择菜单key
      const selectedKeys = computed<string[]>(()=>{
        const selectedKey = getSelectLeftMenuPath(routeItem.value);
        return [selectedKey]
      });

      // 左侧展开菜单keys
      const leftOpenKeys = ref<string[]>(routeParentPaths.value);
      watch([routeParentPaths,collapsed],()=>{
        if (routeParentPaths.value.length && !collapsed.value) {
          leftOpenKeys.value = ArrayMergeUnique<string>(leftOpenKeys.value, routeParentPaths.value);
        } else {
          nextTick(()=> {
            leftOpenKeys.value =[];
          })
        }
      })

      // 左侧菜单展开收起
      const onOpenChange = (key: any) => {
        leftOpenKeys.value = key;
      }

      // 面包屑导航
      const breadCrumbs = computed<BreadcrumbType[]>(
          () => getBreadcrumbRoutes(routeItem.value,routeParentPaths.value, menuData));

      // 设置title
      useTitle(routeItem);


      return {
        collapsed,
        toggleCollapsed,
        topNavEnable,
        belongTopMenu,
        headFixed, 
        selectedKeys,
        leftOpenKeys,
        breadCrumbs,
        permissionMenuData,
        onOpenChange,
        routeItem
      }
    }
})
</script>
<style lang="less">
@import '../../assets/css/variables.less';
#indexlayout {
  display: flex;
  height: 100vh;
  overflow: hidden;
}
#indexlayout-right {
  position: relative;
  flex: 1;
  overflow: hidden;
  background-color: @mainBgColor;
  &.fiexd-header {
    display: flex;
    flex-direction: column;
    .indexlayout-right-main {
      flex: 1;
      overflow: auto;
    }
  }
}
.indexlayout-main-conent {
  margin: 24px;
  position: relative;
}
</style>