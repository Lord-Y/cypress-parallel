import { RouteRecordRaw } from 'vue-router'

const routes: Array<RouteRecordRaw> = [
  {
    path: '/projects/create',
    component: () => import('@/views/projects/Create.vue'),
    meta: {
      root: 'projects',
      activeLink: '/projects/create',
    },
  },
  {
    path: '/projects/edit/:id',
    component: () => import('@/views/projects/Edit.vue'),
    meta: {
      root: 'projects',
      activeLink: '/projects/edit',
    },
  },
  {
    path: '/projects/list',
    component: () => import('@/views/projects/List.vue'),
    meta: {
      root: 'projects',
      activeLink: '/projects/list',
    },
    children: [
      {
        path: ':page',
        component: () => import('@/views/projects/List.vue'),
        meta: {
          root: 'projects',
          activeLink: '/projects/list',
        },
      },
    ],
  },
]

export default routes
