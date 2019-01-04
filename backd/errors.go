package backd

import (
	"errors"
	"net/http"
)

// HTTP Errors returns if found
var (
	ErrBadRequest       = errors.New(http.StatusText(http.StatusBadRequest))
	ErrUnauthorized     = errors.New(http.StatusText(http.StatusUnauthorized))
	ErrNotFound         = errors.New(http.StatusText(http.StatusNotFound))
	ErrMethodNotAllowed = errors.New(http.StatusText(http.StatusMethodNotAllowed))
	ErrConflict         = errors.New(http.StatusText(http.StatusConflict))
)
