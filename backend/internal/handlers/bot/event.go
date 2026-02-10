package bot

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/luvvano/football-bot/backend/internal/models"
	tele "gopkg.in/telebot.v3"
)

func (h *BotHandler) handleEvent(c tele.Context) error {
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

	// Parse arguments: /event [date] [time]
	args := strings.Fields(c.Text())
	
	var eventDate time.Time
	var eventTime string

	now := time.Now()

	if len(args) >= 2 {
		// Try to parse date
		parsedDate, err := parseDate(args[1])
		if err != nil {
			return c.Reply("❌ Неверный формат даты. Используйте: ДД.ММ, ДД.ММ.ГГГГ или +N (дней)")
		}
		eventDate = parsedDate
	} else {
		// Default to today or tomorrow
		eventDate = now.Truncate(24 * time.Hour)
		if now.Hour() >= 18 {
			eventDate = eventDate.AddDate(0, 0, 1)
		}
	}

	if len(args) >= 3 {
		// Parse time
		eventTime = args[2]
		if !isValidTime(eventTime) {
			return c.Reply("❌ Неверный формат времени. Используйте: ЧЧ:ММ")
		}
	} else {
		eventTime = "19:00"
	}

	// Create event
	event, err := h.eventService.Create(ctx, models.CreateEventInput{
		CommunityID:    community.ID,
		EventDate:      eventDate,
		EventTime:      eventTime,
		TeamSize:       community.Settings.TeamSize,
		MaxSubsPerTeam: community.Settings.MaxSubsPerTeam,
		CreatedBy:      user.ID,
	})
	if err != nil {
		return c.Reply("❌ Ошибка создания игры")
	}

	// Get venues for selection
	venues, _ := h.venueService.GetByCommunity(ctx, community.ID)

	msg := fmt.Sprintf(`⚽ Игра создана!

📅 Дата: %s
🕐 Время: %s
👥 Размер команды: %d
🔄 Замен: %d

`, formatDate(eventDate), eventTime, event.TeamSize, event.MaxSubsPerTeam)

	if len(venues) > 0 {
		msg += "🏟 Выберите площадку:"
		
		btn := &tele.ReplyMarkup{}
		var rows []tele.Row
		
		for _, v := range venues {
			venueBtn := btn.Data(v.Name, "venue", fmt.Sprintf("%d_%d", event.ID, v.ID))
			rows = append(rows, btn.Row(venueBtn))
		}
		
		skipBtn := btn.Data("⏩ Пропустить", "skip_venue", fmt.Sprint(event.ID))
		rows = append(rows, btn.Row(skipBtn))
		
		btn.Inline(rows...)
		return c.Send(msg, btn)
	}

	msg += "Используйте /add для записи на игру"
	return c.Reply(msg)
}

func (h *BotHandler) handleVenueCallback(c tele.Context) error {
	ctx := context.Background()
	data := c.Callback().Data
	
	parts := strings.Split(data, "_")
	if len(parts) != 2 {
		return c.Respond(&tele.CallbackResponse{Text: "Ошибка"})
	}

	eventID, _ := strconv.Atoi(parts[0])
	venueID, _ := strconv.Atoi(parts[1])

	err := h.eventService.SetVenue(ctx, eventID, venueID)
	if err != nil {
		return c.Respond(&tele.CallbackResponse{Text: "Ошибка выбора площадки"})
	}

	venue, _ := h.venueService.GetByID(ctx, venueID)
	event, _ := h.eventService.GetByID(ctx, eventID)

	msg := fmt.Sprintf(`⚽ Игра обновлена!

📅 Дата: %s
🕐 Время: %s
🏟 Площадка: %s

Используйте /add для записи`, 
		formatDate(event.EventDate), 
		event.EventTime,
		venue.Name)

	c.Respond(&tele.CallbackResponse{Text: "✅ Площадка выбрана"})
	return c.Edit(msg)
}

func (h *BotHandler) handleSkipVenueCallback(c tele.Context) error {
	ctx := context.Background()
	eventID, _ := strconv.Atoi(c.Callback().Data)

	event, _ := h.eventService.GetByID(ctx, eventID)

	msg := fmt.Sprintf(`⚽ Игра создана!

📅 Дата: %s
🕐 Время: %s

Площадка будет объявлена позже.
Используйте /add для записи`, 
		formatDate(event.EventDate), 
		event.EventTime)

	c.Respond(&tele.CallbackResponse{Text: "⏩ Площадка пропущена"})
	return c.Edit(msg)
}

func parseDate(s string) (time.Time, error) {
	now := time.Now()
	
	// Check for +N format
	if strings.HasPrefix(s, "+") {
		days, err := strconv.Atoi(s[1:])
		if err != nil {
			return time.Time{}, err
		}
		return now.Truncate(24 * time.Hour).AddDate(0, 0, days), nil
	}

	// Try DD.MM format
	parts := strings.Split(s, ".")
	if len(parts) >= 2 {
		day, _ := strconv.Atoi(parts[0])
		month, _ := strconv.Atoi(parts[1])
		year := now.Year()
		if len(parts) >= 3 {
			year, _ = strconv.Atoi(parts[2])
			if year < 100 {
				year += 2000
			}
		}
		return time.Date(year, time.Month(month), day, 0, 0, 0, 0, now.Location()), nil
	}

	return time.Time{}, fmt.Errorf("invalid date format")
}

func isValidTime(s string) bool {
	parts := strings.Split(s, ":")
	if len(parts) != 2 {
		return false
	}
	hours, err1 := strconv.Atoi(parts[0])
	mins, err2 := strconv.Atoi(parts[1])
	return err1 == nil && err2 == nil && hours >= 0 && hours <= 23 && mins >= 0 && mins <= 59
}

func formatDate(t time.Time) string {
	weekdays := []string{"Вс", "Пн", "Вт", "Ср", "Чт", "Пт", "Сб"}
	return fmt.Sprintf("%s, %02d.%02d", weekdays[t.Weekday()], t.Day(), t.Month())
}
