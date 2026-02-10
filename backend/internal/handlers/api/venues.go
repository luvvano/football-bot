package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/luvvano/football-bot/backend/internal/middleware"
	"github.com/luvvano/football-bot/backend/internal/models"
	"github.com/luvvano/football-bot/backend/internal/services"
)

type VenuesHandler struct {
	venueService     *services.VenueService
	communityService *services.CommunityService
	userService      *services.UserService
}

func NewVenuesHandler(
	venueService *services.VenueService,
	communityService *services.CommunityService,
	userService *services.UserService,
) *VenuesHandler {
	return &VenuesHandler{
		venueService:     venueService,
		communityService: communityService,
		userService:      userService,
	}
}

type CreateVenueRequest struct {
	Name         string  `json:"name"`
	Address      *string `json:"address,omitempty"`
	ContactType  *string `json:"contact_type,omitempty"`
	ContactValue *string `json:"contact_value,omitempty"`
}

// ListVenues GET /api/communities/:id/venues
func (h *VenuesHandler) ListVenues(w http.ResponseWriter, r *http.Request) {
	communityIDStr := chi.URLParam(r, "id")
	communityID, err := strconv.Atoi(communityIDStr)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid community id"})
		return
	}

	venues, err := h.venueService.GetByCommunity(r.Context(), communityID)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to get venues"})
		return
	}

	if venues == nil {
		venues = []models.Venue{}
	}

	respondJSON(w, http.StatusOK, venues)
}

// CreateVenue POST /api/communities/:id/venues
func (h *VenuesHandler) CreateVenue(w http.ResponseWriter, r *http.Request) {
	twaUser := middleware.GetTWAUser(r.Context())
	if twaUser == nil {
		respondJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	communityIDStr := chi.URLParam(r, "id")
	communityID, err := strconv.Atoi(communityIDStr)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid community id"})
		return
	}

	// Get user
	user, err := h.userService.GetByTelegramID(r.Context(), twaUser.ID)
	if err != nil {
		respondJSON(w, http.StatusNotFound, map[string]string{"error": "user not found"})
		return
	}

	// Check if user is admin
	isAdmin, err := h.communityService.IsAdmin(r.Context(), user.ID, communityID)
	if err != nil || !isAdmin {
		respondJSON(w, http.StatusForbidden, map[string]string{"error": "admin access required"})
		return
	}

	var req CreateVenueRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	if req.Name == "" {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "name is required"})
		return
	}

	venue, err := h.venueService.Create(r.Context(), models.CreateVenueInput{
		CommunityID:  communityID,
		Name:         req.Name,
		Address:      req.Address,
		ContactType:  req.ContactType,
		ContactValue: req.ContactValue,
		CreatedBy:    user.ID,
	})
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to create venue"})
		return
	}

	respondJSON(w, http.StatusCreated, venue)
}
