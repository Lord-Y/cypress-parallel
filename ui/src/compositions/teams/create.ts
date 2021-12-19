import { reactive, toRefs } from 'vue'
import { useI18n } from 'vue-i18n'
import { useHead } from '@vueuse/head'
import teamsService from '@/api/teamsService'

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
  })

  const { t } = useI18n({
    useScope: 'global',
  })

  state.meta.title = t('teams.create')
  state.meta.description = t('teams.create')

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

  function submit() {
    state.loading.loading.active = true
    teamsService
      .create({
        name: state.form.name,
      })
      .then((response: any) => {
        if (response.status === 201) {
          state.alert.class = 'green'
          state.alert.message = t('alert.http.team.created', {
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
