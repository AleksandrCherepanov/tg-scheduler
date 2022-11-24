package command

import (
	"strings"

	"github.com/AleksandrCherepanov/go_telegram/pkg/telegram"
	"github.com/AleksandrCherepanov/go_telegram/pkg/telegram/client"
	"github.com/AleksandrCherepanov/tg-scheduler/internal/notification"
	"github.com/AleksandrCherepanov/tg-scheduler/internal/scheduler"
	"github.com/AleksandrCherepanov/tg-scheduler/internal/storage"
	"github.com/AleksandrCherepanov/tg-scheduler/internal/text"
	"github.com/AleksandrCherepanov/tg-scheduler/internal/user"
)

type CommandSet struct {
	chatId              int64
	message             *telegram.Message
	userStorage         user.Storage
	notificationStorage notification.Storage
}

func NewCommandSet(chatId int64, message *telegram.Message) *CommandSet {
	return &CommandSet{
		chatId:              chatId,
		message:             message,
		userStorage:         user.GetUserStorage(storage.GetRedisClient()),
		notificationStorage: notification.GetNotificationRedisStorage(storage.GetRedisClient()),
	}
}

func (c *CommandSet) Handle(command string, args []string) (interface{}, error) {
	userId, err := c.message.GetChatId()
	if err != nil {
		return nil, err
	}

	user, err := c.userStorage.GetById(userId)
	if err != nil {
		return nil, err
	}

	_, err = c.notificationStorage.UpdateByUserId(
		user.Id,
		strings.Join(args, " "),
		scheduler.NotificationDay,
		scheduler.NotificationUtcHour,
	)

	if err != nil {
		return nil, err
	}

	return client.NewTelegramResponse(c.chatId, text.NOTIFICATION_IS_SET, false), nil
}
