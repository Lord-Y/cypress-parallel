import { reactive, toRefs } from 'vue'
import { useI18n } from 'vue-i18n'
import { useHead } from '@vueuse/head'
import { useRoute } from 'vue-router'
import ProjectsService, { Project } from '@/api/projectsService'
import AnnotationsService, { Annotation } from '@/api/annotationsService'
import EnvironmentsService, { Environment } from '@/api/environmentsService'

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
      hook: {
        active: false,
      },
      duplicate: {
        active: false,
      },
    },
    alert: {
      class: '',
      message: '',
    },
    isOpen: false,
    projects: {
      projects: [] as Project[],
      byFilter: [] as Project[],
      duplicate: {
        data: {} as Project,
        annoations: {} as Annotation[],
        environments: {} as Environment[],
      },
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
      aLinks: 'hover:text-green-500 hover:font-extrabold',
    },
    status: {
      annoations: {
        count: 0,
      },
      environments: {
        count: 0,
      },
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

  state.meta.title = t('projects.list')
  state.meta.description = t('projects.list')

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
  ProjectsService.list(page)
    .then((response: any) => {
      switch (response.status) {
        case 200:
          state.projects.projects = response.data
          total = state.projects.projects[0].total
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
    })
    .catch((error: any) => {
      state.alert.class = 'red'
      state.alert.message = t('alert.http.errorOccured')
      throw error
    })
  state.loading.loading.active = false

  const { projects, loading, alert, pagination, search, classes } =
    toRefs(state)

  function deleteItem(index: number, type: string, id: number): void {
    if (confirm(t('confirm.sure'))) {
      state.alert.message = ''
      state.loading.delete.active = true
      ProjectsService.delete(id)
        .then(() => {
          if (type === 'projects') {
            state.projects.projects.splice(index, 1)
            if (state.projects.projects.length == 0) {
              state.alert.class = 'mute'
              state.alert.message = t('alert.http.noDataFound')
            }
          } else {
            state.projects.byFilter.splice(index, 1)
            if (state.projects.byFilter.length == 0) {
              state.alert.class = 'mute'
              state.alert.message = t('alert.http.noDataFound')
            }
          }
        })
        .catch((error: any) => {
          state.alert.class = 'red'
          state.alert.message = t('alert.http.errorOccured')
          throw error
        })
      state.loading.delete.active = false
    }
  }

  function launch(project: Project): void {
    state.alert.message = ''
    state.loading.hook.active = true
    ProjectsService.hook({
      project_name: project.project_name,
      cypress_docker_version: project.cypress_docker_version,
      browser: project.browser,
      config_file: project.config_file,
      max_pods: project.max_pods,
    })
      .then((response: any) => {
        if (response.status === 201) {
          state.alert.class = 'green'
          state.alert.message = t('alert.http.hook.created', {
            field: project.project_name,
          })
        } else {
          state.alert.class = 'red'
          state.alert.message = t('alert.http.errorOccured')
        }
        state.loading.hook.active = false
      })
      .catch((error: any) => {
        state.alert.class = 'red'
        state.alert.message = t('alert.http.errorOccured')
        state.loading.hook.active = false
        throw error
      })
  }

  function duplicate(project_id: number) {
    state.loading.duplicate.active = true
    let new_project_id: number
    getProjectByID(project_id).then((response) => {
      if (response) {
        duplicateProject().then((response) => {
          new_project_id = response
          if (response > 0) {
            getAnnotationsByProjectID(project_id).then((response) => {
              if (response && state.projects.duplicate.annoations.length > 0) {
                duplicateAnnotations(new_project_id)
              }
            })
            getEnvironmentsByProjectID(project_id).then((response) => {
              if (
                response &&
                state.projects.duplicate.environments.length > 0
              ) {
                duplicateEnvionments(new_project_id)
              }
            })
            state.loading.duplicate.active = false
          }
        })
      }
    })
  }

  async function getProjectByID(project_id: number): Promise<boolean> {
    let status: boolean
    return await ProjectsService.get(project_id)
      .then((response: any) => {
        switch (response.status) {
          case 200:
            state.projects.duplicate.data = response.data
            if (response.data.scheduling_enabled === 'false') {
              state.projects.duplicate.data.scheduling_enabled = false
            } else {
              state.projects.duplicate.data.scheduling_enabled = true
            }
            status = true
            break
          default:
            state.alert.class = 'red'
            state.alert.message = t('alert.http.errorOccured')
            state.loading.duplicate.active = false
            status = false
            break
        }
        return status
      })
      .catch((error: any) => {
        if (error.response.status === 404) {
          state.alert.class = 'mute'
          state.alert.message = t('alert.http.pageNotFound')
        } else {
          state.alert.class = 'red'
          state.alert.message = t('alert.http.errorOccured')
        }
        state.loading.duplicate.active = false
        return false
      })
  }

  async function duplicateProject(): Promise<number> {
    return await ProjectsService.create({
      teamId: Number(state.projects.duplicate.data.team_id),
      name: state.projects.duplicate.data.project_name + '_duplicated',
      repository: state.projects.duplicate.data.repository,
      branch: state.projects.duplicate.data.branch,
      username: state.projects.duplicate.data.username,
      password: state.projects.duplicate.data.password,
      specs: state.projects.duplicate.data.specs,
      scheduling: state.projects.duplicate.data.scheduling,
      schedulingEnabled: state.projects.duplicate.data.scheduling_enabled,
      maxPods: Number(state.projects.duplicate.data.max_pods),
      cypress_docker_version:
        state.projects.duplicate.data.cypress_docker_version,
      timeout: Number(state.projects.duplicate.data.timeout),
      browser: state.projects.duplicate.data.browser,
      config_file: state.projects.duplicate.data.config_file,
    })
      .then((response: any) => {
        if (response.status === 201) {
          state.alert.class = 'green'
          state.alert.message = t('alert.http.project.duplicate', {
            field: state.projects.duplicate.data.project_name + '_duplicated',
          })
          return Number(response.data['projectId'])
        } else {
          state.alert.class = 'red'
          state.alert.message = t('alert.http.errorOccured')
          state.loading.duplicate.active = false
          return 0
        }
      })
      .catch(() => {
        state.alert.class = 'red'
        state.alert.message = t('alert.http.errorOccured')
        state.loading.duplicate.active = false
        return 0
      })
  }

  async function getAnnotationsByProjectID(
    project_id: number,
  ): Promise<boolean> {
    let status: boolean
    return await AnnotationsService.projectID(project_id)
      .then((response: any) => {
        switch (response.status) {
          case 200:
            state.projects.duplicate.annoations = response.data
            status = true
            break
          default:
            state.alert.class = 'red'
            state.alert.message = t('alert.http.errorOccured')
            state.loading.duplicate.active = false
            status = false
            break
        }
        return status
      })
      .catch((error: any) => {
        if (error.response.status === 404) {
          return true
        } else {
          state.alert.class = 'red'
          state.alert.message = t('alert.http.errorOccured')
        }
        state.loading.duplicate.active = false
        return false
      })
  }

  function duplicateAnnotations(project_id: number): boolean {
    for (let i = 0; i < state.projects.duplicate.annoations.length; i++) {
      AnnotationsService.create({
        projectId: Number(project_id),
        key: state.projects.duplicate.annoations[i].key,
        value: state.projects.duplicate.annoations[i].value,
      })
        .then((response: any) => {
          if (response.status === 201) {
            state.status.annoations.count++
          } else {
            state.alert.class = 'red'
            state.alert.message = t('alert.http.errorOccured')
            state.loading.duplicate.active = false
          }
        })
        .catch((error: any) => {
          state.alert.class = 'red'
          state.alert.message = t('alert.http.errorOccured')
          state.loading.duplicate.active = false
          throw error
        })
    }
    if (
      state.status.annoations.count ===
      state.projects.duplicate.annoations.length
    ) {
      return true
    } else {
      return false
    }
  }

  async function getEnvironmentsByProjectID(
    project_id: number,
  ): Promise<boolean> {
    let status: boolean
    return await EnvironmentsService.projectID(project_id)
      .then((response: any) => {
        switch (response.status) {
          case 200:
            state.projects.duplicate.environments = response.data
            status = true
            break
          default:
            state.alert.class = 'red'
            state.alert.message = t('alert.http.errorOccured')
            state.loading.duplicate.active = false
            status = false
            break
        }
        return status
      })
      .catch((error: any) => {
        if (error.response.status === 404) {
          return true
        } else {
          state.alert.class = 'red'
          state.alert.message = t('alert.http.errorOccured')
        }
        state.loading.duplicate.active = false
        return false
      })
  }

  function duplicateEnvionments(project_id: number): boolean {
    for (let i = 0; i < state.projects.duplicate.environments.length; i++) {
      EnvironmentsService.create({
        projectId: Number(project_id),
        key: state.projects.duplicate.environments[i].key,
        value: state.projects.duplicate.environments[i].value,
      })
        .then((response: any) => {
          if (response.status === 201) {
            state.status.environments.count++
          } else {
            state.alert.class = 'red'
            state.alert.message = t('alert.http.errorOccured')
            state.loading.duplicate.active = false
          }
        })
        .catch((error: any) => {
          state.alert.class = 'red'
          state.alert.message = t('alert.http.errorOccured')
          state.loading.duplicate.active = false
          throw error
        })
    }
    if (
      state.status.environments.count ===
      state.projects.duplicate.environments.length
    ) {
      return true
    } else {
      return false
    }
  }

  return {
    deleteItem,
    launch,
    duplicate,
    ...toRefs(state),
  }
}
