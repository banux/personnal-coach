<template>
  <div>
    <!-- Hero -->
    <div class="text-center mb-10">
      <h1 class="text-4xl font-bold text-gray-900 mb-3">Votre Programme Personnalisé</h1>
      <p class="text-lg text-gray-500">Généré par IA selon votre profil, objectifs et ressenti</p>
    </div>

    <!-- Error -->
    <div v-if="store.error" class="mb-6 bg-red-50 border border-red-200 text-red-700 rounded-lg p-4">
      {{ store.error }}
    </div>

    <!-- Form -->
    <div class="card mb-8">
      <h2 class="text-xl font-bold text-gray-800 mb-6">Votre profil</h2>
      <form @submit.prevent="handleSubmit" class="space-y-6">

        <!-- Personal info -->
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          <div>
            <label class="label">Prénom *</label>
            <input v-model="form.person.name" type="text" class="input-field" placeholder="Ex: Marie" required />
          </div>
          <div>
            <label class="label">Âge</label>
            <input v-model.number="form.person.age" type="number" class="input-field" placeholder="30" min="10" max="100" />
          </div>
          <div>
            <label class="label">Poids (kg)</label>
            <input v-model.number="form.person.weight" type="number" class="input-field" placeholder="70" min="30" max="300" step="0.1" />
          </div>
          <div>
            <label class="label">Taille (cm)</label>
            <input v-model.number="form.person.height" type="number" class="input-field" placeholder="170" min="100" max="250" />
          </div>
          <div>
            <label class="label">Niveau</label>
            <select v-model="form.person.level" class="input-field">
              <option value="beginner">Débutant</option>
              <option value="intermediate">Intermédiaire</option>
              <option value="advanced">Avancé</option>
            </select>
          </div>
        </div>

        <!-- Goals -->
        <div>
          <label class="label">Objectifs</label>
          <div class="flex flex-wrap gap-2">
            <button
              v-for="goal in availableGoals"
              :key="goal.value"
              type="button"
              @click="toggleGoal(goal.value)"
              :class="['tag py-1.5 px-3 text-sm cursor-pointer transition-colors', form.person.goals.includes(goal.value)
                ? 'bg-blue-100 text-blue-800 border border-blue-300'
                : 'bg-gray-100 text-gray-600 border border-gray-200 hover:bg-gray-200']"
            >
              {{ goal.label }}
            </button>
          </div>
        </div>

        <!-- Equipment -->
        <div>
          <label class="label">Équipement disponible</label>
          <div class="flex flex-wrap gap-2">
            <button
              v-for="eq in availableEquipment"
              :key="eq.value"
              type="button"
              @click="toggleEquipment(eq.value)"
              :class="['tag py-1.5 px-3 text-sm cursor-pointer transition-colors', form.person.equipment.includes(eq.value)
                ? 'bg-green-100 text-green-800 border border-green-300'
                : 'bg-gray-100 text-gray-600 border border-gray-200 hover:bg-gray-200']"
            >
              {{ eq.label }}
            </button>
          </div>
        </div>

        <!-- Program settings -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label class="label">Jours par semaine</label>
            <select v-model.number="form.days_per_week" class="input-field">
              <option v-for="d in [2,3,4,5,6]" :key="d" :value="d">{{ d }} jours</option>
            </select>
          </div>
          <div>
            <label class="label">Durée du programme</label>
            <select v-model.number="form.weeks" class="input-field">
              <option v-for="w in [4,6,8,12]" :key="w" :value="w">{{ w }} semaines</option>
            </select>
          </div>
        </div>

        <!-- Weekly feedback (optional) -->
        <div class="border border-gray-200 rounded-lg p-4">
          <button type="button" @click="showFeedback = !showFeedback" class="flex items-center justify-between w-full text-left">
            <span class="font-medium text-gray-700">Ressenti de la semaine précédente (optionnel)</span>
            <span class="text-gray-400">{{ showFeedback ? '▲' : '▼' }}</span>
          </button>

          <div v-if="showFeedback" class="mt-4 space-y-4">
            <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
              <div>
                <label class="label">Énergie (1-10)</label>
                <input v-model.number="form.feedback.energy_level" type="range" min="1" max="10" class="w-full" />
                <div class="text-center text-sm text-gray-500 mt-1">{{ form.feedback.energy_level }}/10</div>
              </div>
              <div>
                <label class="label">Courbatures (1-10)</label>
                <input v-model.number="form.feedback.soreness_level" type="range" min="1" max="10" class="w-full" />
                <div class="text-center text-sm text-gray-500 mt-1">{{ form.feedback.soreness_level }}/10</div>
              </div>
              <div>
                <label class="label">Motivation (1-10)</label>
                <input v-model.number="form.feedback.motivation_level" type="range" min="1" max="10" class="w-full" />
                <div class="text-center text-sm text-gray-500 mt-1">{{ form.feedback.motivation_level }}/10</div>
              </div>
            </div>
            <div>
              <label class="label">Jours complétés</label>
              <input v-model.number="form.feedback.completed_days" type="number" min="0" max="7" class="input-field" />
            </div>
            <div>
              <label class="label">Notes / Commentaires</label>
              <textarea v-model="form.feedback.notes" class="input-field" rows="2" placeholder="Blessure, fatigue, commentaire..."></textarea>
            </div>
          </div>
        </div>

        <button type="submit" class="btn-primary w-full py-3 text-lg" :disabled="store.loading">
          <span v-if="store.loading">⏳ Génération en cours...</span>
          <span v-else>🚀 Générer mon programme</span>
        </button>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useProgramStore } from '../stores/program.js'

