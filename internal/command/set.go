package command

import (
	"strings"

	"github.com/AleksandrCherepanov/go_telegram/pkg/telegram"
	"github.com/AleksandrCherepanov/go_telegram/pkg/telegram/client"
	"github.com/AleksandrCherepanov/tg-scheduler/internal/notification"
	"github.com/AleksandrCherepanov/tg-scheduler/internal/text"
	"github.com/AleksandrCherepanov/tg-scheduler/internal/user"
)

type CommandSet struct {
	chatId      int64
	message     *telegram.Message
	userStorage *user.UserStorage
}

func NewCommandSet(chatId int64, message *telegram.Message) *CommandSet {
	return &CommandSet{
		chatId:      chatId,
		message:     message,
		userStorage: user.GetUserStorage(),
	}
}

func (c *CommandSet) Handle(command string, args []string) (interface{}, error) {
	userId, err := c.message.GetChatId()
	if err != nil {
		return err, nil
	}

	user, err := c.userStorage.GetUserById(userId)
	if err != nil {
		return err, nil
	}
	ns := notification.GetNotificationStorage()
	ns.Set(user.Id, strings.Join(args, " "))
	return client.NewTelegramResponse(c.chatId, text.NOTIFICATION_IS_SET, false), nil
}
