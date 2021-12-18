import { RouteRecordRaw } from 'vue-router'

const routes: Array<RouteRecordRaw> = [
  {
    path: '/environments/create',
    component: () => import('@views/environments/Create.vue'),
    meta: {
      root: 'environments',
      activeLink: '/environments/create',
    },
  },
  {
    path: '/environments/edit/:id',
    component: () => import('@views/environments/Edit.vue'),
    meta: {
      root: 'environments',
      activeLink: '/environments/edit',
    },
  },
  {
    path: '/environments/list',
    component: () => import('@views/environments/List.vue'),
    meta: {
      root: 'environments',
      activeLink: '/environments/list',
    },
    children: [
      {
        path: ':page',
        component: () => import('@views/environments/List.vue'),
        meta: {
          root: 'environments',
          activeLink: '/environments/list',
        },
      },
    ],
  },
]

export default routes
