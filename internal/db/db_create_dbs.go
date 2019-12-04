package db

import (
	"context"

	"github.com/backd-io/backd/internal/constants"
	"github.com/backd-io/backd/internal/structs"
	"go.mongodb.org/mongo-driver/bson"
)

// CreateCollection creates a collection and its indexes on the required database
//   adding the validation if required
func (db *Mongo) CreateCollection(ctx context.Context, database, collection string, validator map[string]interface{}, indexes []structs.Index) (err error) {

	if validator != nil {
		command := bson.D{{"create", collection}, {"validator", validator}, {"validationLevel", "strict"}}

		err = db.client.Database(database).RunCommand(ctx, command).Err()
		if err != nil {
			return
		}
	}

	for _, index := range indexes {
		err = db.CreateIndex(ctx, database, collection, index.Fields, index.Unique)
		if err != nil {
			return
		}
	}

	return

}

// CreateBackdDatabases is called on bootstrap to create
func (db *Mongo) CreateBackdDatabases(ctx context.Context) (err error) {

	err = db.CreateCollection(ctx, constants.DBBackdApp, constants.ColApplications, structs.ApplicationValidator(), structs.ApplicationIndexes)
	if err != nil {
		return
	}

	err = db.CreateCollection(ctx, constants.DBBackdApp, constants.ColDomains, structs.DomainValidator(), structs.DomainIndexes)
	if err != nil {
		return
	}

	err = db.CreateApplicationDatabase(ctx, constants.DBBackdApp)
	if err != nil {
		return
	}

	err = db.CreateDomainDatabase(ctx, constants.DBBackdDom)
	if err != nil {
		return
	}

	// Add the domain to the _Domains to allow registration
	var thisDomain structs.Domain

	thisDomain.ID = constants.DBBackdDom
	thisDomain.Type = structs.DomainTypeBackd
	thisDomain.Description = "backd domain"
	thisDomain.SetCreate(constants.DBBackdDom, constants.SystemUserID) // system ID (this domain must no mutate over the time)

	_, err = db.Insert(ctx, constants.DBBackdApp, constants.ColDomains, &thisDomain)
	return

}

// CreateApplicationDatabase creates the required collection for an application to be usable
func (db *Mongo) CreateApplicationDatabase(ctx context.Context, name string) (err error) {

	err = db.CreateCollection(ctx, name, constants.ColRelations, structs.RelationValidator(), structs.RelationIndexes)
	if err != nil {
		return
	}

	err = db.CreateCollection(ctx, name, constants.ColFunctions, structs.FunctionValidator(), structs.FunctionIndexes)
	if err != nil {
		return
	}

	err = db.CreateCollection(ctx, name, constants.ColRBAC, structs.RBACValidator(), structs.RBACIndexes)
	return

}

// CreateDomainDatabase creates the required collections on the domain to be usable
func (db *Mongo) CreateDomainDatabase(ctx context.Context, name string) (err error) {

	err = db.CreateCollection(ctx, name, constants.ColRBAC, structs.RBACValidator(), structs.RBACIndexes)
	if err != nil {
		return
	}

	err = db.CreateCollection(ctx, name, constants.ColUsers, structs.UserValidator(), structs.UserIndexes)
	if err != nil {
		return
	}

	err = db.CreateCollection(ctx, name, constants.ColGroups, structs.GroupValidator(), structs.GroupIndexes)
	if err != nil {
		return
	}

	err = db.CreateCollection(ctx, name, constants.ColMembership, structs.MembershipValidator(), structs.MembershipIndexes)
	return

}
