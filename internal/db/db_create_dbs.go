package db

import (
	"fmt"

	"github.com/fernandezvara/backd/internal/constants"
	"github.com/fernandezvara/backd/internal/structs"
	mgo "github.com/globalsign/mgo"
)

// CreateCollection creates a collection and its indexes on the required database
//   adding the validation if required
func (db *Mongo) CreateCollection(database, collectionName string, validator map[string]interface{}, indexes []structs.Index) error {

	var (
		colInfo mgo.CollectionInfo
		err     error
	)

	if validator != nil {
		colInfo.Validator = validator
	}

	err = db.Session().DB(database).C(collectionName).Create(&colInfo)
	if err != nil {
		return err
	}

	for _, index := range indexes {
		err = db.CreateIndex(database, collectionName, index.Fields, index.Unique)
		if err != nil {
			return err
		}
	}

	return nil

}

// CreateBackdDatabases is called on bootstrap to create
func (db *Mongo) CreateBackdDatabases() (err error) {

	err = db.CreateCollection(constants.DBBackdApp, constants.ColApplications, structs.ApplicationValidator(), structs.ApplicationIndexes)
	if err != nil {
		return
	}

	err = db.CreateCollection(constants.DBBackdApp, constants.ColDomains, structs.DomainValidator(), structs.DomainIndexes)
	if err != nil {
		return
	}

	err = db.CreateApplicationDatabase(constants.DBBackdApp)
	if err != nil {
		return
	}

	err = db.CreateDomainDatabase(constants.DBBackdDom)
	if err != nil {
		return
	}

	// Add the domain to the _Domains to allow registration
	var thisDomain structs.Domain

	thisDomain.ID = constants.DBBackdDom
	thisDomain.Type = structs.DomainTypeBackd
	thisDomain.Description = "backd domain"
	thisDomain.SetCreate(constants.DBBackdDom, constants.SystemUserID) // system ID (this domain must no mutate over the time)

	err = db.Insert(constants.DBBackdApp, constants.ColDomains, &thisDomain)
	return

}

// CreateApplicationDatabase creates the required collection for an application to be usable
func (db *Mongo) CreateApplicationDatabase(name string) (err error) {

	err = db.CreateCollection(name, constants.ColRelations, structs.RelationValidator(), structs.RelationIndexes)
	if err != nil {
		return
	}

	err = db.CreateCollection(name, constants.ColFunctions, structs.FunctionValidator(), structs.FunctionIndexes)
	if err != nil {
		fmt.Println("err:", err)
		return
	}

	err = db.CreateCollection(name, constants.ColRBAC, structs.RBACValidator(), structs.RBACIndexes)
	return

}

// CreateDomainDatabase creates the required collections on the domain to be usable
func (db *Mongo) CreateDomainDatabase(name string) (err error) {

	err = db.CreateCollection(name, constants.ColRBAC, structs.RBACValidator(), structs.RBACIndexes)
	if err != nil {
		return
	}

	err = db.CreateCollection(name, constants.ColUsers, structs.UserValidator(), structs.UserIndexes)
	if err != nil {
		return
	}

	err = db.CreateCollection(name, constants.ColGroups, structs.GroupValidator(), structs.GroupIndexes)
	if err != nil {
		return
	}

	err = db.CreateCollection(name, constants.ColMembership, structs.MembershipValidator(), structs.MembershipIndexes)
	return

}
