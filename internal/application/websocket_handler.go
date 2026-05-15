package application

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/infopek/agorio/internal/constants"
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

	// The client should know who they are, one-time thingy
	initMsg := map[string]any{
		"type":         "init",
		"player_id":    player.ID,
		"world_width":  constants.WorldWidth,
		"world_height": constants.WorldHeight,
	}
	conn.WriteJSON(initMsg)

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

			log.Printf("Target: %v, %v", msg.X, msg.Y)
			worldX, worldY := app.canvasToWorld(msg.X, msg.Y)
			app.world.UpdateDirection(playerID, worldX, worldY)
		}
	}
}
