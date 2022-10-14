import { reactive, toRefs } from 'vue'
import { useI18n } from 'vue-i18n'
import { useHead } from '@vueuse/head'
import TeamsService, { Teams } from '@/api/teamsService'
import ProjectsService from '@/api/projectsService'

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
      cypress_docker_version: '10.10.0-0.3.0',
      timeout: 10,
      browser: 'chrome',
      config_file: 'cypress.config.js',
    },
  })

  const { t } = useI18n({
    useScope: 'global',
  })

  state.meta.title = t('projects.create')
  state.meta.description = t('projects.create')

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
        case 204:
          state.alert.class = 'red'
          state.alert.message = t('alert.http.create.team')
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

  function submit() {
    state.loading.loading.active = true
    ProjectsService.create({
      teamId: Number(state.form.team_id),
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
        if (response.status === 201) {
          state.alert.class = 'green'
          state.alert.message = t('alert.http.project.created', {
            field: state.form.project_name,
          })
        } else {
          state.alert.class = 'red'
          state.alert.message = t('alert.http.errorOccured')
        }
        state.loading.loading.active = false
      })
      .catch((error: any) => {
        switch (error.response.status) {
          case 409:
            state.alert.class = 'red'
            state.alert.message = t('alert.http.conflict.project', {
              field: state.form.project_name,
            })
            break
          default:
            state.alert.class = 'red'
            state.alert.message = t('alert.http.errorOccured')
            break
        }
        state.loading.loading.active = false
        throw error
      })
  }

  return {
    submit,
    ...toRefs(state),
  }
}
