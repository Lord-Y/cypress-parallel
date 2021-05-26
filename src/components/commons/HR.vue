<template>
  <hr :class="classes" />
</template>

<script lang="ts">
import { defineComponent } from 'vue'

export default defineComponent({
  props: {
    classes: {
      type: String,
      default: 'my-5 border-2 border-gray-300',
    },
  },
})
</script>
