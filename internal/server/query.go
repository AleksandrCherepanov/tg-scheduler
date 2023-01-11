package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/AleksandrCherepanov/tg-scheduler/internal/text"
	"github.com/gorilla/mux"
)

type Query struct {
}

func (q *Query) GetParamInt64(name string, req *http.Request) (int64, error) {
	vars := mux.Vars(req)
	if len(vars) == 0 {
		return 0, fmt.Errorf("%s: empty", text.INVALID_QUARY_PARAMS)
	}

	if _, ok := vars[name]; !ok {
		return 0, fmt.Errorf("%s: not found", text.INVALID_QUARY_PARAMS)
	}

	param, err := strconv.ParseInt(vars[name], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("%s: not int64", text.INVALID_QUARY_PARAMS)
	}

	return param, nil
}
