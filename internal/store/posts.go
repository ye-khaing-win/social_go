package store

import (
	"context"
	"database/sql"
	"github.com/lib/pq"
	"github.com/ye-khaing-win/social_go/internal/models"
)

type PostStore struct {
	db *sql.DB
}

func (s *PostStore) Create(ctx context.Context, post *models.Post) error {
	query := `
INSERT INTO posts 
(content, title, user_id, tags) VALUES 
($1, $2, $3, $4) RETURNING id, created_at, updated_at
`

	err := s.db.QueryRowContext(
		ctx,
		query,
		post.Content,
		post.Title,
		post.UserID,
		pq.Array(post.Tags),
	).Scan(
		&post.ID,
		&post.CreatedAt,
		&post.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}
