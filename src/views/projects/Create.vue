<template>
  <div class="lg:flex min-h-screen w-full bg-gray-50">
    <Menu />
    <div class="w-full p-4">
      <div class="block">
        <Title :title="$t('projects.create')" />
        <SpinnerCommon v-if="loading.loading.active" />
        <AlertMessage :message="alert.message" :classes="alert.class" />
        <div
          class="mx-auto px-3 mt-20 w-full sm:max-w-2xl"
          v-if="!loading.loading.active && teams.length > 0"
        >
          <Form @submit="submit">
            <CreateUpdateTeam
              :teams="teams"
              :team-id="form.team_id"
              v-model:updateTeam="form.team_id"
            />
            <CreateUpdateName
              :name="form.project_name"
              v-model:updateName="form.project_name"
              id="projectName"
              :translation="$t('projects.name')"
            />
            <CreateUpdateRepository
              :repository="form.repository"
              v-model:updateRepository="form.repository"
            />
            <CreateUpdateUsername
              :username="form.username"
              v-model:updateUsername="form.username"
            />
            <CreateUpdatePassword
              :password="form.password"
              v-model:updatePassword="form.password"
            />
            <CreateUpdateBranch
              :branch="form.branch"
              v-model:updateBranch="form.branch"
            />
            <CreateUpdateSpecs
              :specs="form.specs"
              v-model:updateSpecs="form.specs"
            />
            <CreateUpdateSchedulingEnabled
              v-model:updateSchedulingEnabled="form.schedulingEnabled"
            />
            <CreateUpdateScheduling
              v-if="form.schedulingEnabled"
              v-model:updateScheduling="form.scheduling"
            />
            <CreateUpdateMaxPods
              :max-pods="form.maxPods"
              v-model:updateMaxPods="form.maxPods"
            />
            <CreateUpdateCypressDockerVersion
              :docker-version="form.cypress_docker_version"
              v-model:updateDockerVersion="form.cypress_docker_version"
            />
            <CreateUpdateTimeout
              :timeout="form.timeout"
              v-model:updateTimeout="form.timeout"
            />
            <SubmitButton :text="$t('button.submit')" />
          </Form>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, reactive, toRefs } from 'vue'
import { useI18n } from 'vue-i18n'
import Menu from '@views/menu/Menu.vue'
import Title from '@components/commons/Title.vue'
import SpinnerCommon from '@components/commons/SpinnerCommon.vue'
import AlertMessage from '@components/commons/AlertMessage.vue'
import { Form } from 'vee-validate'
import CreateUpdateTeam from '@components/projects/CreateUpdateTeam.vue'
import CreateUpdateName from '@components/commons/CreateUpdateName.vue'
import CreateUpdateRepository from '@components/projects/CreateUpdateRepository.vue'
import CreateUpdateUsername from '@components/projects/CreateUpdateUsername.vue'
import CreateUpdatePassword from '@components/projects/CreateUpdatePassword.vue'
import CreateUpdateBranch from '@components/projects/CreateUpdateBranch.vue'
import CreateUpdateSpecs from '@components/projects/CreateUpdateSpecs.vue'
import CreateUpdateSchedulingEnabled from '@components/projects/CreateUpdateSchedulingEnabled.vue'
import CreateUpdateScheduling from '@components/projects/CreateUpdateScheduling.vue'
import CreateUpdateMaxPods from '@components/projects/CreateUpdateMaxPods.vue'
import CreateUpdateCypressDockerVersion from '@components/projects/CreateUpdateCypressDockerVersion.vue'
import CreateUpdateTimeout from '@components/projects/CreateUpdateTimeout.vue'
import SubmitButton from '@components/buttons/SubmitButton.vue'
import TeamsService, { Teams } from '@api/teamsService'
import ProjectsService from '@api/projectsService'

export default defineComponent({
  components: {
    Menu,
    Title,
    SpinnerCommon,
    AlertMessage,
    Form,
    CreateUpdateTeam,
    CreateUpdateName,
    CreateUpdateRepository,
    CreateUpdateUsername,
    CreateUpdatePassword,
    CreateUpdateBranch,
    CreateUpdateSpecs,
    CreateUpdateSchedulingEnabled,
    CreateUpdateScheduling,
    CreateUpdateMaxPods,
    CreateUpdateCypressDockerVersion,
    CreateUpdateTimeout,
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
      teams: [] as Teams[],
      form: {
        team_id: '',
        project_name: '',
        repository: '',
        branch: '',
        username: '',
        password: '',
        specs: '',
        scheduling: '',
        schedulingEnabled: false,
        maxPods: 10,
        cypress_docker_version: '7.2.0-0.0.2',
        timeout: 10,
      },
    })
    const { t } = useI18n({
      useScope: 'global',
    })

    TeamsService.all()
      .then((response: any) => {
        switch (response.status) {
          case 200:
            state.teams = response.data
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

    function submit() {
      state.loading.loading.active = true
      ProjectsService.create({
        teamId: Number(state.form.team_id),
        name: state.form.project_name,
        repository: state.form.repository,
        branch: state.form.branch,
        username: state.form.username,
        password: state.form.password,
        specs: state.form.specs,
        scheduling: state.form.scheduling,
        schedulingEnabled: state.form.schedulingEnabled,
        maxPods: state.form.maxPods,
        cypress_docker_version: state.form.cypress_docker_version,
        timeout: state.form.timeout,
      })
        .then((response: any) => {
          if (response.status === 201) {
            state.alert.class = 'green'
            state.alert.message = t('alert.http.team.created', {
              field: state.form.project_name,
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

    let { loading, alert, teams, form } = toRefs(state)

    return {
      loading,
      alert,
      teams,
      form,
      submit,
    }
  },
})
</script>
