package api

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/luvvano/football-bot/backend/internal/services"
)

type CommunitiesHandler struct {
	communityService *services.CommunityService
}

func NewCommunitiesHandler(communityService *services.CommunityService) *CommunitiesHandler {
	return &CommunitiesHandler{
		communityService: communityService,
	}
}

// GetCommunity GET /api/communities/:id
func (h *CommunitiesHandler) GetCommunity(w http.ResponseWriter, r *http.Request) {
	communityIDStr := chi.URLParam(r, "id")
	communityID, err := strconv.Atoi(communityIDStr)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid community id"})
		return
	}

	community, err := h.communityService.GetByID(r.Context(), communityID)
	if err != nil {
		respondJSON(w, http.StatusNotFound, map[string]string{"error": "community not found"})
		return
	}

	respondJSON(w, http.StatusOK, community)
}
