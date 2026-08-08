import { createRouter, createWebHistory, createMemoryHistory } from 'vue-router'
import LoginView from '../views/LoginView.vue'
import DashboardView from '../views/DashboardView.vue'
import EmployeesView from '../views/EmployeesView.vue'
import TasksView from '../views/TasksView.vue'
import WorksitesView from '../views/WorksitesView.vue'
import ClientsView from '../views/ClientsView.vue'
import ReportsView from '../views/ReportsView.vue'
import SettingsView from '../views/SettingsView.vue'
import AttendanceManagementView from '../views/AttendanceManagementView.vue'
import NotesView from '../views/NotesView.vue'

// صفحة 404 بسيطة
const NotFoundView = {
  template: `
    <div class="not-found">
      <h1>404</h1>
      <p>الصفحة غير موجودة</p>
      <router-link to="/dashboard">العودة للوحة التحكم</router-link>
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
  { path: '/', redirect: '/dashboard' },
  { path: '/login', component: LoginView, meta: { public: true } },
  { path: '/dashboard', component: DashboardView, meta: { requiresAuth: true } },
  { path: '/employees', component: EmployeesView, meta: { requiresAuth: true } },
  { path: '/tasks', component: TasksView, meta: { requiresAuth: true } },
  { path: '/worksites', component: WorksitesView, meta: { requiresAuth: true } },
  { path: '/reports', component: ReportsView, meta: { requiresAuth: true } },
  { path: '/settings', component: SettingsView, meta: { requiresAuth: true } },
  { path: '/attendance-management', component: AttendanceManagementView, meta: { requiresAuth: true } },
  { path: '/notes', component: NotesView, meta: { requiresAuth: true } },
  { path: '/:pathMatch(.*)*', component: NotFoundView }, // صفحة 404 بدلاً من redirect
]

// استخدام MemoryHistory لـ Electron و WebHistory للويب
const isElectron = typeof window !== 'undefined' && (window.isElectron || (window.process && window.process.type === 'renderer'))
const router = createRouter({ 
  history: isElectron ? createMemoryHistory() : createWebHistory(), 
  routes 
})

// حماية المسارات
router.beforeEach((to, from, next) => {
  // التحقق من المستخدم بدلاً من التوكن (نستخدم httpOnly cookies)
  const userStr = localStorage.getItem('worktrack_admin_user')
  let isAuthed = false
  
  try {
    // التحقق من صحة البيانات في localStorage
    if (userStr) {
      const user = JSON.parse(userStr)
      isAuthed = !!user && (user.id || user.email)
    }
  } catch (e) {
    console.error('❌ خطأ في قراءة بيانات المستخدم:', e)
    localStorage.removeItem('worktrack_admin_user')
    isAuthed = false
  }
  
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
