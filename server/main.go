package main

import (
	"log"
	"fmt"
	"net/http"
)

const (
	port int = 8080
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("../client")))
	http.HandleFunc("/ws", server.WebsocketHandler)

	go server.Run()

	log.Printf("Websocket server started on port %d", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	log.Printf("Error: %v\n", err)
}
