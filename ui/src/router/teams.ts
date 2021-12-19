import { RouteRecordRaw } from 'vue-router'
import Create from '@/views/teams/Create.vue'
import Edit from '@/views/teams/Edit.vue'
import List from '@/views/teams/List.vue'

const routes: Array<RouteRecordRaw> = [
  {
    path: '/teams/create',
    component: Create,
    meta: {
      root: 'teams',
      activeLink: '/teams/create',
    },
    props: {
      url: {
        api: {
          default: '/api/v1/cypress-parallel-api/teams',
        },
      },
    },
  },
  {
    path: '/teams/edit/:id',
    component: Edit,
    meta: {
      root: 'teams',
      activeLink: '/teams/edit',
    },
    props: {
      url: {
        api: {
          default: '/api/v1/cypress-parallel-api/teams',
        },
      },
    },
  },
  {
    path: '/teams/list',
    component: List,
    meta: {
      root: 'teams',
      activeLink: '/teams/list',
    },
    props: {
      url: {
        api: {
          default: '/api/v1/cypress-parallel-api/teams',
        },
      },
    },
    children: [
      {
        path: ':page',
        component: List,
        meta: {
          root: 'teams',
          activeLink: '/teams/list',
        },
        props: {
          url: {
            api: {
              default: '/api/v1/cypress-parallel-api/teams',
            },
          },
        },
      },
    ],
  },
]

export default routes
