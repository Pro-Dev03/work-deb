import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()],
  server: { 
    port: 3001,
    cors: true
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
