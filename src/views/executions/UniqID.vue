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
                <td class="px-2 py-3">
                  {{ execution.uniq_id }}
                </td>
                <td class="px-2 py-3">{{ execution.branch }}</td>
                <td class="px-2 py-3">
                  <router-link
                    :class="['cursor-pointer', classes.aLinks]"
                    :title="$t('see.by.spec')"
                    :to="'/executions/' + execution.execution_id"
                  >
                    {{ execution.spec }}
                  </router-link>
                </td>
                <td class="px-2 py-3" :class="classes.status.global">
                  {{ getGlobalStatus(execution) }}
                </td>
                <td class="px-2 py-3" :class="classes.status.system">
                  {{ getSystemStatus(execution.execution_status) }}
                </td>
                <td class="px-2 py-3" :class="classes.status.global">
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

<script lang="ts">
import { defineComponent, reactive, toRefs } from 'vue'
import { useRoute } from 'vue-router'
import Menu from '@views/menu/Menu.vue'
import Title from '@components/commons/Title.vue'
import SpinnerCommon from '@components/commons/SpinnerCommon.vue'
import AlertMessage from '@components/commons/AlertMessage.vue'
import ExecutionsService, { Execution } from '@api/executionsService'
import Statuses from '@/tools/status'
import { useI18n } from 'vue-i18n'

export default defineComponent({
  components: {
    Menu,
    Title,
    SpinnerCommon,
    AlertMessage,
  },
  setup() {
    let state = reactive({
      executions: [] as Execution[],
      alert: {
        class: '',
        message: '',
      },
      isOpen: false,
      loading: {
        loading: {
          active: true,
        },
      },
      classes: {
        aLinks: 'hover:text-green-500 hover:font-extrabold',
        status: {
          global: '',
          system: '',
        },
      },
    })
    const route = useRoute()
    const { t } = useI18n({
      useScope: 'global',
    })
    let id: string
    id = String(route.params.id)

    state.loading.loading.active = true

    ExecutionsService.uniqid(id)
      .then((response: any) => {
        switch (response.status) {
          case 200:
            state.executions = response.data
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

    function getSpecStatus(execution: Execution): string {
      return Statuses.tests(execution)
    }

    function getSystemStatus(s: string): string {
      switch (s) {
        case 'DONE':
          state.classes.status.system = 'text-green-500 font-semibold'
          break
        case 'NOT_STARTED':
          state.classes.status.system = 'text-gray-500 font-semibold'
          break
        default:
          state.classes.status.system = 'text-red-500 font-semibold'
          break
      }
      return s
    }

    function getGlobalStatus(execution: Execution): string {
      let status: string
      status = Statuses.global(execution)
      switch (status) {
        case 'PASSED':
          state.classes.status.global = 'text-green-500 font-semibold'
          break
        case 'NOT_STARTED':
        case 'SCHEDULED':
        case 'CANCELLED':
          state.classes.status.global = 'text-gray-500 font-semibold'
          break
        default:
          state.classes.status.global = 'text-red-500 font-semibold'
          break
      }
      return Statuses.global(execution)
    }

    state.loading.loading.active = false

    let { executions, loading, alert, classes } = toRefs(state)

    return {
      executions,
      loading,
      alert,
      classes,
      getSpecStatus,
      getSystemStatus,
      getGlobalStatus,
    }
  },
})
</script>
