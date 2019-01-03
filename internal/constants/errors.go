package constants

import "errors"

var (
	// ErrBadConfiguration is returned when the configuration is not properly filled
	ErrBadConfiguration = errors.New("bad configuration")
)

// Data errors
var (
	// ErrItemWithoutID is returned when an item does not have an ID, that must have it
	ErrItemWithoutID = errors.New("item without id")
)
