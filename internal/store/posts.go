package store

import (
	"context"
	"database/sql"
	"errors"
	"github.com/lib/pq"
	"github.com/ye-khaing-win/social_go/internal/models"
	q "github.com/ye-khaing-win/social_go/internal/query"
)

type PostStore struct {
	db *sql.DB
}

func (s *PostStore) GetUserFeed(ctx context.Context, userID int64) ([]*models.PostWithMetadata, error) {
	pg := q.GetPgFromContext(ctx)

	query := `
		SELECT p.id, p.user_id, p.title, p.content, p.created_at, p.version, p.tags,
		u.username,
		COUNT(c.id) AS comments_count
		FROM posts p
		LEFT JOIN comments c ON p.id = c.post_id
		LEFT JOIN users u ON p.user_id = u.id
		INNER JOIN followers f ON p.user_id = f.follower_id OR p.user_id = $1
		WHERE p.user_id = $1 OR f.user_id = $1
		GROUP BY p.id, p.created_at, u.username
		ORDER BY p.created_at DESC
		LIMIT $2 OFFSET $3
		`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, userID, pg.Limit, pg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feed []*models.PostWithMetadata
	for rows.Next() {
		var post models.PostWithMetadata

		if err := rows.Scan(
			&post.ID,
			&post.UserID,
			&post.Title,
			&post.Content,
			&post.CreatedAt,
			&post.Version,
			pq.Array(&post.Tags),
			&post.User.Username,
			&post.CommentCount,
		); err != nil {
			return nil, err
		}

		feed = append(feed, &post)
	}

	return feed, nil
}

func (s *PostStore) Create(ctx context.Context, post *models.Post) error {
	query := `
			INSERT INTO posts 
			(content, title, user_id, tags) VALUES 
			($1, $2, $3, $4) RETURNING id, created_at, updated_at
			`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

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

func (s *PostStore) GetByID(ctx context.Context, postID int64) (*models.Post, error) {
	query := `
			SELECT id, user_id, title, content, tags, created_at, updated_at, version
			FROM posts
			WHERE id = $1
			`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var post models.Post
	if err := s.db.QueryRowContext(ctx, query, postID).Scan(
		&post.ID,
		&post.UserID,
		&post.Title,
		&post.Content,
		pq.Array(&post.Tags),
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.Version,
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
			SET title = $1, content = $2, version = version + 1
			WHERE id = $3 AND version = $4
			RETURNING version
			`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		post.Title,
		post.Content,
		post.ID,
		post.Version,
	).Scan(&post.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrNotFound
		default:
			return err
		}

	}

	return nil
}

func (s *PostStore) Delete(ctx context.Context, postID int64) error {
	query := `
			DELETE FROM posts WHERE id= $1
			`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

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
