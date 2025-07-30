package store

import (
	"context"
	"database/sql"
)

type FollowerStore struct {
	db *sql.DB
}

func (s *FollowerStore) Follow(ctx context.Context, followerID, userID int64) error {
	query := `
		INSERT INTO followers
		(user_id, follower_id) VALUES
		($1, $2)
		`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, userID, followerID)
	if err != nil {
		return err
	}

	return nil
}

func (s *FollowerStore) Unfollow(ctx context.Context, followerID, userID int64) error {
	query := `
		DELETE FROM followers
		WHERE user_id = $1 AND follower_id = $2
		`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, userID, followerID)
	if err != nil {
		return err
	}

	return nil
}
