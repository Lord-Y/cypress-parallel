<template>
  <div class="lg:flex min-h-screen w-full bg-gray-50">
    <Menu />
    <div class="w-full p-4">
      <div class="block">
        <Title :title="$t('environments.create')" />
        <SpinnerCommon v-if="loading.loading.active" />
        <AlertMessage :message="alert.message" :classes="alert.class" />
        <div
          class="mx-auto px-3 mt-20 w-full sm:max-w-2xl"
          v-if="!loading.loading.active && projects.length > 0"
        >
          <Form @submit="submit">
            <CreateUpdateProject
              :projects="projects"
              :project-id="form.project_id"
              v-model:updateProject="form.project_id"
            />
            <CreateUpdateKey
              :kkey="form.key"
              :translation="$t('environments.key')"
              v-model:updateKkey="form.key"
            />
            <CreateUpdateValue
              :vvalue="form.value"
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

<script lang="ts">
import { defineComponent, reactive, toRefs } from 'vue'
import { useI18n } from 'vue-i18n'
import Menu from '@views/menu/Menu.vue'
import Title from '@components/commons/Title.vue'
import SpinnerCommon from '@components/commons/SpinnerCommon.vue'
import AlertMessage from '@components/commons/AlertMessage.vue'
import { Form } from 'vee-validate'
import CreateUpdateProject from '@components/commons/CreateUpdateProject.vue'
import CreateUpdateKey from '@/components/commons/CreateUpdateKey.vue'
import CreateUpdateValue from '@/components/commons/CreateUpdateValue.vue'
import SubmitButton from '@components/buttons/SubmitButton.vue'
import ProjectsService, { ProjectOnly } from '@api/projectsService'
import EnvironmentsService from '@api/environmentsService'

export default defineComponent({
  components: {
    Menu,
    Title,
    SpinnerCommon,
    AlertMessage,
    Form,
    CreateUpdateProject,
    CreateUpdateKey,
    CreateUpdateValue,
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
      projects: [] as ProjectOnly[],
      form: {
        project_id: '',
        key: '',
        value: '',
      },
    })
    const { t } = useI18n({
      useScope: 'global',
    })

    ProjectsService.all()
      .then((response: any) => {
        switch (response.status) {
          case 200:
            state.projects = response.data
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
      EnvironmentsService.create({
        projectId: Number(state.form.project_id),
        key: state.form.key,
        value: state.form.value,
      })
        .then((response: any) => {
          if (response.status === 201) {
            state.alert.class = 'green'
            state.alert.message = t('alert.http.environment.created', {
              field: state.form.key,
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

    let { loading, alert, projects, form } = toRefs(state)

    return {
      loading,
      alert,
      projects,
      form,
      submit,
    }
  },
})
</script>
