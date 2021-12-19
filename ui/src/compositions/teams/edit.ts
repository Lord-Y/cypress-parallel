import { reactive, toRefs } from 'vue'
import { useI18n } from 'vue-i18n'
import { useHead } from '@vueuse/head'
import { useRoute } from 'vue-router'
import TeamsService, { Team } from '@/api/teamsService'

export default function () {
  const state = reactive({
    meta: {
      title: '',
      description: '',
    },
    loading: {
      loading: {
        active: false,
      },
    },
    alert: {
      class: '',
      message: '',
    },
    form: {
      name: '',
    },
    team: {} as Team,
  })

  const route = useRoute()
  const { t } = useI18n({
    useScope: 'global',
  })

  state.meta.title = t('teams.edit')
  state.meta.description = t('teams.edit')

  useHead({
    title: state.meta.title,
    meta: [
      {
        name: 'description',
        content: state.meta.description,
      },
      {
        property: 'og:title',
        content: state.meta.title,
      },
      {
        property: 'og:description',
        content: state.meta.description,
      },
    ],
  })

  state.loading.loading.active = true
  TeamsService.get(Number(route.params.id))
    .then((response: any) => {
      switch (response.status) {
        case 200:
          state.team = response.data
          break
        default:
          state.alert.class = 'red'
          state.alert.message = t('alert.http.errorOccured')
          break
      }
      state.loading.loading.active = false
    })
    .catch((error: any) => {
      if (error.response.status === 404) {
        state.alert.class = 'mute'
        state.alert.message = t('alert.http.pageNotFound')
      } else {
        state.alert.class = 'red'
        state.alert.message = t('alert.http.errorOccured')
      }
      state.loading.loading.active = false
      throw error
    })

  function submit() {
    state.loading.loading.active = true
    TeamsService.update({
      teamId: Number(state.team.team_id),
      name: state.form.name,
    })
      .then((response: any) => {
        if (response.status === 200) {
          state.alert.class = 'green'
          state.alert.message = t('alert.http.team.updated', {
            field: state.form.name,
          })
        } else {
          state.alert.class = 'red'
          state.alert.message = t('alert.http.errorOccured')
        }
        state.loading.loading.active = false
      })
      .catch((error: any) => {
        state.alert.class = 'red'
        state.alert.message = t('alert.http.errorOccured')
        state.loading.loading.active = false
        throw error
      })
  }

  return {
    submit,
    ...toRefs(state),
  }
}
