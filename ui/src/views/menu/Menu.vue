<template>
  <div class="w-full lg:w-80 bg-gray-900">
    <header class="block md:px-4 md:py-1">
      <div
        class="flex justify-between lg:justify-start items-center px-4 py-1 md:p-0"
      >
        <div>
          <img class="h-8 items-center" src="/logo.png" alt="Logo" />
        </div>
        <div class="px-2 pt-2 pb-4">
          <router-link
            class="block text-green-500 font-semibold hover:text-white px-2 py-1 uppercase"
            to="/"
            :title="$t('brand.name')"
            >{{ $t('brand.name') }}</router-link
          >
        </div>
        <div class="lg:hidden">
          <button
            type="button"
            @click="isOpen = !isOpen"
            class="block text-green-500 focus:text-white focus:outline-none"
          >
            <svg class="h-6 w-6 fill-current" viewBox="0 0 24 24">
              <path
                v-if="isOpen"
                fill-rule="evenodd"
                d="M18.278 16.864a1 1 0 0 1-1.414 1.414l-4.829-4.828-4.828 4.828a1 1 0 0 1-1.414-1.414l4.828-4.829-4.828-4.828a1 1 0 0 1 1.414-1.414l4.829 4.828 4.828-4.828a1 1 0 1 1 1.414 1.414l-4.828 4.829 4.828 4.828z"
              />
              <path
                v-if="!isOpen"
                fill-rule="evenodd"
                d="M4 5h16a1 1 0 0 1 0 2H4a1 1 0 1 1 0-2zm0 6h16a1 1 0 0 1 0 2H4a1 1 0 0 1 0-2zm0 6h16a1 1 0 0 1 0 2H4a1 1 0 0 1 0-2z"
              />
            </svg>
          </button>
        </div>
      </div>
      <div :class="isOpen ? 'block' : 'hidden'" class="pt-2 pb-4 lg:flex">
        <nav class="w-full">
          <div>
            <button
              type="button"
              class="block text-white w-full text-left p-3 border-t-2 border-b-2 border-black"
              :class="getSelectedMenu('teams')"
              @click="menu.isOpen.teams = !menu.isOpen.teams"
            >
              {{ $t('teams.teams') }}
            </button>
            <div class="bg-black" v-if="menu.isOpen.teams">
              <router-link
                class="block text-white hover:text-green-500 hover:font-extrabold p-2"
                :class="getActiveLink('/teams/create')"
                to="/teams/create"
                :title="$t('teams.create')"
                >{{ $t('teams.create') }}</router-link
              >
              <router-link
                class="block text-white hover:text-green-500 hover:font-extrabold p-2"
                :class="getActiveLink('/teams/list')"
                to="/teams/list"
                :title="$t('teams.list')"
                >{{ $t('teams.list') }}</router-link
              >
            </div>
          </div>
          <div>
            <button
              type="button"
              class="block text-white w-full text-left p-3 border-b-2 border-black"
              :class="getSelectedMenu('projects')"
              @click="menu.isOpen.projects = !menu.isOpen.projects"
            >
              {{ $t('projects.projects') }}
            </button>
            <div class="bg-black" v-if="menu.isOpen.projects">
              <router-link
                class="block text-white hover:text-green-500 hover:font-extrabold p-2"
                :class="getActiveLink('/projects/create')"
                to="/projects/create"
                :title="$t('projects.create')"
                >{{ $t('projects.create') }}</router-link
              >
              <router-link
                class="block text-white hover:text-green-500 hover:font-extrabold p-2"
                :class="getActiveLink('/projects/list')"
                to="/projects/list"
                :title="$t('projects.list')"
                >{{ $t('projects.list') }}</router-link
              >
            </div>
          </div>
          <div>
            <button
              type="button"
              class="block text-white w-full text-left p-3 border-b-2 border-black"
              :class="getSelectedMenu('annotations')"
              @click="menu.isOpen.annotations = !menu.isOpen.annotations"
            >
              {{ $t('annotations.annotations') }}
            </button>
            <div class="bg-black" v-if="menu.isOpen.annotations">
              <router-link
                class="block text-white hover:text-green-500 hover:font-extrabold p-2"
                :class="getActiveLink('/annotations/create')"
                to="/annotations/create"
                :title="$t('annotations.create')"
                >{{ $t('annotations.create') }}</router-link
              >
              <router-link
                class="block text-white hover:text-green-500 hover:font-extrabold p-2"
                :class="getActiveLink('/annotations/list')"
                to="/annotations/list"
                :title="$t('annotations.list')"
                >{{ $t('annotations.list') }}</router-link
              >
            </div>
          </div>
          <div>
            <button
              type="button"
              class="block text-white w-full text-left p-3 border-b-2 border-black"
              :class="getSelectedMenu('environments')"
              @click="menu.isOpen.environments = !menu.isOpen.environments"
            >
              {{ $t('environments.environments') }}
            </button>
            <div class="bg-black" v-if="menu.isOpen.environments">
              <router-link
                class="block text-white hover:text-green-500 hover:font-extrabold p-2"
                :class="getActiveLink('/environments/create')"
                to="/environments/create"
                :title="$t('environments.create')"
                >{{ $t('environments.create') }}</router-link
              >
              <router-link
                class="block text-white hover:text-green-500 hover:font-extrabold p-2"
                :class="getActiveLink('/environments/list')"
                to="/environments/list"
                :title="$t('environments.list')"
                >{{ $t('environments.list') }}</router-link
              >
            </div>
          </div>
          <div>
            <button
              type="button"
              class="block text-white w-full text-left p-3 border-b-2 border-black"
              :class="getSelectedMenu('executions')"
              @click="menu.isOpen.executions = !menu.isOpen.executions"
            >
              {{ $t('executions.executions', 2) }}
            </button>
            <div class="bg-black" v-if="menu.isOpen.executions">
              <router-link
                class="block text-white hover:text-green-500 hover:font-extrabold p-2"
                :class="getActiveLink('/executions/list')"
                to="/executions/list"
                :title="$t('executions.list')"
                >{{ $t('executions.list') }}</router-link
              >
            </div>
          </div>
        </nav>
      </div>
    </header>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'

export default defineComponent({
  data() {
    return {
      isOpen: false,
      menu: {
        isOpen: {
          teams: false,
          projects: false,
          annotations: false,
          environments: false,
          executions: false,
        },
      },
    }
  },
  methods: {
    getSelectedMenu(data: string): string | undefined {
      if (this.$route.meta.root === data) {
        switch (this.$route.meta.root) {
          case 'teams':
            this.menu.isOpen.teams = true
            break
          case 'projects':
            this.menu.isOpen.projects = true
            break
          case 'annotations':
            this.menu.isOpen.annotations = true
            break
          case 'environments':
            this.menu.isOpen.environments = true
            break
          case 'executions':
            this.menu.isOpen.executions = true
            break
        }
        return 'bg-green-500'
      }
    },
    getActiveLink(data: string): string | undefined {
      if (this.$route.meta.activeLink === data) {
        return 'border-l-4 border-green-500'
      }
    },
  },
})
</script>
