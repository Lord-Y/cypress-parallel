<template>
  <div class="mt-4">
    <label :for="id" class="block w-full pb-1">{{
      $t('projects.maxPods')
    }}</label>
    <Field
      as=""
      :name="name"
      v-model="local"
      :label="$t('projects.maxPods').toLowerCase()"
      rules="required"
      v-slot="{ field, meta, errorMessage }"
      type="number"
    >
      <input
        v-bind="field"
        type="number"
        :id="id"
        :placeholder="$t('projects.maxPods').toLowerCase()"
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
      default: 'maxPods',
    },
    name: {
      type: String,
      default: 'maxPods',
    },
    maxPods: {
      type: Number,
      default: 10,
    },
  },
  components: {
    Field,
  },
  emits: ['update:updateMaxPods'],
  setup(props, { emit }) {
    const local = computed({
      get: () => {
        emit('update:updateMaxPods', props.maxPods)
        return props.maxPods
      },
      set: (value: number) => emit('update:updateMaxPods', value),
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
