<template>
  <div class="lg:flex min-h-screen w-full bg-gray-50">
    <Menu />
    <div class="w-full p-4">
      <div class="block">
        <Title :title="$t('environments.edit')" />
        <SpinnerCommon v-if="loading.loading.active" />
        <AlertMessage :message="alert.message" :classes="alert.class" />
        <div
          class="mx-auto px-3 mt-20 w-full sm:max-w-2xl"
          v-if="
            !loading.loading.active &&
            projects.length > 0 &&
            Object.keys(environment).length > 0
          "
        >
          <Form @submit="submit">
            <CreateUpdateProject
              :projects="projects"
              :project-id="String(environment.project_id)"
              v-model:updateProject="form.project_id"
            />
            <CreateUpdateKey
              :kkey="environment.key"
              :translation="$t('environments.key')"
              v-model:updateKkey="form.key"
            />
            <CreateUpdateValue
              :vvalue="environment.value"
              :translation="$t('environments.value')"
              v-model:updateVvalue="form.value"
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
import CreateUpdateProject from '@/components/commons/CreateUpdateProject.vue'
import CreateUpdateKey from '@/components/commons/CreateUpdateKey.vue'
import CreateUpdateValue from '@/components/commons/CreateUpdateValue.vue'
import SubmitButton from '@/components/buttons/SubmitButton.vue'
import edit from '@/compositions/environments/edit'

const { loading, alert, projects, environment, form, submit } = edit()
</script>
