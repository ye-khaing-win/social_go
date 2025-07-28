package store

import (
	"context"
	"database/sql"
	"errors"
	"github.com/ye-khaing-win/social_go/internal/models"
)

var (
	ErrNotFound = errors.New("resource not found")
)

type Storage struct {
	Posts interface {
		Create(ctx context.Context, post *models.Post) error
		GetByID(ctx context.Context, postID int) (*models.Post, error)
		Delete(ctx context.Context, postID int) error
		Update(ctx context.Context, post *models.Post) error
	}
	Users interface {
		Create(ctx context.Context, user *models.User) error
	}
	Comments interface {
		GetByPostID(ctx context.Context, postID int64) ([]*models.Comment, error)
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts:    &PostStore{db},
		Users:    &UserStore{db},
		Comments: &CommentStore{db},
	}
}
