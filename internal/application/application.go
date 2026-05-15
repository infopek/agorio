package application

import (
	"encoding/json"
	_"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"github.com/infopek/agorio/internal/constants"
	"github.com/infopek/agorio/internal/game"
	"github.com/infopek/agorio/internal/types"
)

type ApplicationConfig struct {
	ApplicationName string
	WorldWidth      types.Real
	WorldHeight     types.Real
	CanvasWidth     types.Real
	CanvasHeight    types.Real
}

type Application struct {
	config ApplicationConfig
	world  *game.World

	clients map[*websocket.Conn]uuid.UUID // player IDs

	mu sync.RWMutex
}

func NewApplication(config ApplicationConfig) *Application {
	world := game.NewWorld(game.WorldConfig{
		Width:  config.WorldWidth,
		Height: config.WorldHeight,
	})
	return &Application{
		config:  config,
		world:   world,
		clients: make(map[*websocket.Conn]uuid.UUID),
	}
}

func (app *Application) Run() {
	ticker := time.NewTicker(time.Second / time.Duration(constants.TickRate))
	defer ticker.Stop()

	for range ticker.C {
		app.world.Update(constants.Dt)

		app.BroadcastState()
	}
}

func (app *Application) BroadcastState() error {
	state := app.world.Snapshot()

	data, err := json.Marshal(state)
	if err != nil {
		return err
	}

	for conn, _ := range app.clients {
		err := conn.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			conn.Close()
			delete(app.clients, conn)
		}
	}

	return nil
}

func (app *Application) canvasToWorld(canvasX, canvasY types.Real) (worldX, worldY types.Real) {
	worldX = (canvasX / app.config.CanvasWidth) * app.config.WorldWidth
	worldY = (canvasY / app.config.CanvasHeight) * app.config.WorldHeight
	return
}
