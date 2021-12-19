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
              v-model:updateName="form.name"
              :translation="$t('teams.name')"
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
import CreateUpdateName from '@/components/commons/CreateUpdateName.vue'
import SubmitButton from '@/components/buttons/SubmitButton.vue'
import edit from '@/compositions/teams/edit'

const { loading, alert, form, team, submit } = edit()
</script>
