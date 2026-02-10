package models

import "time"

type PlayerRating struct {
	ID              int       `json:"id"`
	CommunityID     int       `json:"community_id"`
	UserID          int64     `json:"user_id"`
	Speed           float64   `json:"speed"`
	Dribbling       float64   `json:"dribbling"`
	Shot            float64   `json:"shot"`
	GamesPlayed     int       `json:"games_played"`
	RatingsReceived int       `json:"ratings_received"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func (r *PlayerRating) OverallRating() float64 {
	return (r.Speed + r.Dribbling + r.Shot) / 3
}

type RatingVote struct {
	ID         int       `json:"id"`
	EventID    int       `json:"event_id"`
	FromUserID int64     `json:"from_user_id"`
	ToUserID   int64     `json:"to_user_id"`
	Speed      int       `json:"speed"`
	Dribbling  int       `json:"dribbling"`
	Shot       int       `json:"shot"`
	CreatedAt  time.Time `json:"created_at"`
}
