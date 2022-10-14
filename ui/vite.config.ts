import path from 'path'
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueI18n from '@intlify/vite-plugin-vue-i18n'

let api_url: string
if (process.env.CYPRESS_PARALLEL_PORT) {
  api_url = `http://localhost:${process.env.CYPRESS_PARALLEL_PORT}`
} else {
  api_url = 'http://localhost:8080'
}

// https://vitejs.dev/config/
export default defineConfig({
  base: '/ui/',
  resolve: {
    alias: {
      'vue-i18n': 'vue-i18n/dist/vue-i18n.runtime.esm-bundler.js',
      '@': path.resolve(__dirname, 'src/'),
    },
  },
  plugins: [
    vue(),
    vueI18n({
      include: path.resolve(__dirname, './src/i18n/locales/**'),
    }),
  ],
  server: {
    port: 3000,
    strictPort: true,
    proxy: {
      '/api': {
        target: api_url,
      },
    },
  },
})
