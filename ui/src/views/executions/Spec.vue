<template>
  <div class="lg:flex min-h-screen w-full bg-gray-50">
    <Menu />
    <div class="w-full p-4">
      <div class="block">
        <Title :title="$t('executions.by.spec')" />
        <SpinnerCommon v-if="loading.loading.active" />
        <AlertMessage :message="alert.message" :classes="alert.class" />
        <div
          class="mx-auto px-3 mt-20 w-full overflow-auto"
          v-if="!loading.loading.active && Object.keys(execution).length > 0"
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
              <tr class="hover:bg-gray-100 hover:font-semibold">
                <td class="px-2 py-3">
                  <router-link
                    :class="['cursor-pointer', classes.aLinks]"
                    :title="$t('projects.name') + ' ' + execution.project_name"
                    :to="'/projects/edit/' + execution.project_id"
                    >{{ execution.project_name }}</router-link
                  >
                </td>
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
                <td class="px-2 py-3">{{ execution.spec }}</td>
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
          <div
            class="max-w-5xl mx-auto border-2 p-2"
            v-if="execution.execution_error_output"
          >
            <HR :classes="'my-5 border border-gray-300'" />
            <div class="block">
              <h4 class="text-center font-semibold">
                {{ $t('executions.status.failed') }}
              </h4>
              <div class="break-words border-2 p-2 bg-gray-200 shadow-inner">
                {{ execution.execution_error_output }}
              </div>
            </div>
          </div>
          <div v-if="Object.keys(JSON.parse(execution.result)).length > 0">
            <HR :classes="'my-5 border border-gray-300'" />
            <template
              v-for="(result, index) in JSON.parse(execution.result).results[0]
                .suites[0].tests"
              :key="index"
            >
              <div class="justify-center">
                <div class="flex max-w-5xl mx-auto border-2">
                  <div class="w-5/6 p-1">
                    <div class="flex items-center">
                      <i>
                        <svg
                          class="w-8 h-8"
                          :class="
                            result.state === 'passed'
                              ? 'text-green-500'
                              : 'text-red-500'
                          "
                          fill="none"
                          stroke="currentColor"
                          viewBox="0 0 24 24"
                        >
                          <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
                          />
                        </svg>
                      </i>
                      <h4 class="ml-2">{{ result.title }}</h4>
                    </div>
                  </div>
                  <div class="w-1/6 p-1">
                    <span class="flex items-center justify-end">
                      {{ convertDuration(Number(result.duration)) }}
                      <i class="ml-2">
                        <svg
                          class="w-8 h-8"
                          fill="none"
                          stroke="currentColor"
                          viewBox="0 0 24 24"
                        >
                          <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"
                          />
                        </svg>
                      </i>
                    </span>
                  </div>
                </div>
                <div class="max-w-5xl mx-auto border-2 p-2">
                  <div class="block">
                    <h4 class="text-center font-semibold">
                      {{ $t('executions.status.code') }}
                    </h4>
                    <div
                      class="break-words border-2 p-2 bg-gray-200 shadow-inner"
                    >
                      {{ result.code }}
                    </div>
                    <template v-if="result.err.message">
                      <h4 class="text-center font-semibold">
                        {{ $t('executions.status.err.message') }}
                      </h4>
                      <div
                        class="break-words mt-2 border-2 p-2 bg-gray-200 shadow-inner"
                      >
                        {{ result.err.message }}
                      </div>
                    </template>
                  </div>
                </div>
              </div>
            </template>
          </div>
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
import HR from '@/components/commons/HR.vue'
import spec from '@/compositions/executions/spec'

const {
  execution,
  loading,
  alert,
  classes,
  convertDuration,
  getSpecStatus,
  getSystemStatus,
  getGlobalStatus,
} = spec()
</script>
