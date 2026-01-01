package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/mabidoli/gravity-bff/internal/domain/model"
	"github.com/mabidoli/gravity-bff/internal/domain/repository"
)

// Ensure PgUserRepository implements the UserRepository interface.
var _ repository.UserRepository = (*PgUserRepository)(nil)

// PgUserRepository implements UserRepository using PostgreSQL.
type PgUserRepository struct {
	db *pgxpool.Pool
}

// NewPgUserRepository creates a new PostgreSQL user repository.
func NewPgUserRepository(db *pgxpool.Pool) *PgUserRepository {
	return &PgUserRepository{db: db}
}

// GetUserByID retrieves a user by their ID.
func (r *PgUserRepository) GetUserByID(ctx context.Context, userID string) (*model.User, error) {
	query := `
		SELECT id, name, email, avatar_url
		FROM users
		WHERE id = $1
	`

	var user model.User
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.AvatarURL,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// GetUsersByIDs retrieves multiple users by their IDs.
func (r *PgUserRepository) GetUsersByIDs(ctx context.Context, userIDs []string) ([]model.User, error) {
	if len(userIDs) == 0 {
		return []model.User{}, nil
	}

	query := `
		SELECT id, name, email, avatar_url
		FROM users
		WHERE id = ANY($1)
	`

	rows, err := r.db.Query(ctx, query, userIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.AvatarURL)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return users, nil
}
