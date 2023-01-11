package api

import (
	"encoding/json"
	"net/http"

	"github.com/AleksandrCherepanov/tg-scheduler/internal/notification"
	"github.com/AleksandrCherepanov/tg-scheduler/internal/scheduler"
	"github.com/AleksandrCherepanov/tg-scheduler/internal/server"
	"github.com/AleksandrCherepanov/tg-scheduler/internal/text"
	validator "github.com/go-playground/validator/v10"
)

type NotificationAPI struct {
	ns notification.Storage
	server.Query
}

var notificationAPI *NotificationAPI

func GetNotificationAPI(ns notification.Storage) *NotificationAPI {
	if notificationAPI == nil {
		notificationAPI = &NotificationAPI{ns: ns}
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
	userId, err := api.GetParamInt64("user_id", req)
	if err != nil {
		status := http.StatusUnprocessableEntity
		server.ResponseWithError(res, server.GetResponseError(text.INVALID_QUARY_PARAMS, status))
		return
	}

	result, err := api.ns.GetByUserId(userId)
	if err != nil {
		server.ResponseWithError(res, server.GetResponseError("Notification not found", http.StatusNotFound))
		return
	}

	server.ResponseJson(res, result)
}

func (api *NotificationAPI) CreateNotification(res http.ResponseWriter, req *http.Request) {
	errStatus := http.StatusUnprocessableEntity
	body, ok := server.GetParsedBody(req)
	if !ok {
		server.ResponseWithError(res, server.GetResponseError(text.CANT_GET_BODY, errStatus))
	}

	newNotification := &notification.Notification{}
	err := json.Unmarshal(body, newNotification)
	if err != nil {
		server.ResponseWithError(res, server.GetResponseError(text.INVALID_BODY, errStatus))
		return
	}

	if !checkNotificationBody(res, req, newNotification) {
		return
	}

	_, err = api.ns.GetByUserId(newNotification.UserId)
	if err != nil {
		server.ResponseWithError(res, server.GetResponseError("User not found", http.StatusNotFound))
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

func (api *NotificationAPI) UpdateNotification(res http.ResponseWriter, req *http.Request) {
	body, ok := server.GetParsedBody(req)
	if !ok {
		server.ResponseWithError(
			res,
			server.GetResponseError(text.CANT_GET_BODY, http.StatusUnprocessableEntity),
		)
	}

	userId, err := api.GetParamInt64("user_id", req)
	if err != nil {
		return
	}

	_, err = api.ns.GetByUserId(userId)
	if err != nil {
		server.ResponseWithError(
			res,
			server.GetResponseError("Notification not found", http.StatusNotFound),
		)
		return
	}

	newNotification := &notification.Notification{}
	err = json.Unmarshal(body, newNotification)
	if err != nil {
		server.ResponseWithError(
			res,
			server.GetResponseError(text.INVALID_BODY, http.StatusUnprocessableEntity),
		)
		return
	}

	if !checkNotificationBody(res, req, newNotification) {
		return
	}

	api.ns.UpdateByUserId(
		newNotification.UserId,
		newNotification.Value,
		newNotification.Day,
		newNotification.Hour,
	)
	res.WriteHeader(http.StatusNoContent)
}

func (api NotificationAPI) DeleteNotification(res http.ResponseWriter, req *http.Request) {
	userId, err := api.GetParamInt64("user_id", req)
	if err != nil {
		return
	}

	err = api.ns.DeleteByUserId(userId)
	if err != nil {
		server.ResponseWithError(
			res,
			server.GetResponseError(err.Error(), http.StatusInternalServerError),
		)
		return
	}
	res.WriteHeader(http.StatusNoContent)
}

func checkNotificationBody(
	res http.ResponseWriter,
	req *http.Request,
	n *notification.Notification,
) bool {
	jsonValildator := validator.New()
	err := jsonValildator.Struct(n)
	if err != nil {
		errorList := err.(validator.ValidationErrors)
		errTexts := make([]string, 0, len(errorList))
		for _, e := range errorList {
			errTexts = append(errTexts, e.Error())
		}
		server.ResponseWithError(
			res,
			server.GetResponseError(errTexts, http.StatusUnprocessableEntity),
		)
		return false
	}

	return true
}
