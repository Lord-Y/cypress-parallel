import { reactive, toRefs } from 'vue'
import { useI18n } from 'vue-i18n'
import { useHead } from '@vueuse/head'

// eslint-disable-next-line @typescript-eslint/explicit-module-boundary-types
export default function () {
  const state = reactive({
    meta: {
      title: '',
      description: '',
    },
  })

  const { t } = useI18n({
    useScope: 'global',
  })

  state.meta.title =
    t('brand.name').toUpperCase() +
    ' | ' +
    t('alert.http.maintenance.maintenance')
  state.meta.description = t('alert.http.maintenance.maintenance')

  useHead({
    title: state.meta.title,
    meta: [
      {
        name: 'description',
        content: state.meta.description,
      },
      {
        property: 'og:title',
        content: t('brand.name').toUpperCase() + ' | ' + state.meta.title,
      },
      {
        property: 'og:description',
        content: t('metaTags.estateAnnounce') + ' - ' + state.meta.description,
      },
    ],
  })

  return {
    ...toRefs(state),
  }
}
