package net

import (
	"log"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"github.com/infopek/agorio/server/game"
)

type Session struct {
	playerID   uuid.UUID
	conn       *websocket.Conn
	inputChan  chan<- game.PlayerMessage
	outputChan chan game.ServerMessage
}

func NewSession(conn *websocket.Conn, inputChan chan<- game.PlayerMessage) *Session {
	return &Session{
		playerID:   uuid.New(),
		conn:       conn,
		inputChan:  inputChan,
		outputChan: make(chan game.ServerMessage),
	}
}

func (s *Session) readLoop() {
	defer close(s.inputChan)
	for {
		_, msg, err := s.conn.ReadMessage()
		if err != nil {
			return
		}

		playerMsg, err := parsePlayerMessage(msg)
		if err != nil {
			return
		}

		switch m := playerMsg.(type) {
		case game.JoinMessage:
			m.PlayerID = s.playerID
			m.OutputChan = s.outputChan
			s.inputChan <- m
		case game.MoveMessage:
			m.PlayerID = s.playerID
			s.inputChan <- m
		case game.SplitMessage:
			m.PlayerID = s.playerID
			s.inputChan <- m
		case game.FeedMessage:
			m.PlayerID = s.playerID
			s.inputChan <- m
		}
	}
}

func (s *Session) writeLoop() {
	defer close(s.outputChan)
	for state := range s.outputChan {
		data, err := toJSON(state)
		if err != nil {
			log.Printf("error while converting to json: %v\n", err)
			continue
		}
		err = s.conn.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			log.Printf("error while sending data: %v\n", err)
			continue
		}
	}
}
