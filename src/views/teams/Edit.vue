<template>
  <div class="lg:flex min-h-screen w-full bg-gray-50">
    <Menu />
    <div class="w-full p-4">
      <div class="block">
        <Title :title="$t('teams.edit')" />
        <SpinnerCommon v-if="loading.loading.active" />
        <AlertMessage :message="alert.message" :classes="alert.class" />
        <div
          class="mx-auto px-3 mt-20 w-full sm:max-w-2xl"
          v-if="!loading.loading.active && Object.keys(team).length > 0"
        >
          <Form @submit="submit">
            <CreateUpdateName
              :name="team.team_name"
              v-model:nameUpdate="form.name"
              :translation="$t('teams.name')"
            />
            <SubmitButton :text="$t('button.create')" />
          </Form>
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
import { Form } from 'vee-validate'
import CreateUpdateName from '@/components/commons/CreateUpdateName.vue'
import SubmitButton from '@components/buttons/SubmitButton.vue'
import TeamsService, { Team } from '@api/teamsService'
import { useI18n } from 'vue-i18n'

export default defineComponent({
  components: {
    Menu,
    Title,
    SpinnerCommon,
    AlertMessage,
    Form,
    CreateUpdateName,
    SubmitButton,
  },
  setup() {
    let state = reactive({
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
      })
      .catch((error: any) => {
        if (error.response.status === 404) {
          state.alert.class = 'mute'
          state.alert.message = t('alert.http.pageNotFound')
        } else {
          state.alert.class = 'red'
          state.alert.message = t('alert.http.errorOccured')
        }
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

    let { loading, alert, team, form } = toRefs(state)

    return {
      loading,
      alert,
      team,
      form,
      submit,
    }
  },
})
</script>
