package websocket

import (
	"sync"

	"github.com/gorilla/websocket"
	"github.com/kust1q/Zapp/main/internal/domain/entity"
)

type hub struct {
	clients    map[*client]bool
	users      map[int]*client
	broadcast  chan *entity.Notification
	register   chan *client
	unregister chan *client
	mu         sync.RWMutex
}

func NewHub() *hub {
	return &hub{
		clients:    make(map[*client]bool),
		users:      make(map[int]*client),
		broadcast:  make(chan *entity.Notification, 256),
		register:   make(chan *client),
		unregister: make(chan *client),
	}
}

func (h *hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.users[client.userID] = client
			h.mu.Unlock()

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				delete(h.users, client.userID)
				close(client.send)
			}
			h.mu.Unlock()

		case notification := <-h.broadcast:
			h.mu.RLock()
			if client, ok := h.users[notification.RecipientID]; ok {
				select {
				case client.send <- notification:
				default:
					close(client.send)
					delete(h.clients, client)
					delete(h.users, client.userID)
				}
			}
			h.mu.RUnlock()
		}
	}
}

func (h *hub) SendNotification(notification *entity.Notification) {
	h.broadcast <- notification
}

func (h *hub) HandleNewConnection(conn *websocket.Conn, userID int) {
	client := NewClient(h, conn, userID)
	h.register <- client

	go client.WritePump()
	go client.ReadPump()
}
