import { RoutesDataItem } from "@/utils/routes";
import BlankLayout from '@/layouts/BlankLayout.vue';

const IndexLayoutRoutes: Array<RoutesDataItem> = [
  {
    icon: 'script',
    title: 'index-layout.menu.script',
    path: '/script',
    redirect: '/script/list',
    component: BlankLayout,
    children: [
      {
        title: 'index-layout.menu.script.list',
        path: 'list',
        component: () => import('@/views/script/index/main.vue'),
        hidden: true,
      },
      {
        title: 'index-layout.menu.script.view',
        path: 'view/:id',
        component: () => import('@/views/script/view/index.vue'),
        hidden: true,
      },
    ],
  },

  {
    icon: 'execution',
    title: 'index-layout.menu.execution',
    path: '/exec',
    redirect: '/exec/history',
    component: BlankLayout,
    children: [
      {
        title: 'index-layout.menu.execution.history',
        path: 'history',
        component: () => import('@/views/exec/history/index.vue'),
        hidden: true,
      },
      {
        title: 'index-layout.menu.execution',
        path: 'exec',
        component: BlankLayout,
        hidden: true,
        children: [
          {
            title: 'index-layout.menu.execution.execCase',
            path: 'case',
            component: () => import('@/views/exec/exec/case.vue'),
            hidden: true,
          },
          {
            title: 'index-layout.menu.execution.execSuite',
            path: 'suite',
            component: () => import('@/views/exec/exec/suite.vue'),
            hidden: true,
          },
          {
            title: 'index-layout.menu.execution.execTask',
            path: 'task',
            component: () => import('@/views/exec/exec/task.vue'),
            hidden: true,
          },
        ]
      },
      {
        title: 'index-layout.menu.execution.result',
        path: 'result/:id',
        component: () => import('@/views/exec/result/index.vue'),
        hidden: true,
      },
      {
        title: 'index-layout.menu.execution.result',
        path: 'result/:id',
        component: () => import('@/views/exec/result/index.vue'),
        hidden: true,
      },
      {
        title: 'index-layout.menu.execution.result',
        path: 'result/:id',
        component: () => import('@/views/exec/result/index.vue'),
        hidden: true,
      },
    ],
  },

  {
    icon: 'config',
    title: 'index-layout.menu.config',
    path: '/config',
    component: () => import('@/views/config/index.vue'),
  },

  {
    icon: 'sync',
    title: 'index-layout.menu.sync',
    path: '/sync',
    component: () => import('@/views/sync/index.vue'),
  },

];

export default IndexLayoutRoutes;