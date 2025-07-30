package store

import (
	"context"
	"database/sql"
	"errors"
	"github.com/ye-khaing-win/social_go/internal/models"
	"time"
)

var (
	ErrNotFound          = errors.New("resource not found")
	QueryTimeoutDuration = 5 * time.Second
)

type Storage struct {
	Posts interface {
		Create(ctx context.Context, post *models.Post) error
		GetByID(ctx context.Context, postID int64) (*models.Post, error)
		Delete(ctx context.Context, postID int64) error
		Update(ctx context.Context, post *models.Post) error
		GetUserFeed(ctx context.Context, userID int64) ([]*models.PostWithMetadata, error)
	}
	Users interface {
		Create(ctx context.Context, tx *sql.Tx, user *models.User) error
		GetByID(ctx context.Context, userID int64) (*models.User, error)
		CreateAndInvite(ctx context.Context, user *models.User, token string, expiry time.Duration) error
	}
	Comments interface {
		GetByPostID(ctx context.Context, postID int64) ([]*models.Comment, error)
		Create(ctx context.Context, comment *models.Comment) error
	}
	Followers interface {
		Follow(ctx context.Context, followerID, userID int64) error
		Unfollow(ctx context.Context, followerID, userID int64) error
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts:     &PostStore{db},
		Users:     &UserStore{db},
		Comments:  &CommentStore{db},
		Followers: &FollowerStore{db: db},
	}
}

func withTx(db *sql.DB, ctx context.Context, fn func(tx *sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}
