<template>
  <div class="mt-4">
    <button type="submit" :class="classes">
      {{ text }}
    </button>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'

export default defineComponent({
  props: {
    text: {
      type: String,
      required: true,
    },
    classes: {
      type: String,
      default: 'w-full p-3 bg-gray-900 text-white hover:text-green-500',
    },
  },
})
</script>
