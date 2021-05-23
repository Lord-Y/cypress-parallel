import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import { createI18n } from 'vue-i18n'
import messages from '@intlify/vite-plugin-vue-i18n/messages'
import './tailwind.css'
import { defineRule, configure } from 'vee-validate'
import { localize } from '@vee-validate/i18n'
import AllRules from '@vee-validate/rules'
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

Object.keys(AllRules).forEach((rule) => {
  defineRule(rule, AllRules[rule])
})

const app = createApp(App)
app.use(router)
app.use(i18n)
app.mount('#app')
