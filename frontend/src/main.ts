import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { router } from './router'
import { canDirective } from './directives/can'
import { useAuthStore, setupAuthCrossTabSync } from './stores/authStore'
import App from './App.vue'
import './style.css'

const app = createApp(App)
const pinia = createPinia()

app.use(pinia)
app.use(router)
app.directive('can', canDirective)

// Initialize cross-tab authentication synchronization (BroadcastChannel + storage events)
const authStore = useAuthStore()
setupAuthCrossTabSync(authStore, router)

app.mount('#app')
