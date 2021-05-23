<template>
  <div class="lg:flex min-h-screen w-full bg-gray-50">
    <Menu />
    <div class="w-full p-4">
      <div class="block">
        <Title :title="$t('teams.create')" />
        <SpinnerCommon v-if="loading.loading.active" />
        <AlertMessage :message="alert.message" :classes="alert.class" />
        <div
          class="mx-auto px-3 mt-20 w-full sm:max-w-2xl"
          v-if="!loading.loading.active"
        >
          <Form @submit="submit">
            <CreateUpdateName
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
import Menu from '@views/menu/Menu.vue'
import Title from '@components/commons/Title.vue'
import SpinnerCommon from '@components/commons/SpinnerCommon.vue'
import AlertMessage from '@components/commons/AlertMessage.vue'
import { Form } from 'vee-validate'
import CreateUpdateName from '@/components/commons/CreateUpdateName.vue'
import SubmitButton from '@components/buttons/SubmitButton.vue'
import teamsService from '@api/teamsService'
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
    })
    const { t } = useI18n({
      useScope: 'global',
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

    let { loading, alert, form } = toRefs(state)

    return {
      loading,
      alert,
      form,
      submit,
    }
  },
})
</script>
