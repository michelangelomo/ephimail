import { createApp } from 'vue'
import App from './App.vue'
import 'bootstrap/dist/css/bootstrap.min.css'
import { library } from '@fortawesome/fontawesome-svg-core'
import { faChevronLeft, faChevronRight, faDownload, faLock, faCopy, faExclamationTriangle } from "@fortawesome/free-solid-svg-icons"
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'
import Notifications from '@kyvg/vue3-notification'

// Add all the icons we need
library.add(
  faChevronLeft, 
  faChevronRight, 
  faDownload, 
  faLock, 
  faCopy, 
  faExclamationTriangle
)

const app = createApp(App)
app.component('font-awesome-icon', FontAwesomeIcon)
app.use(Notifications)
app.mount('#app')