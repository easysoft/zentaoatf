import { RoutesDataItem } from "@/utils/routes";
import BlankLayout from '@/layouts/BlankLayout.vue';

const IndexLayoutRoutes: Array<RoutesDataItem> = [
  {
    icon: 'script',
    title: 'test_script',
    path: '/script',
    redirect: '/script/index',
    component: BlankLayout,
    children: [
      {
        title: 'test_script',
        path: 'index',
        component: () => import('@/views/script/index.vue'),
        hidden: true,
      },
      {
        title: 'test_script',
        path: 'index/:workspace/:seq/:scope',
        component: () => import('@/views/script/index.vue'),
        hidden: true,
      },
    ],
  },

  {
    icon: 'result',
    title: 'test_result',
    path: '/result',
    redirect: '/result/list',
    component: BlankLayout,
    children: [
      {
        title: 'test_result',
        path: 'list',
        component: () => import('@/views/result/index.vue'),
        hidden: true,
      },
      {
        title: 'execution.result.func',
        path: 'func/:workspaceId/:seq',
        component: () => import('@/views/result/result-func.vue'),
        hidden: true,
      },
      {
        title: 'execution.result.unit',
        path: 'unit/:workspaceId/:seq',
        component: () => import('@/views/result/result-unit.vue'),
        hidden: true,
      },

     /* {
        title: 'execution',
        path: 'run',
        component: BlankLayout,
        hidden: true,
        children: [
          {
            title: 'execution.execCase',
            path: 'case/:seq/:scope',
            component: () => import('@/views/exec/exec/case.vue'),
            hidden: true,
          },
          {
            title: 'execution.execModule',
            path: 'module/:productId/:moduleId/:seq/:scope',
            component: () => import('@/views/exec/exec/module.vue'),
            hidden: true,
          },
          {
            title: 'execution.execSuite',
            path: 'suite/:productId/:suiteId/:seq/:scope',
            component: () => import('@/views/exec/exec/suite.vue'),
            hidden: true,
          },
          {
            title: 'execution.execTask',
            path: 'task/:productId/:taskId/:seq/:scope',
            component: () => import('@/views/exec/exec/task.vue'),
            hidden: true,
          },
          {
            title: 'execution.execUnit',
            path: 'unit',
            component: () => import('@/views/exec/exec/unit.vue'),
            hidden: true,
          },
        ]
      },*/
    ],
  },

  {
    icon: 'empty',
    title: 'workspace',
    path: '/workspace',
    redirect: '/workspace/list',
    component: BlankLayout,
    children: [
      {
        title: 'workspace',
        path: 'list',
        component: () => import('@/views/workspace/index.vue'),
        hidden: true,
      },
      {
        title: 'create_workspace',
        path: 'edit/:id',
        component: () => import('@/views/workspace/edit.vue'),
        hidden: true,
      },
      {
        title: 'edit_workspace',
        path: 'edit/:id',
        component: () => import('@/views/workspace/edit.vue'),
        hidden: true,
      },
    ],
  },

  {
    icon: 'empty',
    title: 'site',
    path: '/site',
    redirect: '/site/list',
    component: BlankLayout,
    hidden: true,
    children: [
      {
        title: 'zentao_site',
        path: 'list',
        component: () => import('@/views/site/index.vue'),
        hidden: true,
      },
      {
        title: 'create_site',
        path: 'create/:id',
        component: () => import('@/views/site/edit.vue'),
        hidden: true,
      },
      {
        title: 'edit_site',
        path: 'edit/:id',
        component: () => import('@/views/site/edit.vue'),
        hidden: true,
      },
    ],
  },

  {
    icon: 'empty',
    title: 'interpreter',
    path: '/interpreter',
    redirect: '/interpreter/list',
    component: BlankLayout,
    hidden: true,
    children: [
      {
        title: 'interpreter',
        path: 'list',
        component: () => import('@/views/interpreter/index.vue'),
        hidden: true,
      },
    ],
  },
];

export default IndexLayoutRoutes;