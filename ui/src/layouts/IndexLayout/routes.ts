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
        title: 'index-layout.menu.execution.result.func',
        path: 'history/func/:seq',
        component: () => import('@/views/exec/history/result-func.vue'),
        hidden: true,
      },
      {
        title: 'index-layout.menu.execution.result.unit',
        path: 'history/unit',
        component: () => import('@/views/exec/history/result-unit.vue'),
        hidden: true,
      },

      {
        title: 'index-layout.menu.execution',
        path: 'run',
        component: BlankLayout,
        hidden: true,
        children: [
          {
            title: 'index-layout.menu.execution.execCase',
            path: 'case/:seq/:scope',
            component: () => import('@/views/exec/exec/case.vue'),
            hidden: true,
          },
          {
            title: 'index-layout.menu.execution.execModule',
            path: 'module/:productId/:moduleId/:seq/:scope',
            component: () => import('@/views/exec/exec/module.vue'),
            hidden: true,
          },
          {
            title: 'index-layout.menu.execution.execSuite',
            path: 'suite/:productId/:suiteId/:seq/:scope',
            component: () => import('@/views/exec/exec/suite.vue'),
            hidden: true,
          },
          {
            title: 'index-layout.menu.execution.execTask',
            path: 'task/:productId/:taskId/:seq/:scope',
            component: () => import('@/views/exec/exec/task.vue'),
            hidden: true,
          },
          {
            title: 'index-layout.menu.execution.execUnit',
            path: 'unit',
            component: () => import('@/views/exec/exec/unit.vue'),
            hidden: true,
          },
        ]
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