const router = useRouter()
const store = useProgramStore()
const showFeedback = ref(false)

const form = reactive({
  person: {
    name: '',
    age: null,
    weight: null,
    height: null,
    level: 'intermediate',
    goals: [],
    equipment: ['bodyweight'],
  },
  days_per_week: 3,
  weeks: 4,
  feedback: {
    energy_level: 7,
    soreness_level: 3,
    motivation_level: 7,
    completed_days: 0,
    notes: '',
  }
})

const availableGoals = [
  { value: 'weight_loss', label: '⚖️ Perte de poids' },
  { value: 'muscle_gain', label: '💪 Prise de masse' },
  { value: 'strength', label: '🏋️ Force' },
  { value: 'endurance', label: '🏃 Endurance' },
  { value: 'flexibility', label: '🧘 Souplesse' },
  { value: 'general_fitness', label: '❤️ Forme générale' },
]

const availableEquipment = [
  { value: 'bodyweight', label: '🤸 Poids du corps' },
  { value: 'dumbbell', label: '🏋️ Haltères' },
  { value: 'barbell', label: '⚡ Barre' },
  { value: 'machine', label: '🔧 Machines' },
  { value: 'kettlebell', label: '🔔 Kettlebell' },
  { value: 'bands', label: '↔️ Élastiques' },
  { value: 'pullup_bar', label: '⬆️ Barre de traction' },
]

function toggleGoal(value) {
  const idx = form.person.goals.indexOf(value)
  if (idx >= 0) form.person.goals.splice(idx, 1)
  else form.person.goals.push(value)
}

function toggleEquipment(value) {
  const idx = form.person.equipment.indexOf(value)
  if (idx >= 0) form.person.equipment.splice(idx, 1)
  else form.person.equipment.push(value)
}

async function handleSubmit() {
  if (!form.person.name.trim()) return

  const payload = {
    person: { ...form.person },
    days_per_week: form.days_per_week,
    weeks: form.weeks,
  }

  if (showFeedback.value) {
    payload.feedback = { ...form.feedback }
  }

  try {
    const result = await store.generateProgram(payload)
    router.push(`/program/${result.program.id}`)
  } catch {
    // error already set in store
  }
}
</script>
