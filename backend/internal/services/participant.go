package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/luvvano/football-bot/backend/internal/models"
)

var (
	ErrAlreadyRegistered = errors.New("already registered for this event")
	ErrEventFull         = errors.New("event is full")
	ErrNotParticipant    = errors.New("not a participant of this event")
	ErrInvalidToken      = errors.New("invalid or expired claim token")
)

type ParticipantService struct {
	pool *pgxpool.Pool
}

func NewParticipantService(pool *pgxpool.Pool) *ParticipantService {
	return &ParticipantService{pool: pool}
}

func (s *ParticipantService) GetByEvent(ctx context.Context, eventID int) ([]models.Participant, error) {
	query := `
		SELECT p.id, p.event_id, p.user_id, p.guest_name, p.guest_claim_token, p.added_by, p.status, p.registered_at,
		       u.id, u.telegram_id, u.username, u.first_name, u.last_name
		FROM participants p
		LEFT JOIN users u ON p.user_id = u.id
		WHERE p.event_id = $1 AND p.status = 'confirmed'
		ORDER BY p.registered_at
	`

	rows, err := s.pool.Query(ctx, query, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var participants []models.Participant
	for rows.Next() {
		var p models.Participant
		var userID, userTelegramID *int64
		var userName, userFirstName, userLastName *string

		if err := rows.Scan(
			&p.ID, &p.EventID, &p.UserID, &p.GuestName, &p.GuestClaimToken, &p.AddedBy, &p.Status, &p.RegisteredAt,
			&userID, &userTelegramID, &userName, &userFirstName, &userLastName,
		); err != nil {
			return nil, err
		}

		if userID != nil {
			p.User = &models.User{
				ID:         *userID,
				TelegramID: *userTelegramID,
				Username:   userName,
				FirstName:  userFirstName,
				LastName:   userLastName,
			}
		}

		participants = append(participants, p)
	}

	return participants, rows.Err()
}

func (s *ParticipantService) GetByEventAndUser(ctx context.Context, eventID int, userID int64) (*models.Participant, error) {
	query := `
		SELECT id, event_id, user_id, guest_name, guest_claim_token, added_by, status, registered_at
		FROM participants
		WHERE event_id = $1 AND user_id = $2 AND status = 'confirmed'
	`

	var p models.Participant
	err := s.pool.QueryRow(ctx, query, eventID, userID).Scan(
		&p.ID, &p.EventID, &p.UserID, &p.GuestName, &p.GuestClaimToken, &p.AddedBy, &p.Status, &p.RegisteredAt,
	)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (s *ParticipantService) GetByEventAndGuest(ctx context.Context, eventID int, guestName string) (*models.Participant, error) {
	query := `
		SELECT id, event_id, user_id, guest_name, guest_claim_token, added_by, status, registered_at
		FROM participants
		WHERE event_id = $1 AND guest_name = $2 AND status = 'confirmed'
	`

	var p models.Participant
	err := s.pool.QueryRow(ctx, query, eventID, guestName).Scan(
		&p.ID, &p.EventID, &p.UserID, &p.GuestName, &p.GuestClaimToken, &p.AddedBy, &p.Status, &p.RegisteredAt,
	)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (s *ParticipantService) CountConfirmed(ctx context.Context, eventID int) (int, error) {
	query := `SELECT COUNT(*) FROM participants WHERE event_id = $1 AND status = 'confirmed'`
	var count int
	err := s.pool.QueryRow(ctx, query, eventID).Scan(&count)
	return count, err
}

func (s *ParticipantService) AddUser(ctx context.Context, eventID int, userID int64, addedBy int64) (*models.Participant, error) {
	// Check if already registered
	existing, err := s.GetByEventAndUser(ctx, eventID, userID)
	if err == nil && existing != nil {
		return nil, ErrAlreadyRegistered
	}
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	query := `
		INSERT INTO participants (event_id, user_id, added_by)
		VALUES ($1, $2, $3)
		RETURNING id, event_id, user_id, guest_name, guest_claim_token, added_by, status, registered_at
	`

	var p models.Participant
	err = s.pool.QueryRow(ctx, query, eventID, userID, addedBy).Scan(
		&p.ID, &p.EventID, &p.UserID, &p.GuestName, &p.GuestClaimToken, &p.AddedBy, &p.Status, &p.RegisteredAt,
	)
	if err != nil {
		return nil, fmt.Errorf("add participant: %w", err)
	}
	return &p, nil
}

func (s *ParticipantService) AddGuest(ctx context.Context, eventID int, guestName string, addedBy int64) (*models.Participant, string, error) {
	// Generate claim token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return nil, "", fmt.Errorf("generate token: %w", err)
	}
	claimToken := hex.EncodeToString(tokenBytes)

	query := `
		INSERT INTO participants (event_id, guest_name, guest_claim_token, added_by)
		VALUES ($1, $2, $3, $4)
		RETURNING id, event_id, user_id, guest_name, guest_claim_token, added_by, status, registered_at
	`

	var p models.Participant
	err := s.pool.QueryRow(ctx, query, eventID, guestName, claimToken, addedBy).Scan(
		&p.ID, &p.EventID, &p.UserID, &p.GuestName, &p.GuestClaimToken, &p.AddedBy, &p.Status, &p.RegisteredAt,
	)
	if err != nil {
		return nil, "", fmt.Errorf("add guest: %w", err)
	}
	return &p, claimToken, nil
}

func (s *ParticipantService) RemoveUser(ctx context.Context, eventID int, userID int64) error {
	query := `
		UPDATE participants SET status = 'cancelled'
		WHERE event_id = $1 AND user_id = $2 AND status = 'confirmed'
	`
	result, err := s.pool.Exec(ctx, query, eventID, userID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrNotParticipant
	}
	return nil
}

func (s *ParticipantService) RemoveGuest(ctx context.Context, eventID int, guestName string) error {
	query := `
		UPDATE participants SET status = 'cancelled'
		WHERE event_id = $1 AND guest_name = $2 AND status = 'confirmed'
	`
	result, err := s.pool.Exec(ctx, query, eventID, guestName)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrNotParticipant
	}
	return nil
}

func (s *ParticipantService) GetByClaimToken(ctx context.Context, token string) (*models.Participant, error) {
	query := `
		SELECT id, event_id, user_id, guest_name, guest_claim_token, added_by, status, registered_at
		FROM participants
		WHERE guest_claim_token = $1 AND status = 'confirmed'
	`

	var p models.Participant
	err := s.pool.QueryRow(ctx, query, token).Scan(
		&p.ID, &p.EventID, &p.UserID, &p.GuestName, &p.GuestClaimToken, &p.AddedBy, &p.Status, &p.RegisteredAt,
	)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (s *ParticipantService) ClaimGuest(ctx context.Context, token string, userID int64) (*models.Participant, error) {
	query := `
		UPDATE participants
		SET user_id = $1, guest_name = NULL, guest_claim_token = NULL
		WHERE guest_claim_token = $2 AND status = 'confirmed'
		RETURNING id, event_id, user_id, guest_name, guest_claim_token, added_by, status, registered_at
	`
	var p models.Participant
	err := s.pool.QueryRow(ctx, query, userID, token).Scan(
		&p.ID, &p.EventID, &p.UserID, &p.GuestName, &p.GuestClaimToken, &p.AddedBy, &p.Status, &p.RegisteredAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrInvalidToken
		}
		return nil, err
	}
	return &p, nil
}
