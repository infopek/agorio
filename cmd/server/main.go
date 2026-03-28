package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/infopek/agorio/internal/application"
	"github.com/infopek/agorio/internal/constants"
)

func main() {

	app := application.NewApplication(application.ApplicationConfig{
		ApplicationName: "Agorio",
		WorldWidth:      constants.DefaultWorldWidth,
		WorldHeight:     constants.DefaultWorldHeight,
	})

	http.Handle("/", http.FileServer(http.Dir("./cmd/server/static")))
	http.HandleFunc("/ws", app.WebsocketHandler)

	go app.Run()

	log.Printf("WebSocket server started on port %d", constants.Port)
	log.Println(http.ListenAndServe(fmt.Sprintf(":%d", constants.Port), nil))
}
