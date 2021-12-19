<template>
  <div class="lg:flex min-h-screen w-full bg-gray-50">
    <Menu />
    <div class="w-full p-4">
      <div class="block">
        <Title :title="$t('projects.list')" />
        <SpinnerCommon v-if="loading.loading.active" />
        <SpinnerCommon v-if="loading.delete.active" />
        <SpinnerCommon v-if="loading.hook.active" />
        <SpinnerCommon v-if="loading.duplicate.active" />
        <AlertMessage :message="alert.message" :classes="alert.class" />
        <SearchProjectsByFilter
          v-if="search.bar.enabled"
          v-model:loading="loading.loading.active"
          v-model:alertClass="alert.class"
          v-model:alertMessage="alert.message"
          v-model:searchTable="search.table.enabled"
          v-model:byFilter="projects.byFilter"
        />
        <div
          class="mx-auto px-3 mt-20 w-full overflow-auto"
          v-if="
            !loading.loading.active &&
            !search.table.enabled &&
            projects.projects.length > 0
          "
        >
          <table
            class="table-auto w-full text-left border-collapse divide-y border-t-2"
          >
            <thead>
              <tr>
                <th class="py-3">{{ $t('teams.name') }}</th>
                <th class="py-3">{{ $t('projects.projects') }}</th>
                <th class="py-3">{{ $t('projects.branch') }}</th>
                <th class="py-3">{{ $t('projects.hooks.launch') }}</th>
                <th class="py-3">{{ $t('projects.duplicate') }}</th>
                <th class="py-3">{{ $t('edit.edit') }}</th>
                <th class="py-3">{{ $t('delete.delete') }}</th>
              </tr>
            </thead>
            <tbody class="divide-y">
              <tr
                class="hover:bg-gray-100 hover:font-semibold"
                v-for="(project, index) in projects.projects"
                :id="project.project_id.toString()"
                :key="index"
              >
                <td class="px-2 py-3">
                  <router-link
                    :class="['cursor-pointer', classes.aLinks]"
                    :title="$t('edit.edit')"
                    :to="'/teams/edit/' + project.team_id"
                    >{{ project.project_name }}</router-link
                  >
                </td>
                <td class="px-2 py-3">{{ project.project_name }}</td>
                <td class="px-2 py-3">{{ project.branch }}</td>
                <td class="px-2 py-3">
                  <router-link
                    :class="['cursor-pointer', classes.aLinks]"
                    :title="$t('projects.hooks.launch')"
                    to
                    @click="launch(project)"
                  >
                    <svg
                      class="w-10 h-10"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                    >
                      <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z"
                      />
                      <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                      />
                    </svg>
                  </router-link>
                </td>
                <td class="px-2 py-3">
                  <router-link
                    :class="['cursor-pointer', classes.aLinks]"
                    :title="$t('projects.duplicate')"
                    to
                    @click="duplicate(project.project_id)"
                  >
                    <svg
                      class="w-10 h-10"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                    >
                      <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M8 7v8a2 2 0 002 2h6M8 7V5a2 2 0 012-2h4.586a1 1 0 01.707.293l4.414 4.414a1 1 0 01.293.707V15a2 2 0 01-2 2h-2M8 7H6a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2v-2"
                      />
                    </svg>
                  </router-link>
                </td>
                <td class="px-2 py-3">
                  <router-link
                    :class="['cursor-pointer', classes.aLinks]"
                    :title="$t('edit.edit')"
                    :to="'/projects/edit/' + project.project_id"
                  >
                    <svg
                      class="w-10 h-10 hover:text-green-500"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                    >
                      <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"
                      />
                    </svg>
                  </router-link>
                </td>
                <td class="px-2 py-3">
                  <router-link
                    class="cursor-pointer"
                    :title="$t('delete.delete')"
                    to
                    @click="deleteItem(index, 'projects', project.project_id)"
                  >
                    <svg
                      class="w-10 h-10 hover:text-green-500"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                    >
                      <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                      />
                    </svg>
                  </router-link>
                </td>
              </tr>
            </tbody>
          </table>
          <Pagination v-if="pagination.enabled" :pagination="pagination.data" />
        </div>
        <div
          class="mx-auto px-3 mt-20 w-full overflow-auto"
          v-if="
            !loading.loading.active &&
            search.table.enabled &&
            projects.byFilter.length > 0
          "
        >
          <table
            class="table-auto w-full text-left border-collapse divide-y border-t-2"
          >
            <thead>
              <tr>
                <th class="py-3">{{ $t('teams.name') }}</th>
                <th class="py-3">{{ $t('projects.projects') }}</th>
                <th class="py-3">{{ $t('projects.branch') }}</th>
                <th class="py-3">{{ $t('projects.hooks.launch') }}</th>
                <th class="py-3">{{ $t('projects.duplicate') }}</th>
                <th class="py-3">{{ $t('edit.edit') }}</th>
                <th class="py-3">{{ $t('delete.delete') }}</th>
              </tr>
            </thead>
            <tbody class="divide-y">
              <tr
                class="hover:bg-gray-100 hover:font-semibold"
                v-for="(project, index) in projects.byFilter"
                :id="project.project_id.toString()"
                :key="index"
              >
                <td class="px-2 py-3">
                  <router-link
                    :class="['cursor-pointer', classes.aLinks]"
                    :title="$t('edit.edit')"
                    :to="'/teams/edit/' + project.team_id"
                    >{{ project.project_name }}</router-link
                  >
                </td>
                <td class="px-2 py-3">{{ project.project_name }}</td>
                <td class="px-2 py-3">{{ project.branch }}</td>
                <td class="px-2 py-3">
                  <router-link
                    :class="['cursor-pointer', classes.aLinks]"
                    :title="$t('projects.hooks.launch')"
                    to
                    @click="launch(project)"
                  >
                    <svg
                      class="w-10 h-10"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                    >
                      <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z"
                      />
                      <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                      />
                    </svg>
                  </router-link>
                </td>
                <td class="px-2 py-3">
                  <router-link
                    :class="['cursor-pointer', classes.aLinks]"
                    :title="$t('projects.duplicate')"
                    to
                    @click="duplicate(project.project_id)"
                  >
                    <svg
                      class="w-10 h-10"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                    >
                      <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M8 7v8a2 2 0 002 2h6M8 7V5a2 2 0 012-2h4.586a1 1 0 01.707.293l4.414 4.414a1 1 0 01.293.707V15a2 2 0 01-2 2h-2M8 7H6a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2v-2"
                      />
                    </svg>
                  </router-link>
                </td>
                <td class="px-2 py-3">
                  <router-link
                    :class="['cursor-pointer', classes.aLinks]"
                    :title="$t('edit.edit')"
                    :to="'/projects/edit/' + project.project_id"
                  >
                    <svg
                      class="w-10 h-10 hover:text-green-500"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                    >
                      <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"
                      />
                    </svg>
                  </router-link>
                </td>
                <td class="px-2 py-3">
                  <router-link
                    class="cursor-pointer"
                    :title="$t('delete.delete')"
                    to
                    @click="deleteItem(index, 'projectsf', project.project_id)"
                  >
                    <svg
                      class="w-10 h-10 hover:text-green-500"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                    >
                      <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                      />
                    </svg>
                  </router-link>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, reactive, toRefs } from 'vue'
