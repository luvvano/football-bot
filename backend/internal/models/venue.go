package models

import "time"

type Venue struct {
	ID           int       `json:"id"`
	CommunityID  int       `json:"community_id"`
	Name         string    `json:"name"`
	Address      *string   `json:"address,omitempty"`
	ContactType  *string   `json:"contact_type,omitempty"` // telegram, whatsapp, phone
	ContactValue *string   `json:"contact_value,omitempty"`
	CreatedBy    *int64    `json:"created_by,omitempty"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
}

type CreateVenueInput struct {
	CommunityID  int
	Name         string
	Address      *string
	ContactType  *string
	ContactValue *string
	CreatedBy    int64
}
