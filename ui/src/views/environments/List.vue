<template>
  <div class="lg:flex min-h-screen w-full bg-gray-50">
    <Menu />
    <div class="w-full p-4">
      <div class="block">
        <Title :title="$t('environments.list')" />
        <SpinnerCommon v-if="loading.loading.active" />
        <SpinnerCommon v-if="loading.delete.active" />
        <AlertMessage :message="alert.message" :classes="alert.class" />
        <SearchEnvironmentsByFilter
          v-if="search.bar.enabled"
          v-model:loading="loading.loading.active"
          v-model:alertClass="alert.class"
          v-model:alertMessage="alert.message"
          v-model:searchTable="search.table.enabled"
          v-model:byFilter="environments.byFilter"
        />
        <div
          class="mx-auto px-3 mt-20 w-full overflow-auto"
          v-if="
            !loading.loading.active &&
            !search.table.enabled &&
            environments.environments.length > 0
          "
        >
          <table
            class="table-auto w-full text-left border-collapse divide-y border-t-2"
          >
            <thead>
              <tr>
                <th class="py-3">{{ $t('projects.name') }}</th>
                <th class="py-3">{{ $t('environments.key') }}</th>
                <th class="py-3">{{ $t('environments.value') }}</th>
                <th class="py-3">{{ $t('edit.edit') }}</th>
                <th class="py-3">{{ $t('delete.delete') }}</th>
              </tr>
            </thead>
            <tbody class="divide-y">
              <tr
                class="hover:bg-gray-100 hover:font-semibold"
                v-for="(environment, index) in environments.environments"
                :id="environment.environment_id.toString()"
                :key="index"
              >
                <td class="px-2 py-3">
                  <router-link
                    :class="['cursor-pointer', classes.aLinks]"
                    :title="$t('edit.edit')"
                    :to="'/projects/edit/' + environment.project_id"
                    >{{ environment.project_name }}</router-link
                  >
                </td>
                <td class="px-2 py-3">{{ environment.key }}</td>
                <td class="px-2 py-3">{{ environment.value }}</td>
                <td class="px-2 py-3">
                  <router-link
                    :class="['cursor-pointer', classes.aLinks]"
                    :title="$t('edit.edit')"
                    :to="'/environments/edit/' + environment.environment_id"
                  >
                    <svg
                      class="w-10 h-10 hover:text-emerald-500"
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
                    to=""
                    @click="
                      deleteItem(
                        index,
                        'environments',
                        environment.environment_id,
                      )
                    "
                  >
                    <svg
                      class="w-10 h-10 hover:text-emerald-500"
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
            environments.byFilter.length > 0
          "
        >
          <table
            class="table-auto w-full text-left border-collapse divide-y border-t-2"
          >
            <thead>
              <tr>
                <th class="py-3">{{ $t('projects.name') }}</th>
                <th class="py-3">{{ $t('environments.key') }}</th>
                <th class="py-3">{{ $t('environments.value') }}</th>
                <th class="py-3">{{ $t('edit.edit') }}</th>
                <th class="py-3">{{ $t('delete.delete') }}</th>
              </tr>
            </thead>
            <tbody class="divide-y">
              <tr
                class="hover:bg-gray-100 hover:font-semibold"
                v-for="(environment, index) in environments.byFilter"
                :id="environment.environment_id.toString()"
                :key="index"
              >
                <td class="px-2 py-3">
                  <router-link
                    :class="['cursor-pointer', classes.aLinks]"
                    :title="$t('edit.edit')"
                    :to="'/projects/edit/' + environment.project_id"
                    >{{ environment.project_name }}</router-link
                  >
                </td>
                <td class="px-2 py-3">{{ environment.key }}</td>
                <td class="px-2 py-3">{{ environment.value }}</td>
                <td class="px-2 py-3">
                  <router-link
                    :class="['cursor-pointer', classes.aLinks]"
                    :title="$t('edit.edit')"
                    :to="'/environments/edit/' + environment.environment_id"
                  >
                    <svg
                      class="w-10 h-10 hover:text-emerald-500"
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
                    to=""
                    @click="
                      deleteItem(
                        index,
                        'environmentsf',
                        environment.environment_id,
                      )
                    "
                  >
                    <svg
                      class="w-10 h-10 hover:text-emerald-500"
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

<script setup lang="ts">
import Menu from '@/views/menu/Menu.vue'
import Title from '@/components/commons/Title.vue'
import SpinnerCommon from '@/components/commons/SpinnerCommon.vue'
import AlertMessage from '@/components/commons/AlertMessage.vue'
import SearchEnvironmentsByFilter from '@/components/search/SearchEnvironmentsByFilter.vue'
import Pagination from '@/components/commons/Pagination.vue'
import list from '@/compositions/environments/list'

const {
  environments,
  loading,
  alert,
  pagination,
  search,
  classes,
  deleteItem,
} = list()
</script>