import { useRoute } from 'vue-router'
import Menu from '@/views/menu/Menu.vue'
import Title from '@/components/commons/Title.vue'
import SpinnerCommon from '@/components/commons/SpinnerCommon.vue'
import AlertMessage from '@/components/commons/AlertMessage.vue'
import SearchProjectsByFilter from '@/components/search/SearchProjectsByFilter.vue'
import Pagination from '@/components/commons/Pagination.vue'
import ProjectsService, { Project } from '@/api/projectsService'
import AnnotationsService, { Annotation } from '@/api/annotationsService'
import EnvironmentsService, { Environment } from '@/api/environmentsService'
import { useI18n } from 'vue-i18n'

export default defineComponent({
  components: {
    Menu,
    Title,
    SpinnerCommon,
    AlertMessage,
    SearchProjectsByFilter,
    Pagination,
  },
  setup() {
    let state = reactive({
      projects: {
        projects: [] as Project[],
        byFilter: [] as Project[],
        duplicate: {
          data: {} as Project,
          annoations: {} as Annotation[],
          environments: {} as Environment[],
        },
      },
      alert: {
        class: '',
        message: '',
      },
      isOpen: false,
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

    let { projects, loading, alert, pagination, search, classes } =
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
                if (
                  response &&
                  state.projects.duplicate.annoations.length > 0
                ) {
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
      projects,
      loading,
      alert,
      pagination,
      search,
      classes,
      deleteItem,
      launch,
      duplicate,
    }
  },
})
</script>
