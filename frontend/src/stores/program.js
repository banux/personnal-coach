import { defineStore } from 'pinia'
import { ref } from 'vue'
import axios from 'axios'

// When served by Go, API is on same origin. Dev mode uses proxy or explicit URL.
const API_BASE = import.meta.env.VITE_API_URL || ''

// All requests use credentials (cookie-based auth)
const api = axios.create({ baseURL: API_BASE, withCredentials: true })

export const useProgramStore = defineStore('program', () => {
  const currentProgram = ref(null)
  const programs = ref([])
  const loading = ref(false)
  const error = ref(null)
  const timer = ref(null)

  async function generateProgram(formData) {
    loading.value = true
    error.value = null
    try {
      const res = await api.post('/api/programs/generate', formData)
      currentProgram.value = res.data.program
      return res.data
    } catch (err) {
      error.value = err.response?.data?.error || 'Erreur lors de la génération du programme'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function fetchProgram(id) {
    loading.value = true
    error.value = null
    try {
      const res = await api.get(`/api/programs/${id}`)
      currentProgram.value = res.data
      return res.data
    } catch (err) {
      error.value = err.response?.data?.error || 'Erreur lors du chargement du programme'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function fetchTimer(programId, dayIndex) {
    loading.value = true
    error.value = null
    try {
      const res = await api.get(`/api/programs/${programId}/timer/${dayIndex}`)
      timer.value = res.data
      return res.data
    } catch (err) {
      error.value = err.response?.data?.error || 'Erreur lors du chargement du timer'
      throw err
    } finally {
      loading.value = false
    }
  }

  function downloadPDF(programId) {
    window.open(`${API_BASE}/api/programs/${programId}/pdf`, '_blank')
  }

  return { currentProgram, programs, loading, error, timer, generateProgram, fetchProgram, fetchTimer, downloadPDF }
})
