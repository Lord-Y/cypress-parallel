import { RouteRecordRaw } from 'vue-router'

const routes: Array<RouteRecordRaw> = [
  {
    path: '/executions/list',
    component: () => import('@views/executions/List.vue'),
    meta: {
      root: 'executions',
      activeLink: '/executions/list',
    },
    children: [
      {
        path: ':page',
        component: () => import('@views/executions/List.vue'),
        meta: {
          root: 'executions',
          activeLink: '/executions/list',
        },
      },
    ],
  },
  {
    path: '/executions/:id',
    component: () => import('@views/executions/Spec.vue'),
    meta: {
      root: 'executions',
      activeLink: '/executions/list',
    },
  },
  {
    path: '/executions/uniqid/:id',
    component: () => import('@views/executions/UniqID.vue'),
    meta: {
      root: 'executions',
      activeLink: '/executions/list',
    },
  },
]

export default routes
