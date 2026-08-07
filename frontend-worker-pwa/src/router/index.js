import { createRouter, createWebHistory, createMemoryHistory } from 'vue-router'
import LoginView from '../views/LoginView.vue'
import AttendanceView from '../views/AttendanceView.vue'
import ProfileView from '../views/ProfileView.vue'
import NotesView from '../views/NotesView.vue'
import HistoryView from '../views/HistoryView.vue'

const routes = [
  { path: '/', redirect: '/attendance' },
  { path: '/login', component: LoginView, meta: { public: true } },
  { path: '/attendance', component: AttendanceView },
  { path: '/profile', component: ProfileView },
  { path: '/notes', component: NotesView },
  { path: '/history', component: HistoryView },
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
