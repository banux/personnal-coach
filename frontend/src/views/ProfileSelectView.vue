<template>
  <div class="min-h-screen bg-gray-50 flex items-center justify-center px-4">
    <div class="w-full max-w-md">
      <div class="text-center mb-8">
        <span class="text-5xl">👤</span>
        <h1 class="text-2xl font-bold text-gray-900 mt-3">Choisissez votre profil</h1>
        <p class="text-gray-500 mt-1">Tous les utilisateurs partagent le même mot de passe</p>
      </div>

      <!-- Error -->
      <div v-if="error" class="mb-4 bg-red-50 border border-red-200 text-red-700 rounded-lg p-3 text-sm">
        {{ error }}
      </div>

      <!-- Existing profiles -->
      <div v-if="profiles.length > 0" class="card mb-6">
        <h2 class="text-sm font-semibold text-gray-500 uppercase tracking-wide mb-3">Profils existants</h2>
        <div class="space-y-2">
          <button
            v-for="p in profiles"
            :key="p.id"
            @click="handleSelect(p.id)"
            class="w-full flex items-center gap-3 px-4 py-3 rounded-lg border border-gray-200 hover:border-blue-400 hover:bg-blue-50 transition-colors text-left"
            :class="{ 'opacity-50 pointer-events-none': selecting }"
          >
            <span class="text-xl">👤</span>
            <div>
              <div class="font-medium text-gray-900">{{ p.name }}</div>
              <div class="text-xs text-gray-400">Créé le {{ formatDate(p.created_at) }}</div>
            </div>
          </button>
        </div>
      </div>

      <!-- Create new profile -->
      <div class="card">
        <h2 class="text-sm font-semibold text-gray-500 uppercase tracking-wide mb-3">Nouveau profil</h2>
        <form @submit.prevent="handleCreate" class="flex gap-2">
          <input
            v-model="newName"
            type="text"
            class="input-field flex-1"
            placeholder="Ex: Marie, Pierre..."
            maxlength="50"
            required
          />
          <button type="submit" class="btn-primary px-4" :disabled="creating || !newName.trim()">
            <span v-if="creating">...</span>
            <span v-else>Créer</span>
          </button>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'
import { useAuthStore } from '../stores/auth.js'

const router = useRouter()
const authStore = useAuthStore()

const profiles = ref([])
const newName = ref('')
const error = ref(null)
const creating = ref(false)
const selecting = ref(false)

const API_BASE = import.meta.env.VITE_API_URL || ''

onMounted(async () => {
  try {
    const res = await axios.get(`${API_BASE}/api/profiles`, { withCredentials: true })
    profiles.value = res.data.profiles || []
  } catch (err) {
    error.value = 'Impossible de charger les profils'
  }
})

async function handleSelect(id) {
  selecting.value = true
  error.value = null
  try {
    await authStore.selectProfile(id)
    router.push('/')
  } catch (err) {
    error.value = err.response?.data?.error || 'Erreur lors de la sélection'
  } finally {
    selecting.value = false
  }
}

async function handleCreate() {
  const name = newName.value.trim()
  if (!name) return
  creating.value = true
  error.value = null
  try {
    const res = await axios.post(`${API_BASE}/api/profiles`, { name }, { withCredentials: true })
    const created = res.data
    profiles.value.push(created)
    await authStore.selectProfile(created.id)
    router.push('/')
  } catch (err) {
    error.value = err.response?.data?.error || 'Erreur lors de la création'
  } finally {
    creating.value = false
    newName.value = ''
  }
}

function formatDate(dateStr) {
  if (!dateStr) return ''
  try {
    return new Date(dateStr).toLocaleDateString('fr-FR', { day: 'numeric', month: 'short', year: 'numeric' })
  } catch {
    return dateStr
  }
}
</script>
