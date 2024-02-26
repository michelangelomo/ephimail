import { createApp } from 'vue'
import App from './App.vue'
import 'bootstrap/dist/css/bootstrap.min.css';
import { library } from '@fortawesome/fontawesome-svg-core';
import { faChevronLeft, faChevronRight, faDownload } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome';
import Notifications from '@kyvg/vue3-notification'

library.add(faChevronLeft, faChevronRight, faDownload)

createApp(App)
.component('font-awesome-icon', FontAwesomeIcon)
.use(Notifications)
.mount('#app')
