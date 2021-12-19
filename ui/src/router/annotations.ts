import { RouteRecordRaw } from 'vue-router'

const routes: Array<RouteRecordRaw> = [
  {
    path: '/annotations/create',
    component: () => import('@/views/annotations/Create.vue'),
    meta: {
      root: 'annotations',
      activeLink: '/annotations/create',
    },
  },
  {
    path: '/annotations/edit/:id',
    component: () => import('@/views/annotations/Edit.vue'),
    meta: {
      root: 'annotations',
      activeLink: '/annotations/edit',
    },
  },
  {
    path: '/annotations/list',
    component: () => import('@/views/annotations/List.vue'),
    meta: {
      root: 'annotations',
      activeLink: '/annotations/list',
    },
    children: [
      {
        path: ':page',
        component: () => import('@/views/annotations/List.vue'),
        meta: {
          root: 'annotations',
          activeLink: '/annotations/list',
        },
      },
    ],
  },
]

export default routes
