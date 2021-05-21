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
]

export default routes
