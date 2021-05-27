<template>
  <div class="mt-4">
    <label :for="id" class="block w-full pb-1">{{
      $t('projects.specs')
    }}</label>
    <Field
      as=""
      :name="name"
      v-model="local"
      :label="$t('projects.specs').toLowerCase()"
      rules="required"
      v-slot="{ field, meta, errorMessage }"
      type="text"
    >
      <input
        v-bind="field"
        type="text"
        :id="id"
        :placeholder="$t('projects.specs').toLowerCase()"
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
import { defineRule } from 'vee-validate'
import { required } from '@vee-validate/rules'
defineRule('required', required)

export default defineComponent({
  props: {
    id: {
      type: String,
      default: 'specs',
    },
    name: {
      type: String,
      default: 'specs',
    },
    specs: {
      type: String,
      default: '',
    },
  },
  components: {
    Field,
  },
  emits: ['update:updateSpecs'],
  setup(props, { emit }) {
    const local = computed({
      get: () => {
        emit('update:updateSpecs', props.specs)
        return props.specs
      },
      set: (value: string) => emit('update:updateSpecs', value),
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
