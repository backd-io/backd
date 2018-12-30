package rest

import (
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// HTTP errors
var (
	ErrBadRequest   = errors.New("400")
	ErrUnauthorized = errors.New("401")
	ErrNotFound     = errors.New("404")
	ErrConflict     = errors.New("409")
)

// HTTPErrorStatus is the struct that user will receive as body. It can contain a reason
//  if we are controlling what to explain about the error.
type HTTPErrorStatus struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Reason  string `json:"reason,omitempty"`
}

// APIHandler defines the handler function that will be used for each API endpoint
type APIHandler func(w http.ResponseWriter, r *http.Request, ps httprouter.Params)

// APIMatcher is an array of rules that must be applied for each request to ensure
//   the required input is being passed to the API
type APIMatcher []string
