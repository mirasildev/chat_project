package domain

import "time"

type Message struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	FileURL   string    `json:"file_url,omitempty"`
	UserID    string    `json:"user_id"`
	ChatID    int64     `json:"chat_id"`
	User      *User     `json:"user,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
