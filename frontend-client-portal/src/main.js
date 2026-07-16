import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import './styles/tokens.css'

// التأكد من تحميل الصفحة بشكل كامل قبل تعريف التطبيق
document.addEventListener('DOMContentLoaded', () => {
  createApp(App).use(router).mount('#app')
})
