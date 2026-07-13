import { reactive } from 'vue'
import api from '../services/api'

// بيانات تجريبية تُستخدم للعرض إلى حين توفر استجابة فعلية من الـ API
const mockTasks = [
  { id: '1', title: 'صيانة مكيّفات — برج الأمل', worksite: 'برج الأمل، عمّان', time: '09:00 - 11:00', status: 'in_progress' },
  { id: '2', title: 'فحص دوري — فيلا الشميساني', worksite: 'الشميساني', time: '12:30 - 13:30', status: 'pending' },
  { id: '3', title: 'تركيب كاميرات — معرض السيارات', worksite: 'طريق المطار', time: 'أمس، 16:00', status: 'completed' },
  { id: '4', title: 'إصلاح تسريب مياه — مجمع الزهور', worksite: 'الجبيهة', time: 'أمس، 10:00', status: 'late' },
]

export const tasksStore = reactive({
  items: [],
  loading: false,
  async fetchMine() {
    this.loading = true
    try {
      const { data } = await api.get('/tasks/mine')
      this.items = data
    } catch (e) {
      this.items = mockTasks
    } finally {
      this.loading = false
    }
  },
  find(id) {
    return this.items.find((t) => String(t.id) === String(id))
  },
})
