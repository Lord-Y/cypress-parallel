import { reactive, toRefs } from 'vue'
import { useRoute } from 'vue-router'

export default function () {
  const state = reactive({
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
  })

  const route = useRoute()

  function getSelectedMenu(data: string): string | undefined {
    if (route.meta.root === data) {
      switch (route.meta.root) {
        case 'teams':
          state.menu.isOpen.teams = true
          break
        case 'projects':
          state.menu.isOpen.projects = true
          break
        case 'annotations':
          state.menu.isOpen.annotations = true
          break
        case 'environments':
          state.menu.isOpen.environments = true
          break
        case 'executions':
          state.menu.isOpen.executions = true
          break
      }
      return 'bg-emerald-500'
    }
  }

  function getActiveLink(data: string): string | undefined {
    if (route.meta.activeLink === data) {
      return 'border-l-4 border-emerald-500'
    }
  }

  return {
    getSelectedMenu,
    getActiveLink,
    ...toRefs(state),
  }
}
