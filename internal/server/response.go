package server

import (
	"encoding/json"
	"log"
	"net/http"
)

type ResponseError struct {
	Error      interface{} `json:"error"`
	StatusCode int         `json:"-"`
}

func GetResponseError(err interface{}, statusCode int) *ResponseError {
	return &ResponseError{
		err,
		statusCode,
	}
}

func ResponseJson(w http.ResponseWriter, data interface{}) {
	jsonData, err := json.Marshal(data)
	log.Println(string(jsonData))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	statusCode := 200
	res, ok := data.(ResponseError)
	if ok {
		statusCode = res.StatusCode
	}

	w.WriteHeader(statusCode)
	w.Write(jsonData)
}

func ResponseWithError(w http.ResponseWriter, res *ResponseError) {
	ResponseJson(w, res)
}
