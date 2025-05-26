import { createApp } from 'vue'
import App from './App.vue'
import 'bootstrap/dist/css/bootstrap.min.css'
import Notifications from '@kyvg/vue3-notification'
import configService from './services/config'

// Load configuration before mounting the app
configService.loadConfig().then(() => {
  const app = createApp(App)
  app.use(Notifications)
  app.mount('#app')
}).catch(error => {
  console.error('Failed to load configuration:', error)
  // Still mount the app with defaults
  const app = createApp(App)
  app.use(Notifications)
  app.mount('#app')
})