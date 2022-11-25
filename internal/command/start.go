package command

import (
	"github.com/AleksandrCherepanov/go_telegram/pkg/telegram"
	"github.com/AleksandrCherepanov/go_telegram/pkg/telegram/client"
	"github.com/AleksandrCherepanov/tg-scheduler/internal/storage"
	"github.com/AleksandrCherepanov/tg-scheduler/internal/text"
	"github.com/AleksandrCherepanov/tg-scheduler/internal/user"
)

type CommandStart struct {
	chatId      int64
	message     *telegram.Message
	userStorage user.Storage
}

func NewCommandStart(chatId int64, message *telegram.Message) *CommandStart {
	return &CommandStart{
		chatId:      chatId,
		message:     message,
		userStorage: user.GetUserStorage(storage.GetRedisClient()),
	}
}

func (c *CommandStart) Handle(command string, args []string) (interface{}, error) {
	c.userStorage.Create(c.chatId, c.message.Chat.GetName())
	return client.NewTelegramResponse(c.chatId, text.USER_CREATED, false), nil
}
