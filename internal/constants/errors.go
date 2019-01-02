package constants

import "errors"

var (
	// ErrBadConfiguration is returned when the configuration is not properly filled
	ErrBadConfiguration = errors.New("bad configuration")
)
