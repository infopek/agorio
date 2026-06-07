package net

import (
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"github.com/infopek/agorio/server/game"
)

type Session struct {
	playerID   uuid.UUID
	conn       *websocket.Conn
	inputChan  chan<- game.PlayerEvent
	outputChan chan game.ServerEvent

	lastSplit time.Time // for rate-limit
	lastFeed  time.Time // for rate-limit
	lastMove  time.Time // for rate-limit

	closeOnce sync.Once
	onClose   func() // server cleanup
}

func NewSession(conn *websocket.Conn, inputChan chan<- game.PlayerEvent, onClose func()) *Session {
	conn.SetReadLimit(game.MaxPlayerMessageBytes)

	return &Session{
		playerID:   uuid.New(),
		conn:       conn,
		inputChan:  inputChan,
		outputChan: make(chan game.ServerEvent, 8),

		onClose: onClose,
	}
}

/** Session.close
 *
 * Centralized session close
 *
 */
func (s *Session) close() {
	s.closeOnce.Do(func() {
		s.conn.Close()
		s.onClose() // server callback
		s.inputChan <- game.PlayerDisconnectEvent{
			PlayerID: s.playerID,
		}
	})
}

/** Session.readLoop
 * Responsible for sending the disconnect, and
 *  forwarding player inputs to the game
 *
 * Stamps every message with the session's player ID
 */
func (s *Session) readLoop() {
	defer s.close()

	s.conn.SetReadDeadline(time.Now().Add(game.PongWait))
	s.conn.SetPongHandler(func(string) error {
		return s.conn.SetReadDeadline(time.Now().Add(game.PongWait))
	})

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
		case game.PlayerJoinEvent:
			m.PlayerID = s.playerID
			m.OutputChan = s.outputChan
			if !s.sendInput(m) {
				return
			}
		case game.PlayerMoveEvent:
			now := time.Now()
			if now.Sub(s.lastMove) < game.MoveCooldown {
				continue
			}
			s.lastMove = now

			m.PlayerID = s.playerID
			s.sendInput(m)
		case game.PlayerSplitEvent:
			now := time.Now()
			if now.Sub(s.lastSplit) < game.ActionCooldown {
				continue
			}
			s.lastSplit = now

			m.PlayerID = s.playerID
			s.sendInput(m)
		case game.PlayerFeedEvent:
			now := time.Now()
			if now.Sub(s.lastFeed) < game.ActionCooldown {
				continue
			}
			s.lastFeed = now

			m.PlayerID = s.playerID
			s.sendInput(m)
		}
	}
}

/** Session.writeLoop
 * Serializes each server message, and
 *  sends it towards the client (player)
 */
func (s *Session) writeLoop() {
	ticker := time.NewTicker(game.PingPeriod)
	defer ticker.Stop()
	defer s.close()

	for {
		select {
		case serverMsg, ok := <-s.outputChan:
			if !ok {
				return
			}

			data, err := toJSON(serverMsg)
			if err != nil {
				log.Printf("error while converting to json: %v\n", err)
				continue
			}

			s.conn.SetWriteDeadline(time.Now().Add(game.WriteWait))
			err = s.conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				log.Printf("error while sending data: %v\n", err)
				return
			}
		case <-ticker.C:
			// AFK connections are kicked after a while
			deadline := time.Now().Add(game.WriteWait)
			if err := s.conn.WriteControl(websocket.PingMessage, nil, deadline); err != nil {
				log.Printf("error while sending ping: %v\n", err)
				return
			}
		}
	}
}

/** Session.sendInput
 *
 * Helper for sending inputs towards server in
 *  an unblocking fashion
 *
 */
func (s *Session) sendInput(event game.PlayerEvent) bool {
	select {
	case s.inputChan <- event:
		return true
	default:
		return false
	}
}
