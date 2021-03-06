import { createRouter, createWebHistory } from 'vue-router'
import home from './home'
import teams from './teams'
import projects from './projects'
import executions from './executions'
import annotations from './annotations'
import environments from './environments'

const router = createRouter({
  history: createWebHistory('/ui/'),
  routes: [
    ...home,
    ...teams,
    ...projects,
    ...annotations,
    ...environments,
    ...executions,
  ],
})

export default router
