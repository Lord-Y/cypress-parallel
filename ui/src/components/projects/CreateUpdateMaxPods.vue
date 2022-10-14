<template>
  <div class="mt-4">
    <label :for="id" class="block w-full pb-1">{{
      $t('projects.maxPods')
    }}</label>
    <Field
      as=""
      :name="name"
      v-model.number="local"
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
        class="block w-full border-gray-300 focus:outline-none focus:border-emerald-500 focus:ring-emerald-500"
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
    const local = ref(props.maxPods)
    emit('update:updateMaxPods', local)

    function getValidationClass(meta: any): string {
      if (meta.valid && meta.validated && meta.dirty) {
        return 'outline-none border-emerald-500'
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
