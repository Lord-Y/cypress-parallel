import { reactive, toRefs } from 'vue'
import { useI18n } from 'vue-i18n'
import { useHead } from '@vueuse/head'
import ProjectsService, { ProjectOnly } from '@/api/projectsService'
import EnvironmentsService from '@/api/environmentsService'

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
    form: {
      project_id: '',
      key: '',
      value: '',
    },
  })

  const { t } = useI18n({
    useScope: 'global',
  })

  state.meta.title = t('environments.create')
  state.meta.description = t('environments.create')

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

  ProjectsService.all()
    .then((response: any) => {
      switch (response.status) {
        case 200:
          state.projects = response.data
          break
        case 204:
          state.alert.class = 'red'
          state.alert.message = t('alert.http.create.project')
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

  function submit() {
    state.loading.loading.active = true
    EnvironmentsService.create({
      projectId: Number(state.form.project_id),
      key: state.form.key,
      value: state.form.value,
    })
      .then((response: any) => {
        if (response.status === 201) {
          state.alert.class = 'green'
          state.alert.message = t('alert.http.environment.created', {
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
