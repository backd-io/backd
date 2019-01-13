package lua

import "errors"

// common errors
var (
	ErrApplicationNotEspecified = errors.New("Application ID not especified")
)

const (
	noAppID = "no_app"
)
