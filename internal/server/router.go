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
		ResponseWithError(w, GetResponseError(text.CANT_GET_BODY, 422))
		return
	}

	update := &telegram.Update{}
	err := json.Unmarshal(body, update)
	if err != nil {
		ResponseWithError(w, GetResponseError(err.Error(), 422))
		return
	}

	message := update.Message
	if message == nil {
		ResponseWithError(w, GetResponseError(text.CANT_PROCESS_MESSAGE, 422))
		return
	}

	if message.Entities == nil {
		ResponseWithError(w, GetResponseError(text.CANT_PROCESS_MESSAGE, 422))
		return
	}

	chatId, err := update.Message.GetChatId()
	if err != nil {
		ResponseWithError(w, GetResponseError(text.CANT_PROCESS_MESSAGE, 422))
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
				ResponseWithError(w, GetResponseError(err.Error(), 500))
				return
			}
		}
		ResponseWithError(w, GetResponseError(handleError.Error(), 500))
		return
	}

	if tgResult, ok := result.(client.TelegramResponse); ok {
		res, err := tgResult.Send()
		log.Printf("%v\n", string(res.([]byte)))
		if err != nil {
			ResponseWithError(w, GetResponseError(err.Error(), 500))
			return
		}
	}
	ResponseJson(w, result)
}
