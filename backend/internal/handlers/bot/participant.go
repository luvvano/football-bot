package bot

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/luvvano/football-bot/backend/internal/models"
	"github.com/luvvano/football-bot/backend/internal/services"
	tele "gopkg.in/telebot.v3"
)

func (h *BotHandler) handleAdd(c tele.Context) error {
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

	// Get next event
	event, err := h.eventService.GetNextEvent(ctx, community.ID)
	if err != nil {
		return c.Reply("❌ Нет предстоящих игр. Создайте игру командой /event")
	}

	// Get user
	user, err := h.userService.Upsert(ctx, models.CreateUserInput{
		TelegramID: sender.ID,
		Username:   stringPtr(sender.Username),
		FirstName:  stringPtr(sender.FirstName),
		LastName:   stringPtr(sender.LastName),
	})
	if err != nil {
		return c.Reply("❌ Ошибка регистрации пользователя")
	}

	// Ensure user is community member
	h.communityService.AddMember(ctx, user.ID, community.ID, "player")

	args := strings.TrimSpace(strings.TrimPrefix(c.Text(), "/add"))
	
	if args == "" {
		// Add self
		return h.addSelf(c, ctx, event, user)
	}

	if strings.HasPrefix(args, "@") {
		// Add another user by username
		return h.addByUsername(c, ctx, event, user, strings.TrimPrefix(args, "@"))
	}

	// Add guest
	return h.addGuest(c, ctx, event, user, args)
}

func (h *BotHandler) addSelf(c tele.Context, ctx context.Context, event *models.Event, user *models.User) error {
	// Check capacity
	count, _ := h.participantService.CountConfirmed(ctx, event.ID)
	if count >= event.MaxParticipants() {
		return c.Reply("❌ Игра заполнена")
	}

	_, err := h.participantService.AddUser(ctx, event.ID, user.ID, user.ID)
	if err != nil {
		if err == services.ErrAlreadyRegistered {
			return c.Reply("⚠️ Вы уже записаны на эту игру")
		}
		return c.Reply("❌ Ошибка записи")
	}

	return h.showParticipantsList(c, ctx, event, fmt.Sprintf("✅ %s записан(а)", safeUsername(user)))
}

func (h *BotHandler) addByUsername(c tele.Context, ctx context.Context, event *models.Event, addedBy *models.User, username string) error {
	// Find user by username
	targetUser, err := h.userService.GetByUsername(ctx, username)
	if err != nil {
		return c.Reply(fmt.Sprintf("❌ Пользователь @%s не найден. Он должен сначала использовать /start", username))
	}

	// Check capacity
	count, _ := h.participantService.CountConfirmed(ctx, event.ID)
	if count >= event.MaxParticipants() {
		return c.Reply("❌ Игра заполнена")
	}

	_, err = h.participantService.AddUser(ctx, event.ID, targetUser.ID, addedBy.ID)
	if err != nil {
		if err == services.ErrAlreadyRegistered {
			return c.Reply(fmt.Sprintf("⚠️ @%s уже записан(а) на эту игру", username))
		}
		return c.Reply("❌ Ошибка записи")
	}

	return h.showParticipantsList(c, ctx, event, fmt.Sprintf("✅ @%s записан(а)", username))
}

func (h *BotHandler) addGuest(c tele.Context, ctx context.Context, event *models.Event, addedBy *models.User, guestName string) error {
	// Check capacity
	count, _ := h.participantService.CountConfirmed(ctx, event.ID)
	if count >= event.MaxParticipants() {
		return c.Reply("❌ Игра заполнена")
	}

	participant, claimToken, err := h.participantService.AddGuest(ctx, event.ID, guestName, addedBy.ID)
	if err != nil {
		return c.Reply("❌ Ошибка записи гостя")
	}

	_ = participant
	msg := fmt.Sprintf(`✅ Гость "%s" записан!

🔗 Токен для привязки:
%s

📲 Чтобы привязать аккаунт, гость должен написать боту в личку:
/claim %s`, guestName, claimToken, claimToken)

	return h.showParticipantsList(c, ctx, event, msg)
}

func (h *BotHandler) handleRemove(c tele.Context) error {
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

	// Get next event
	event, err := h.eventService.GetNextEvent(ctx, community.ID)
	if err != nil {
		return c.Reply("❌ Нет предстоящих игр")
	}

	// Get user
	user, err := h.userService.GetByTelegramID(ctx, sender.ID)
	if err != nil {
		return c.Reply("❌ Вы не зарегистрированы")
	}

	args := strings.TrimSpace(strings.TrimPrefix(c.Text(), "/remove"))

	if args == "" {
		// Remove self
		return h.removeSelf(c, ctx, event, user)
	}

	// Check if admin for removing others
	isAdmin, _ := h.communityService.IsAdmin(ctx, user.ID, community.ID)

	if strings.HasPrefix(args, "@") {
		// Remove another user
		return h.removeByUsername(c, ctx, event, user, strings.TrimPrefix(args, "@"), isAdmin)
	}

	// Remove guest
	return h.removeGuest(c, ctx, event, user, args, isAdmin)
}

