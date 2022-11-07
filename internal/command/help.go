package command

import (
	"github.com/AleksandrCherepanov/go_telegram/pkg/telegram/client"
)

type CommandHelp struct {
	chatId int64
}

func NewCommandHelp(chatId int64) *CommandHelp {
	return &CommandHelp{
		chatId: chatId,
	}
}

func (c *CommandHelp) Handle(command string, args []string) (interface{}, error) {
	text := "This is notification bot\\.\n\n"
	text += "Command list:\n"
	text += "1\\. /start \\- use this command for user initialization\\.\n"
	text += "2\\. /set \\- use this command for setting some value for notification\\.\n"
	text += "Format is: `/set some value`\\. This `some value` will be sent to you every Monday at 04:00 GMT\\+0"
	return client.NewTelegramResponse(c.chatId, text, false), nil
}
