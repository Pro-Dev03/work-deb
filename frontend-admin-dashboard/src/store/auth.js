import { reactive } from 'vue'
import { currentUser } from '../services/auth'

export const authStore = reactive({
  user: currentUser(),
  
  setUser(user) {
    this.user = user
  },
  
  clear() {
    this.user = null
  }
})
