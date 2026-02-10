package services

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/luvvano/football-bot/backend/internal/models"
)

type EventService struct {
	pool *pgxpool.Pool
}

func NewEventService(pool *pgxpool.Pool) *EventService {
	return &EventService{pool: pool}
}

func (s *EventService) GetByID(ctx context.Context, id int) (*models.Event, error) {
	query := `
		SELECT id, community_id, venue_id, event_date, event_time, team_size, max_subs_per_team,
		       status, team_a, team_b, subs_a, subs_b, created_by, created_at, finished_at
		FROM events WHERE id = $1
	`

	var e models.Event
	err := s.pool.QueryRow(ctx, query, id).Scan(
		&e.ID, &e.CommunityID, &e.VenueID, &e.EventDate, &e.EventTime, &e.TeamSize, &e.MaxSubsPerTeam,
		&e.Status, &e.TeamA, &e.TeamB, &e.SubsA, &e.SubsB, &e.CreatedBy, &e.CreatedAt, &e.FinishedAt,
	)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (s *EventService) GetByCommunity(ctx context.Context, communityID int, includeFinished bool) ([]models.Event, error) {
	query := `
		SELECT e.id, e.community_id, e.venue_id, e.event_date, e.event_time, e.team_size, e.max_subs_per_team,
		       e.status, e.team_a, e.team_b, e.subs_a, e.subs_b, e.created_by, e.created_at, e.finished_at
		FROM events e
		WHERE e.community_id = $1
	`

	if !includeFinished {
		query += ` AND e.status NOT IN ('finished', 'cancelled', 'not_held')`
	}
	query += ` ORDER BY e.event_date, e.event_time`

	rows, err := s.pool.Query(ctx, query, communityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []models.Event
	for rows.Next() {
		var e models.Event
		if err := rows.Scan(
			&e.ID, &e.CommunityID, &e.VenueID, &e.EventDate, &e.EventTime, &e.TeamSize, &e.MaxSubsPerTeam,
			&e.Status, &e.TeamA, &e.TeamB, &e.SubsA, &e.SubsB, &e.CreatedBy, &e.CreatedAt, &e.FinishedAt,
		); err != nil {
			return nil, err
		}
		events = append(events, e)
	}

	return events, rows.Err()
}

func (s *EventService) GetNextEvent(ctx context.Context, communityID int) (*models.Event, error) {
	query := `
		SELECT id, community_id, venue_id, event_date, event_time, team_size, max_subs_per_team,
		       status, team_a, team_b, subs_a, subs_b, created_by, created_at, finished_at
		FROM events
		WHERE community_id = $1 AND status IN ('open', 'full')
		  AND (event_date > CURRENT_DATE OR (event_date = CURRENT_DATE AND event_time >= CURRENT_TIME))
		ORDER BY event_date, event_time
		LIMIT 1
	`

	var e models.Event
	err := s.pool.QueryRow(ctx, query, communityID).Scan(
		&e.ID, &e.CommunityID, &e.VenueID, &e.EventDate, &e.EventTime, &e.TeamSize, &e.MaxSubsPerTeam,
		&e.Status, &e.TeamA, &e.TeamB, &e.SubsA, &e.SubsB, &e.CreatedBy, &e.CreatedAt, &e.FinishedAt,
	)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (s *EventService) Create(ctx context.Context, input models.CreateEventInput) (*models.Event, error) {
	query := `
		INSERT INTO events (community_id, venue_id, event_date, event_time, team_size, max_subs_per_team, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, community_id, venue_id, event_date, event_time, team_size, max_subs_per_team,
		          status, team_a, team_b, subs_a, subs_b, created_by, created_at, finished_at
	`

	var e models.Event
	err := s.pool.QueryRow(ctx, query,
		input.CommunityID, input.VenueID, input.EventDate, input.EventTime,
		input.TeamSize, input.MaxSubsPerTeam, input.CreatedBy,
	).Scan(
		&e.ID, &e.CommunityID, &e.VenueID, &e.EventDate, &e.EventTime, &e.TeamSize, &e.MaxSubsPerTeam,
		&e.Status, &e.TeamA, &e.TeamB, &e.SubsA, &e.SubsB, &e.CreatedBy, &e.CreatedAt, &e.FinishedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create event: %w", err)
	}
	return &e, nil
}

func (s *EventService) UpdateStatus(ctx context.Context, eventID int, status models.EventStatus) error {
	query := `UPDATE events SET status = $1 WHERE id = $2`
	_, err := s.pool.Exec(ctx, query, status, eventID)
	return err
}

func (s *EventService) SetVenue(ctx context.Context, eventID int, venueID int) error {
	query := `UPDATE events SET venue_id = $1 WHERE id = $2`
	_, err := s.pool.Exec(ctx, query, venueID, eventID)
	return err
}

func (s *EventService) Cancel(ctx context.Context, eventID int) error {
	return s.UpdateStatus(ctx, eventID, models.EventStatusCancelled)
}

func (s *EventService) Finish(ctx context.Context, eventID int) error {
	query := `UPDATE events SET status = 'finished', finished_at = $1 WHERE id = $2`
	_, err := s.pool.Exec(ctx, query, time.Now(), eventID)
	return err
}

func (s *EventService) GetWithVenue(ctx context.Context, eventID int) (*models.Event, error) {
	event, err := s.GetByID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	if event.VenueID != nil {
		venueQuery := `
			SELECT id, community_id, name, address, contact_type, contact_value, created_by, is_active, created_at
			FROM venues WHERE id = $1
		`
		var v models.Venue
		err := s.pool.QueryRow(ctx, venueQuery, *event.VenueID).Scan(
			&v.ID, &v.CommunityID, &v.Name, &v.Address, &v.ContactType,
			&v.ContactValue, &v.CreatedBy, &v.IsActive, &v.CreatedAt,
		)
		if err == nil {
			event.Venue = &v
		}
	}

	return event, nil
}
