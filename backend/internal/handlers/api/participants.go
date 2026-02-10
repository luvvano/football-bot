package api

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/luvvano/football-bot/backend/internal/services"
)

type ParticipantsHandler struct {
	participantService *services.ParticipantService
}

func NewParticipantsHandler(participantService *services.ParticipantService) *ParticipantsHandler {
	return &ParticipantsHandler{
		participantService: participantService,
	}
}

// ListParticipants GET /api/events/:id/participants
func (h *ParticipantsHandler) ListParticipants(w http.ResponseWriter, r *http.Request) {
	eventIDStr := chi.URLParam(r, "id")
	eventID, err := strconv.Atoi(eventIDStr)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid event id"})
		return
	}

	participants, err := h.participantService.GetByEvent(r.Context(), eventID)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to get participants"})
		return
	}

	respondJSON(w, http.StatusOK, participants)
}
