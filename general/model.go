package general

import "net/http"

type Endpoint func(w http.ResponseWriter, r *http.Request)

type ResponseStatus string

const (
	StatusOK    ResponseStatus = "ok"
	StatusError ResponseStatus = "error"
)
