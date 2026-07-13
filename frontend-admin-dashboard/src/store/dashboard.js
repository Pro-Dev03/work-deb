import { reactive } from 'vue'
import api from '../services/api'

// بيانات تجريبية للعرض إلى حين توفر استجابة فعلية من /reports/daily-summary
const mockSummary = {
  total_employees: 18,
  completed: 27,
  in_progress: 6,
  pending: 9,
  late: 2,
}

export const dashboardStore = reactive({
  summary: mockSummary,
  loading: false,
  async fetchSummary() {
    this.loading = true
    try {
      const { data } = await api.get('/reports/daily-summary')
      this.summary = data
    } catch (e) {
      this.summary = mockSummary
    } finally {
      this.loading = false
    }
  },
})
