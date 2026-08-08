import { createRouter, createWebHistory, createMemoryHistory } from 'vue-router'
import LoginView from '../views/LoginView.vue'
import AttendanceView from '../views/AttendanceView.vue'
import ProfileView from '../views/ProfileView.vue'
import NotesView from '../views/NotesView.vue'
import HistoryView from '../views/HistoryView.vue'

// صفحة 404 بسيطة
const NotFoundView = {
  template: `
    <div class="not-found">
      <h1>404</h1>
      <p>الصفحة غير موجودة</p>
      <router-link to="/attendance">العودة للحضور</router-link>
    </div>
  `,
  style: `
    .not-found {
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      height: 100vh;
      text-align: center;
      font-family: system-ui, sans-serif;
    }
    .not-found h1 {
      font-size: 72px;
      margin: 0;
      color: #666;
    }
    .not-found p {
      font-size: 18px;
      color: #999;
      margin: 20px 0;
    }
    .not-found a {
      color: #667eea;
      text-decoration: none;
      font-size: 16px;
    }
    .not-found a:hover {
      text-decoration: underline;
    }
  `
}

const routes = [
  { path: '/', redirect: '/attendance' },
  { path: '/login', component: LoginView, meta: { public: true } },
  { path: '/attendance', component: AttendanceView },
  { path: '/profile', component: ProfileView },
  { path: '/notes', component: NotesView },
  { path: '/history', component: HistoryView },
  { path: '/:pathMatch(.*)*', component: NotFoundView }, // صفحة 404
]

// استخدام MemoryHistory لـ Electron و WebHistory للويب
const isElectron = typeof window !== 'undefined' && (window.isElectron || (window.process && window.process.type === 'renderer'))
const router = createRouter({ 
  history: isElectron ? createMemoryHistory() : createWebHistory(), 
  routes 
})

router.beforeEach((to) => {
  if (!to.meta.public && !localStorage.getItem('worktrack_token')) return '/login'
  return true
})
export default router
