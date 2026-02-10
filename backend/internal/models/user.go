package models

import "time"

type User struct {
	ID            int64     `json:"id"`
	TelegramID    int64     `json:"telegram_id"`
	Username      *string   `json:"username,omitempty"`
	FirstName     *string   `json:"first_name,omitempty"`
	LastName      *string   `json:"last_name,omitempty"`
	IsGlobalAdmin bool      `json:"is_global_admin"`
	CreatedAt     time.Time `json:"created_at"`
	LastActiveAt  time.Time `json:"last_active_at"`
}

type CreateUserInput struct {
	TelegramID int64
	Username   *string
	FirstName  *string
	LastName   *string
}
