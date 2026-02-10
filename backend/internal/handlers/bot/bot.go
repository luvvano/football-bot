package bot

import (
	"github.com/luvvano/football-bot/backend/internal/services"
	tele "gopkg.in/telebot.v3"
)

type BotHandler struct {
	bot              *tele.Bot
	webAppURL        string
	userService      *services.UserService
	communityService *services.CommunityService
	venueService     *services.VenueService
	eventService     *services.EventService
	participantService *services.ParticipantService
}

func NewBotHandler(
	bot *tele.Bot,
	webAppURL string,
	userService *services.UserService,
	communityService *services.CommunityService,
	venueService *services.VenueService,
	eventService *services.EventService,
	participantService *services.ParticipantService,
) *BotHandler {
	return &BotHandler{
		bot:              bot,
		webAppURL:        webAppURL,
		userService:      userService,
		communityService: communityService,
		venueService:     venueService,
		eventService:     eventService,
		participantService: participantService,
	}
}

func (h *BotHandler) RegisterHandlers() {
	// Basic commands
	h.bot.Handle("/start", h.handleStart)
	h.bot.Handle("/event", h.handleEvent)
	h.bot.Handle("/add", h.handleAdd)
	h.bot.Handle("/remove", h.handleRemove)
	h.bot.Handle("/info", h.handleInfo)
	h.bot.Handle("/events", h.handleEvents)
	h.bot.Handle("/app", h.handleApp)
	h.bot.Handle("/claim", h.handleClaim)

	// Admin commands
	h.bot.Handle("/venue", h.handleVenue)
	h.bot.Handle("/cancel", h.handleCancel)
	h.bot.Handle("/finish", h.handleFinish)
	h.bot.Handle("/notheld", h.handleNotHeld)

	// Callback queries for venue selection
	h.bot.Handle(&tele.InlineButton{Unique: "venue"}, h.handleVenueCallback)
	h.bot.Handle(&tele.InlineButton{Unique: "skip_venue"}, h.handleSkipVenueCallback)
}

func stringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
