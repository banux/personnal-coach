<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Navigation (hide on login page) -->
    <nav v-if="!isLoginPage" class="bg-white shadow-sm border-b border-gray-100">
      <div class="max-w-6xl mx-auto px-4 py-3 flex items-center justify-between">
        <router-link to="/" class="flex items-center space-x-2">
          <span class="text-2xl">💪</span>
          <span class="text-xl font-bold text-blue-700">Coach Personnel IA</span>
        </router-link>
        <div class="flex items-center gap-4">
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
    <main :class="isLoginPage ? '' : 'max-w-6xl mx-auto px-4 py-8'">
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

const isLoginPage = computed(() => route.path === '/login')

async function handleLogout() {
  await authStore.logout()
  router.push('/login')
}
</script>
