package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/AleksandrCherepanov/tg-scheduler/internal/server"
	"github.com/AleksandrCherepanov/tg-scheduler/internal/text"
	"github.com/AleksandrCherepanov/tg-scheduler/internal/user"
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
	vars := mux.Vars(req)
	if len(vars) == 0 {
		res.WriteHeader(http.StatusUnprocessableEntity)
		server.ResponseError(res, text.INVALID_QUARY_PARAMS)
		return
	}

	if _, ok := vars["id"]; !ok {
		res.WriteHeader(http.StatusUnprocessableEntity)
		server.ResponseError(res, text.INVALID_QUARY_PARAMS)
		return
	}

	userId, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		res.WriteHeader(http.StatusUnprocessableEntity)
		server.ResponseError(res, text.INVALID_QUARY_PARAMS)
		return
	}

	result, err := api.us.GetById(userId)
	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		server.ResponseError(res, err.Error())
		return
	}

	server.ResponseJson(res, result)
}

func (api *UserAPI) CreateUser(res http.ResponseWriter, req *http.Request) {
	body, ok := server.GetParsedBody(req)
	if !ok {
		server.ResponseError(res, text.CANT_GET_BODY)
	}

	newUser := &user.User{}
	err := json.Unmarshal(body, newUser)
	if err != nil {
		res.WriteHeader(http.StatusUnprocessableEntity)
		server.ResponseError(res, text.INVALID_BODY)
		return
	}

	api.us.Create(newUser.Id, newUser.Name)
	res.WriteHeader(http.StatusCreated)
}
