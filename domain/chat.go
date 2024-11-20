package domain

import "time"

type Chat struct {
	ID        int64        `json:"id"`
	Type      string       `json:"type"`
	Name      string       `json:"name"`
	CreatedBy string       `json:"created_by"`
	CreatedAt time.Time    `json:"created_at"`
	Members   []ChatMember `json:"members,omitempty"`
}

type ChatMember struct {
	User User   `json:"user"`
	Role string `json:"role"` // admin or member
}

type GetAllChats struct {
	Data  []*Chat
	Count int64
}

type GetAllChatMembersResponse struct {
	Data  []*ChatMember
	Count int64
}
