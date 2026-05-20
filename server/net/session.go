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

/** readLoop
 * Responsible for sending the disconnect, and
 *  forwarding player inputs to the game
 *
 * Stamps every message with the session's player ID
 */
func (s *Session) readLoop() {
	defer s.conn.Close()
	defer func() {
		s.inputChan <- game.DisconnectMessage{
			PlayerID: s.playerID,
		}
	}()

	for {
		_, msg, err := s.conn.ReadMessage()
		if err != nil {
			log.Printf("error reading message: %v\n", err)
			return
		}

		playerMsg, err := parsePlayerMessage(msg)
		if err != nil {
			log.Printf("player msg parse error: %v\n", err)
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

/** writeLoop
 * Serializes each server message, and
 *  sends it towards the client (player)
 */
func (s *Session) writeLoop() {
	for serverMsg := range s.outputChan {
		data, err := toJSON(serverMsg)
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
