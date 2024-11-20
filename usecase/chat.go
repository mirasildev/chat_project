package usecase

import (
	"errors"

	"github.com/mirasildev/chat_task/domain"
)

type ChatRepository interface {
	Store(chat *domain.Chat) error
	GetByID(id int64) (*domain.Chat, error)
	GetAll(limit, page int64, userID string) (*domain.GetAllChats, error)
	AddMember(chatID int64, userID string) error
	RemoveMember(chatID int64, userID string) error
	GetUserChats(userID string) ([]domain.Chat, error)
	GetChatMembers(chatID int64) (*domain.GetAllChatMembersResponse, error)
}

type ChatService struct {
	chatRepo ChatRepository
}

func NewChatService(g ChatRepository) *ChatService {
	return &ChatService{
		chatRepo: g,
	}
}

func (c *ChatService) CreateChat(chat *domain.Chat) error {
	if chat.Name == "" {
		return errors.New("chat name is required")
	}

	return c.chatRepo.Store(chat)
}

func (c *ChatService) JoinChat(chatID int64, userID string) error {
	return c.chatRepo.AddMember(chatID, userID)
}

func (c *ChatService) LeaveChat(chatID int64, userID string) error {
	chat, err := c.chatRepo.GetByID(chatID)
	if err != nil {
		return err
	}

	// Check if user is the creator
	if chat.CreatedBy == userID {
		return errors.New("chat creator cannot leave the chat")
	}

	return c.chatRepo.RemoveMember(chatID, userID)
}

func (c *ChatService) GetChat(id int64) (*domain.Chat, error) {
	return c.chatRepo.GetByID(id)
}

func (c *ChatService) GetUserChats(userID string) ([]domain.Chat, error) {
	return c.chatRepo.GetUserChats(userID)
}

func (c *ChatService) GetChatMembers(chatID int64) (*domain.GetAllChatMembersResponse, error) {
	return c.chatRepo.GetChatMembers(chatID)
}
