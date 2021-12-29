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
              :team-id="form.team_id.toString()"
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
            <CreateUpdateBrowser
              :browser="form.browser"
              v-model:updateBrowser="form.browser"
            />
            <CreateUpdateConfigFile
              :config-file="form.config_file"
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
import create from '@/compositions/projects/create'

const { loading, alert, teams, form, submit } = create()
</script>
