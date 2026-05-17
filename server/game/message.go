package game

import (
	"github.com/google/uuid"
)

type PlayerMessage interface {
	playerMessage() // to seal the interface
}

type JoinMessage struct {
	PlayerID uuid.UUID
	Name     string
	OutputChan chan<- ServerMessage
}

type MoveMessage struct {
	PlayerID uuid.UUID
	Target   Vec2
}

type SplitMessage struct {
	PlayerID uuid.UUID
}

type FeedMessage struct {
	PlayerID uuid.UUID
}

func (msg JoinMessage) playerMessage()  {}
func (msg MoveMessage) playerMessage()  {}
func (msg SplitMessage) playerMessage() {}
func (msg FeedMessage) playerMessage()  {}
