<template>
  <div class="lg:flex min-h-screen w-full bg-gray-50">
    <Menu />
    <div class="w-full p-4">
      <div class="block">
        <Title :title="$t('teams.list')" />
        <SpinnerCommon v-if="loading.loading.active" />
        <AlertMessage :message="alert.message" :classes="alert.class" />
        <SearchTeamsByFilter
          v-if="search.bar.enabled"
          @update:loading="loading.loading.active = $event"
          @update:alertClass="alert.class = $event"
          @update:alertMessage="alert.message = $event"
          @update:searchTable="search.table.enabled = $event"
          @update:byFilter="teams.byFilter = $event"
        />
        <div
          class="mx-auto px-3 mt-20 w-full overflow-auto"
          v-if="
            !loading.loading.active &&
            !search.table.enabled &&
            teams.teams.length > 0
          "
        >
          <table
            class="
              table-auto
              w-full
              text-left
              border-collapse
              divide-y
              border-t-2
            "
          >
            <thead>
              <tr>
                <th class="py-3">{{ $t('teams.teams') }}</th>
                <th class="py-3">{{ $t('edit.edit') }}</th>
                <th class="py-3">{{ $t('delete.delete') }}</th>
              </tr>
            </thead>
            <tbody class="divide-y">
              <tr
                class="hover:bg-gray-100 hover:font-semibold"
                v-for="(team, index) in teams.teams"
                :id="team.team_id.toString()"
                :key="index"
              >
                <td class="px-2 py-3">{{ team.team_name }}</td>
                <td class="px-2 py-3">
                  <router-link
                    :class="['cursor-pointer', classes.aLinks]"
                    :title="$t('edit.edit')"
                    :to="'/teams/edit/' + team.team_id"
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
                      ></path>
                    </svg>
                  </router-link>
                </td>

                <td class="px-2 py-3">
                  <router-link
                    class="cursor-pointer"
                    :title="$t('delete.delete')"
                    to=""
                    @click="deleteItem(index, 'teams', team.team_id)"
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
                      ></path>
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
            teams.byFilter.length > 0
          "
        >
          <table
            class="
              table-auto
              w-full
              text-left
              border-collapse
              divide-y
              border-t-2
            "
          >
            <thead>
              <tr>
                <th class="py-3">{{ $t('teams.teams') }}</th>
                <th class="py-3">{{ $t('edit.edit') }}</th>
                <th class="py-3">{{ $t('delete.delete') }}</th>
              </tr>
            </thead>
            <tbody class="divide-y">
              <tr
                class="hover:bg-gray-100 hover:font-semibold"
                v-for="(team, index) in teams.byFilter"
                :id="team.team_id.toString()"
                :key="index"
              >
                <td class="px-2 py-3">{{ team.team_name }}</td>
                <td class="px-2 py-3">
                  <router-link
                    :class="['cursor-pointer', classes.aLinks]"
                    :title="$t('edit.edit')"
                    :to="'/teams/edit/' + team.team_id"
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
                      ></path>
                    </svg>
                  </router-link>
                </td>

                <td class="px-2 py-3">
                  <router-link
                    class="cursor-pointer"
                    :title="$t('delete.delete')"
                    to=""
                    @click="deleteItem(index, 'teamsf', team.team_id)"
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
                      ></path>
                    </svg>
                  </router-link>
                </td>
              </tr>
            </tbody>
          </table>
          <Pagination v-if="pagination.enabled" :pagination="pagination.data" />
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, reactive, toRefs } from 'vue'
import { useRoute } from 'vue-router'
import Menu from '@views/menu/Menu.vue'
import Title from '@components/commons/Title.vue'
import SpinnerCommon from '@components/commons/SpinnerCommon.vue'
import AlertMessage from '@components/commons/AlertMessage.vue'
import SearchTeamsByFilter from '@components/search/SearchTeamsByFilter.vue'
import Pagination from '@components/commons/Pagination.vue'
import TeamsService, { Teams } from '@api/teamsService'
import { useI18n } from 'vue-i18n'

export default defineComponent({
  components: {
    Menu,
    Title,
    SpinnerCommon,
    AlertMessage,
    SearchTeamsByFilter,
    Pagination,
  },
  setup() {
    let state = reactive({
      teams: {
        teams: [] as Teams[],
        byFilter: [] as Teams[],
      },
      alert: {
        message: '',
        class: '',
      },
      isOpen: false,
      loading: {
        loading: {
          active: true,
        },
        delete: {
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
          enabled: true,
        },
        table: {
          enabled: false,
        },
      },
      classes: {
        aLinks: 'hover:text-green-500 hover:font-extrabold',
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

    TeamsService.get(page)
      .then((response: any) => {
        switch (response.status) {
          case 200:
            state.teams.teams = response.data
            total = state.teams.teams[0].total
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

    let { teams, loading, alert, pagination, search, classes } = toRefs(state)

    function deleteItem(index: number, type: string, id: number): void {
      if (confirm(t('confirm.sure'))) {
        state.alert.message = ''
        state.loading.delete.active = true
        TeamsService.delete(id)
          .then(() => {
            if (type === 'teams') {
              state.teams.teams.splice(index, 1)
              if (state.teams.teams.length == 0) {
                state.alert.class = 'mute'
                state.alert.message = t('alert.http.noDataFound')
              }
            } else {
              state.teams.byFilter.splice(index, 1)
              if (state.teams.byFilter.length == 0) {
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

    return {
      teams,
      loading,
      alert,
      pagination,
      search,
      classes,
      deleteItem,
    }
  },
})
</script>
