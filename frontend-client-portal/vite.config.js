import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { readFileSync } from 'fs'
import { resolve } from 'path'

// قراءة الإصدار من package.json
const packageJson = JSON.parse(readFileSync(resolve(__dirname, 'package.json'), 'utf-8'))
const appVersion = packageJson.version

// Vite plugin مخصص لمعالجة service-worker.js
function serviceWorkerPlugin() {
  return {
    name: 'service-worker-plugin',
    generateBundle() {
      const serviceWorkerPath = resolve(__dirname, 'public/service-worker.js')
      let serviceWorkerContent = readFileSync(serviceWorkerPath, 'utf-8')
      serviceWorkerContent = serviceWorkerContent.replace('__APP_VERSION__', appVersion)
      
      this.emitFile({
        type: 'asset',
        fileName: 'service-worker.js',
        source: serviceWorkerContent
      })
    }
  }
}

export default defineConfig({
  plugins: [vue(), serviceWorkerPlugin()],
  server: { port: 3002 },
  define: {
    __APP_VERSION__: JSON.stringify(appVersion)
  },
  build: {
    cssCodeSplit: true,
    rollupOptions: {
      output: {
        manualChunks: {
          'vue-vendor': ['vue', 'vue-router']
        }
      }
    }
  }
})
