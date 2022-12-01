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
	vars := mux.Vars(req)
	resErr := server.GetResponseError(text.INVALID_QUARY_PARAMS, 422)
	if len(vars) == 0 {
		res.WriteHeader(http.StatusUnprocessableEntity)
		server.ResponseWithError(res, resErr)
		return
	}

	if _, ok := vars["id"]; !ok {
		res.WriteHeader(http.StatusUnprocessableEntity)
		server.ResponseWithError(res, resErr)
		return
	}

	userId, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		res.WriteHeader(http.StatusUnprocessableEntity)
		server.ResponseWithError(res, resErr)
		return
	}

	result, err := api.us.GetById(userId)
	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		server.ResponseWithError(res, server.GetResponseError(err.Error(), 500))
		return
	}

	server.ResponseJson(res, result)
}

func (api *UserAPI) CreateUser(res http.ResponseWriter, req *http.Request) {
	body, ok := server.GetParsedBody(req)
	if !ok {
		server.ResponseWithError(res, server.GetResponseError(text.CANT_GET_BODY, 422))
	}

	newUser := &user.User{}
	err := json.Unmarshal(body, newUser)
	if err != nil {
		res.WriteHeader(http.StatusUnprocessableEntity)
		server.ResponseWithError(res, server.GetResponseError(text.INVALID_BODY, 422))
		return
	}

	jsonValildator := validator.New()
	err = jsonValildator.Struct(newUser)
	if err != nil {
		errorList := err.(validator.ValidationErrors)
		errTexts := make([]string, 0, len(errorList))
		for _, e := range errorList {
			errTexts = append(errTexts, e.Error())
		}
		server.ResponseWithError(res, server.GetResponseError(errTexts, 422))
		return
	}
	api.us.Create(newUser.Id, newUser.Name)
	res.WriteHeader(http.StatusCreated)
}
