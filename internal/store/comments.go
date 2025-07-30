package store

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/ye-khaing-win/social_go/internal/models"
)

type CommentStore struct {
	db *sql.DB
}

func (s *CommentStore) GetByPostID(ctx context.Context, postID int64) ([]*models.Comment, error) {

	fmt.Println("PostID: ", postID)
	query := `
			SELECT c.id, c.post_id, c.user_id, c.content, c.created_at, u.id, u.username
			FROM comments AS c
			INNER JOIN users AS u
			ON c.user_id = u.id
			WHERE c.post_id = $1
			ORDER BY c.created_at DESC;
			`
	rows, err := s.db.QueryContext(ctx, query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*models.Comment
	for rows.Next() {
		var comment models.Comment
		if err := rows.Scan(
			&comment.ID,
			&comment.PostID,
			&comment.UserID,
			&comment.Content,
			&comment.CreatedAt,
			&comment.User.ID,
			&comment.User.Username,
		); err != nil {
			return nil, err
		}
		comments = append(comments, &comment)
	}

	return comments, err
}

func (s *CommentStore) Create(ctx context.Context, comment *models.Comment) error {
	query := `
	INSERT INTO comments
	(post_id, user_id, content) VALUES 
	($1, $2, $3)
	RETURNING id, created_at
`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	if err := s.db.QueryRowContext(
		ctx,
		query,
		comment.PostID,
		comment.UserID,
		comment.Content,
	).Scan(&comment.ID, &comment.CreatedAt); err != nil {
		return err
	}

	return nil
}
