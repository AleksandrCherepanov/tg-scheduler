package command

import (
	"strings"

	"github.com/AleksandrCherepanov/go_telegram/pkg/telegram"
	"github.com/AleksandrCherepanov/go_telegram/pkg/telegram/client"
	"github.com/AleksandrCherepanov/tg-scheduler/internal/text"
)

const unknownCommand = "/unknown"

type HandlerInterface interface {
	Handle(chatId int64, message *telegram.Message) (interface{}, error)
}

type CommandHandlerInterface interface {
	Handle(command string, args []string) (interface{}, error)
}

type CommandHandler struct {
	handlers map[string]CommandHandlerInterface
}

func NewCommandHandler(chatId int64, message *telegram.Message) *CommandHandler {
	commandHandler := &CommandHandler{}
	commandHandler.handlers = map[string]CommandHandlerInterface{
		"/start": NewCommandStart(chatId, message),
		"/set":   NewCommandSet(chatId, message),
		"/help":  NewCommandHelp(chatId),
	}

	return commandHandler
}

func (commandHandler *CommandHandler) Handle(chatId int64, message *telegram.Message) (interface{}, error) {
	commandWithArgs := strings.Split(*message.Text, " ")

	if len(commandWithArgs) == 0 {
		return nil, client.NewTelegramResponse(chatId, text.INVALID_COMMAND, true)
	}

	handler, ok := commandHandler.handlers[commandWithArgs[0]]
	if !ok {
		return nil, client.NewTelegramResponse(chatId, text.INVALID_COMMAND, true)
	}

	return handler.Handle(commandWithArgs[0], commandWithArgs[1:])
}
