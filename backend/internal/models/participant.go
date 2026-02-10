package models

import "time"

type ParticipantStatus string

const (
	ParticipantStatusConfirmed ParticipantStatus = "confirmed"
	ParticipantStatusCancelled ParticipantStatus = "cancelled"
)

type Participant struct {
	ID              int               `json:"id"`
	EventID         int               `json:"event_id"`
	UserID          *int64            `json:"user_id,omitempty"`
	GuestName       *string           `json:"guest_name,omitempty"`
	GuestClaimToken *string           `json:"guest_claim_token,omitempty"`
	AddedBy         *int64            `json:"added_by,omitempty"`
	Status          ParticipantStatus `json:"status"`
	RegisteredAt    time.Time         `json:"registered_at"`

	// Joined data
	User *User `json:"user,omitempty"`
}

func (p *Participant) DisplayName() string {
	if p.GuestName != nil {
		return *p.GuestName + " (гость)"
	}
	if p.User != nil {
		if p.User.FirstName != nil {
			name := *p.User.FirstName
			if p.User.LastName != nil {
				name += " " + *p.User.LastName
			}
			return name
		}
		if p.User.Username != nil {
			return "@" + *p.User.Username
		}
	}
	return "Unknown"
}

func (p *Participant) IsGuest() bool {
	return p.GuestName != nil
}

type AddParticipantInput struct {
	EventID   int
	UserID    *int64
	GuestName *string
	AddedBy   int64
}
