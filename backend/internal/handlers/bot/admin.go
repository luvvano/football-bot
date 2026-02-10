package bot

import (
	"context"
	"fmt"

	"github.com/luvvano/football-bot/backend/internal/models"
	tele "gopkg.in/telebot.v3"
)

func (h *BotHandler) handleCancel(c tele.Context) error {
	ctx := context.Background()
	chat := c.Chat()
	sender := c.Sender()

	if chat.Type == tele.ChatPrivate {
		return c.Reply("⚠️ Эта команда работает только в группах")
	}

	// Get community
	community, err := h.communityService.GetByTelegramChatID(ctx, chat.ID)
	if err != nil {
		return c.Reply("❌ Сообщество не зарегистрировано")
	}

	// Get user and check admin
	user, err := h.userService.GetByTelegramID(ctx, sender.ID)
	if err != nil {
		return c.Reply("❌ Вы не зарегистрированы")
	}

	isAdmin, _ := h.communityService.IsAdmin(ctx, user.ID, community.ID)
	if !isAdmin {
		return c.Reply("❌ Только администратор может отменять игры")
	}

	// Get next event
	event, err := h.eventService.GetNextEvent(ctx, community.ID)
	if err != nil {
		return c.Reply("❌ Нет предстоящих игр для отмены")
	}

	// Update status to cancelled
	err = h.eventService.UpdateStatus(ctx, event.ID, models.EventStatusCancelled)
	if err != nil {
		return c.Reply("❌ Ошибка отмены игры")
	}

	return c.Reply(fmt.Sprintf("🚫 Игра на %s в %s отменена", formatDate(event.EventDate), formatTime(event.EventTime)))
}

func (h *BotHandler) handleFinish(c tele.Context) error {
	ctx := context.Background()
	chat := c.Chat()
	sender := c.Sender()

	if chat.Type == tele.ChatPrivate {
		return c.Reply("⚠️ Эта команда работает только в группах")
	}

	// Get community
	community, err := h.communityService.GetByTelegramChatID(ctx, chat.ID)
	if err != nil {
		return c.Reply("❌ Сообщество не зарегистрировано")
	}

	// Get user and check admin
	user, err := h.userService.GetByTelegramID(ctx, sender.ID)
	if err != nil {
		return c.Reply("❌ Вы не зарегистрированы")
	}

	isAdmin, _ := h.communityService.IsAdmin(ctx, user.ID, community.ID)
	if !isAdmin {
		return c.Reply("❌ Только администратор может завершать игры")
	}

	// Get next event (or current in_progress)
	event, err := h.eventService.GetNextEvent(ctx, community.ID)
	if err != nil {
		return c.Reply("❌ Нет активных игр для завершения")
	}

	// Update status to finished
	err = h.eventService.UpdateStatus(ctx, event.ID, models.EventStatusFinished)
	if err != nil {
		return c.Reply("❌ Ошибка завершения игры")
	}

	participants, _ := h.participantService.GetByEvent(ctx, event.ID)

	msg := fmt.Sprintf("✅ Игра на %s завершена!\n\n👥 Участвовали: %d человек",
		formatDate(event.EventDate), len(participants))

	// TODO: In future, add rating prompt here

	return c.Reply(msg)
}

func (h *BotHandler) handleNotHeld(c tele.Context) error {
	ctx := context.Background()
	chat := c.Chat()
	sender := c.Sender()

	if chat.Type == tele.ChatPrivate {
		return c.Reply("⚠️ Эта команда работает только в группах")
	}

	// Get community
	community, err := h.communityService.GetByTelegramChatID(ctx, chat.ID)
	if err != nil {
		return c.Reply("❌ Сообщество не зарегистрировано")
	}

	// Get user and check admin
	user, err := h.userService.GetByTelegramID(ctx, sender.ID)
	if err != nil {
		return c.Reply("❌ Вы не зарегистрированы")
	}

	isAdmin, _ := h.communityService.IsAdmin(ctx, user.ID, community.ID)
	if !isAdmin {
		return c.Reply("❌ Только администратор может отмечать несостоявшиеся игры")
	}

	// Get next event
	event, err := h.eventService.GetNextEvent(ctx, community.ID)
	if err != nil {
		return c.Reply("❌ Нет активных игр")
	}

	// Update status to not_held
	err = h.eventService.UpdateStatus(ctx, event.ID, models.EventStatusNotHeld)
	if err != nil {
		return c.Reply("❌ Ошибка")
	}

	return c.Reply(fmt.Sprintf("⚪ Игра на %s отмечена как несостоявшаяся", formatDate(event.EventDate)))
}
