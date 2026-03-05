<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Navigation (hide on login/profile pages) -->
    <nav v-if="showNav" class="bg-white shadow-sm border-b border-gray-100">
      <div class="max-w-6xl mx-auto px-4 py-3 flex items-center justify-between">
        <router-link to="/" class="flex items-center space-x-2">
          <span class="text-2xl">💪</span>
          <span class="text-xl font-bold text-blue-700">Coach Personnel IA</span>
        </router-link>
        <div class="flex items-center gap-4">
          <!-- Current profile badge -->
          <div v-if="authStore.profile" class="flex items-center gap-2">
            <span class="text-sm text-gray-500">👤 {{ authStore.profile.name }}</span>
            <router-link
              to="/profiles"
              class="text-xs text-gray-400 hover:text-blue-600 transition-colors border border-gray-200 rounded px-2 py-0.5 hover:border-blue-300"
            >
              Changer
            </router-link>
          </div>
          <router-link
            to="/programs"
            class="text-sm text-gray-500 hover:text-blue-600 transition-colors"
            :class="route.path === '/programs' ? 'text-blue-600 font-medium' : ''"
          >
            Historique
          </router-link>
          <router-link
            to="/"
            class="text-sm text-gray-500 hover:text-blue-600 transition-colors"
            :class="route.path === '/' ? 'text-blue-600 font-medium' : ''"
          >
            Nouveau
          </router-link>
          <button
            @click="handleLogout"
            class="text-sm text-gray-400 hover:text-red-500 transition-colors"
          >
            Déconnexion
          </button>
        </div>
      </div>
    </nav>

    <!-- Main content -->
    <main :class="showNav ? 'max-w-6xl mx-auto px-4 py-8' : ''">
      <router-view />
    </main>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from './stores/auth.js'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const showNav = computed(() => route.path !== '/login' && route.path !== '/profiles')

async function handleLogout() {
  await authStore.logout()
  router.push('/login')
}
</script>
