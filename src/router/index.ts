import { createRouter, createWebHistory } from 'vue-router'
import home from './home'
import teams from './teams'
import projects from './projects'
import executions from './executions'

const router = createRouter({
  // history: createWebHistory(process.env.BASE_URL),
  history: createWebHistory(),
  routes: [...home, ...teams, ...projects, ...executions],
})

export default router
