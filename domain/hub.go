package domain

import "sync"

type Hub struct {
	sync.RWMutex

	Clients map[*Client]bool

	Broadcast  chan *Message
	Register   chan *Client
	Unregister chan *Client

	Messages []*Message
}
