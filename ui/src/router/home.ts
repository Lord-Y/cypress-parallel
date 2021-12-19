import { RouteRecordRaw } from 'vue-router'

const routes: Array<RouteRecordRaw> = [
  {
    path: '/',
    redirect: '/executions/list',
  },
  {
    path: '/:pathMatch(.*)*',
    component: () => import('@/views/pageStatus/404.vue'),
  },
  {
    path: '/404',
    component: () => import('@/views/pageStatus/404.vue'),
  },
  {
    path: '/500',
    component: () => import('@/views/pageStatus/500.vue'),
  },
]

export default routes
