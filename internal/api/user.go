package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/AleksandrCherepanov/tg-scheduler/internal/server"
	"github.com/AleksandrCherepanov/tg-scheduler/internal/text"
	"github.com/AleksandrCherepanov/tg-scheduler/internal/user"
	validator "github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type UserAPI struct {
	us user.Storage
}

var userAPI *UserAPI

func GetUserAPI(us user.Storage) *UserAPI {
	if userAPI == nil {
		userAPI = &UserAPI{us}
	}

	return userAPI
}

func (api *UserAPI) GetUserList(res http.ResponseWriter, req *http.Request) {
	users := api.us.Get()
	if len(users) == 0 {
		server.ResponseJson(res, users)
		return
	}

	userSlice := make([]user.User, 0, len(users))
	for _, u := range users {
		userSlice = append(userSlice, u)
	}

	server.ResponseJson(res, userSlice)
}

func (api *UserAPI) GetUser(res http.ResponseWriter, req *http.Request) {
	userId := getQueryUserId(res, req)
	if userId < 0 {
		return
	}

	result, err := api.us.GetById(userId)
	if err != nil {
		server.ResponseWithError(res, server.GetResponseError(err.Error(), http.StatusNotFound))
		return
	}

	server.ResponseJson(res, result)
}

func (api *UserAPI) CreateUser(res http.ResponseWriter, req *http.Request) {
	body, ok := server.GetParsedBody(req)
	if !ok {
		server.ResponseWithError(
			res,
			server.GetResponseError(text.CANT_GET_BODY, http.StatusUnprocessableEntity),
		)
	}

	newUser := &user.User{}
	err := json.Unmarshal(body, newUser)
	if err != nil {
		server.ResponseWithError(
			res,
			server.GetResponseError(text.INVALID_BODY, http.StatusUnprocessableEntity),
		)
		return
	}

	if !checkUserBody(res, req, newUser) {
		return
	}

	api.us.Create(newUser.Id, newUser.Name)
	res.WriteHeader(http.StatusCreated)
}

func (api *UserAPI) UpdateUser(res http.ResponseWriter, req *http.Request) {
	body, ok := server.GetParsedBody(req)
	if !ok {
		server.ResponseWithError(
			res,
			server.GetResponseError(text.CANT_GET_BODY, http.StatusUnprocessableEntity),
		)
	}

	userId := getQueryUserId(res, req)
	if userId < 0 {
		return
	}

	_, err := api.us.GetById(userId)
	if err != nil {
		server.ResponseWithError(res, server.GetResponseError(err.Error(), http.StatusNotFound))
		return
	}

	newUser := &user.User{}
	err = json.Unmarshal(body, newUser)
	if err != nil {
		server.ResponseWithError(
			res,
			server.GetResponseError(text.INVALID_BODY, http.StatusUnprocessableEntity),
		)
		return
	}

	if !checkUserBody(res, req, newUser) {
		return
	}

	api.us.Update(newUser.Id, newUser.Name)
	res.WriteHeader(http.StatusNoContent)
}

func (api UserAPI) DeleteUser(res http.ResponseWriter, req *http.Request) {
	userId := getQueryUserId(res, req)
	if userId < 0 {
		return
	}

	err := api.us.Delete(userId)
	if err != nil {
		server.ResponseWithError(
			res,
			server.GetResponseError(err.Error(), http.StatusInternalServerError),
		)
		return
	}
	res.WriteHeader(http.StatusNoContent)
}

func getQueryUserId(res http.ResponseWriter, req *http.Request) int64 {
	vars := mux.Vars(req)
	resErr := server.GetResponseError(text.INVALID_QUARY_PARAMS, http.StatusUnprocessableEntity)
	if len(vars) == 0 {
		res.WriteHeader(http.StatusUnprocessableEntity)
		server.ResponseWithError(res, resErr)
		return -1
	}

	if _, ok := vars["id"]; !ok {
		res.WriteHeader(http.StatusUnprocessableEntity)
		server.ResponseWithError(res, resErr)
		return -1
	}

	userId, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		res.WriteHeader(http.StatusUnprocessableEntity)
		server.ResponseWithError(res, resErr)
		return -1
	}
	return userId
}

func checkUserBody(res http.ResponseWriter, req *http.Request, u *user.User) bool {
	jsonValildator := validator.New()
	err := jsonValildator.Struct(u)
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
