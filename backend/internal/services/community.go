package services

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/luvvano/football-bot/backend/internal/models"
)

type CommunityService struct {
	pool *pgxpool.Pool
}

func NewCommunityService(pool *pgxpool.Pool) *CommunityService {
	return &CommunityService{pool: pool}
}

func (s *CommunityService) GetByTelegramChatID(ctx context.Context, chatID int64) (*models.Community, error) {
	query := `
		SELECT id, telegram_chat_id, name, created_by, settings, created_at
		FROM communities WHERE telegram_chat_id = $1
	`

	var c models.Community
	var settingsJSON []byte
	err := s.pool.QueryRow(ctx, query, chatID).Scan(
		&c.ID, &c.TelegramChatID, &c.Name, &c.CreatedBy, &settingsJSON, &c.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(settingsJSON, &c.Settings); err != nil {
		c.Settings = models.DefaultCommunitySettings()
	}

	return &c, nil
}

func (s *CommunityService) GetByID(ctx context.Context, id int) (*models.Community, error) {
	query := `
		SELECT id, telegram_chat_id, name, created_by, settings, created_at
		FROM communities WHERE id = $1
	`

	var c models.Community
	var settingsJSON []byte
	err := s.pool.QueryRow(ctx, query, id).Scan(
		&c.ID, &c.TelegramChatID, &c.Name, &c.CreatedBy, &settingsJSON, &c.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(settingsJSON, &c.Settings); err != nil {
		c.Settings = models.DefaultCommunitySettings()
	}

	return &c, nil
}

func (s *CommunityService) Create(ctx context.Context, chatID int64, name string, createdBy int64) (*models.Community, error) {
	settings := models.DefaultCommunitySettings()
	settingsJSON, _ := json.Marshal(settings)

	query := `
		INSERT INTO communities (telegram_chat_id, name, created_by, settings)
		VALUES ($1, $2, $3, $4)
		RETURNING id, telegram_chat_id, name, created_by, settings, created_at
	`

	var c models.Community
	var retSettingsJSON []byte
	err := s.pool.QueryRow(ctx, query, chatID, name, createdBy, settingsJSON).Scan(
		&c.ID, &c.TelegramChatID, &c.Name, &c.CreatedBy, &retSettingsJSON, &c.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create community: %w", err)
	}

	c.Settings = settings
	return &c, nil
}

func (s *CommunityService) AddMember(ctx context.Context, userID int64, communityID int, role string) (*models.CommunityMember, error) {
	query := `
		INSERT INTO community_members (user_id, community_id, role)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, community_id) DO UPDATE SET role = EXCLUDED.role
		RETURNING id, user_id, community_id, role, joined_at
	`

	var m models.CommunityMember
	err := s.pool.QueryRow(ctx, query, userID, communityID, role).Scan(
		&m.ID, &m.UserID, &m.CommunityID, &m.Role, &m.JoinedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("add member: %w", err)
	}
	return &m, nil
}

func (s *CommunityService) GetMember(ctx context.Context, userID int64, communityID int) (*models.CommunityMember, error) {
	query := `
		SELECT id, user_id, community_id, role, joined_at
		FROM community_members
		WHERE user_id = $1 AND community_id = $2
	`

	var m models.CommunityMember
	err := s.pool.QueryRow(ctx, query, userID, communityID).Scan(
		&m.ID, &m.UserID, &m.CommunityID, &m.Role, &m.JoinedAt,
	)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (s *CommunityService) IsAdmin(ctx context.Context, userID int64, communityID int) (bool, error) {
	member, err := s.GetMember(ctx, userID, communityID)
	if err != nil {
		return false, err
	}
	return member.Role == "admin", nil
}

func (s *CommunityService) GetMembers(ctx context.Context, communityID int) ([]models.CommunityMember, error) {
	query := `
		SELECT id, user_id, community_id, role, joined_at
		FROM community_members
		WHERE community_id = $1
		ORDER BY joined_at
	`

	rows, err := s.pool.Query(ctx, query, communityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []models.CommunityMember
	for rows.Next() {
		var m models.CommunityMember
		if err := rows.Scan(&m.ID, &m.UserID, &m.CommunityID, &m.Role, &m.JoinedAt); err != nil {
			return nil, err
		}
		members = append(members, m)
	}

	return members, rows.Err()
}
