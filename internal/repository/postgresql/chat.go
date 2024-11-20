package postgresql

import (
	"database/sql"
	"fmt"

	"github.com/mirasildev/chat_task/domain"
)

type ChatRepository struct {
	db *sql.DB
}

func NewChatRepository(db *sql.DB) *ChatRepository {
	return &ChatRepository{db}
}

func (p *ChatRepository) Store(chat *domain.Chat) error {
	query := `
        INSERT INTO chats (name, created_by, type)
        VALUES ($1, $2, $3)
    `

	// Begin transaction
	tx, err := p.db.Begin()
	if err != nil {
		return err
	}

	// Insert chat
	_, err = tx.Exec(query,
		chat.Name,
		chat.CreatedBy,
		chat.Type,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Add creator as admin member
	memberQuery := `
        INSERT INTO chat_members (chat_id, user_id, role)
        VALUES ($1, $2, 'admin')
    `
	_, err = tx.Exec(memberQuery, chat.ID, chat.CreatedBy)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (p *ChatRepository) AddMember(chatID int64, userID, role string) error {
	query := `
        INSERT INTO chat_members (chat_id, user_id, role)
        VALUES ($1, $2, $3)
    `
	_, err := p.db.Exec(query, chatID, userID, role)
	return err
}

func (p *ChatRepository) RemoveMember(chatID int64, userID string) error {
	query := `DELETE FROM chat_members WHERE chat_id = $1 AND user_id = $2`
	_, err := p.db.Exec(query, chatID, userID)
	return err
}

func (p *ChatRepository) GetByID(id string) (*domain.Chat, error) {
	query := `
        SELECT id, name, created_by, type, created_at
        FROM chats WHERE id = $1
    `

	group := &domain.Chat{}
	err := p.db.QueryRow(query, id).Scan(
		&group.ID,
		&group.Name,
		&group.Type,
		&group.CreatedBy,
		&group.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	membersQuery := `
        SELECT u.id, u.username, gm.role
        FROM users u
        JOIN chat_members cm ON u.id = cm.user_id
        WHERE cm.chat_id = $1
    `

	rows, err := p.db.Query(membersQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var member domain.ChatMember
		err := rows.Scan(&member.User.ID, &member.User.Username,
			&member.Role)
		if err != nil {
			return nil, err
		}
		group.Members = append(group.Members, member)
	}

	return group, nil
}

func (pr *ChatRepository) GetAll(limit, page int64, userID string) (*domain.GetAllChats, error) {
	result := domain.GetAllChats{
		Data: make([]*domain.Chat, 0),
	}

	offset := (page - 1) * limit
	limit2 := fmt.Sprintf(" LIMIT %d OFFSET %d ", limit, offset)

	filter := ""
	if userID != "" {
		filter = fmt.Sprintf(" WHERE user_id = %s ", userID)
	}

	query := `
		SELECT
			id,
			name,
			created_by,
			type
		FROM chats
	` + filter + limit2

	rows, err := pr.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var chat domain.Chat

		err := rows.Scan(
			&chat.ID,
			&chat.Name,
			&chat.CreatedBy,
			&chat.Type,
		)
		if err != nil {
			return nil, err
		}

		if err != nil {
			return nil, err
		}

		result.Data = append(result.Data, &chat)
	}

	queryCount := `SELECT count(1) FROM chats` + filter
	err = pr.db.QueryRow(queryCount).Scan(&result.Count)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (p *ChatRepository) GetUserChats(userID string) ([]domain.Chat, error) {
	query := `
        SELECT
	        c.id,
	        c.name,
	        c.created_by,
	        c.type,
	        c.created_at,
	        gm.role
        FROM chats c
        JOIN chat_members cm ON c.id = cm.chat_id
        WHERE cm.user_id = $1
        ORDER BY c.created_at DESC
    `

	rows, err := p.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chats []domain.Chat
	for rows.Next() {
		var group domain.Chat
		var role string
		err := rows.Scan(
			&group.ID,
			&group.Name,
			&group.CreatedBy,
			&group.Type,
			&group.CreatedAt,
			&role,
		)
		if err != nil {
			return nil, err
		}
		chats = append(chats, group)
	}
	return chats, nil
}

func (ur *ChatRepository) GetChatMembers(chatID int64) (*domain.GetAllChatMembersResponse, error) {
	result := domain.GetAllChatMembersResponse{
		Data: make([]*domain.ChatMember, 0),
	}

	query := `
		SELECT
			u.id,
			u.email,
			u.username,
			u.created_at,
			cm.role
		FROM users u
		INNER JOIN chat_members cm ON cm.user_id=u.id
		WHERE cm.chat_id=$1
	`

	row := ur.db.QueryRow(query, chatID)

	var (
		u = domain.ChatMember{
			User: domain.User{},
		}
	)

	err := row.Scan(
		&u.User.ID,
		&u.User.Email,
		&u.User.Username,
		&u.User.CreatedAt,
		&u.Role,
	)
	if err != nil {
		return nil, err
	}

	result.Data = append(result.Data, &u)

	queryCount := `
		SELECT count(1) FROM users u
		INNER JOIN chat_members cm ON cm.user_id=u.id
		WHERE cm.chat_id=$1
	`
	err = ur.db.QueryRow(queryCount, chatID).Scan(&result.Count)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
