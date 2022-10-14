import { reactive, toRefs } from 'vue'
import { useI18n } from 'vue-i18n'
import { useHead } from '@vueuse/head'
import { useRoute } from 'vue-router'
import AnnotationsService, { Annotations } from '@/api/annotationsService'

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
    annotations: {
      annotations: [] as Annotations[],
      byFilter: [] as Annotations[],
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

  state.meta.title = t('annotations.list')
  state.meta.description = t('annotations.list')

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
  AnnotationsService.list(page)
    .then((response: any) => {
      switch (response.status) {
        case 200:
          state.annotations.annotations = response.data
          total = state.annotations.annotations[0].total
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

  function deleteItem(index: number, type: string, id: number): void {
    if (confirm(t('confirm.sure'))) {
      state.alert.message = ''
      state.loading.delete.active = true
      AnnotationsService.delete(id)
        .then(() => {
          if (type === 'annotations') {
            state.annotations.annotations.splice(index, 1)
            if (state.annotations.annotations.length == 0) {
              state.alert.class = 'mute'
              state.alert.message = t('alert.http.noDataFound')
            }
          } else {
            state.annotations.byFilter.splice(index, 1)
            if (state.annotations.byFilter.length == 0) {
              state.alert.class = 'mute'
              state.alert.message = t('alert.http.noDataFound')
            }
          }
          state.loading.delete.active = false
        })
        .catch((error: any) => {
          state.alert.class = 'red'
          state.alert.message = t('alert.http.errorOccured')
          state.loading.delete.active = false
          throw error
        })
    }
  }

  return {
    deleteItem,
    ...toRefs(state),
  }
}
