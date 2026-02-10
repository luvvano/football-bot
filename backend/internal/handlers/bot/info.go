package bot

import (
	"context"
	"fmt"

	tele "gopkg.in/telebot.v3"
)

func (h *BotHandler) handleInfo(c tele.Context) error {
	ctx := context.Background()
	chat := c.Chat()

	if chat.Type == tele.ChatPrivate {
		return c.Reply("⚠️ Эта команда работает только в группах")
	}

	// Get community
	community, err := h.communityService.GetByTelegramChatID(ctx, chat.ID)
	if err != nil {
		return c.Reply("❌ Сообщество не зарегистрировано. Используйте /start")
	}

	// Get next event
	event, err := h.eventService.GetNextEvent(ctx, community.ID)
	if err != nil {
		return c.Reply("ℹ️ Нет предстоящих игр. Создайте игру командой /event")
	}

	// Get venue if set
	var venueName string
	if event.VenueID != nil {
		venue, err := h.venueService.GetByID(ctx, *event.VenueID)
		if err == nil {
			venueName = venue.Name
			if venue.Address != nil {
				venueName += fmt.Sprintf("\n   📍 %s", *venue.Address)
			}
		}
	}

	// Get participants
	participants, _ := h.participantService.GetByEvent(ctx, event.ID)

	msg := fmt.Sprintf(`⚽ Ближайшая игра

📅 Дата: %s
🕐 Время: %s
👥 Размер команды: %d (+%d замены)
`, formatDate(event.EventDate), formatTime(event.EventTime), event.TeamSize, event.MaxSubsPerTeam)

	if venueName != "" {
		msg += fmt.Sprintf("🏟 Площадка: %s\n", venueName)
	}

	msg += fmt.Sprintf("\n📋 Участники (%d/%d):\n", len(participants), event.MaxParticipants())

	if len(participants) == 0 {
		msg += "Пока никого нет. /add для записи"
	} else {
		for i, p := range participants {
			msg += fmt.Sprintf("%d. %s\n", i+1, p.DisplayName())
		}
		
		remaining := event.MaxParticipants() - len(participants)
		if remaining > 0 {
			msg += fmt.Sprintf("\n🟢 Свободных мест: %d", remaining)
		} else {
			msg += "\n🔴 Мест нет"
		}
	}

	return c.Reply(msg)
}

func (h *BotHandler) handleEvents(c tele.Context) error {
	ctx := context.Background()
	chat := c.Chat()
	sender := c.Sender()

	// Private chat - show events from all user's communities
	if chat.Type == tele.ChatPrivate {
		return h.handlePrivateEvents(c, ctx, sender.ID)
	}

	// Group chat - existing logic
	community, err := h.communityService.GetByTelegramChatID(ctx, chat.ID)
	if err != nil {
		return c.Reply("❌ Сообщество не зарегистрировано. Используйте /start")
	}

	// Get upcoming events
	events, err := h.eventService.GetByCommunity(ctx, community.ID, false)
	if err != nil || len(events) == 0 {
		return c.Reply("ℹ️ Нет предстоящих игр. Создайте игру командой /event")
	}

	msg := "📅 Предстоящие игры:\n\n"

	for _, e := range events {
		count, _ := h.participantService.CountConfirmed(ctx, e.ID)
		
		statusEmoji := "🟢"
		if count >= e.MaxParticipants() {
			statusEmoji = "🔴"
		} else if count >= e.MaxParticipants()/2 {
			statusEmoji = "🟡"
		}

		msg += fmt.Sprintf("%s %s в %s — %d/%d игроков\n",
			statusEmoji,
			formatDate(e.EventDate),
			formatTime(e.EventTime),
			count,
			e.MaxParticipants())
	}

	msg += "\nИспользуйте /info для подробностей о ближайшей игре"

	return c.Reply(msg)
}

func (h *BotHandler) handlePrivateEvents(c tele.Context, ctx context.Context, telegramID int64) error {
	user, err := h.userService.GetByTelegramID(ctx, telegramID)
	if err != nil {
		return c.Reply("❌ Вы не зарегистрированы. Используйте /start в группе с ботом")
	}

	communities, err := h.communityService.GetUserCommunities(ctx, user.ID)
	if err != nil || len(communities) == 0 {
		return c.Reply("ℹ️ Вы не состоите ни в одном сообществе")
	}

	msg := "📅 Ваши предстоящие игры:\n"
	hasEvents := false

	for _, comm := range communities {
		events, err := h.eventService.GetByCommunity(ctx, comm.ID, false)
		if err != nil || len(events) == 0 {
			continue
		}

		hasEvents = true
		msg += fmt.Sprintf("\n⚽ %s:\n", comm.Name)

		for _, e := range events {
			count, _ := h.participantService.CountConfirmed(ctx, e.ID)
			
			statusEmoji := "🟢"
			if count >= e.MaxParticipants() {
				statusEmoji = "🔴"
			} else if count >= e.MaxParticipants()/2 {
				statusEmoji = "🟡"
			}

			msg += fmt.Sprintf("  %s %s в %s — %d/%d\n",
				statusEmoji,
				formatDate(e.EventDate),
				formatTime(e.EventTime),
				count,
				e.MaxParticipants())
		}
	}

	if !hasEvents {
		return c.Reply("ℹ️ Нет предстоящих игр в ваших сообществах")
	}

	return c.Reply(msg)
}

func (h *BotHandler) handleApp(c tele.Context) error {
	ctx := context.Background()
	chat := c.Chat()

	var url string
	if chat.Type != tele.ChatPrivate {
		community, err := h.communityService.GetByTelegramChatID(ctx, chat.ID)
		if err == nil {
			url = fmt.Sprintf("%s?community=%d", h.webAppURL, community.ID)
		}
	}

	if url == "" {
		url = h.webAppURL
	}

	btn := &tele.ReplyMarkup{}
	webAppBtn := btn.WebApp("⚽ Открыть приложение", &tele.WebApp{URL: url})
	btn.Inline(btn.Row(webAppBtn))

	return c.Send("Нажмите кнопку, чтобы открыть приложение:", btn)
}
