import { RouteRecordRaw } from 'vue-router'
import Home from '@views/Home.vue'

const routes: Array<RouteRecordRaw> = [
  {
    path: '/',
    name: 'Home',
    component: Home,
    meta: {
      root: 'home',
      activeLink: '/',
    },
  },
]

export default routes
