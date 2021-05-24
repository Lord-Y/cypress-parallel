<template>
  <div class="mt-4">
    <label :for="id" class="block w-full pb-1">{{
      $t('projects.branch')
    }}</label>
    <Field
      as=""
      :name="name"
      v-model="local"
      :label="$t('projects.branch').toLowerCase()"
      rules="required"
      v-slot="{ field, meta, errorMessage }"
      type="text"
    >
      <input
        v-bind="field"
        type="text"
        :id="id"
        :placeholder="$t('projects.branch').toLowerCase()"
        autocomplete="off"
        class="
          block
          w-full
          border-gray-300
          focus:outline-none
          focus:border-green-500
          focus:ring-green-500
        "
        :class="getValidationClass(meta)"
      />
      <span v-if="errorMessage" class="text-red-500">{{ errorMessage }}</span>
    </Field>
  </div>
</template>

<script lang="ts">
import { defineComponent, computed } from 'vue'
import { Field } from 'vee-validate'

export default defineComponent({
  props: {
    id: {
      type: String,
      default: 'branch',
    },
    name: {
      type: String,
      default: 'branch',
    },
    branch: {
      type: String,
      default: '',
    },
  },
  components: {
    Field,
  },
  emits: ['update:updateBranch'],
  setup(props, { emit }) {
    const local = computed({
      get: () => {
        emit('update:updateBranch', props.branch)
        return props.branch
      },
      set: (value: string) => emit('update:updateBranch', value),
    })
    function getValidationClass(meta: any): string {
      if (meta.valid && meta.validated && meta.dirty) {
        return 'outline-none border-green-500'
      }
      if (!meta.valid && !meta.validated && meta.dirty) {
        return 'outline-none border-red-500'
      }
      return ''
    }

    return {
      local,
      getValidationClass,
    }
  },
})
</script>
