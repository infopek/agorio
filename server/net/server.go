package net

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"

	"github.com/infopek/agorio/server/game"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			origin := r.Header.Get("Origin")
			for _, o := range game.AllowedOrigins {
				if origin == o {
					return true
				}
			}
			return false
		},
	}
)

type Server struct {
	world *game.World

	mu          sync.Mutex
	connections uint64
}

func NewServer(world *game.World) *Server {
	return &Server{
		world: world,
	}
}

func (s *Server) WebsocketHandler(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	if s.connections >= game.MaxConnections {
		s.mu.Unlock()
		http.Error(w, "server full :(", http.StatusServiceUnavailable)
		return
	}
	s.connections++
	s.mu.Unlock()

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.mu.Lock()
		s.connections--
		s.mu.Unlock()

		log.Printf("Error upgrading connection: %v\n", err)
		return
	}

	session := NewSession(conn, s.world.InputChan, func() {
		s.mu.Lock()
		s.connections--
		s.mu.Unlock()
	})
	go session.readLoop()
	go session.writeLoop()
}
