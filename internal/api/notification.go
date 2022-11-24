package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/AleksandrCherepanov/tg-scheduler/internal/notification"
	"github.com/AleksandrCherepanov/tg-scheduler/internal/scheduler"
	"github.com/AleksandrCherepanov/tg-scheduler/internal/server"
	"github.com/AleksandrCherepanov/tg-scheduler/internal/text"
	"github.com/gorilla/mux"
)

type NotificationAPI struct {
	ns notification.Storage
}

var notificationAPI *NotificationAPI

func GetNotificationAPI(ns notification.Storage) *NotificationAPI {
	if notificationAPI == nil {
		notificationAPI = &NotificationAPI{ns}
	}

	return notificationAPI
}

func (api *NotificationAPI) GetNotificationList(res http.ResponseWriter, req *http.Request) {
	notifications := api.ns.Get()
	if len(notifications) == 0 {
		server.ResponseJson(res, notifications)
		return
	}

	notificationSlice := make([]notification.Notification, 0, len(notifications))
	for _, u := range notifications {
		notificationSlice = append(notificationSlice, u)
	}

	server.ResponseJson(res, notificationSlice)
}

func (api *NotificationAPI) GetNotification(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	if len(vars) == 0 {
		res.WriteHeader(http.StatusUnprocessableEntity)
		server.ResponseError(res, text.INVALID_QUARY_PARAMS)
		return
	}

	if _, ok := vars["user_id"]; !ok {
		res.WriteHeader(http.StatusUnprocessableEntity)
		server.ResponseError(res, text.INVALID_QUARY_PARAMS)
		return
	}

	userId, err := strconv.ParseInt(vars["user_id"], 10, 64)
	if err != nil {
		res.WriteHeader(http.StatusUnprocessableEntity)
		server.ResponseError(res, text.INVALID_QUARY_PARAMS)
		return
	}

	result, err := api.ns.GetByUserId(userId)
	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		server.ResponseError(res, err.Error())
		return
	}

	server.ResponseJson(res, result)
}

func (api *NotificationAPI) CreateNotification(res http.ResponseWriter, req *http.Request) {
	body, ok := server.GetParsedBody(req)
	if !ok {
		server.ResponseError(res, text.CANT_GET_BODY)
	}

	newNotification := &notification.Notification{}
	err := json.Unmarshal(body, newNotification)
	if err != nil {
		res.WriteHeader(http.StatusUnprocessableEntity)
		server.ResponseError(res, text.INVALID_BODY)
		return
	}

	api.ns.CreateByUserId(
		newNotification.UserId,
		newNotification.Value,
		scheduler.NotificationDay,
		scheduler.NotificationUtcHour,
	)

	res.WriteHeader(http.StatusCreated)
}
