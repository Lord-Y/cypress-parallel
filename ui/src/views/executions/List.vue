<template>
  <div class="lg:flex min-h-screen w-full bg-gray-50">
    <Menu />
    <div class="w-full p-4">
      <div class="block">
        <Title :title="$t('executions.list')" />
        <SpinnerCommon v-if="loading.loading.active" />
        <AlertMessage :message="alert.message" :classes="alert.class" />
        <SearchExecutionsByFilter
          v-if="search.bar.enabled"
          v-model:loading="loading.loading.active"
          v-model:alertClass="alert.class"
          v-model:alertMessage="alert.message"
          v-model:searchTable="search.table.enabled"
          v-model:byFilter="executions.byFilter"
        />
        <div
          class="mx-auto px-3 mt-20 w-full overflow-auto"
          v-if="
            !loading.loading.active &&
            !search.table.enabled &&
            executions.executions.length > 0
          "
        >
          <table
            class="table-auto w-full text-left border-collapse divide-y border-t-2"
          >
            <thead>
              <tr>
                <th class="py-3">{{ $t('projects.name') }}</th>
                <th class="py-3">{{ $t('executions.uniqId') }}</th>
                <th class="py-3">{{ $t('projects.branch') }}</th>
                <th class="py-3">{{ $t('executions.spec') }}</th>
                <th class="py-3">{{ $t('date.date') }}</th>
              </tr>
            </thead>
            <tbody class="divide-y">
              <tr
                class="hover:bg-gray-100 hover:font-semibold"
                v-for="(execution, index) in executions.executions"
                :id="execution.execution_id.toString()"
                :key="index"
              >
                <td class="px-2 py-3">{{ execution.project_name }}</td>
                <td class="px-2 py-3">
                  <router-link
                    :class="['cursor-pointer', classes.aLinks]"
                    :title="$t('executions.uniqId')"
                    :to="'/executions/uniqid/' + execution.uniq_id"
                  >
                    {{ $t('see.by.uniqId') }}
                    {{ execution.uniq_id }}
                  </router-link>
                </td>
                <td class="px-2 py-3">{{ execution.branch }}</td>
                <td class="px-2 py-3">
                  <router-link
                    :class="['cursor-pointer', classes.aLinks]"
                    :title="$t('see.by.spec')"
                    :to="'/executions/' + execution.execution_id"
                    >{{ execution.spec }}</router-link
                  >
                </td>
                <td class="px-2 py-3">{{ execution.date }}</td>
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
            executions.byFilter.length > 0
          "
        >
          <table
            class="table-auto w-full text-left border-collapse divide-y border-t-2"
          >
            <thead>
              <tr>
                <th class="py-3">{{ $t('projects.name') }}</th>
                <th class="py-3">{{ $t('executions.uniqId') }}</th>
                <th class="py-3">{{ $t('projects.branch') }}</th>
                <th class="py-3">{{ $t('executions.spec') }}</th>
                <th class="py-3">{{ $t('date.date') }}</th>
              </tr>
            </thead>
            <tbody class="divide-y">
              <tr
                class="hover:bg-gray-100 hover:font-semibold"
                v-for="(execution, index) in executions.byFilter"
                :id="execution.execution_id.toString()"
                :key="index"
              >
                <td class="px-2 py-3">{{ execution.project_name }}</td>
                <td class="px-2 py-3">
                  <router-link
                    :class="['cursor-pointer', classes.aLinks]"
                    :title="$t('executions.uniqId')"
                    :to="'/executions/uniqid/' + execution.uniq_id"
                  >
                    {{ $t('see.by.uniqId') }}
                    {{ execution.uniq_id }}
                  </router-link>
                </td>
                <td class="px-2 py-3">{{ execution.branch }}</td>
                <td class="px-2 py-3">
                  <router-link
                    :class="['cursor-pointer', classes.aLinks]"
                    :title="$t('see.by.spec')"
                    :to="'/executions/' + execution.execution_id"
                    >{{ execution.spec }}</router-link
                  >
                </td>
                <td class="px-2 py-3">{{ execution.date }}</td>
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
import SearchExecutionsByFilter from '@/components/search/SearchExecutionsByFilter.vue'
import Pagination from '@/components/commons/Pagination.vue'
import list from '@/compositions/executions/list'

const { executions, loading, alert, pagination, search, classes } = list()
</script>
