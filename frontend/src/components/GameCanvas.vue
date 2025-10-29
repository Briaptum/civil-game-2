<template>
  <canvas
    ref="canvasRef"
    class="w-full h-full"
    @click="handleCanvasClick"
  ></canvas>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useWebSocket } from '../services/websocket'

const props = defineProps({
  playerId: {
    type: String,
    required: true
  }
})

const emit = defineEmits(['connection-status', 'player-count'])

const canvasRef = ref(null)
let ctx = null
let animationFrameId = null
const keys = {}
const player = {
  id: props.playerId,
  x: 400,
  y: 300,
  velocity: 0,
  angle: 0,
  speed: 3
}

const ws = useWebSocket(props.playerId)

ws.on('open', () => {
  emit('connection-status', 'connected')
})

ws.on('close', () => {
  emit('connection-status', 'disconnected')
})

ws.on('error', () => {
  emit('connection-status', 'error')
})

ws.on('message', (gameState) => {
  renderGame(gameState)
  emit('player-count', Object.keys(gameState.players || {}).length)
})

const handleKeyDown = (e) => {
  keys[e.key.toLowerCase()] = true
  keys[e.code] = true
}

const handleKeyUp = (e) => {
  keys[e.key.toLowerCase()] = false
  keys[e.code] = false
}

const handleCanvasClick = (e) => {
  const rect = canvasRef.value.getBoundingClientRect()
  const clickX = e.clientX - rect.left
  const clickY = e.clientY - rect.top
  
  player.x = clickX
  player.y = clickY
  updatePlayerPosition()
}

const updatePlayer = () => {
  let moved = false

  if (keys['w'] || keys['arrowup']) {
    player.y -= player.speed
    moved = true
  }
  if (keys['s'] || keys['arrowdown']) {
    player.y += player.speed
    moved = true
  }
  if (keys['a'] || keys['arrowleft']) {
    player.x -= player.speed
    moved = true
  }
  if (keys['d'] || keys['arrowright']) {
    player.x += player.speed
    moved = true
  }

  if (moved) {
    updatePlayerPosition()
  }
}

const updatePlayerPosition = () => {
  if (ws.isConnected()) {
    ws.send(JSON.stringify(player))
  }
}

const renderGame = (gameState) => {
  if (!ctx || !canvasRef.value) return

  const canvas = canvasRef.value
  ctx.clearRect(0, 0, canvas.width, canvas.height)

  // Draw background grid
  drawGrid()

  // Draw all players
  if (gameState.players) {
    Object.values(gameState.players).forEach(p => {
      drawPlayer(p, p.id === props.playerId)
    })
  }
}

const drawGrid = () => {
  const canvas = canvasRef.value
  ctx.strokeStyle = '#374151'
  ctx.lineWidth = 1

  const gridSize = 50
  for (let x = 0; x < canvas.width; x += gridSize) {
    ctx.beginPath()
    ctx.moveTo(x, 0)
    ctx.lineTo(x, canvas.height)
    ctx.stroke()
  }

  for (let y = 0; y < canvas.height; y += gridSize) {
    ctx.beginPath()
    ctx.moveTo(0, y)
    ctx.lineTo(canvas.width, y)
    ctx.stroke()
  }
}

const drawPlayer = (p, isCurrentPlayer) => {
  ctx.save()
  
  if (isCurrentPlayer) {
    ctx.fillStyle = '#3b82f6'
    ctx.strokeStyle = '#60a5fa'
    ctx.lineWidth = 3
  } else {
    ctx.fillStyle = '#ef4444'
    ctx.strokeStyle = '#f87171'
    ctx.lineWidth = 2
  }

  ctx.beginPath()
  ctx.arc(p.x, p.y, 15, 0, Math.PI * 2)
  ctx.fill()
  ctx.stroke()

  // Draw direction indicator
  ctx.strokeStyle = isCurrentPlayer ? '#60a5fa' : '#f87171'
  ctx.lineWidth = 2
  ctx.beginPath()
  ctx.moveTo(p.x, p.y)
  ctx.lineTo(p.x + Math.cos(p.angle) * 20, p.y + Math.sin(p.angle) * 20)
  ctx.stroke()

  ctx.restore()
}

const gameLoop = () => {
  updatePlayer()
  animationFrameId = requestAnimationFrame(gameLoop)
}

const resizeCanvas = () => {
  if (canvasRef.value) {
    const canvas = canvasRef.value
    const rect = canvas.getBoundingClientRect()
    canvas.width = rect.width
    canvas.height = rect.height
  }
}

onMounted(() => {
  const canvas = canvasRef.value
  ctx = canvas.getContext('2d')
  
  resizeCanvas()
  window.addEventListener('resize', resizeCanvas)
  window.addEventListener('keydown', handleKeyDown)
  window.addEventListener('keyup', handleKeyUp)
  
  gameLoop()
  
  // Send initial position
  setTimeout(() => {
    updatePlayerPosition()
  }, 500)
})

onUnmounted(() => {
  window.removeEventListener('resize', resizeCanvas)
  window.removeEventListener('keydown', handleKeyDown)
  window.removeEventListener('keyup', handleKeyUp)
  
  if (animationFrameId) {
    cancelAnimationFrame(animationFrameId)
  }
  
  ws.close()
})
</script>

