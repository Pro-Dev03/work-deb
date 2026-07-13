import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import './styles/tokens.css'
import i18n from './services/i18n'

import 'leaflet/dist/leaflet.css'
import L from 'leaflet'

delete L.Icon.Default.prototype._getIconUrl
L.Icon.Default.mergeOptions({
  iconRetinaUrl: new URL('leaflet/dist/images/marker-icon-2x.png', import.meta.url).href,
  iconUrl: new URL('leaflet/dist/images/marker-icon.png', import.meta.url).href,
  shadowUrl: new URL('leaflet/dist/images/marker-shadow.png', import.meta.url).href,
})

const app = createApp(App)
app.use(i18n)
app.use(router)
app.mount('#app')
