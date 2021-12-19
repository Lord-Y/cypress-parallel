import { reactive, toRefs } from 'vue'
import { useI18n } from 'vue-i18n'
import { useHead } from '@vueuse/head'
import { useRoute } from 'vue-router'
import TeamsService, { Teams } from '@/api/teamsService'
import ProjectsService, { Project } from '@/api/projectsService'

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
    teams: [] as Teams[],
    project: {} as Project,
    form: {
      team_id: '',
      project_name: '',
      repository: '',
      branch: '',
      username: '',
      password: '',
      specs: '',
      scheduling: '',
      schedulingEnabled: false,
      maxPods: 10,
      cypress_docker_version: '7.2.0-0.0.3',
      timeout: 10,
      browser: 'chrome',
      config_file: 'cypress.json',
    },
  })

  const route = useRoute()
  const { t } = useI18n({
    useScope: 'global',
  })

  state.meta.title = t('projects.edit')
  state.meta.description = t('projects.edit')

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
  TeamsService.all()
    .then((response: any) => {
      switch (response.status) {
        case 200:
          state.teams = response.data
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

  state.loading.loading.active = true
  ProjectsService.get(Number(route.params.id))
    .then((response: any) => {
      switch (response.status) {
        case 200:
          state.project = response.data
          if (response.data.scheduling_enabled === 'false') {
            state.project.scheduling_enabled = false
          } else {
            state.project.scheduling_enabled = true
          }
          break
        default:
          state.alert.class = 'red'
          state.alert.message = t('alert.http.errorOccured')
          break
      }
      state.loading.loading.active = false
    })
    .catch((error: any) => {
      if (error.response.status === 404) {
        state.alert.class = 'mute'
        state.alert.message = t('alert.http.pageNotFound')
      } else {
        state.alert.class = 'red'
        state.alert.message = t('alert.http.errorOccured')
      }
      state.loading.loading.active = false
      throw error
    })

  function submit() {
    state.loading.loading.active = true
    ProjectsService.update({
      teamId: Number(state.form.team_id),
      projectId: Number(state.project.project_id),
      name: state.form.project_name,
      repository: state.form.repository,
      branch: state.form.branch,
      username: state.form.username,
      password: state.form.password,
      specs: state.form.specs,
      scheduling: state.form.scheduling,
      schedulingEnabled: state.form.schedulingEnabled,
      maxPods: state.form.maxPods,
      cypress_docker_version: state.form.cypress_docker_version,
      timeout: state.form.timeout,
      browser: state.form.browser,
      config_file: state.form.config_file,
    })
      .then((response: any) => {
        if (response.status === 200) {
          state.alert.class = 'green'
          state.alert.message = t('alert.http.project.updated', {
            field: state.form.project_name,
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
