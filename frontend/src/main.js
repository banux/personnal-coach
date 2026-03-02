import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { createRouter, createWebHistory } from 'vue-router'
import './style.css'
import App from './App.vue'
import LoginView from './views/LoginView.vue'
import HomeView from './views/HomeView.vue'
import ProgramView from './views/ProgramView.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/login', component: LoginView, meta: { public: true } },
    { path: '/', component: HomeView },
    { path: '/program/:id', component: ProgramView, props: true },
  ]
})

const pinia = createPinia()
const app = createApp(App)
app.use(pinia)
app.use(router)

// Auth guard: check session before mounting
import { useAuthStore } from './stores/auth.js'

router.beforeEach(async (to) => {
  const authStore = useAuthStore()

  if (to.meta.public) return true

  await authStore.checkStatus()
  if (!authStore.authenticated) {
    return '/login'
  }
  return true
})

app.mount('#app')
