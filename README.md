# Civil Game

A multiplayer 2D game built with Go backend, Vue.js frontend, Tailwind CSS, and Docker.

## Tech Stack

- **Backend**: Go 1.21 with Gorilla Mux and WebSockets
- **Frontend**: Vue.js 3 with Vite, Tailwind CSS
- **Containerization**: Docker and Docker Compose
- **Real-time Communication**: WebSocket for game state synchronization

## Project Structure

```
civil-game/
├── backend/              # Go backend server
│   ├── main.go          # Main server application
│   ├── go.mod           # Go dependencies
│   └── Dockerfile       # Backend Docker image
├── frontend/            # Vue.js frontend application
│   ├── src/
│   │   ├── components/  # Vue components
│   │   ├── services/    # WebSocket service
│   │   ├── App.vue      # Main app component
│   │   └── main.js      # Entry point
│   ├── package.json     # Frontend dependencies
│   ├── vite.config.js   # Vite configuration
│   ├── tailwind.config.js  # Tailwind CSS configuration
│   └── Dockerfile       # Frontend Docker image
├── docker-compose.yml   # Docker Compose configuration
└── README.md           # This file
```

## Prerequisites

- Docker and Docker Compose installed
- OR Node.js 20+ and Go 1.21+ for local development

## Quick Start with Docker

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd civil-game
   ```

2. **Build and run with Docker Compose**
   ```bash
   docker-compose up --build
   ```

3. **Access the game**
   - Frontend: http://localhost
   - Backend API: http://localhost:8080/api/health

## Local Development

### Backend Setup

1. **Navigate to backend directory**
   ```bash
   cd backend
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Run the server**
   ```bash
   go run main.go
   ```

   The backend will start on `http://localhost:8080`

### Frontend Setup

1. **Navigate to frontend directory**
   ```bash
   cd frontend
   ```

2. **Install dependencies**
   ```bash
   npm install
   ```

3. **Start development server**
   ```bash
   npm run dev
   ```

   The frontend will start on `http://localhost:5173` with Vite's hot module replacement.

4. **Build for production**
   ```bash
   npm run build
   ```

## Game Controls

- **WASD** or **Arrow Keys**: Move your player
- **Click**: Move player to clicked position

## Features

- Real-time multiplayer using WebSockets
- Multiple players can join and interact
- Responsive UI with Tailwind CSS
- Canvas-based 2D rendering
- Grid background for visual reference

## API Endpoints

- `GET /api/health` - Health check endpoint
- `WS /ws?playerId=<id>` - WebSocket connection for game state

## Development Notes

- The backend serves game state updates at 10 FPS (100ms intervals)
- Frontend renders at browser refresh rate (typically 60 FPS)
- Player position updates are sent immediately on movement
- WebSocket connection automatically reconnects on disconnect

## Docker Commands

- **Build and start**: `docker-compose up --build`
- **Start in background**: `docker-compose up -d`
- **Stop containers**: `docker-compose down`
- **View logs**: `docker-compose logs -f`
- **Rebuild specific service**: `docker-compose build backend` or `docker-compose build frontend`

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test locally
5. Submit a pull request

## License

MIT License
