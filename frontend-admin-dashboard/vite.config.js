import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { readFileSync } from 'fs'
import { resolve } from 'path'

// قراءة الإصدار من package.json
const packageJson = JSON.parse(readFileSync(resolve(__dirname, 'package.json'), 'utf-8'))
const appVersion = packageJson.version

export default defineConfig({
  plugins: [vue()],
  base: './', // استخدام مسارات نسبية لـ Electron
  server: { 
    port: 3001,
    cors: true
  },
  define: {
    __APP_VERSION__: JSON.stringify(appVersion)
  },
  optimizeDeps: {
    include: ['vue', 'vue-router', 'leaflet']
  },
  assetsInclude: ['**/*.json'],
  build: {
    cssCodeSplit: true,
    rollupOptions: {
      output: {
        manualChunks: {
          'vue-vendor': ['vue', 'vue-router'],
          'leaflet-vendor': ['leaflet']
        }
      }
    }
  }
})
