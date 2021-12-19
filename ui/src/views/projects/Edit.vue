<template>
  <div class="lg:flex min-h-screen w-full bg-gray-50">
    <Menu />
    <div class="w-full p-4">
      <div class="block">
        <Title :title="$t('projects.edit')" />
        <SpinnerCommon v-if="loading.loading.active" />
        <AlertMessage :message="alert.message" :classes="alert.class" />
        <div
          class="mx-auto px-3 mt-20 w-full sm:max-w-2xl"
          v-if="
            !loading.loading.active &&
            teams.length > 0 &&
            Object.keys(project).length > 0
          "
        >
          <Form @submit="submit">
            <CreateUpdateTeam
              :teams="teams"
              :team-id="String(project.team_id)"
              v-model:updateTeam="form.team_id"
            />
            <CreateUpdateName
              :name="project.project_name"
              v-model:updateName="form.project_name"
              id="projectName"
              :translation="$t('projects.name')"
            />
            <CreateUpdateRepository
              :repository="project.repository"
              v-model:updateRepository="form.repository"
            />
            <CreateUpdateUsername
              :username="project.username"
              v-model:updateUsername="form.username"
            />
            <CreateUpdatePassword
              :password="project.password"
              v-model:updatePassword="form.password"
            />
            <CreateUpdateBranch
              :branch="project.branch"
              v-model:updateBranch="form.branch"
            />
            <CreateUpdateSpecs
              :specs="project.specs"
              v-model:updateSpecs="form.specs"
            />
            <CreateUpdateSchedulingEnabled
              :scheduling-enabled="project.scheduling_enabled"
              v-model:updateSchedulingEnabled="form.schedulingEnabled"
            />
            <CreateUpdateScheduling
              v-if="form.schedulingEnabled"
              :scheduling="project.scheduling"
              v-model:updateScheduling="form.scheduling"
            />
            <CreateUpdateMaxPods
              :max-pods="Number(project.max_pods)"
              v-model:updateMaxPods="form.maxPods"
            />
            <CreateUpdateCypressDockerVersion
              :docker-version="project.cypress_docker_version"
              v-model:updateDockerVersion="form.cypress_docker_version"
            />
            <CreateUpdateTimeout
              :timeout="Number(project.timeout)"
              v-model:updateTimeout="form.timeout"
            />
            <CreateUpdateBrowser
              :browser="project.browser"
              v-model:updateBrowser="form.browser"
            />
            <CreateUpdateConfigFile
              :config-file="project.config_file"
              v-model:updateConfigFile="form.config_file"
            />
            <SubmitButton :text="$t('button.submit')" />
          </Form>
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
import { Form } from 'vee-validate'
import CreateUpdateTeam from '@/components/projects/CreateUpdateTeam.vue'
import CreateUpdateName from '@/components/commons/CreateUpdateName.vue'
import CreateUpdateRepository from '@/components/projects/CreateUpdateRepository.vue'
import CreateUpdateUsername from '@/components/projects/CreateUpdateUsername.vue'
import CreateUpdatePassword from '@/components/projects/CreateUpdatePassword.vue'
import CreateUpdateBranch from '@/components/projects/CreateUpdateBranch.vue'
import CreateUpdateSpecs from '@/components/projects/CreateUpdateSpecs.vue'
import CreateUpdateSchedulingEnabled from '@/components/projects/CreateUpdateSchedulingEnabled.vue'
import CreateUpdateScheduling from '@/components/projects/CreateUpdateScheduling.vue'
import CreateUpdateMaxPods from '@/components/projects/CreateUpdateMaxPods.vue'
import CreateUpdateCypressDockerVersion from '@/components/projects/CreateUpdateCypressDockerVersion.vue'
import CreateUpdateTimeout from '@/components/projects/CreateUpdateTimeout.vue'
import CreateUpdateBrowser from '@/components/projects/CreateUpdateBrowser.vue'
import CreateUpdateConfigFile from '@/components/projects/CreateUpdateConfigFile.vue'
import SubmitButton from '@/components/buttons/SubmitButton.vue'
import edit from '@/compositions/projects/edit'

const { loading, alert, teams, project, form, submit } = edit()
</script>
