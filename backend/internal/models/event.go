package models

import "time"

type EventStatus string

const (
	EventStatusOpen       EventStatus = "open"
	EventStatusFull       EventStatus = "full"
	EventStatusInProgress EventStatus = "in_progress"
	EventStatusFinished   EventStatus = "finished"
	EventStatusCancelled  EventStatus = "cancelled"
	EventStatusNotHeld    EventStatus = "not_held"
)

type Event struct {
	ID             int         `json:"id"`
	CommunityID    int         `json:"community_id"`
	VenueID        *int        `json:"venue_id,omitempty"`
	EventDate      time.Time   `json:"event_date"`
	EventTime      string      `json:"event_time"` // HH:MM format
	TeamSize       int         `json:"team_size"`
	MaxSubsPerTeam int         `json:"max_subs_per_team"`
	Status         EventStatus `json:"status"`
	TeamA          []int       `json:"team_a"`
	TeamB          []int       `json:"team_b"`
	SubsA          []int       `json:"subs_a"`
	SubsB          []int       `json:"subs_b"`
	CreatedBy      *int64      `json:"created_by,omitempty"`
	CreatedAt      time.Time   `json:"created_at"`
	FinishedAt     *time.Time  `json:"finished_at,omitempty"`
	
	// Joined data
	Venue        *Venue        `json:"venue,omitempty"`
	Participants []Participant `json:"participants,omitempty"`
}

type CreateEventInput struct {
	CommunityID    int
	VenueID        *int
	EventDate      time.Time
	EventTime      string
	TeamSize       int
	MaxSubsPerTeam int
	CreatedBy      int64
}

func (e *Event) MaxParticipants() int {
	return (e.TeamSize + e.MaxSubsPerTeam) * 2
}
