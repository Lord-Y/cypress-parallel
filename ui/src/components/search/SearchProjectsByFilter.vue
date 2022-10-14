<template>
  <div class="mt-4">
    <div class="relative w-10/12 mx-auto mt-16">
      <input
        type="text"
        :id="id"
        name="searchByFilter"
        class="block w-full border-gray-300 focus:outline-none focus:border-emerald-500 focus:ring-emerald-500 text-center"
        :placeholder="$t('projects.search')"
        autocomplete="off"
        v-model.trim="search.byFilter"
        @keyup="fetchDataByFilter"
      />
      <button
        type="button"
        class="absolute inset-y-0 right-0 flex items-center px-4 font-bold text-white bg-emerald-500"
        @click="clear"
      >
        <svg
          class="w-6 h-6"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M6 18L18 6M6 6l12 12"
          />
        </svg>
      </button>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, reactive, toRefs } from 'vue'
import ProjectsService, { Projects } from '@/api/projectsService'
import { useI18n } from 'vue-i18n'

export default defineComponent({
  props: {
    id: {
      type: String,
      default: 'searchByFilter',
    },
  },
  emits: [
    'update:loading',
    'update:alertClass',
    'update:alertMessage',
    'update:searchTable',
    'update:byFilter',
  ],
  setup(props, { emit }) {
    let state = reactive({
      search: {
        byFilter: '',
      },
    })
    const { t } = useI18n({
      useScope: 'global',
    })
    let projects: Projects[]

    function fetchDataByFilter() {
      if (state.search.byFilter.length >= 2) {
        emit('update:loading', true)
        ProjectsService.search(state.search.byFilter.trim(), 1)
          .then((response) => {
            switch (response.status) {
              case 200:
                projects = response.data
                emit('update:byFilter', projects)
                emit('update:searchTable', true)
                break
              case 204:
                emit('update:alertClass', 'mute')
                emit('update:alertMessage', t('alert.http.noDataFound'))
                break
              default:
                emit('update:alertClass', 'red')
                emit('update:alertMessage', t('alert.http.errorOccured'))
                break
            }
            emit('update:loading', false)
          })
          .catch((error: any) => {
            emit('update:alertClass', 'red')
            emit('update:alertMessage', t('alert.http.errorOccured'))
            emit('update:loading', false)
            throw error
          })
      } else {
        emit('update:byFilter', projects)
        emit('update:searchTable', false)
        emit('update:loading', false)
      }
    }

    function clear() {
      state.search.byFilter = ''
      emit('update:byFilter', projects)
      emit('update:searchTable', false)
      emit('update:loading', false)
    }

    let { search } = toRefs(state)
    return {
      search,
      fetchDataByFilter,
      clear,
    }
  },
})
</script>
