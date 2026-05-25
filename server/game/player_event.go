package game

import (
	"github.com/google/uuid"
)

type PlayerEvent interface {
	playerEvent() // to seal the interface
}

type PlayerJoinEvent struct {
	PlayerID   uuid.UUID
	Name       string
	OutputChan chan<- ServerEvent
}

type PlayerDisconnectEvent struct {
	PlayerID uuid.UUID
}

type PlayerMoveEvent struct {
	PlayerID uuid.UUID
	Target   Vec2
}

type PlayerSplitEvent struct {
	PlayerID uuid.UUID
}

type PlayerFeedEvent struct {
	PlayerID uuid.UUID
}

func (msg PlayerJoinEvent) playerEvent()       {}
func (msg PlayerDisconnectEvent) playerEvent() {}
func (msg PlayerMoveEvent) playerEvent()       {}
func (msg PlayerSplitEvent) playerEvent()      {}
func (msg PlayerFeedEvent) playerEvent()       {}
