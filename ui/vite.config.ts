import path from 'path'
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueI18n from '@intlify/vite-plugin-vue-i18n'

// https://vitejs.dev/config/
export default defineConfig({
  base: '/ui/',
  resolve: {
    alias: {
      'vue-i18n': 'vue-i18n/dist/vue-i18n.runtime.esm-bundler.js',
      '@': path.resolve(__dirname, 'src/'),
      '@router': path.resolve(__dirname, 'src/router'),
      '@components': path.resolve(__dirname, 'src/components'),
      '@views': path.resolve(__dirname, 'src/views'),
      '@api': path.resolve(__dirname, 'src/api'),
    },
  },
  plugins: [
    vue(),
    vueI18n({
      include: path.resolve(__dirname, './src/i18n/locales/**'),
    }),
  ],
  server: {
    strictPort: true,
    proxy: {
      '/api': {
        target: process.env.API_GATEWAY || 'http://localhost:8080',
      },
    }
  }
})
