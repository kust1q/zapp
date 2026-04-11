package websocket

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

type service struct {
	hub      hubProvider
	upgrader websocket.Upgrader
}

func NewWebSocketService(hub hubProvider) *service {
	return &service{
		hub: hub,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
	}
}

func (s *service) HandleConnection(w http.ResponseWriter, r *http.Request, userID int) error {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return fmt.Errorf("failed to do websocket upgrade: %w", err)
	}
	s.hub.HandleNewConnection(conn, userID)

	return nil
}
