import { createRouter, createWebHistory, createMemoryHistory } from 'vue-router'
import LoginView from '../views/LoginView.vue'
import DashboardView from '../views/DashboardView.vue'
import EmployeesView from '../views/EmployeesView.vue'
import TasksView from '../views/TasksView.vue'
import WorksitesView from '../views/WorksitesView.vue'
import ClientsView from '../views/ClientsView.vue'
import ReportsView from '../views/ReportsView.vue'
import SettingsView from '../views/SettingsView.vue'
import ServiceRequestsView from '../views/ServiceRequestsView.vue'

const routes = [
  { path: '/', redirect: '/dashboard' },
  { path: '/login', component: LoginView, meta: { public: true } },
  { path: '/dashboard', component: DashboardView, meta: { requiresAuth: true } },
  { path: '/service-requests', component: ServiceRequestsView, meta: { requiresAuth: true } },
  { path: '/employees', component: EmployeesView, meta: { requiresAuth: true } },
  { path: '/tasks', component: TasksView, meta: { requiresAuth: true } },
  { path: '/worksites', component: WorksitesView, meta: { requiresAuth: true } },
  { path: '/reports', component: ReportsView, meta: { requiresAuth: true } },
  { path: '/settings', component: SettingsView, meta: { requiresAuth: true } },
]

// استخدام MemoryHistory لـ Electron و WebHistory للويب
const isElectron = typeof window !== 'undefined' && (window.isElectron || (window.process && window.process.type === 'renderer'))
const router = createRouter({ 
  history: isElectron ? createMemoryHistory() : createWebHistory(), 
  routes 
})

// حماية المسارات
router.beforeEach((to, from, next) => {
  const isAuthed = !!localStorage.getItem('worktrack_admin_token')
  
  // إذا كان المسار يتطلب مصادقة والمستخدم غير مسجل
  if (to.meta.requiresAuth && !isAuthed) {
    next('/login')
  } 
  // إذا كان المستخدم مسجل ويحاول الذهاب إلى login
  else if (to.path === '/login' && isAuthed) {
    next('/dashboard')
  }
  // السماح بالوصول
  else {
    next()
  }
})

export default router
