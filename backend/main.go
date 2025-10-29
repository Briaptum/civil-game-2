package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins in development
	},
}

// Player represents a player in the game
type Player struct {
	ID       string  `json:"id"`
	X        float64 `json:"x"`
	Y        float64 `json:"y"`
	Velocity float64 `json:"velocity"`
	Angle    float64 `json:"angle"`
}

// GameState represents the current state of the game
type GameState struct {
	Players map[string]*Player `json:"players"`
	mu      sync.RWMutex
}

// Client represents a WebSocket client connection
type Client struct {
	conn   *websocket.Conn
	player *Player
	game   *GameState
}

// Hub manages active clients and broadcasts messages
type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	game       *GameState
}

func newHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		game: &GameState{
			Players: make(map[string]*Player),
		},
	}
}

func (h *Hub) run() {
	ticker := time.NewTicker(100 * time.Millisecond) // 10 FPS game loop
	defer ticker.Stop()

	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			client.game.mu.Lock()
			client.game.Players[client.player.ID] = client.player
			client.game.mu.Unlock()
			log.Printf("Client registered: %s", client.player.ID)

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				client.game.mu.Lock()
				delete(client.game.Players, client.player.ID)
				client.game.mu.Unlock()
				client.conn.Close()
				log.Printf("Client unregistered: %s", client.player.ID)
			}

		case message := <-h.broadcast:
			for client := range h.clients {
				err := client.conn.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					client.conn.Close()
					delete(h.clients, client)
				}
			}

		case <-ticker.C:
			// Broadcast game state to all clients
			h.game.mu.RLock()
			state, err := json.Marshal(h.game)
			h.game.mu.RUnlock()
			if err == nil {
				h.broadcast <- state
			}
		}
	}
}

func (c *Client) readPump() {
	defer func() {
		c.game.mu.Lock()
		if _, ok := c.game.Players[c.player.ID]; ok {
			delete(c.game.Players, c.player.ID)
		}
		c.game.mu.Unlock()
		c.conn.Close()
	}()

	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		var playerUpdate Player
		if err := json.Unmarshal(message, &playerUpdate); err == nil {
			c.game.mu.Lock()
			if p, ok := c.game.Players[playerUpdate.ID]; ok {
				p.X = playerUpdate.X
				p.Y = playerUpdate.Y
				p.Velocity = playerUpdate.Velocity
				p.Angle = playerUpdate.Angle
			}
			c.game.mu.Unlock()
		}
	}
}

func handleWebSocket(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	playerID := r.URL.Query().Get("playerId")
	if playerID == "" {
		playerID = generatePlayerID()
	}

	player := &Player{
		ID:       playerID,
		X:        400,
		Y:        300,
		Velocity: 0,
		Angle:    0,
	}

	client := &Client{
		conn:   conn,
		player: player,
		game:   hub.game,
	}

	hub.register <- client

	go client.readPump()
}

func generatePlayerID() string {
	return time.Now().Format("20060102150405") + "-" + string(rune(time.Now().UnixNano()%10000+65))
}

func apiHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func main() {
	hub := newHub()
	go hub.run()

	r := mux.NewRouter()

	// API routes
	r.HandleFunc("/api/health", apiHealth).Methods("GET")

	// WebSocket route
	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handleWebSocket(hub, w, r)
	})

	// Serve static files (for production, use nginx in front)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("../frontend/dist/")))

	port := ":8080"
	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(port, r))
}

