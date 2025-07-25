package store

import (
	"context"
	"database/sql"
	"github.com/ye-khaing-win/social_go/internal/models"
)

type UserStore struct {
	db *sql.DB
}

func (s *UserStore) Create(ctx context.Context, user *models.User) error {
	query := `
INSERT INTO users
(username, password, email) VALUES
($1, $2, $3)
RETURNING id, created_at
`
	if err := s.db.QueryRowContext(
		ctx,
		query,
		user.Username,
		user.Password,
		user.Email).Scan(
		&user.ID,
		&user.CreatedAt,
	); err != nil {
		return err
	}
	return nil
}
