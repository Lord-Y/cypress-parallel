<template>
  <div class="mt-4">
    <label :for="id" class="block w-full pb-1">{{
      $t('projects.password')
    }}</label>
    <input
      v-model.lazy.trim="local"
      type="text"
      :id="id"
      :name="name"
      :placeholder="$t('projects.password').toLowerCase()"
      autocomplete="off"
      class="block w-full border-gray-300 focus:outline-none focus:border-green-500 focus:ring-green-500"
    />
  </div>
</template>

<script lang="ts">
import { defineComponent, ref } from 'vue'

export default defineComponent({
  props: {
    id: {
      type: String,
      default: 'password',
    },
    name: {
      type: String,
      default: 'password',
    },
    password: {
      type: String,
      default: '',
    },
  },
  emits: ['update:updatePassword'],
  setup(props, { emit }) {
    const local = ref(props.password)
    emit('update:updatePassword', local)

    return {
      local,
    }
  },
})
</script>
