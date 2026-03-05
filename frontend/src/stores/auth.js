import { defineStore } from 'pinia'
import { ref } from 'vue'
import axios from 'axios'

const API_BASE = import.meta.env.VITE_API_URL || ''

export const useAuthStore = defineStore('auth', () => {
  const authenticated = ref(false)
  const loading = ref(false)
  const error = ref(null)
  const profile = ref(null) // { id, name } or null

  async function checkStatus() {
    try {
      const res = await axios.get(`${API_BASE}/auth/status`, { withCredentials: true })
      authenticated.value = res.data.authenticated === true
      if (res.data.profile_id) {
        profile.value = { id: res.data.profile_id, name: res.data.profile_name }
      } else {
        profile.value = null
      }
    } catch {
      authenticated.value = false
      profile.value = null
    }
  }

  async function login(password) {
    loading.value = true
    error.value = null
    try {
      await axios.post(`${API_BASE}/auth/login`, { password }, { withCredentials: true })
      authenticated.value = true
      profile.value = null
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
      profile.value = null
    }
  }

  async function selectProfile(id) {
    const res = await axios.post(`${API_BASE}/api/profiles/select`, { id }, { withCredentials: true })
    profile.value = { id: res.data.profile_id, name: res.data.profile_name }
  }

  return { authenticated, loading, error, profile, checkStatus, login, logout, selectProfile }
})
