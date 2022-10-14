<template>
  <div class="mt-4">
    <label :for="id" class="block w-full pb-1">{{
      $t('projects.browsers', 2)
    }}</label>
    <Field
      as=""
      :name="name"
      v-model="local"
      :label="$t('projects.browsers', 2).toLowerCase()"
      rules="required"
      v-slot="{ field, errorMessage, meta }"
    >
      <select
        v-bind="field"
        :id="id"
        class="block w-full border-gray-300 focus:outline-none focus:border-emerald-500 focus:ring-emerald-500"
        :class="getValidationClass(meta)"
      >
        <option value="" disabled>
          {{ $t('select.selectYourChoice') }}
        </option>
        <option
          v-for="(browser_name, index) in browsers"
          :value="browser_name"
          :key="index"
        >
          {{ browser_name }}
        </option>
      </select>
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
      default: 'browser',
    },
    name: {
      type: String,
      default: 'browser',
    },
    browser: {
      type: String,
      default: 'chrome',
    },
  },
  components: {
    Field,
  },
  emits: ['update:updateBrowser'],
  setup(props, { emit }) {
    const browsers = ref(['chrome', 'firefox'])
    const local = ref(props.browser)
    emit('update:updateBrowser', local)

    function getValidationClass(meta: any): string {
      if (meta.valid && meta.validated && meta.dirty) {
        return 'outline-none border-emerald-500'
      }
      if (!meta.valid && meta.validated && !meta.dirty) {
        return 'outline-none border-red-500'
      }
      return ''
    }

    return {
      browsers,
      local,
      getValidationClass,
    }
  },
})
</script>
