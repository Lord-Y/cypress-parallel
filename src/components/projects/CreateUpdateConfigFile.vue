<template>
  <div class="mt-4">
    <label :for="id" class="block w-full pb-1">{{
      $t('projects.configFile')
    }}</label>
    <input
      v-model.lazy.trim="local"
      type="text"
      :id="id"
      :name="name"
      :placeholder="$t('projects.configFile').toLowerCase()"
      autocomplete="off"
      class="
        block
        w-full
        border-gray-300
        focus:outline-none
        focus:border-green-500
        focus:ring-green-500
      "
    />
  </div>
</template>

<script lang="ts">
import { defineComponent, ref } from 'vue'

export default defineComponent({
  props: {
    id: {
      type: String,
      default: 'configFile',
    },
    name: {
      type: String,
      default: 'configFile',
    },
    configFile: {
      type: String,
      default: 'cypress.json',
    },
  },
  emits: ['update:updateConfigFile'],
  setup(props, { emit }) {
    const local = ref(props.configFile)
    emit('update:updateConfigFile', local)

    return {
      local,
    }
  },
})
</script>
