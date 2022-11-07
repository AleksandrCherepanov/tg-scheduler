package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/AleksandrCherepanov/go_telegram/pkg/telegram"
	"github.com/AleksandrCherepanov/go_telegram/pkg/telegram/client"
	"github.com/AleksandrCherepanov/tg-scheduler/internal/command"
	"github.com/AleksandrCherepanov/tg-scheduler/internal/text"
)

type Router struct {
	handlers map[string]command.HandlerInterface
}

func NewRouter() *Router {
	return &Router{}
}

func (router *Router) WithHandlers(chatId int64, message *telegram.Message) *Router {
	router.handlers = map[string]command.HandlerInterface{
		"command": command.NewCommandHandler(chatId, message),
	}
	return router
}

func (router *Router) Resolve(w http.ResponseWriter, req *http.Request) {
	body, ok := GetParsedBody(req)
	if !ok {
		ResponseError(w, text.CANT_GET_BODY)
		return
	}

	update := &telegram.Update{}
	err := json.Unmarshal(body, update)
	if err != nil {
		ResponseError(w, err.Error())
		return
	}

	message := update.Message
	if message == nil {
		ResponseError(w, text.CANT_PROCESS_MESSAGE)
		return
	}

	if message.Entities == nil {
		ResponseError(w, text.CANT_PROCESS_MESSAGE)
		return
	}

	chatId, err := update.Message.GetChatId()
	if err != nil {
		ResponseError(w, text.CANT_PROCESS_MESSAGE)
	}

	var result interface{}
	var handleError error
	for _, entity := range *&message.Entities {
		if entity.IsCommand() {
			router = router.WithHandlers(chatId, message)
			result, handleError = router.handlers["command"].Handle(chatId, update.Message)
		}
	}

	if handleError != nil {
		tgResponse, ok := handleError.(client.TelegramResponse)
		if ok {
			res, err := tgResponse.Send()
			log.Printf("%v\n", string(res.([]byte)))
			if err != nil {
				ResponseError(w, err.Error())
				return
			}
		}
		ResponseError(w, handleError.Error())
		return
	}

	if tgResult, ok := result.(client.TelegramResponse); ok {
		res, err := tgResult.Send()
		log.Printf("%v\n", string(res.([]byte)))
		if err != nil {
			ResponseError(w, err.Error())
			return
		}
	}
	ResponseJson(w, result)
}
