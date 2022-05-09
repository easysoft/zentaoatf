/**
 * 路由入口
 * @author LiQingSong
 */
import NProgress from 'nprogress'; // progress bar
import 'nprogress/nprogress.css'; // progress bar style
NProgress.configure({ showSpinner: false, easing: 'ease', speed: 1000 }); // NProgress Configuration

import { createRouter, createWebHashHistory } from 'vue-router';
import { RoutesDataItem } from "@/utils/routes";

import IndexLayoutRoutes from '@/layouts/IndexLayout/routes';
import IndexLayout from '@/layouts/IndexLayout/index.vue';
import MainLayout from '@/layouts/MainLayout/Main.vue';
import BlankLayout from "@/layouts/BlankLayout.vue";

const routes: RoutesDataItem[] = [
  {
    title: 'empty',
    path: '/',
    component: BlankLayout,
    children: [
      {
        title: 'main',
        path: '/main',
        component: MainLayout
      },
      {
        title: 'empty',
        path: '/',
        redirect: '/script/index',
        component: IndexLayout,
        children: IndexLayoutRoutes
      },
      {
        title: 'empty',
        path: '/refresh',
        component: () => import('@/views/refresh/index.vue'),
      },
    ]
  },
  {
    title: 'app.global.menu.notfound',
    path: '/:pathMatch(.*)*',
    component: () => import('@/views/404/index.vue'),
  }
]

const router = createRouter({
  scrollBehavior(/* to, from, savedPosition */) {
    return { top: 0 }
  },
  history: createWebHashHistory(process.env.BASE_URL),
  routes: routes,
});

router.beforeEach((/* to, from */) => {
  // start progress bar
  NProgress.start();
});

router.afterEach(() => {
  // finish progress bar
  NProgress.done();
});

export default router;
