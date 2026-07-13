import { createRouter, createWebHistory } from 'vue-router'
import LoginView from '../views/LoginView.vue'
import AttendanceView from '../views/AttendanceView.vue'
import ProfileView from '../views/ProfileView.vue'

const routes = [
  { path: '/', redirect: '/attendance' },
  { path: '/login', component: LoginView, meta: { public: true } },
  { path: '/attendance', component: AttendanceView },
  { path: '/profile', component: ProfileView },
]

const router = createRouter({ history: createWebHistory(), routes })
router.beforeEach((to) => {
  if (!to.meta.public && !localStorage.getItem('worktrack_token')) return '/login'
  return true
})
export default router
