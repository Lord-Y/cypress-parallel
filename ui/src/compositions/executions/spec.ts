import { reactive, toRefs } from 'vue'
import { useI18n } from 'vue-i18n'
import { useHead } from '@vueuse/head'
import { useRoute } from 'vue-router'
import ExecutionsService, { Execution } from '@/api/executionsService'
import Statuses from '@/tools/status'
import moment from 'moment'

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
    },
    alert: {
      class: '',
      message: '',
    },
    isOpen: false,
    execution: {} as Execution,
    classes: {
      aLinks: 'hover:text-emerald-500 hover:font-extrabold',
    },
  })

  const route = useRoute()
  const { t } = useI18n({
    useScope: 'global',
  })

  const id = String(route.params.id)

  state.meta.title = t('executions.by.spec')
  state.meta.description = t('executions.by.spec')

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
  ExecutionsService.get(Number(id))
    .then((response: any) => {
      switch (response.status) {
        case 200:
          state.execution = response.data
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

  function convertDuration(s: number) {
    if (s > 60000) {
      return moment.duration(s).minutes() + 'min'
    } else {
      return moment.duration(s).seconds() + 's'
    }
  }

  function getSpecStatus(execution: Execution): string {
    return Statuses.tests(execution)
  }

  function getSystemStatus(s: string): string {
    let classes: string
    switch (s) {
      case 'DONE':
        classes = 'text-emerald-500 font-semibold'
        break
      case 'NOT_STARTED':
      case 'QUEUED':
      case 'SCHEDULED':
      case 'CANCELLED':
        classes = 'text-gray-500 font-semibold'
        break
      default:
        classes = 'text-red-500 font-semibold'
        break
    }
    return classes
  }

  function getGlobalStatus(execution: Execution, mode: string): string {
    let classes: string
    const status = Statuses.global(execution)
    switch (status) {
      case 'PASSED':
        classes = 'text-emerald-500 font-semibold'
        break
      case 'NOT_STARTED':
      case 'QUEUED':
      case 'SCHEDULED':
      case 'CANCELLED':
        classes = 'text-gray-500 font-semibold'
        break
      default:
        classes = 'text-red-500 font-semibold'
        break
    }
    if (mode === 'classes') {
      return classes
    } else {
      return status
    }
  }

  return {
    convertDuration,
    getSpecStatus,
    getSystemStatus,
    getGlobalStatus,
    ...toRefs(state),
  }
}
