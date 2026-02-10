package models

import (
	"encoding/json"
	"time"
)

type CommunitySettings struct {
	TeamSize       int      `json:"team_size"`
	MaxSubsPerTeam int      `json:"max_subs_per_team"`
	Timezone       string   `json:"timezone"`
	RatingParams   []string `json:"rating_params"`
}

func DefaultCommunitySettings() CommunitySettings {
	return CommunitySettings{
		TeamSize:       5,
		MaxSubsPerTeam: 1,
		Timezone:       "UTC",
		RatingParams:   []string{"speed", "dribbling", "shot"},
	}
}

type Community struct {
	ID             int               `json:"id"`
	TelegramChatID int64             `json:"telegram_chat_id"`
	Name           string            `json:"name"`
	CreatedBy      *int64            `json:"created_by,omitempty"`
	Settings       CommunitySettings `json:"settings"`
	CreatedAt      time.Time         `json:"created_at"`
}

func (c *Community) SettingsJSON() ([]byte, error) {
	return json.Marshal(c.Settings)
}

type CommunityMember struct {
	ID          int       `json:"id"`
	UserID      int64     `json:"user_id"`
	CommunityID int       `json:"community_id"`
	Role        string    `json:"role"` // admin, player
	JoinedAt    time.Time `json:"joined_at"`
}
