<template>
  <div class="mt-4">
    <label :for="id" class="block w-full pb-1">{{ translation }}</label>
    <Field
      as=""
      name="name"
      v-model="local"
      :label="translation.toLowerCase()"
      rules="required"
      v-slot="{ field, meta, errorMessage }"
      type="text"
      :id="id"
      :placeholder="translation.toLowerCase()"
      autocomplete="off"
    >
      <input
        v-bind="field"
        type="text"
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
    name: {
      type: String,
      default: '',
    },
    id: {
      type: String,
      default: 'name',
    },
    translation: {
      type: String,
      required: true,
    },
  },
  components: {
    Field,
  },
  emits: ['update:updateName'],
  setup(props, { emit }) {
    const local = computed({
      get: () => {
        emit('update:updateName', props.name)
        return props.name
      },
      set: (value: string) => emit('update:updateName', value),
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
