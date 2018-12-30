package db

import "github.com/rs/xid"

// Interface is the interface that define all methods to operate the databases
type Interface interface {
	IsInitialized() error
	Insert(database, collection string, this interface{}) error
	Count(database, collection string, query interface{}) (int, error)
	GetMany(database, collection string, query interface{}, sort []string, this interface{}, skip, limit int) error
	GetOne(database, collection string, query, this interface{}) error
	GetOneByID(database, collection, id string, this interface{}) error
	Update(database, collection string, selector, to interface{}) error
	UpdateByID(database, collection, id string, to interface{}) error
	Delete(database, collection string, selector interface{}) error
	DeleteByID(database, collection, id string) error
	CreateIndex(database, collection string, fields []string, unique bool) error
	CreateDefaultAppIndexes() error
	CreateDefaultDomainIndexes() error
}

// NewXID returns a new secure ID using the rs/xid librady
func NewXID() xid.ID {
	return xid.New()
}
