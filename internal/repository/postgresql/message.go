package postgresql

import (
	"database/sql"

	"github.com/mirasildev/chat_task/domain"
)

type MessageRepository struct {
	db *sql.DB
}

func NewMessageRepository(db *sql.DB) *MessageRepository {
	return &MessageRepository{db}
}

func (m *MessageRepository) Store(msg *domain.Message) error {
	query := `INSERT INTO messages (id, content, file_url, user_id, chat_id)
              VALUES ($1, $2, $3, $4, $5)`

	_, err := m.db.Exec(query,
		msg.ID,
		msg.Content,
		msg.FileURL,
		msg.UserID,
		msg.ChatID,
	)

	return err
}

func (p *MessageRepository) GetChatMessages(chatID int64, limit, offset int) ([]domain.Message, error) {
    query := `
        SELECT m.id, m.content, m.file_url, m.user_id, m.chat_id, m.created_at,
               u.username
        FROM messages m
        JOIN users u ON m.user_id = u.id
        WHERE m.chat_id = $1
        ORDER BY m.created_at DESC
        LIMIT $2 OFFSET $3
    `
    
    rows, err := p.db.Query(query, chatID, limit, offset)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var messages []domain.Message
    for rows.Next() {
        var msg domain.Message
        var user domain.User
        
        err := rows.Scan(
            &msg.ID,
            &msg.Content,
            &msg.FileURL,
            &msg.UserID,
            &msg.ChatID,
            &msg.CreatedAt,
            &user.Username,
        )
        if err != nil {
            return nil, err
        }
        
        msg.User = &user
        messages = append(messages, msg)
    }
    
    return messages, nil
}
