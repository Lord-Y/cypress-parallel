import { reactive, toRefs } from 'vue'
import { useI18n } from 'vue-i18n'
import { useHead } from '@vueuse/head'
import { useRoute } from 'vue-router'
import ExecutionsService, { Execution } from '@/api/executionsService'

export default function () {
  const state = reactive({
    meta: {
      title: '',
      description: '',
    },
    loading: {
      loading: {
        active: true,
      },
      delete: {
        active: false,
      },
    },
    alert: {
      class: '',
      message: '',
    },
    isOpen: false,
    executions: {
      executions: [] as Execution[],
      byFilter: [] as Execution[],
    },
    pagination: {
      enabled: false,
      data: {
        url: '',
        actualPage: 1,
        total: 0,
      },
    },
    search: {
      byFilter: '',
      bar: {
        enabled: false,
      },
      table: {
        enabled: false,
      },
    },
    classes: {
      aLinks: 'hover:text-emerald-500 hover:font-extrabold',
    },
  })

  const route = useRoute()
  const { t } = useI18n({
    useScope: 'global',
  })

  let page: number, total: number

  if (!route.params.page) {
    page = 1
  } else {
    page = Number(route.params.page)
  }

  state.meta.title = t('executions.list')
  state.meta.description = t('executions.list')

  useHead({
    title: state.meta.title,
    meta: [
      {
        name: 'description',
        content: state.meta.description,
      },
      {
        property: 'og:title',
        content: state.meta.title,
      },
      {
        property: 'og:description',
        content: state.meta.description,
      },
    ],
  })

  state.loading.loading.active = true
  ExecutionsService.list(page)
    .then((response: any) => {
      switch (response.status) {
        case 200:
          state.executions.executions = response.data
          total = state.executions.executions[0].total
          if (total > 25) {
            state.search.bar.enabled = true
            state.pagination.data.url = route.path.replace('/' + page, '')
            state.pagination.data.actualPage = page
            state.pagination.data.total = total
            state.pagination.enabled = true
          }
          break
        case 204:
          state.alert.class = 'mute'
          state.alert.message = t('alert.http.noDataFound')
          break
        default:
          state.alert.class = 'red'
          state.alert.message = t('alert.http.errorOccured')
          break
      }
      state.loading.loading.active = false
    })
    .catch((error: any) => {
      state.alert.class = 'red'
      state.alert.message = t('alert.http.errorOccured')
      state.loading.loading.active = false
      throw error
    })

  return {
    ...toRefs(state),
  }
}
