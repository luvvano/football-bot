package services

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/luvvano/football-bot/backend/internal/models"
)

type UserService struct {
	pool *pgxpool.Pool
}

func NewUserService(pool *pgxpool.Pool) *UserService {
	return &UserService{pool: pool}
}

func (s *UserService) GetByTelegramID(ctx context.Context, telegramID int64) (*models.User, error) {
	query := `
		SELECT id, telegram_id, username, first_name, last_name, is_global_admin, created_at, last_active_at
		FROM users WHERE telegram_id = $1
	`

	var u models.User
	err := s.pool.QueryRow(ctx, query, telegramID).Scan(
		&u.ID, &u.TelegramID, &u.Username, &u.FirstName, &u.LastName,
		&u.IsGlobalAdmin, &u.CreatedAt, &u.LastActiveAt,
	)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *UserService) GetByID(ctx context.Context, id int64) (*models.User, error) {
	query := `
		SELECT id, telegram_id, username, first_name, last_name, is_global_admin, created_at, last_active_at
		FROM users WHERE id = $1
	`

	var u models.User
	err := s.pool.QueryRow(ctx, query, id).Scan(
		&u.ID, &u.TelegramID, &u.Username, &u.FirstName, &u.LastName,
		&u.IsGlobalAdmin, &u.CreatedAt, &u.LastActiveAt,
	)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *UserService) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	query := `
		SELECT id, telegram_id, username, first_name, last_name, is_global_admin, created_at, last_active_at
		FROM users WHERE username = $1
	`

	var u models.User
	err := s.pool.QueryRow(ctx, query, username).Scan(
		&u.ID, &u.TelegramID, &u.Username, &u.FirstName, &u.LastName,
		&u.IsGlobalAdmin, &u.CreatedAt, &u.LastActiveAt,
	)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *UserService) Create(ctx context.Context, input models.CreateUserInput) (*models.User, error) {
	query := `
		INSERT INTO users (telegram_id, username, first_name, last_name)
		VALUES ($1, $2, $3, $4)
		RETURNING id, telegram_id, username, first_name, last_name, is_global_admin, created_at, last_active_at
	`

	var u models.User
	err := s.pool.QueryRow(ctx, query, input.TelegramID, input.Username, input.FirstName, input.LastName).Scan(
		&u.ID, &u.TelegramID, &u.Username, &u.FirstName, &u.LastName,
		&u.IsGlobalAdmin, &u.CreatedAt, &u.LastActiveAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	return &u, nil
}

func (s *UserService) Upsert(ctx context.Context, input models.CreateUserInput) (*models.User, error) {
	query := `
		INSERT INTO users (telegram_id, username, first_name, last_name)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (telegram_id) DO UPDATE SET
			username = EXCLUDED.username,
			first_name = EXCLUDED.first_name,
			last_name = EXCLUDED.last_name,
			last_active_at = CURRENT_TIMESTAMP
		RETURNING id, telegram_id, username, first_name, last_name, is_global_admin, created_at, last_active_at
	`

	var u models.User
	err := s.pool.QueryRow(ctx, query, input.TelegramID, input.Username, input.FirstName, input.LastName).Scan(
		&u.ID, &u.TelegramID, &u.Username, &u.FirstName, &u.LastName,
		&u.IsGlobalAdmin, &u.CreatedAt, &u.LastActiveAt,
	)
	if err != nil {
		return nil, fmt.Errorf("upsert user: %w", err)
	}
	return &u, nil
}

func (s *UserService) UpdateActivity(ctx context.Context, userID int64) error {
	query := `UPDATE users SET last_active_at = CURRENT_TIMESTAMP WHERE id = $1`
	_, err := s.pool.Exec(ctx, query, userID)
	return err
}
