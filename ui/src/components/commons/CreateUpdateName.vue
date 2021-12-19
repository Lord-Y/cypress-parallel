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
    >
      <input
        v-bind="field"
        type="text"
        :id="id"
        :placeholder="translation.toLowerCase()"
        autocomplete="off"
        class="block w-full border-gray-300 focus:outline-none focus:border-green-500 focus:ring-green-500"
        :class="getValidationClass(meta)"
      />
      <span v-if="errorMessage" class="text-red-500">{{ errorMessage }}</span>
    </Field>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref } from 'vue'
import { Field } from 'vee-validate'
import { defineRule } from 'vee-validate'
import { required } from '@vee-validate/rules'
defineRule('required', required)

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
    const local = ref(props.name)
    emit('update:updateName', local)

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
