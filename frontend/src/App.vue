<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Navigation -->
    <nav class="bg-white shadow-sm border-b border-gray-100">
      <div class="max-w-6xl mx-auto px-4 py-3 flex items-center justify-between">
        <router-link to="/" class="flex items-center space-x-2">
          <span class="text-2xl">💪</span>
          <span class="text-xl font-bold text-blue-700">Coach Personnel IA</span>
        </router-link>
        <router-link to="/" class="text-sm text-gray-500 hover:text-blue-600 transition-colors">
          Nouveau programme
        </router-link>
      </div>
    </nav>

    <!-- Main content -->
    <main class="max-w-6xl mx-auto px-4 py-8">
      <router-view />
    </main>
  </div>
</template>

<script setup>
</script>
