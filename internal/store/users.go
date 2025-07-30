package store

import (
	"context"
	"database/sql"
	"errors"
	"github.com/lib/pq"
	"github.com/ye-khaing-win/social_go/internal/models"
	"time"
)

var ErrDuplicateEmail = errors.New("user with this email already exists")

type UserStore struct {
	db *sql.DB
}

func (s *UserStore) Create(ctx context.Context, tx *sql.Tx, user *models.User) error {
	query := `
			INSERT INTO users
			(username, password, email) VALUES
			($1, $2, $3)
			RETURNING id, created_at
			`
	if err := tx.QueryRowContext(
		ctx,
		query,
		user.Username,
		user.Password,
		user.Email).Scan(
		&user.ID,
		&user.CreatedAt,
	); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23505" {
				return ErrDuplicateEmail
			}
		}
		return err
	}
	return nil
}

func (s *UserStore) GetByID(ctx context.Context, userID int64) (*models.User, error) {
	query := `
		SELECT id, username, email, password, created_at
		FROM users
		WHERE id = $1
		`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var user models.User
	if err := s.db.QueryRowContext(
		ctx,
		query,
		userID,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}

	}

	return &user, nil
}

func (s *UserStore) CreateAndInvite(ctx context.Context, user *models.User, token string, exp time.Duration) error {

	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		if err := s.Create(ctx, tx, user); err != nil {
			return err
		}

		if err := s.invite(ctx, tx, token, exp, user.ID); err != nil {
			return err
		}

		return nil
	})
}

func (s *UserStore) invite(ctx context.Context, tx *sql.Tx, token string, exp time.Duration, userID int64) error {
	query := `
		INSERT INTO user_invitations
		(token, user_id, expiry) VALUES
		($1, $2, $3)
		`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, token, userID, time.Now().Add(exp))
	if err != nil {
		return err
	}

	return nil
}
