package bot

import (
	"context"
	"fmt"

	"github.com/luvvano/football-bot/backend/internal/models"
	tele "gopkg.in/telebot.v3"
)

func (h *BotHandler) handleStart(c tele.Context) error {
	ctx := context.Background()
	chat := c.Chat()
	sender := c.Sender()

	// Upsert the user
	user, err := h.userService.Upsert(ctx, models.CreateUserInput{
		TelegramID: sender.ID,
		Username:   stringPtr(sender.Username),
		FirstName:  stringPtr(sender.FirstName),
		LastName:   stringPtr(sender.LastName),
	})
	if err != nil {
		return c.Reply("❌ Ошибка регистрации пользователя")
	}

	// Handle private chat
	if chat.Type == tele.ChatPrivate {
		return h.handlePrivateStart(c, user)
	}

	// Handle group chat
	return h.handleGroupStart(c, ctx, chat, user)
}

func (h *BotHandler) handlePrivateStart(c tele.Context, user *models.User) error {
	name := "Игрок"
	if user.FirstName != nil {
		name = *user.FirstName
	}

	msg := fmt.Sprintf(`👋 Привет, %s!

⚽ Я Football Bot — помогаю организовывать футбольные матчи.

Чтобы начать:
1. Добавьте меня в группу
2. Напишите /start в группе
3. Первый, кто это сделает, станет администратором

Или используйте MiniApp для управления играми:`, name)

	// Create WebApp button
	btn := &tele.ReplyMarkup{}
	webAppBtn := btn.WebApp("⚽ Открыть приложение", &tele.WebApp{URL: h.webAppURL})
	btn.Inline(btn.Row(webAppBtn))

	return c.Send(msg, btn)
}

func (h *BotHandler) handleGroupStart(c tele.Context, ctx context.Context, chat *tele.Chat, user *models.User) error {
	// Check if community exists
	community, err := h.communityService.GetByTelegramChatID(ctx, chat.ID)
	
	if err != nil {
		// Create new community
		community, err = h.communityService.Create(ctx, chat.ID, chat.Title, user.ID)
		if err != nil {
			return c.Reply("❌ Ошибка создания сообщества")
		}

		// Add user as admin
		_, err = h.communityService.AddMember(ctx, user.ID, community.ID, "admin")
		if err != nil {
			return c.Reply("❌ Ошибка добавления администратора")
		}

		msg := fmt.Sprintf(`✅ Сообщество "%s" зарегистрировано!

👑 @%s назначен администратором.

📋 Основные команды:
/event [дата] [время] — создать игру
/add [имя] — записаться (или добавить гостя)
/remove — отписаться от игры
/info — текущая игра
/events — список игр

🔧 Админ-команды:
/venue [название | адрес] — создать площадку
/cancel — отменить игру
/finish — завершить игру

/app — открыть приложение`, 
			community.Name, 
			safeUsername(user))

		return c.Reply(msg)
	}

	// Community exists, just add user as player if not member
	member, _ := h.communityService.GetMember(ctx, user.ID, community.ID)
	if member == nil {
		_, err = h.communityService.AddMember(ctx, user.ID, community.ID, "player")
		if err != nil {
			return c.Reply("❌ Ошибка добавления в сообщество")
		}
	}

	msg := fmt.Sprintf(`👋 Добро пожаловать в "%s"!

Используйте /info для информации о текущей игре.`, community.Name)

	btn := &tele.ReplyMarkup{}
	webAppBtn := btn.WebApp("⚽ Открыть приложение", &tele.WebApp{URL: h.webAppURL + "?community=" + fmt.Sprint(community.ID)})
	btn.Inline(btn.Row(webAppBtn))

	return c.Send(msg, btn)
}

func safeUsername(user *models.User) string {
	if user.Username != nil {
		return *user.Username
	}
	if user.FirstName != nil {
		return *user.FirstName
	}
	return "Unknown"
}

func (h *BotHandler) handleHelp(c tele.Context) error {
	chat := c.Chat()

	if chat.Type == tele.ChatPrivate {
		return c.Reply(`📋 Доступные команды:

👤 В личке:
/events — ваши предстоящие игры
/venue — площадки ваших сообществ
/claim <токен> — привязать гостевую запись
/app — открыть приложение

👥 В группе:
/event [дата] [время] — создать игру
/add [имя] — записаться (или добавить гостя)
/remove — отписаться от игры
/info — текущая игра
/events — список игр

🔧 Админ-команды (в группе):
/venue [название | адрес] — создать площадку
/cancel — отменить игру
/finish — завершить игру`)
	}

	return c.Reply(`📋 Команды:

/event [дата] [время] — создать игру
/add [имя] — записаться (или добавить гостя)
/remove — отписаться от игры
/info — текущая игра
/events — список игр
/venue — список площадок

🔧 Админ:
/venue [название | адрес] — создать площадку
/cancel — отменить игру
/finish — завершить игру

/app — открыть приложение`)
}
