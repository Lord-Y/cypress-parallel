<template>
  <div class="mt-4">
    <label class="block w-full pb-1">{{
      $t('projects.schedulingEnabled')
    }}</label>
    <label :for="id" class="flex items-center cursor-pointer">
      <div class="relative">
        <input
          type="checkbox"
          :id="id"
          class="hidden"
          v-model="local"
          @click="update($event)"
        />
        <div
          id="englobe"
          class="w-12 h-8 flex items-center bg-gray-300 rounded-2xl p-1 duration-200 ease-in-out"
          :class="{ 'bg-emerald-500': local }"
        >
          <div
            id="inner"
            class="bg-white w-6 h-6 rounded-2xl shadow-md transform duration-200 ease-in-out"
            :class="{
              'translate-x-4': local,
            }"
          ></div>
        </div>
      </div>
    </label>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref } from 'vue'

export default defineComponent({
  props: {
    id: {
      type: String,
      default: 'schedulingEnabled',
    },
    name: {
      type: String,
      default: 'schedulingEnabled',
    },
    schedulingEnabled: {
      type: Boolean,
      default: false,
    },
  },
  emits: ['update:updateSchedulingEnabled'],
  setup(props, { emit }) {
    const local = ref(props.schedulingEnabled)
    emit('update:updateSchedulingEnabled', local)

    function update(event: any): void {
      let englobe: Element | null, inner: Element | null
      englobe = document.querySelector('#englobe')
      inner = document.querySelector('#inner')
      if (englobe !== null && inner !== null) {
        if (event.target.checked) {
          englobe.classList.add('bg-emerald-500')
          inner.classList.add('translate-x-4')
        } else {
          englobe.classList.remove('bg-emerald-500')
          inner.classList.remove('translate-x-4')
        }
        emit('update:updateSchedulingEnabled', event.target.checked)
      }
    }

    return {
      local,
      update,
    }
  },
})
</script>
