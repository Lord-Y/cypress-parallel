<template>
  <div class="flex h-screen">
    <div class="m-auto text-center">
      <p class="text-gray-600 font-black text-7xl">
        {{ $t('alert.http.maintenance.msg') }}
      </p>
      <p class="text-gray-600 font-black text-7xl">
        {{ $t('alert.http.maintenance.back') }}
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import maintenance from '@/compositions/pageStatus/maintenance'
// eslint-disable-next-line @typescript-eslint/no-unused-vars
const { meta } = maintenance()
</script>
