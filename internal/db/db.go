package db

import (
	"github.com/rs/xid"
)

// NewXID returns a new secure ID using the rs/xid librady
func NewXID() xid.ID {
	return xid.New()
}
