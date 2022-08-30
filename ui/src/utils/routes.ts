import { RouteRecordRaw,  RouteLocationNormalizedLoaded } from 'vue-router';

/**
 * Route utils
 */
/**
 * 面包屑类型
 */
 export interface BreadcrumbType {
  // 标题，路由在菜单、浏览器title 或 面包屑中展示的文字，目前可以使用locales
  title: string;
  // 路由地址或外链
  path: string;
}

/**
 * tab导航存储规则类型
 */
export type TabNavType = 'path' | 'querypath';

declare module 'vue-router' {
  /**
   * 自定义补充扩展 - 路由 - 类型字段
   */
  interface _RouteRecordBase {
      // 菜单中是否隐藏
      hidden?: boolean;
      // 图标的名称，显示在菜单标题前
      icon?: string;
      // 权限控制，页面角色(您可以设置多个角色)
      roles?: string[];
      // 标题，路由在菜单、浏览器title 或 面包屑中展示的文字，目前可以使用locales
      title: string;
      /**
       * 设置tab导航存储规则类型
       *    1、默认不配置按照path(route.path)规则
       *    2、querypath：path + query (route.path+route.query) 规则
       */
      tabNavType?: TabNavType ;
      /**
       * 设置该字段，则在关闭当前tab页时，作为关闭前的钩子函数
       * @param close 关闭回调函数
       */
      tabNavCloseBefore?: (close: ()=>void)=> void;
      /**
        * 左侧菜单选中，如果设置路径，侧栏将突出显示你设置的路径对应的侧栏导航
        *   1、（默认 route.path），此参数是为了满足特殊页面特殊需求，
        *   2、如：详情页等选中侧栏导航或在模块A下面的页面，想选模块B为导航选中状态
        */
      selectLeftMenu?: string;
      /**
        * 所属顶级菜单,当顶级菜单存在时，用于选中顶部菜单，与侧栏菜单切换
        *   1、三级路由此参数的作用是选中顶级菜单
        *   2、二级路由此参数的作用是所属某个顶级菜单的下面，两个层级的必须同时填写一致，如果path设置的是外链，此参数必填
        *   3、(默认不设置 path.split('/')[0])，此参数是为了满足特殊页面特殊需求
        */
      belongTopMenu?: string;
  }
}

/**
 * 自定义重命名 - 路由类型
 */
export type RoutesDataItem = RouteRecordRaw;

import { isExternal } from './validate';
import { equalObject } from "./object";

/**
 * 获取 route
 * @param pathname path
 * @param routesData routes
 */
export const getRouteItem = (pathname: string, routesData: RoutesDataItem[]): RoutesDataItem => {

  let item: RoutesDataItem = { title: '', path: '', redirect: '', roles: [] };

  for (let index = 0, len = routesData.length; index < len; index += 1) {
    const element = routesData[index];

    if (element.path === pathname) {
      item = element;
      break;
    } else if (element.path.indexOf(':') > 0)  {
      const reg = /(:[^/]+)/gmi
      const path = element.path.replaceAll(reg,".+")
      const pass = new RegExp(path).test(pathname)

      if (pass) {
        // console.log('----', element.path, pathname)
        item = element;
        break;
      }
    }

    if (element.children) {
      item = getRouteItem(pathname, element.children);
      if (item.path !== '') {
        break;
      }
    }
  }

  return item;
};

/**
 * 根据 hidden 判断是否有数据子集
 * @param children RoutesDataItem[]
 */
export const hasChildRoute = (children: RoutesDataItem[]): boolean => {
  const showChildren = children.filter(item => {
    if (item.hidden) {
      return false;
    }
    return true;
  });
  if (showChildren.length > 0) {
    return true;
  }
  return false;
};

/**
 * 根据路由 path 格式化 - 获取 父path
 * @param pathname path
 * @param separator 路由分割符- 默认 /
 */
export const formatRoutePathTheParents = (pathname: string, separator = '/'): string[] => {
  const arr: string[] = [];
  if (!pathname || pathname === '') {
    return arr;
  }

  const pathArr = pathname.split(separator);
  for (let index = 1, len = pathArr.length - 1; index < len; index += 1) {
    arr.push(pathArr.slice(0, index + 1).join(separator));
  }

  return arr;
};

/**
 * 根据父path 设置当前项 path
 * @param pathname path
 * @param parentPath 父path - 默认 /
 * @param headStart 路由起始头 - 默认 /
 */
export const setRoutePathForParent = (pathname: string, parentPath = '/', headStart = '/'): string => {
  if (isExternal(pathname)) {
    return pathname;
  }

  return pathname.substring(0, headStart.length) === headStart
    ? pathname
    : `${parentPath}/${pathname}`;
};

/**
 * 根据路由 pathname 数组 - 返回对应的 route 数组
 * @param pathname path[]
 * @param routesData routes
 */
