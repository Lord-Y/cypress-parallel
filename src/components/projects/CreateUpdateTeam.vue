<template>
  <div class="mt-4">
    <label :for="id" class="block w-full pb-1">{{
      $t('teams.teams', 2)
    }}</label>
    <Field
      as=""
      :name="name"
      v-model="local"
      :label="$t('teams.teams', 2).toLowerCase()"
      rules="required"
      v-slot="{ field, errorMessage, meta }"
    >
      <select
        v-bind="field"
        :id="id"
        class="
          block
          w-full
          border-gray-300
          focus:outline-none
          focus:border-green-500
          focus:ring-green-500
        "
        :class="getValidationClass(meta)"
      >
        <option value="" disabled>
          {{ $t('select.selectYourChoice') }}
        </option>
        <option
          v-for="(team, index) in teams"
          :value="team.team_id"
          :key="index"
        >
          {{ team.team_name }}
        </option>
      </select>
      <span v-if="errorMessage" class="text-red-500">{{ errorMessage }}</span>
    </Field>
  </div>
</template>

<script lang="ts">
import { defineComponent, PropType, computed } from 'vue'
import { Field } from 'vee-validate'
import { Teams } from '@api/teamsService'
import { defineRule } from 'vee-validate';
import { required } from '@vee-validate/rules';
defineRule('required', required)

export default defineComponent({
  props: {
    id: {
      type: String,
      default: 'team',
    },
    name: {
      type: String,
      default: 'team',
    },
    teams: {
      type: Array as PropType<Teams[]>,
      required: true,
    },
    teamId: {
      type: String,
      default: '',
    },
  },
  components: {
    Field,
  },
  emits: ['update:updateTeam'],
  setup(props, { emit }) {
    const local = computed({
      get: () => {
        emit('update:updateTeam', props.teamId)
        return props.teamId
      },
      set: (value: string) => emit('update:updateTeam', value),
    })
    function getValidationClass(meta: any): string {
      if (meta.valid && meta.validated && meta.dirty) {
        return 'outline-none border-green-500'
      }
      if (!meta.valid && meta.validated && !meta.dirty) {
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
