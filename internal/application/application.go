package application

import (
	"encoding/json"
	_"fmt"
	"time"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/google/uuid"

	"github.com/infopek/agorio/internal/constants"
	"github.com/infopek/agorio/internal/models"
	"github.com/infopek/agorio/internal/types"
)


type ApplicationConfig struct {
	ApplicationName string
	WorldWidth      types.Real
	WorldHeight     types.Real
}

type Application struct {
	config ApplicationConfig
	world  *models.World

	clients map[*websocket.Conn]uuid.UUID	// player IDs

	mu sync.RWMutex
}

func NewApplication(config ApplicationConfig) *Application {
	world := models.NewWorld(models.WorldConfig{
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
		app.world.Update()

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

