package store

import (
	"context"
	"database/sql"
	"errors"
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

func (s *PostStore) GetByID(ctx context.Context, postID int) (*models.Post, error) {
	query := `
			SELECT id, user_id, title, content, tags, created_at, updated_at
			FROM posts
			WHERE id = $1
			`
	var post models.Post
	if err := s.db.QueryRowContext(ctx, query, postID).Scan(
		&post.ID,
		&post.UserID,
		&post.Title,
		&post.Content,
		pq.Array(&post.Tags),
		&post.CreatedAt,
		&post.UpdatedAt,
	); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}

	}

	return &post, nil
}

func (s *PostStore) Update(ctx context.Context, post *models.Post) error {
	query := `
			UPDATE posts
			SET title = $1, content = $2
			WHERE id = $3
			`
	_, err := s.db.ExecContext(ctx, query, post.Title, post.Content, post.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostStore) Delete(ctx context.Context, postID int) error {
	query := `
			DELETE FROM posts WHERE id= $1
			`
	res, err := s.db.ExecContext(ctx, query, postID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrNotFound
	}

	return nil
}