func (h *BotHandler) removeSelf(c tele.Context, ctx context.Context, event *models.Event, user *models.User) error {
	err := h.participantService.RemoveUser(ctx, event.ID, user.ID)
	if err != nil {
		if err == services.ErrNotParticipant {
			return c.Reply("⚠️ Вы не записаны на эту игру")
		}
		return c.Reply("❌ Ошибка")
	}

	return h.showParticipantsList(c, ctx, event, fmt.Sprintf("✅ %s отписан(а)", safeUsername(user)))
}

func (h *BotHandler) removeByUsername(c tele.Context, ctx context.Context, event *models.Event, remover *models.User, username string, isAdmin bool) error {
	if !isAdmin {
		return c.Reply("❌ Только администратор может удалять других игроков")
	}

	targetUser, err := h.userService.GetByUsername(ctx, username)
	if err != nil {
		return c.Reply(fmt.Sprintf("❌ Пользователь @%s не найден", username))
	}

	err = h.participantService.RemoveUser(ctx, event.ID, targetUser.ID)
	if err != nil {
		if err == services.ErrNotParticipant {
			return c.Reply(fmt.Sprintf("⚠️ @%s не записан(а) на эту игру", username))
		}
		return c.Reply("❌ Ошибка")
	}

	return h.showParticipantsList(c, ctx, event, fmt.Sprintf("✅ @%s отписан(а)", username))
}

func (h *BotHandler) removeGuest(c tele.Context, ctx context.Context, event *models.Event, remover *models.User, guestName string, isAdmin bool) error {
	// Check if remover is the one who added the guest
	guest, err := h.participantService.GetByEventAndGuest(ctx, event.ID, guestName)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.Reply(fmt.Sprintf("❌ Гость \"%s\" не найден", guestName))
		}
		return c.Reply("❌ Ошибка")
	}

	if !isAdmin && (guest.AddedBy == nil || *guest.AddedBy != remover.ID) {
		return c.Reply("❌ Вы можете удалить только гостей, которых добавили сами")
	}

	err = h.participantService.RemoveGuest(ctx, event.ID, guestName)
	if err != nil {
		return c.Reply("❌ Ошибка удаления гостя")
	}

	return h.showParticipantsList(c, ctx, event, fmt.Sprintf("✅ Гость \"%s\" удалён", guestName))
}

func (h *BotHandler) handleClaim(c tele.Context) error {
	ctx := context.Background()
	chat := c.Chat()
	sender := c.Sender()

	// Only works in private chat
	if chat.Type != tele.ChatPrivate {
		return c.Reply("⚠️ Эта команда работает только в личных сообщениях с ботом")
	}

	args := strings.TrimSpace(strings.TrimPrefix(c.Text(), "/claim"))
	if args == "" {
		return c.Reply("⚠️ Укажите токен: /claim <токен>")
	}

	// Get or create user
	user, err := h.userService.Upsert(ctx, models.CreateUserInput{
		TelegramID: sender.ID,
		Username:   stringPtr(sender.Username),
		FirstName:  stringPtr(sender.FirstName),
		LastName:   stringPtr(sender.LastName),
	})
	if err != nil {
		return c.Reply("❌ Ошибка регистрации")
	}

	// Claim the guest spot
	participant, err := h.participantService.ClaimGuest(ctx, args, user.ID)
	if err != nil {
		if err == services.ErrInvalidToken {
			return c.Reply("❌ Неверный или уже использованный токен")
		}
		return c.Reply("❌ Ошибка привязки")
	}

	// Get event info
	event, _ := h.eventService.GetByID(ctx, participant.EventID)
	if event != nil {
		return c.Reply(fmt.Sprintf("✅ Аккаунт привязан!\n\n📅 Вы записаны на игру %s в %s", formatDate(event.EventDate), event.EventTime))
	}

	return c.Reply("✅ Аккаунт успешно привязан!")
}

func (h *BotHandler) showParticipantsList(c tele.Context, ctx context.Context, event *models.Event, header string) error {
	participants, _ := h.participantService.GetByEvent(ctx, event.ID)
	
	msg := header + "\n\n"
	msg += fmt.Sprintf("📅 %s в %s\n", formatDate(event.EventDate), event.EventTime)
	msg += fmt.Sprintf("👥 Участники (%d/%d):\n", len(participants), event.MaxParticipants())

	for i, p := range participants {
		msg += fmt.Sprintf("%d. %s\n", i+1, p.DisplayName())
	}

	if len(participants) == 0 {
		msg += "Пока никого нет. Будьте первым! /add"
	}

	return c.Reply(msg)
}
