import { createRouter, createWebHistory } from 'vue-router'
import ServiceReportView from '../views/ServiceReportView.vue'
import ServiceHistoryView from '../views/ServiceHistoryView.vue'
import RatingView from '../views/RatingView.vue'
import NewRequestView from '../views/NewRequestView.vue'

const routes = [
  { path: '/', component: ServiceHistoryView },
  { path: '/new', component: NewRequestView },
  { path: '/report/:id', component: ServiceReportView },
  { path: '/rate/:id', component: RatingView },
]

const router = createRouter({ history: createWebHistory(), routes })
export default router
