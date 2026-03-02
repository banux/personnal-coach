<template>
  <div class="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 flex items-center justify-center p-4">
    <div class="bg-white rounded-2xl shadow-lg w-full max-w-sm p-8">
      <!-- Logo -->
      <div class="text-center mb-8">
        <div class="text-5xl mb-3">💪</div>
        <h1 class="text-2xl font-bold text-gray-900">Coach Personnel IA</h1>
        <p class="text-gray-500 text-sm mt-1">Connectez-vous pour accéder à votre coach</p>
      </div>

      <!-- Error -->
      <div v-if="authStore.error" class="mb-4 bg-red-50 border border-red-200 text-red-700 rounded-lg p-3 text-sm">
        {{ authStore.error }}
      </div>

      <!-- Form -->
      <form @submit.prevent="handleLogin" class="space-y-4">
        <div>
          <label class="label">Mot de passe</label>
          <input
            v-model="password"
            type="password"
            class="input-field"
            placeholder="••••••••"
            autocomplete="current-password"
            required
            autofocus
          />
        </div>

        <button
          type="submit"
          class="btn-primary w-full py-3"
          :disabled="authStore.loading || !password"
        >
          <span v-if="authStore.loading">⏳ Connexion...</span>
          <span v-else>Se connecter</span>
        </button>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth.js'

const router = useRouter()
const authStore = useAuthStore()
const password = ref('')

async function handleLogin() {
  try {
    await authStore.login(password.value)
    router.push('/')
  } catch {
    password.value = ''
  }
}
</script>
