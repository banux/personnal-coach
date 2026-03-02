<template>
  <div class="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-center justify-center p-4">
    <div class="bg-white rounded-2xl shadow-2xl w-full max-w-md max-h-[90vh] flex flex-col">
      <!-- Header -->
      <div class="bg-blue-600 text-white rounded-t-2xl p-4 flex items-center justify-between">
        <div>
          <h2 class="text-lg font-bold">⏱ Timer</h2>
          <p class="text-blue-200 text-sm">{{ timer?.day_name }}</p>
        </div>
        <button @click="$emit('close')" class="text-white hover:text-blue-200 text-xl font-bold">✕</button>
      </div>

      <!-- Timer display -->
      <div class="p-6 text-center border-b border-gray-100">
        <div class="text-6xl font-mono font-bold mb-2" :class="isResting ? 'text-orange-500' : 'text-blue-600'">
          {{ formatTime(timeRemaining) }}
        </div>
        <div class="text-sm font-medium" :class="isResting ? 'text-orange-500' : 'text-blue-600'">
          {{ isResting ? '😮‍💨 REPOS' : '💪 TRAVAIL' }}
        </div>

        <!-- Current exercise -->
        <div v-if="currentSet" class="mt-4 bg-gray-50 rounded-xl p-4">
          <div class="text-lg font-semibold text-gray-800">{{ currentSet.exercise_name }}</div>
          <div class="text-gray-500 text-sm mt-1">
            Série {{ currentSet.set_number }} • {{ currentSet.reps }} reps
          </div>
        </div>

        <!-- Next up -->
        <div v-if="nextSet" class="mt-2 text-xs text-gray-400">
          Suivant: {{ nextSet.exercise_name }} — Série {{ nextSet.set_number }}
        </div>
      </div>

      <!-- Progress -->
      <div class="px-6 py-3 bg-gray-50">
        <div class="flex justify-between text-xs text-gray-500 mb-1">
          <span>Progression</span>
          <span>{{ currentSetIndex + 1 }} / {{ timer?.sets?.length || 0 }}</span>
        </div>
        <div class="w-full bg-gray-200 rounded-full h-2">
          <div
            class="bg-blue-600 rounded-full h-2 transition-all duration-300"
            :style="{ width: progressPercent + '%' }"
          ></div>
        </div>
      </div>

      <!-- Controls -->
      <div class="p-4 flex gap-3 justify-center">
        <button @click="previous" class="btn-secondary" :disabled="currentSetIndex === 0">
          ⏮
        </button>
        <button @click="togglePause" class="btn-primary px-8">
          {{ isPaused ? '▶ Reprendre' : '⏸ Pause' }}
        </button>
        <button @click="next" class="btn-secondary" :disabled="currentSetIndex >= (timer?.sets?.length || 0) - 1">
          ⏭
        </button>
      </div>

      <!-- Set list (scrollable) -->
      <div class="overflow-y-auto max-h-48 border-t border-gray-100">
        <div
          v-for="(set, idx) in timer?.sets"
          :key="idx"
          :class="['px-4 py-2 text-sm flex justify-between items-center',
            idx === currentSetIndex ? 'bg-blue-50 text-blue-800 font-medium' : 'text-gray-600',
            idx < currentSetIndex ? 'opacity-40 line-through' : '']"
        >
          <span>{{ set.exercise_name }} — Série {{ set.set_number }}</span>
          <span class="text-xs text-gray-400">{{ set.reps }} reps • {{ set.rest_seconds }}s repos</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'

const props = defineProps({ timer: Object })
const emit = defineEmits(['close'])

const currentSetIndex = ref(0)
const timeRemaining = ref(0)
const isResting = ref(false)
const isPaused = ref(false)
let intervalId = null

const currentSet = computed(() => props.timer?.sets?.[currentSetIndex.value])
const nextSet = computed(() => props.timer?.sets?.[currentSetIndex.value + 1])
const progressPercent = computed(() => {
  const total = props.timer?.sets?.length || 1
  return Math.round((currentSetIndex.value / total) * 100)
})

function initTimer() {
  if (!currentSet.value) return
  isResting.value = false
  timeRemaining.value = currentSet.value.work_seconds || 30
}

function startInterval() {
  clearInterval(intervalId)
  intervalId = setInterval(() => {
    if (isPaused.value) return

    if (timeRemaining.value > 0) {
      timeRemaining.value--
    } else {
      // Transition
      if (!isResting.value && currentSet.value?.rest_seconds > 0) {
        isResting.value = true
        timeRemaining.value = currentSet.value.rest_seconds
        // Beep
        playBeep()
      } else {
        // Move to next set
        if (currentSetIndex.value < (props.timer?.sets?.length || 0) - 1) {
          currentSetIndex.value++
          initTimer()
          playBeep()
        } else {
          clearInterval(intervalId)
          timeRemaining.value = 0
          isResting.value = false
        }
      }
    }
  }, 1000)
}

function togglePause() {
  isPaused.value = !isPaused.value
}

function next() {
  if (currentSetIndex.value < (props.timer?.sets?.length || 0) - 1) {
    currentSetIndex.value++
    initTimer()
  }
}

function previous() {
  if (currentSetIndex.value > 0) {
    currentSetIndex.value--
    initTimer()
  }
}

function formatTime(seconds) {
  const m = Math.floor(seconds / 60)
  const s = seconds % 60
  return `${String(m).padStart(2, '0')}:${String(s).padStart(2, '0')}`
}

function playBeep() {
  try {
    const ctx = new (window.AudioContext || window.webkitAudioContext)()
    const osc = ctx.createOscillator()
    const gain = ctx.createGain()
    osc.connect(gain)
    gain.connect(ctx.destination)
    osc.frequency.value = isResting.value ? 440 : 880
    gain.gain.setValueAtTime(0.3, ctx.currentTime)
    gain.gain.exponentialRampToValueAtTime(0.001, ctx.currentTime + 0.3)
    osc.start(ctx.currentTime)
    osc.stop(ctx.currentTime + 0.3)
  } catch {
    // Audio not supported
  }
}

onMounted(() => {
  initTimer()
  startInterval()
})

onUnmounted(() => {
  clearInterval(intervalId)
})

watch(currentSetIndex, () => {
  initTimer()
})
</script>
