import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()],
  server: { 
    port: 3001,
    cors: true
  },
  publicDir: 'public',
  optimizeDeps: {
    include: ['vue', 'vue-router', 'leaflet']
  },
  assetsInclude: ['**/*.json', '**/*.jpg', '**/*.png', '**/*.svg'],
  build: {
    cssCodeSplit: true,
    rollupOptions: {
      output: {
        manualChunks: {
          'vue-vendor': ['vue', 'vue-router'],
          'leaflet-vendor': ['leaflet']
        },
        // إضافة timestamp تلقائي لضمان cache busting
        entryFileNames: 'assets/[name]-[hash].js',
        chunkFileNames: 'assets/[name]-[hash].js',
        assetFileNames: 'assets/[name]-[hash][extname]'
      }
    }
  }
})
