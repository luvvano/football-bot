package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/luvvano/football-bot/backend/internal/middleware"
)

type AuthHandler struct {
	botToken string
	maxAge   time.Duration
}

func NewAuthHandler(botToken string, maxAge time.Duration) *AuthHandler {
	return &AuthHandler{
		botToken: botToken,
		maxAge:   maxAge,
	}
}

type ValidateRequest struct {
	InitData string `json:"init_data"`
}

type ValidateResponse struct {
	Valid bool                 `json:"valid"`
	User  *middleware.TWAUser `json:"user,omitempty"`
	Error string              `json:"error,omitempty"`
}

// Validate POST /api/auth/validate
func (h *AuthHandler) Validate(w http.ResponseWriter, r *http.Request) {
	var req ValidateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, ValidateResponse{
			Valid: false,
			Error: "invalid request body",
		})
		return
	}

	data, err := middleware.ValidateTWAInitData(req.InitData, h.botToken, h.maxAge)
	if err != nil {
		respondJSON(w, http.StatusOK, ValidateResponse{
			Valid: false,
			Error: err.Error(),
		})
		return
	}

	respondJSON(w, http.StatusOK, ValidateResponse{
		Valid: true,
		User:  &data.User,
	})
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
