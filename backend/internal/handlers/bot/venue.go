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

	// Get user
	user, err := h.userService.GetByTelegramID(ctx, sender.ID)
	if err != nil {
		return c.Reply("❌ Вы не зарегистрированы. Используйте /start в группе с ботом")
	}

	// Private chat - show all user's venues
	if chat.Type == tele.ChatPrivate {
		return h.handlePrivateVenue(c, ctx, user)
	}

	// Group chat - existing logic
	community, err := h.communityService.GetByTelegramChatID(ctx, chat.ID)
	if err != nil {
		return c.Reply("❌ Сообщество не зарегистрировано. Используйте /start")
	}

	// Check if admin for creating venues
	isAdmin, _ := h.communityService.IsAdmin(ctx, user.ID, community.ID)

	// Parse: /venue Название | Адрес (optional)
	args := strings.TrimSpace(strings.TrimPrefix(c.Text(), "/venue"))
	if args == "" {
		// List venues
		return h.listVenues(c, ctx, community.ID, isAdmin)
	}

	if !isAdmin {
		return c.Reply("❌ Только администратор может создавать площадки")
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

func (h *BotHandler) handlePrivateVenue(c tele.Context, ctx context.Context, user *models.User) error {
	// Get user's communities
	communities, err := h.communityService.GetUserCommunities(ctx, user.ID)
	if err != nil || len(communities) == 0 {
		return c.Reply("ℹ️ Вы не состоите ни в одном сообществе.\n\nДобавьте бота в группу и используйте /start")
	}

	msg := "🏟 Ваши площадки:\n"

	hasVenues := false
	for _, comm := range communities {
		venues, err := h.venueService.GetByCommunity(ctx, comm.ID)
		if err != nil || len(venues) == 0 {
			continue
		}

		hasVenues = true
		msg += fmt.Sprintf("\n📍 %s:\n", comm.Name)
		for _, v := range venues {
			msg += fmt.Sprintf("  • %s", v.Name)
			if v.Address != nil {
				msg += fmt.Sprintf(" — %s", *v.Address)
			}
			msg += "\n"
		}
	}

	if !hasVenues {
		msg = "ℹ️ Площадки не созданы.\n\nСоздайте площадку командой /venue в группе (только для админов)"
	}

	return c.Reply(msg)
}

func (h *BotHandler) listVenues(c tele.Context, ctx context.Context, communityID int, isAdmin bool) error {
	venues, err := h.venueService.GetByCommunity(ctx, communityID)
	if err != nil || len(venues) == 0 {
		if isAdmin {
			return c.Reply("ℹ️ Площадки не созданы.\n\nДля создания: /venue Название | Адрес")
		}
		return c.Reply("ℹ️ Площадки не созданы")
	}

	msg := "🏟 Площадки:\n\n"
	for i, v := range venues {
		msg += fmt.Sprintf("%d. %s", i+1, v.Name)
		if v.Address != nil {
			msg += fmt.Sprintf(" — %s", *v.Address)
		}
		msg += "\n"
	}

	if isAdmin {
		msg += "\nДля создания: /venue Название | Адрес"
	}

	return c.Reply(msg)
}
