package bot

import (
	"context"
	"fmt"
	"strings"

	"github.com/luvvano/football-bot/backend/internal/models"
	tele "gopkg.in/telebot.v3"
)

func (h *BotHandler) handleVenue(c tele.Context) error {
	ctx := context.Background()
	chat := c.Chat()
	sender := c.Sender()

	if chat.Type == tele.ChatPrivate {
		return c.Reply("⚠️ Эта команда работает только в группах")
	}

	// Get community
	community, err := h.communityService.GetByTelegramChatID(ctx, chat.ID)
	if err != nil {
		return c.Reply("❌ Сообщество не зарегистрировано. Используйте /start")
	}

	// Get user
	user, err := h.userService.GetByTelegramID(ctx, sender.ID)
	if err != nil {
		return c.Reply("❌ Вы не зарегистрированы. Используйте /start")
	}

	// Check if admin
	isAdmin, _ := h.communityService.IsAdmin(ctx, user.ID, community.ID)
	if !isAdmin {
		return c.Reply("❌ Только администратор может создавать площадки")
	}

	// Parse: /venue Название | Адрес (optional)
	args := strings.TrimSpace(strings.TrimPrefix(c.Text(), "/venue"))
	if args == "" {
		// List venues
		return h.listVenues(c, ctx, community.ID)
	}

	parts := strings.SplitN(args, "|", 2)
	name := strings.TrimSpace(parts[0])
	var address *string
	if len(parts) > 1 {
		addr := strings.TrimSpace(parts[1])
		address = &addr
	}

	if name == "" {
		return c.Reply("⚠️ Укажите название площадки: /venue Название | Адрес")
	}

	// Create venue
	venue, err := h.venueService.Create(ctx, models.CreateVenueInput{
		CommunityID: community.ID,
		Name:        name,
		Address:     address,
		CreatedBy:   user.ID,
	})
	if err != nil {
		return c.Reply("❌ Ошибка создания площадки")
	}

	msg := fmt.Sprintf("✅ Площадка создана!\n\n🏟 %s", venue.Name)
	if venue.Address != nil {
		msg += fmt.Sprintf("\n📍 %s", *venue.Address)
	}

	return c.Reply(msg)
}

func (h *BotHandler) listVenues(c tele.Context, ctx context.Context, communityID int) error {
	venues, err := h.venueService.GetByCommunity(ctx, communityID)
	if err != nil || len(venues) == 0 {
		return c.Reply("ℹ️ Площадки не созданы.\n\nДля создания: /venue Название | Адрес")
	}

	msg := "🏟 Площадки:\n\n"
	for i, v := range venues {
		msg += fmt.Sprintf("%d. %s", i+1, v.Name)
		if v.Address != nil {
			msg += fmt.Sprintf(" — %s", *v.Address)
		}
		msg += "\n"
	}

	msg += "\nДля создания: /venue Название | Адрес"

	return c.Reply(msg)
}
