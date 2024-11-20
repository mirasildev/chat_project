package postgresql

import (
	"database/sql"
	"fmt"

	"github.com/mirasildev/chat_task/domain"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (ur *UserRepository) Store(user *domain.User) (*domain.User, error) {
	query := `
		INSERT INTO users(
			email,
			password,
			username
		) VALUES($1, $2, $3)
		RETURNING id, created_at
	`

	row := ur.db.QueryRow(
		query,
		user.Email,
		user.Password,
		user.Username,
	)

	err := row.Scan(
		&user.ID,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (ur *UserRepository) Get(id int64) (*domain.User, error) {
	var (
		result domain.User
	)

	query := `
		SELECT
			id,
			email,
			password,
			username,
			created_at
		FROM users
		WHERE id=$1
	`

	row := ur.db.QueryRow(query, id)
	err := row.Scan(
		&result.ID,
		&result.Email,
		&result.Password,
		&result.Username,
		&result.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (ur *UserRepository) GetAll(limit, page int64, search string) (*domain.GetAllUsersReponse, error) {
	result := domain.GetAllUsersReponse{
		Data: make([]*domain.User, 0),
	}

	offset := (page - 1) * limit
	limit2 := fmt.Sprintf(" LIMIT %d OFFSET %d ", limit, offset)

	filter := ""
	if search != "" {
		str := "%" + search + "%"
		filter += fmt.Sprintf(`
			WHERE email ILIKE '%s' OR username ILIKE '%s'`,
			str, str,
		)
	}

	query := `
		SELECT
			id,
			email,
			password,
			username,
			created_at
		FROM users
		` + filter + `
		ORDER BY created_at desc
		` + limit2

	rows, err := ur.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var (
			u domain.User
		)

		err := rows.Scan(
			&u.ID,
			&u.Email,
			&u.Password,
			&u.Username,
			&u.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		result.Data = append(result.Data, &u)
	}

	queryCount := `SELECT count(1) FROM users ` + filter
	err = ur.db.QueryRow(queryCount).Scan(&result.Count)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (ur *UserRepository) GetByEmail(email string) (*domain.User, error) {
	var (
		result domain.User
	)

	query := `
		SELECT
			id,
			email,
			password,
			username,
			created_at
		FROM users
		WHERE email=$1
	`

	row := ur.db.QueryRow(query, email)
	err := row.Scan(
		&result.ID,
		&result.Email,
		&result.Password,
		&result.Username,
		&result.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (ur *UserRepository) UpdatePassword(userID, password string) error {
	query := `UPDATE users SET password=$1 WHERE id=$2`

	_, err := ur.db.Exec(query, password, userID)
	if err != nil {
		return err
	}

	return nil
}

func (ur *UserRepository) Delete(id string) error {
	query := `DELETE FROM users WHERE id=$1`

	result, err := ur.db.Exec(query, id)
	if err != nil {
		return err
	}

	if count, _ := result.RowsAffected(); count == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (ur *UserRepository) Update(user *domain.User) (*domain.User, error) {
	query := `
		UPDATE users SET
			username=$1
		WHERE id=$2
		RETURNING
			email,
			created_at
	`

	err := ur.db.QueryRow(
		query,
		user.Username,
		user.ID,
	).Scan(
		&user.Email,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}
