<template>
  <div>
    <div class="flex items-center justify-between mb-8">
      <div>
        <h1 class="text-3xl font-bold text-gray-900">Historique des programmes</h1>
        <p class="text-gray-500 mt-1">{{ programs.length }} programme{{ programs.length !== 1 ? 's' : '' }} généré{{ programs.length !== 1 ? 's' : '' }}</p>
      </div>
      <router-link to="/" class="btn-primary flex items-center gap-2">
        + Nouveau programme
      </router-link>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="text-center py-16">
      <div class="text-4xl mb-3">⏳</div>
      <p class="text-gray-400">Chargement...</p>
    </div>

    <!-- Empty state -->
    <div v-else-if="programs.length === 0" class="text-center py-20">
      <div class="text-6xl mb-4">🏋️</div>
      <h2 class="text-xl font-semibold text-gray-700 mb-2">Aucun programme pour l'instant</h2>
      <p class="text-gray-400 mb-6">Créez votre premier programme personnalisé</p>
      <router-link to="/" class="btn-primary">Créer un programme</router-link>
    </div>

    <!-- Programs grid -->
    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
      <div
        v-for="program in programs"
        :key="program.id"
        class="card hover:shadow-md transition-shadow cursor-pointer group"
        @click="openProgram(program)"
      >
        <!-- Header -->
        <div class="flex items-start justify-between mb-3">
          <div>
            <h3 class="text-lg font-bold text-gray-900 group-hover:text-blue-700 transition-colors">
              {{ program.person_name }}
            </h3>
            <p class="text-blue-600 text-sm font-medium">{{ program.objective }}</p>
          </div>
          <span class="tag bg-blue-50 text-blue-700 text-xs">
            S{{ program.week_number }}/{{ program.total_weeks }}
          </span>
        </div>

        <!-- Stats -->
        <div class="flex gap-4 text-sm text-gray-500 mb-4">
          <span>📅 {{ program.days?.length || 0 }} jours/sem</span>
          <span>💪 {{ totalExercises(program) }} exercices</span>
        </div>

        <!-- Day badges -->
        <div class="flex flex-wrap gap-1 mb-4">
          <span
            v-for="day in program.days?.slice(0, 4)"
            :key="day.day"
            class="tag bg-gray-100 text-gray-600 text-xs"
          >{{ day.focus || day.name }}</span>
          <span v-if="(program.days?.length || 0) > 4" class="tag bg-gray-100 text-gray-400 text-xs">
            +{{ program.days.length - 4 }}
          </span>
        </div>

        <!-- Footer -->
        <div class="flex items-center justify-between pt-3 border-t border-gray-100">
          <span class="text-xs text-gray-400">{{ formatDate(program.generated_at) }}</span>
          <div class="flex gap-2">
            <button
              @click.stop="downloadPDF(program.id)"
              class="text-xs text-gray-400 hover:text-blue-600 transition-colors px-2 py-1 rounded hover:bg-blue-50"
              title="Télécharger PDF"
            >
              📥 PDF
            </button>
            <button
              @click.stop="openProgram(program)"
              class="text-xs text-blue-600 hover:text-blue-800 font-medium px-2 py-1 rounded hover:bg-blue-50"
            >
              Voir →
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useProgramStore } from '../stores/program.js'
import axios from 'axios'

const router = useRouter()
const store = useProgramStore()

const programs = ref([])
const loading = ref(true)

const API_BASE = import.meta.env.VITE_API_URL || ''

onMounted(async () => {
  try {
    const res = await axios.get(`${API_BASE}/api/programs`, { withCredentials: true })
    programs.value = res.data.programs || []
  } catch (err) {
    console.error('Failed to load programs', err)
  } finally {
    loading.value = false
  }
})

function openProgram(program) {
  store.currentProgram = program
  router.push(`/program/${program.id}`)
}

function downloadPDF(id) {
  store.downloadPDF(id)
}

function totalExercises(program) {
  return (program.days || []).reduce((total, day) => {
    return total + (day.blocks || []).reduce((t, block) => t + (block.exercises || []).length, 0)
  }, 0)
}

function formatDate(dateStr) {
  if (!dateStr) return ''
  try {
    return new Date(dateStr).toLocaleDateString('fr-FR', {
      day: 'numeric', month: 'short', year: 'numeric'
    })
  } catch {
    return dateStr
  }
}
</script>
