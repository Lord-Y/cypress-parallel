import { createRouter, createWebHistory } from 'vue-router'
import home from './home'
import teams from './teams'
import projects from './projects'
import executions from './executions'
import annotations from './annotations'

const router = createRouter({
  // history: createWebHistory(process.env.BASE_URL),
  history: createWebHistory(),
  routes: [...home, ...teams, ...projects, ...annotations, ...executions],
})

export default router
