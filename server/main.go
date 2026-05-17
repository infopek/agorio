package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/infopek/agorio/server/game"
	"github.com/infopek/agorio/server/net"
)

const (
	port int = 8080
)

func main() {
	world := game.NewWorld()
	server := net.NewServer(world)

	http.Handle("/", http.FileServer(http.Dir("../client")))
	http.HandleFunc("/ws", server.WebsocketHandler)

	go world.Tick()

	log.Printf("Websocket server started on port %d", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	log.Printf("Error: %v\n", err)
}