export const getPathsTheRoutes = ( pathname: string[], routesData: RoutesDataItem[]): RoutesDataItem[] => {
  const routeItem: RoutesDataItem[] = [];

  for (let index = 0, len = pathname.length; index < len; index += 1) {
    const element = pathname[index];
    const item = getRouteItem(element, routesData);
    if (item.path !== '') {
      routeItem.push(item);
    }
  }

  return routeItem;
};

/**
 * 获取当前路由选中的侧边栏菜单path
 * @param route route
 */
export const getSelectLeftMenuPath = (route: RoutesDataItem): string => {
  return route.selectLeftMenu || route.path;
};

/**
 * 获取当前路由对应的顶部菜单path
 * @param route route
 */
export const getRouteBelongTopMenu = (route: RoutesDataItem): string => {
  if (route.belongTopMenu) {
    return route.belongTopMenu;
  }
  return `/${route.path.split('/')[1]}`;
};


/**
 * 格式化返回 vue 路由, 主要处理特殊情况
 * @param routesData routes
 * @param parentPath 父path - 默认 /
 * @param headStart 路由起始头 - 默认 /
 */
export const vueRoutes = (routesData: RoutesDataItem[], parentPath = '/', headStart = '/'): RoutesDataItem[] => {
  return routesData.map(item => {
    const { children, ...other } = item;
    const itemChildren = children || [];
    const newItem: RoutesDataItem = {children: [], ...other };
    newItem.path = setRoutePathForParent(newItem.path, parentPath, headStart);

    if (item.children) {
      newItem.children = [
        ...vueRoutes(itemChildren, newItem.path, headStart),
      ];
    }

    return newItem;
  });
};

/**
 * 批量设置route.meta值
 * @param routesData routes
 */
 export const routesSetMeta = (routesData: RoutesDataItem[]): RoutesDataItem[] => {
  return routesData.map(item => {
    const { children, tabNavType, meta, ...other } = item;
    const newItem: RoutesDataItem = {
      meta: {
        ...meta,

        // 自定义设置的 meta 值 S
        tabNavType: tabNavType || 'path',
        // 自定义设置的 meta 值 E
      },
      children: [],
      ...other
     };

    if (item.children) {
      const itemChildren = children || [];
      newItem.children = [
        ...routesSetMeta(itemChildren),
      ];
    }

    return newItem;
  });

}


/**
 * 根据 自定义传入权限名 判断当前用户是否有权限
 * @param userRoles 用户的权限
 * @param roles 自定义权限名
 */
export const hasPermissionRouteRoles = (userRoles: string[], roles?: string | string[]): boolean => {
  if (userRoles.includes('admin')) {
    return true;
  }

  if(typeof roles === 'undefined' || roles.length == 0) {
    return true;
  }

  if (typeof roles === 'string') {
    return userRoles.includes(roles);
  }

  if(roles instanceof Array && roles.length > 0) {
    return roles.some(role => userRoles.includes(role));
  }

  return false;
};

/**
 * 根据 route.roles 判断当前用户是否有权限
 * @param roles 用户的权限
 * @param route 当前路由
 */
export const hasPermission = (roles: string[], route: RoutesDataItem): boolean => {
  if (roles.includes('admin')) {
    return true;
  }

  if (route.roles) {
    return route.roles.some(role => roles.includes(role));
  }

  return true;
};

/**
 * 根据用户权限 获取 对应权限菜单
 * @param roles 用户的权限
 * @param routes 框架对应路由
 */
export const getPermissionMenuData = ( roles: string[], routes: RoutesDataItem[]): RoutesDataItem[] => {
  const menu: RoutesDataItem[] = [];
  for (let index = 0, len = routes.length; index < len; index += 1) {
    const element = {...routes[index]};
    if (hasPermission(roles, element)) {
      if (element.children) {
        element.children = getPermissionMenuData(roles, element.children);
      }
      menu.push(element);
    }
  }

  return menu;
};


/**
 * 判断tabNav，对应的route是否相等
 * @param route1 vue-route
 * @param route2 vue-route
 * @param type 判断规则
 * @returns
 */
 export const equalTabNavRoute = (route1: RouteLocationNormalizedLoaded,
                                  route2: RouteLocationNormalizedLoaded,
                                  type: TabNavType = 'path'): boolean=> {

  let is = false;
  if (type === 'querypath') { // path + query
    is = equalObject(route1.query, route2.query) && route1.path === route2.path
  } else { // path
    is = route1.path === route2.path

    if (!is) {
      const arr1 = route1.path.split('/')
      const arr2 = route2.path.split('/')

      is = arr1[1] == '~' && arr2[1] == '~' && arr1[2] == arr2[2] // 同源高亮
    }
  }

  return is;
}

