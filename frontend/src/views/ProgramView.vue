<template>
  <div>
    <!-- Loading -->
    <div v-if="store.loading" class="text-center py-20">
      <div class="text-4xl mb-4">⏳</div>
      <p class="text-gray-500 text-lg">Chargement du programme...</p>
    </div>

    <!-- Error -->
    <div v-else-if="store.error" class="bg-red-50 border border-red-200 text-red-700 rounded-lg p-6 text-center">
      <p class="font-semibold mb-2">Erreur</p>
      <p>{{ store.error }}</p>
      <router-link to="/" class="btn-primary mt-4 inline-block">Retour</router-link>
    </div>

    <!-- Program -->
    <div v-else-if="program">
      <!-- Program header -->
      <div class="card mb-6">
        <div class="flex flex-col md:flex-row md:items-center justify-between gap-4">
          <div>
            <h1 class="text-2xl font-bold text-gray-900">
              Programme de {{ program.person_name }}
            </h1>
            <p class="text-gray-500 mt-1">
              Semaine {{ program.week_number }}/{{ program.total_weeks }} •
              {{ program.days.length }} jours/semaine
            </p>
            <p class="text-blue-700 font-medium mt-1">{{ program.objective }}</p>
          </div>
          <div class="flex gap-3">
            <button @click="store.downloadPDF(program.id)" class="btn-secondary flex items-center gap-2">
              📥 Télécharger PDF
            </button>
            <router-link to="/" class="btn-primary flex items-center gap-2">
              + Nouveau
            </router-link>
          </div>
        </div>

        <!-- Feedback summary if present -->
        <div v-if="program.feedback" class="mt-4 pt-4 border-t border-gray-100">
          <p class="text-sm font-medium text-gray-500 mb-2">Basé sur votre ressenti précédent:</p>
          <div class="flex flex-wrap gap-3">
            <span class="tag bg-blue-50 text-blue-700">⚡ Énergie: {{ program.feedback.energy_level }}/10</span>
            <span class="tag bg-orange-50 text-orange-700">😅 Courbatures: {{ program.feedback.soreness_level }}/10</span>
            <span class="tag bg-green-50 text-green-700">🔥 Motivation: {{ program.feedback.motivation_level }}/10</span>
          </div>
        </div>
      </div>

      <!-- Day tabs -->
      <div class="flex gap-2 mb-6 overflow-x-auto pb-2">
        <button
          v-for="(day, idx) in program.days"
          :key="idx"
          @click="activeDay = idx"
          :class="['px-4 py-2 rounded-lg font-medium text-sm whitespace-nowrap transition-colors',
            activeDay === idx
              ? 'bg-blue-600 text-white shadow-sm'
              : 'bg-white text-gray-600 hover:bg-gray-50 border border-gray-200']"
        >
          Jour {{ day.day }}: {{ day.name }}
        </button>
      </div>

      <!-- Active day content -->
      <div v-if="currentDay" class="space-y-4">
        <!-- Day header -->
        <div class="card bg-gradient-to-r from-blue-50 to-indigo-50 border-blue-100">
          <div class="flex flex-col md:flex-row md:items-center justify-between gap-3">
            <div>
              <h2 class="text-xl font-bold text-blue-900">{{ currentDay.name }}</h2>
              <p class="text-blue-700">{{ currentDay.focus }}</p>
              <p class="text-sm text-gray-500 mt-1">⏱ Durée estimée: {{ currentDay.duration }} min</p>
            </div>
            <button @click="openTimer(activeDay)" class="btn-primary flex items-center gap-2">
              ⏱ Timer d'entraînement
            </button>
          </div>
        </div>

        <!-- Warm-up notes -->
        <div v-if="currentDay.warmup_notes" class="card border-l-4 border-yellow-400 bg-yellow-50">
          <h3 class="font-semibold text-yellow-800 mb-1">🔥 Échauffement</h3>
          <p class="text-yellow-700 text-sm">{{ currentDay.warmup_notes }}</p>
        </div>

        <!-- Exercise blocks -->
        <div v-for="block in currentDay.blocks" :key="block.name" class="card">
          <h3 class="font-bold text-gray-800 mb-4 text-lg border-b border-gray-100 pb-2">
            {{ block.name }}
          </h3>

          <!-- Exercise list -->
          <div class="space-y-3">
            <div
              v-for="(exercise, idx) in block.exercises"
              :key="idx"
              class="rounded-lg border border-gray-100 overflow-hidden"
              :class="idx % 2 === 0 ? 'bg-gray-50' : 'bg-white'"
            >
              <div class="p-4">
                <!-- Exercise name and muscle groups -->
                <div class="flex flex-wrap items-start justify-between gap-2 mb-2">
                  <h4 class="font-semibold text-gray-900 text-base">{{ exercise.name }}</h4>
                  <div class="flex flex-wrap gap-1">
                    <span
                      v-for="muscle in exercise.muscle_groups"
                      :key="muscle"
                      class="tag bg-purple-50 text-purple-700 text-xs"
                    >{{ muscle }}</span>
                  </div>
                </div>

                <!-- Exercise details grid -->
                <div class="grid grid-cols-2 md:grid-cols-4 gap-3 mt-3">
                  <div class="text-center bg-blue-50 rounded-lg p-2">
                    <div class="text-lg font-bold text-blue-700">{{ exercise.sets }}</div>
                    <div class="text-xs text-blue-500">Séries</div>
                  </div>
                  <div class="text-center bg-green-50 rounded-lg p-2">
                    <div class="text-lg font-bold text-green-700">{{ exercise.reps }}</div>
                    <div class="text-xs text-green-500">Répétitions</div>
                  </div>
                  <div class="text-center bg-orange-50 rounded-lg p-2">
                    <div class="text-lg font-bold text-orange-700">{{ exercise.intensity }}</div>
                    <div class="text-xs text-orange-500">Intensité</div>
                  </div>
                  <div class="text-center bg-gray-100 rounded-lg p-2">
                    <div class="text-lg font-bold text-gray-700">{{ exercise.rest_seconds }}s</div>
                    <div class="text-xs text-gray-500">Repos</div>
                  </div>
                </div>

                <!-- Tempo and notes -->
                <div v-if="exercise.tempo || exercise.notes" class="mt-3 space-y-1">
                  <div v-if="exercise.tempo" class="text-sm">
                    <span class="font-medium text-gray-600">Tempo: </span>
                    <span class="font-mono bg-gray-100 px-1.5 py-0.5 rounded text-gray-700">{{ exercise.tempo }}</span>
                    <span class="text-xs text-gray-400 ml-1">(excentrique-pause-concentrique-pause)</span>
                  </div>
                  <div v-if="exercise.notes" class="text-sm text-gray-600 bg-blue-50 rounded px-3 py-2 border-l-2 border-blue-300">
                    💡 {{ exercise.notes }}
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Cooldown notes -->
        <div v-if="currentDay.cooldown_notes" class="card border-l-4 border-green-400 bg-green-50">
          <h3 class="font-semibold text-green-800 mb-1">🧘 Retour au calme</h3>
          <p class="text-green-700 text-sm">{{ currentDay.cooldown_notes }}</p>
        </div>
      </div>

      <!-- Program notes -->
      <div v-if="program.notes" class="card mt-6">
        <h3 class="font-bold text-gray-700 mb-2">📋 Notes du programme</h3>
        <p class="text-gray-600 text-sm">{{ program.notes }}</p>
      </div>
    </div>

    <!-- Timer modal -->
    <TimerModal v-if="showTimer" :timer="store.timer" @close="showTimer = false" />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useProgramStore } from '../stores/program.js'
import TimerModal from '../components/TimerModal.vue'

const props = defineProps({ id: String })
const store = useProgramStore()
const activeDay = ref(0)
const showTimer = ref(false)

const program = computed(() => store.currentProgram)
const currentDay = computed(() => program.value?.days?.[activeDay.value])

onMounted(async () => {
  if (props.id && (!store.currentProgram || store.currentProgram.id !== props.id)) {
    await store.fetchProgram(props.id)
  }
})

async function openTimer(dayIndex) {
  if (!program.value?.id) return
  await store.fetchTimer(program.value.id, dayIndex)
  showTimer.value = true
}
</script>
