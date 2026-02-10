package api

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/luvvano/football-bot/backend/internal/middleware"
	"github.com/luvvano/football-bot/backend/internal/models"
	"github.com/luvvano/football-bot/backend/internal/services"
)

type EventsHandler struct {
	eventService       *services.EventService
	participantService *services.ParticipantService
	userService        *services.UserService
}

func NewEventsHandler(
	eventService *services.EventService,
	participantService *services.ParticipantService,
	userService *services.UserService,
) *EventsHandler {
	return &EventsHandler{
		eventService:       eventService,
		participantService: participantService,
		userService:        userService,
	}
}

type EventResponse struct {
	Event        *models.Event       `json:"event"`
	Participants []models.Participant `json:"participants"`
}

// GetEvent GET /api/events/:id
func (h *EventsHandler) GetEvent(w http.ResponseWriter, r *http.Request) {
	eventIDStr := chi.URLParam(r, "id")
	eventID, err := strconv.Atoi(eventIDStr)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid event id"})
		return
	}

	event, err := h.eventService.GetWithVenue(r.Context(), eventID)
	if err != nil {
		respondJSON(w, http.StatusNotFound, map[string]string{"error": "event not found"})
		return
	}

	participants, err := h.participantService.GetByEvent(r.Context(), eventID)
	if err != nil {
		participants = []models.Participant{}
	}

	respondJSON(w, http.StatusOK, EventResponse{
		Event:        event,
		Participants: participants,
	})
}

// Join POST /api/events/:id/join
func (h *EventsHandler) Join(w http.ResponseWriter, r *http.Request) {
	twaUser := middleware.GetTWAUser(r.Context())
	if twaUser == nil {
		respondJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	eventIDStr := chi.URLParam(r, "id")
	eventID, err := strconv.Atoi(eventIDStr)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid event id"})
		return
	}

	// Get or create user
	user, err := h.userService.Upsert(r.Context(), models.CreateUserInput{
		TelegramID: twaUser.ID,
		Username:   stringPtr(twaUser.Username),
		FirstName:  stringPtr(twaUser.FirstName),
		LastName:   stringPtr(twaUser.LastName),
	})
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to create user"})
		return
	}

	// Check event exists
	event, err := h.eventService.GetByID(r.Context(), eventID)
	if err != nil {
		respondJSON(w, http.StatusNotFound, map[string]string{"error": "event not found"})
		return
	}

	// Check event is open
	if event.Status != models.EventStatusOpen {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "event is not open for registration"})
		return
	}

	// Check capacity
	count, _ := h.participantService.CountConfirmed(r.Context(), eventID)
	if count >= event.MaxParticipants() {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "event is full"})
		return
	}

	// Add participant
	participant, err := h.participantService.AddUser(r.Context(), eventID, user.ID, user.ID)
	if err != nil {
		if err == services.ErrAlreadyRegistered {
			respondJSON(w, http.StatusConflict, map[string]string{"error": "already registered"})
			return
		}
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to join event"})
		return
	}

	respondJSON(w, http.StatusOK, participant)
}

// Leave DELETE /api/events/:id/leave
func (h *EventsHandler) Leave(w http.ResponseWriter, r *http.Request) {
	twaUser := middleware.GetTWAUser(r.Context())
	if twaUser == nil {
		respondJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	eventIDStr := chi.URLParam(r, "id")
	eventID, err := strconv.Atoi(eventIDStr)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid event id"})
		return
	}

	// Get user
	user, err := h.userService.GetByTelegramID(r.Context(), twaUser.ID)
	if err != nil {
		respondJSON(w, http.StatusNotFound, map[string]string{"error": "user not found"})
		return
	}

	// Remove participant
	err = h.participantService.RemoveUser(r.Context(), eventID, user.ID)
	if err != nil {
		if err == services.ErrNotParticipant {
			respondJSON(w, http.StatusNotFound, map[string]string{"error": "not a participant"})
			return
		}
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to leave event"})
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// ListEvents GET /api/communities/:id/events
func (h *EventsHandler) ListEvents(w http.ResponseWriter, r *http.Request) {
	communityIDStr := chi.URLParam(r, "id")
	communityID, err := strconv.Atoi(communityIDStr)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid community id"})
		return
	}

	includeFinished := r.URL.Query().Get("include_finished") == "true"
	events, err := h.eventService.GetByCommunity(r.Context(), communityID, includeFinished)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to get events"})
		return
	}

	if events == nil {
		events = []models.Event{}
	}

	respondJSON(w, http.StatusOK, events)
}

func stringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
