package net

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/infopek/agorio/server/game"
	"github.com/infopek/agorio/server/utils"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Server struct {
	world *game.World
}

func NewServer(world *game.World) *Server {
	return &Server{
		world: world,
	}
}

func (s *Server) WebsocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading connection: %v\n", err)
		return
	}
	defer conn.Close()

	session := NewSession(conn, s.world.InputChan)
	go session.readLoop()
	go session.writeLoop()
}

// Later, this should be from user input
func (s *Server) getUsername() string {
	return fmt.Sprintf("random_name%d", utils.RandIntRange(1000, 10000))
}
