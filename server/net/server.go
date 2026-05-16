package net

import (
	"net/http"

	"github.com/infopek/agorio/server/game/world"
)

func (world *World) WebsocketHandler(w http.ResponseWriter, r http.Request)
