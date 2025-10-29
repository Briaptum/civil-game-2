<template>
  <div class="h-screen w-screen bg-gray-900 text-white flex flex-col">
    <header class="bg-gray-800 p-4 shadow-lg">
      <h1 class="text-2xl font-bold text-center">Civil Game</h1>
      <div class="flex justify-center gap-4 mt-2">
        <span class="text-sm">Status: <span :class="connectionStatus === 'connected' ? 'text-green-500' : 'text-red-500'">{{ connectionStatus }}</span></span>
        <span class="text-sm">Players: {{ playerCount }}</span>
      </div>
    </header>
    <main class="flex-1 relative">
      <GameCanvas :player-id="playerId" @connection-status="handleConnectionStatus" @player-count="handlePlayerCount" />
    </main>
    <footer class="bg-gray-800 p-2 text-center text-xs">
      <p>Use WASD or Arrow Keys to move | Click to interact</p>
    </footer>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import GameCanvas from './components/GameCanvas.vue'

const playerId = ref('')
const connectionStatus = ref('disconnected')
const playerCount = ref(0)

onMounted(() => {
  playerId.value = `player-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`
})

const handleConnectionStatus = (status) => {
  connectionStatus.value = status
}

const handlePlayerCount = (count) => {
  playerCount.value = count
}
</script>

