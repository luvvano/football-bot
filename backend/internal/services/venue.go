package services

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/luvvano/football-bot/backend/internal/models"
)

type VenueService struct {
	pool *pgxpool.Pool
}

func NewVenueService(pool *pgxpool.Pool) *VenueService {
	return &VenueService{pool: pool}
}

func (s *VenueService) GetByID(ctx context.Context, id int) (*models.Venue, error) {
	query := `
		SELECT id, community_id, name, address, contact_type, contact_value, created_by, is_active, created_at
		FROM venues WHERE id = $1
	`

	var v models.Venue
	err := s.pool.QueryRow(ctx, query, id).Scan(
		&v.ID, &v.CommunityID, &v.Name, &v.Address, &v.ContactType,
		&v.ContactValue, &v.CreatedBy, &v.IsActive, &v.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func (s *VenueService) GetByCommunity(ctx context.Context, communityID int) ([]models.Venue, error) {
	query := `
		SELECT id, community_id, name, address, contact_type, contact_value, created_by, is_active, created_at
		FROM venues
		WHERE community_id = $1 AND is_active = true
		ORDER BY name
	`

	rows, err := s.pool.Query(ctx, query, communityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var venues []models.Venue
	for rows.Next() {
		var v models.Venue
		if err := rows.Scan(
			&v.ID, &v.CommunityID, &v.Name, &v.Address, &v.ContactType,
			&v.ContactValue, &v.CreatedBy, &v.IsActive, &v.CreatedAt,
		); err != nil {
			return nil, err
		}
		venues = append(venues, v)
	}

	return venues, rows.Err()
}

func (s *VenueService) Create(ctx context.Context, input models.CreateVenueInput) (*models.Venue, error) {
	query := `
		INSERT INTO venues (community_id, name, address, contact_type, contact_value, created_by)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, community_id, name, address, contact_type, contact_value, created_by, is_active, created_at
	`

	var v models.Venue
	err := s.pool.QueryRow(ctx, query,
		input.CommunityID, input.Name, input.Address, input.ContactType, input.ContactValue, input.CreatedBy,
	).Scan(
		&v.ID, &v.CommunityID, &v.Name, &v.Address, &v.ContactType,
		&v.ContactValue, &v.CreatedBy, &v.IsActive, &v.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create venue: %w", err)
	}
	return &v, nil
}

func (s *VenueService) Deactivate(ctx context.Context, id int) error {
	query := `UPDATE venues SET is_active = false WHERE id = $1`
	_, err := s.pool.Exec(ctx, query, id)
	return err
}
