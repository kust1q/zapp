package websocket

import "github.com/gorilla/websocket"

type (
	hubProvider interface {
		HandleNewConnection(conn *websocket.Conn, userID int)
	}
)
