package application

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/infopek/agorio/internal/math"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (app *Application) WebsocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading: ", err)
		return
	}
	defer conn.Close()

	// TODO: Get username from input
	name := fmt.Sprintf("randika%d", math.RandRange(1000, 10000))
	player := app.world.AddPlayer(name)

	app.mu.Lock()
	app.clients[conn] = player.ID
	app.mu.Unlock()

	defer func() {
		app.mu.Lock()
		delete(app.clients, conn)
		app.mu.Unlock()
		app.world.RemovePlayer(player.ID)
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			break
		}

		var msg ClientMessage
		if err = json.Unmarshal(message, &msg); err != nil {
			log.Printf("Invalid JSON: %v", err)
			continue
		}

		switch msg.Type {
		case "move":
			app.mu.RLock()
			playerID := app.clients[conn]
			app.mu.RUnlock()

			app.world.UpdateTarget(playerID, msg.X, msg.Y)
		}
		log.Printf("Message received: %v\n", msg)
	}
}
