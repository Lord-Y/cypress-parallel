<template>
  <div class="lg:flex min-h-screen w-full bg-gray-50">
    <Menu />
    <div class="w-full p-4">
      <div class="block">
        <Title :title="$t('executions.by.uniqId')" />
        <SpinnerCommon v-if="loading.loading.active" />
        <AlertMessage :message="alert.message" :classes="alert.class" />
        <div
          class="mx-auto px-3 mt-20 w-full overflow-auto"
          v-if="!loading.loading.active && executions.length > 0"
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
                <th class="py-3">{{ $t('executions.status.global') }}</th>
                <th class="py-3">{{ $t('executions.status.system') }}</th>
                <th class="py-3">{{ $t('executions.status.spec') }}</th>
                <th class="py-3">{{ $t('date.date') }}</th>
              </tr>
            </thead>
            <tbody class="divide-y">
              <tr
                class="hover:bg-gray-100 hover:font-semibold"
                v-for="(execution, index) in executions"
                :id="execution.execution_id.toString()"
                :key="index"
              >
                <td class="px-2 py-3">{{ execution.project_name }}</td>
                <td class="px-2 py-3">{{ execution.uniq_id }}</td>
                <td class="px-2 py-3">{{ execution.branch }}</td>
                <td class="px-2 py-3">
                  <router-link
                    :class="['cursor-pointer', classes.aLinks]"
                    :title="$t('see.by.spec')"
                    :to="'/executions/' + execution.execution_id"
                    >{{ execution.spec }}</router-link
                  >
                </td>
                <td
                  class="px-2 py-3"
                  :class="getGlobalStatus(execution, 'classes')"
                >
                  {{ getGlobalStatus(execution, '') }}
                </td>
                <td
                  class="px-2 py-3"
                  :class="getSystemStatus(execution.execution_status)"
                >
                  {{ execution.execution_status }}
                </td>
                <td
                  class="px-2 py-3"
                  :class="getGlobalStatus(execution, 'classes')"
                >
                  {{ getSpecStatus(execution) }}
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
import uniqID from '@/compositions/executions/uniqID'

const {
  executions,
  loading,
  alert,
  classes,
  getSpecStatus,
  getSystemStatus,
  getGlobalStatus,
} = uniqID()
</script>
