package command

import (
	"github.com/AleksandrCherepanov/go_telegram/pkg/telegram"
	"github.com/AleksandrCherepanov/go_telegram/pkg/telegram/client"
	"github.com/AleksandrCherepanov/tg-scheduler/internal/user"
)

type CommandStart struct {
	chatId      int64
	message     *telegram.Message
	userStorage *user.UserStorage
}

func NewCommandStart(chatId int64, message *telegram.Message) *CommandStart {
	return &CommandStart{
		chatId:      chatId,
		message:     message,
		userStorage: user.GetUserStorage(),
	}
}

func (c *CommandStart) Handle(command string, args []string) (interface{}, error) {
	c.userStorage.CreateUser(c.chatId, c.message.Chat.GetName())
	return client.NewTelegramResponse(c.chatId, "scheduler is started", false), nil
}
