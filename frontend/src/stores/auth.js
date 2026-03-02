import { defineStore } from 'pinia'
import { ref } from 'vue'
import axios from 'axios'

const API_BASE = import.meta.env.VITE_API_URL || ''

export const useAuthStore = defineStore('auth', () => {
  const authenticated = ref(false)
  const loading = ref(false)
  const error = ref(null)

  async function checkStatus() {
    try {
      const res = await axios.get(`${API_BASE}/auth/status`, { withCredentials: true })
      authenticated.value = res.data.authenticated === true
    } catch {
      authenticated.value = false
    }
  }

  async function login(password) {
    loading.value = true
    error.value = null
    try {
      await axios.post(`${API_BASE}/auth/login`, { password }, { withCredentials: true })
      authenticated.value = true
    } catch (err) {
      error.value = err.response?.data?.error || 'Erreur de connexion'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function logout() {
    try {
      await axios.post(`${API_BASE}/auth/logout`, {}, { withCredentials: true })
    } finally {
      authenticated.value = false
    }
  }

  return { authenticated, loading, error, checkStatus, login, logout }
})
