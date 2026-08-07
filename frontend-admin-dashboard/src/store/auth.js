import { reactive } from 'vue'
import { currentUser } from '../services/auth'

export const authStore = reactive({
  user: currentUser(),
  
  setUser(user) {
    console.log('📝 authStore.setUser called with:', user)
    this.user = user
    console.log('✅ authStore.user updated to:', this.user)
  },
  
  clear() {
    console.log('🗑️ authStore.clear called')
    this.user = null
    console.log('✅ authStore.user cleared')
  }
})
