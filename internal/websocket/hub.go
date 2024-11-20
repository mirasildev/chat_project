package websocket

import (
	"encoding/json"
	"fmt"

	"github.com/mirasildev/chat_task/domain"
	"github.com/mirasildev/chat_task/usecase"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[string]*Client

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	ChatService    usecase.ChatService
	MessageService usecase.MessageService
	AuthService    usecase.AuthService
}

func NewHub(chatService usecase.ChatService, messageService usecase.MessageService) *Hub {
	return &Hub{
		broadcast:      make(chan []byte),
		register:       make(chan *Client),
		unregister:     make(chan *Client),
		clients:        make(map[string]*Client),
		ChatService:    chatService,
		MessageService: messageService,
	}
}

type Message struct {
	UserID   string `json:"user_id"`
	ChatType string `json:"chat_type"`
	ChatID   int64  `json:"chat_id"`
	Content  string `json:"content"`
	FileURL  string `json:"file_url"`
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client.userID] = client
			fmt.Println("New client connected", client.userID)
		case client := <-h.unregister:
			if _, ok := h.clients[client.userID]; ok {
				delete(h.clients, client.userID)
				close(client.send)
			}
			fmt.Println("Client disconnected", client.userID)
		case data := <-h.broadcast:
			var message Message
			err := json.Unmarshal(data, &message)
			if err != nil {
				fmt.Println(err)
				continue
			}

			// Create message
			err = h.MessageService.CreateMessage(&domain.Message{
				Content: message.Content,
				FileURL: message.FileURL,
				ChatID:  message.ChatID,
				UserID:  message.UserID,
			})
			if err != nil {
				fmt.Println(err)
				continue
			}

			// Get Chat members
			result, err := h.ChatService.GetChatMembers(message.ChatID)
			if err != nil {
				fmt.Println(err)
				continue
			}

			for _, user := range result.Data {
				if message.UserID == user.User.ID {
					continue
				}

				client, ok := h.clients[user.User.ID]
				if ok {
					select {
					case client.send <- data:
					default:
						close(client.send)
						delete(h.clients, client.userID)
					}
				}
			}
		}
	}
}
