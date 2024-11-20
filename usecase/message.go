package usecase

import (
	"errors"

	"github.com/google/uuid"
	// "github.com/gorilla/websocket"
	"github.com/mirasildev/chat_task/domain"
)

type MessageRepository interface {
	Store(message *domain.Message) error
	GetChatMessages(groupID int64, limit, offset int) ([]domain.Message, error)
	DeleteMessage(messageID, userID string) error
}

type MessageService struct {
	messageRepo MessageRepository
}

func NewMessageService(g MessageRepository) *MessageService {
	return &MessageService{
		messageRepo: g,
	}
}

func (m *MessageService) CreateMessage(msg *domain.Message) error {
	if msg.Content == "" && msg.FileURL == "" {
		return errors.New("message must have content or file")
	}

	msg.ID = uuid.New().String()
	if err := m.messageRepo.Store(msg); err != nil {
		return err
	}

	return nil
}

func (m *MessageService) GetChatMessages(chatID int64, limit, offset int) ([]domain.Message, error) {
	if limit <= 0 {
		limit = 50 // default limit
	}
	return m.messageRepo.GetChatMessages(chatID, limit, offset)
}

func (m *MessageService) DeleteMessage(messageID, userID string) error {
	return m.messageRepo.DeleteMessage(messageID, userID)
}
