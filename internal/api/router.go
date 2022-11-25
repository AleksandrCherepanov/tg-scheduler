package api

import "net/http"

type Router struct {
}

func NewRouter() *Router {
	return &Router{}
}

func (r *Router) Resolve(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("API Response"))
}
