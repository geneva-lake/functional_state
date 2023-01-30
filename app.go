package main

import (
	"net/http"

	"github.com/geneva-lake/functional_state/service"
	"github.com/gorilla/mux"
)

//   - -------------------------------------------------------------------------------------------------------------------
//     Create handlers for http request
//   - -------------------------------------------------------------------------------------------------------------------
func CreateHandlers(config *service.Config) http.Handler {
	s := service.NewService(config)
	e := service.MakeUserEndpoint(s)
	r := mux.NewRouter()
	r.Methods("GET").Path("/users/{id}").HandlerFunc(e)
	return r
}
