package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/infopek/agorio/server/game"
	"github.com/infopek/agorio/server/net"
)

func main() {
	world := game.NewWorld()
	server := net.NewServer(world)

	http.Handle("/", http.FileServer(http.Dir("../client")))
	http.HandleFunc("/ws", server.WebsocketHandler)

	go world.Tick()

	srv := &http.Server{
		Addr: fmt.Sprintf(":%d", game.Port),

		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	log.Printf("Websocket server started on port %d", game.Port)
	err := srv.ListenAndServe()
	log.Printf("Error: %v\n", err)
}
