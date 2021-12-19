import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import { createI18n } from 'vue-i18n'
import messages from '@intlify/vite-plugin-vue-i18n/messages'
import './tailwind.css'
import { createHead } from '@vueuse/head'
import { configure } from 'vee-validate'
import { localize } from '@vee-validate/i18n'
import en from '@vee-validate/i18n/dist/locale/en.json'
import fr from '@vee-validate/i18n/dist/locale/fr.json'
const i18n = createI18n({
  legacy: false,
  locale: 'en-US',
  fallbackLocale: 'en-US',
  globalInjection: true,
  messages,
})

configure({
  generateMessage: localize({
    en,
    fr,
  }),
})

const app = createApp(App)
const head = createHead()
app.use(router)
app.use(i18n)
app.use(head)
app.mount('#app')
