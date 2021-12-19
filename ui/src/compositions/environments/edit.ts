import { reactive, toRefs } from 'vue'
import { useI18n } from 'vue-i18n'
import { useHead } from '@vueuse/head'
import { useRoute } from 'vue-router'
import ProjectsService, { ProjectOnly } from '@/api/projectsService'
import EnvironmentsService, { Environment } from '@/api/environmentsService'

export default function () {
  const state = reactive({
    meta: {
      title: '',
      description: '',
    },
    loading: {
      loading: {
        active: false,
      },
    },
    alert: {
      class: '',
      message: '',
    },
    projects: [] as ProjectOnly[],
    environment: {} as Environment,
    form: {
      project_id: '',
      key: '',
      value: '',
    },
  })

  const route = useRoute()
  const { t } = useI18n({
    useScope: 'global',
  })

  state.meta.title = t('environments.edit')
  state.meta.description = t('environments.edit')

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
  ProjectsService.all()
    .then((response: any) => {
      switch (response.status) {
        case 200:
          state.projects = response.data
          break
        default:
          state.alert.class = 'red'
          state.alert.message = t('alert.http.errorOccured')
          break
      }
    })
    .catch((error: any) => {
      state.alert.class = 'red'
      state.alert.message = t('alert.http.errorOccured')
      throw error
    })

  EnvironmentsService.get(Number(route.params.id))
    .then((response: any) => {
      switch (response.status) {
        case 200:
          state.environment = response.data
          break
        default:
          state.alert.class = 'red'
          state.alert.message = t('alert.http.errorOccured')
          break
      }
    })
    .catch((error: any) => {
      if (error.response.status === 404) {
        state.alert.class = 'mute'
        state.alert.message = t('alert.http.pageNotFound')
      } else {
        state.alert.class = 'red'
        state.alert.message = t('alert.http.errorOccured')
      }
      throw error
    })

  function submit() {
    state.loading.loading.active = true
    EnvironmentsService.update({
      projectId: Number(state.form.project_id),
      environmentId: Number(state.environment.environment_id),
      key: state.form.key,
      value: state.form.value,
    })
      .then((response: any) => {
        if (response.status === 200) {
          state.alert.class = 'green'
          state.alert.message = t('alert.http.environment.updated', {
            field: state.form.key,
          })
        } else {
          state.alert.class = 'red'
          state.alert.message = t('alert.http.errorOccured')
        }
        state.loading.loading.active = false
      })
      .catch((error: any) => {
        state.alert.class = 'red'
        state.alert.message = t('alert.http.errorOccured')
        state.loading.loading.active = false
        throw error
      })
  }

  return {
    submit,
    ...toRefs(state),
  }
}
