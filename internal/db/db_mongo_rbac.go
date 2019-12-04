package db

import (
	"context"

	"github.com/backd-io/backd/backd"
	"github.com/backd-io/backd/internal/constants"
	"github.com/backd-io/backd/internal/pbsessions"
	"github.com/backd-io/backd/internal/rest"
	"github.com/backd-io/backd/internal/structs"
	"go.mongodb.org/mongo-driver/bson"
)

// Can validates the ability to make something by an user
func (db *Mongo) Can(ctx context.Context, session *pbsessions.Session, isDomain bool, database, collection, id string, perm backd.Permission) bool {

	var (
		rbac  structs.RBAC
		query map[string]interface{}
		cid   []string // array of ID inside collections
		err   error
	)

	if isDomain {
		collection = constants.ColDomains
	}

	// if the user can make the action on * can make it for the id
	cid = append(cid, "*")
	if id != "*" {
		cid = append(cid, id)
	}

	query = map[string]interface{}{
		"domain_id": session.GetDomainId(),
		"identity_id": bson.M{
			"$in": db.getIdentities(session),
		},
		"collection": bson.M{
			"$in": []string{collection, "*"},
		},
		"collection_id": bson.M{
			"$in": cid,
		},
		"perm": bson.M{
			"$in": []backd.Permission{perm, backd.PermissionAdmin},
		},
	}

	err = db.GetOne(ctx, database, constants.ColRBAC, query, &rbac)
	if err != nil {
		return false
	}

	return true
}

// VisibleID returns only those IDs that the user is able to see from a collection
func (db *Mongo) VisibleID(ctx context.Context, session *pbsessions.Session, isDomain bool, database, collection string, perm backd.Permission) (all bool, ids []string, err error) {

	if isDomain {
		collection = constants.ColDomains
	}

	var (
		query map[string]interface{}
	)

	query = map[string]interface{}{
		"domain_id": session.GetDomainId(),
		"identity_id": bson.M{
			"$in": db.getIdentities(session),
		},
		"collection": bson.M{
			"$in": []string{collection, "*"},
		},
		"perm": bson.M{
			"$in": []backd.Permission{perm, backd.PermissionAdmin},
		},
	}

	var result []interface{}
	result, err = db.client.Database(database).Collection(collection).Distinct(ctx, "collection_id", query)

	// see if can see all the items to simplify the query
	for _, id := range result {
		// results must be string....
		thisID, ok := id.(string)
		if ok {
			if thisID == "*" {
				all = true
				break
			}
			ids = append(ids, thisID)
		}
	}
	return

}

// getIdentities returns the identities associated with the user.
//   Identities array contains the UserID and all the GroupID where the user belongs.
func (db *Mongo) getIdentities(session *pbsessions.Session) (identities []string) {

	identities = append(identities, session.GetUserId())
	for _, identity := range session.GetGroups() {
		identities = append(identities, identity)
	}
	return

}

// GetOneRBAC returns one object by query
func (db *Mongo) GetOneRBAC(ctx context.Context, session *pbsessions.Session, isDomain bool, perm backd.Permission, database, collection string, query map[string]interface{}) (map[string]interface{}, error) {

	var (
		data map[string]interface{}
		err  error
	)

	err = db.GetOne(ctx, database, collection, query, &data)
	if err != nil {
		return data, err
	}

	if _, ok := data["_id"]; ok {
		return data, constants.ErrItemWithoutID
	}

	// allowed
	if db.Can(ctx, session, isDomain, database, collection, data["_id"].(string), perm) {
		return data, nil
	}

	return nil, rest.ErrUnauthorized
}

// GetOneByIDRBAC returns one object by ID
func (db *Mongo) GetOneByIDRBAC(ctx context.Context, session *pbsessions.Session, isDomain bool, perm backd.Permission, database, collection, id string) (map[string]interface{}, error) {

	var (
		data map[string]interface{}
		err  error
	)

	// allowed
	if db.Can(ctx, session, isDomain, database, collection, id, perm) {
		err = db.GetOneByID(ctx, database, collection, id, &data)
		if err != nil {
			return data, err
		}

		if _, ok := data["_id"]; !ok {
			return data, constants.ErrItemWithoutID
		}

		return data, nil
	}

	return nil, rest.ErrUnauthorized
}

// GetOneByIDRBACInterface returns one object by ID
func (db *Mongo) GetOneByIDRBACInterface(ctx context.Context, session *pbsessions.Session, isDomain bool, perm backd.Permission, database, collection, id string, this interface{}) error {

	// allowed
	if db.Can(ctx, session, isDomain, database, collection, id, perm) {
		return db.GetOneByID(ctx, database, collection, id, this)
	}

	return rest.ErrUnauthorized
}

// GetManyRBAC returns all records that meets RBAC and the desired filter,
//   skip and limit must be passed to limit the number of results
func (db *Mongo) GetManyRBAC(ctx context.Context, session *pbsessions.Session, isDomain bool, perm backd.Permission, database, collection string, query, sort map[string]interface{}, this interface{}, skip, limit int64) error {

	var (
		all          bool
		accesibleIDs []string
		err          error
	)

	if query == nil {
		query = make(map[string]interface{})
	}

	all, accesibleIDs, err = db.VisibleID(ctx, session, isDomain, database, collection, perm)
	// fmt.Println("all, accesibleIDs, err:", all, accesibleIDs, err)
	if err != nil {
		return err
	}

	// restrict only if the user can see a limited amount of items
	if all == false {
		query["_ids"] = bson.M{
			"$in": accesibleIDs,
		}
	}

	return db.GetMany(ctx, database, collection, query, sort, this, skip, limit)

}

// InsertRBAC a new entry on the collection, adding metadata, returns errors if any
func (db *Mongo) InsertRBAC(ctx context.Context, session *pbsessions.Session, isDomain bool, database, collection string, this map[string]interface{}) (map[string]interface{}, error) {

	var err error

	if db.Can(ctx, session, isDomain, database, collection, "*", backd.PermissionCreate) {
		// set metadata
		var metadata structs.Metadata
		metadata.SetCreate(session.GetDomainId(), session.GetUserId())
		this["_id"] = NewXID().String()
		this["meta"] = metadata
		_, err = db.Insert(ctx, database, collection, this)
		return this, err
	}

	return map[string]interface{}{}, rest.ErrUnauthorized
}

// InsertRBACInterface a new entry on the collection, metadata and ID must be written in advance
func (db *Mongo) InsertRBACInterface(ctx context.Context, session *pbsessions.Session, isDomain bool, database, collection string, this interface{}) error {

	var err error

	if db.Can(ctx, session, isDomain, database, collection, "*", backd.PermissionCreate) {
		_, err = db.Insert(ctx, database, collection, this)
		return err
	}

	return rest.ErrUnauthorized

}

// UpdateByIDRBAC updates the data and metadata on the collections, returning errors if any
func (db *Mongo) UpdateByIDRBAC(ctx context.Context, session *pbsessions.Session, isDomain bool, database, collection, id string, this map[string]interface{}) (map[string]interface{}, error) {

	var (
		thisID  string
		ok      bool
		oldData map[string]interface{}
		err     error
	)

	// ensure id has been passed and its the same as sent on the put
	// if id is incorrect or missing then there is a conflict
	thisID, ok = this["_id"].(string)
	if ok != true || thisID != id {
		return nil, rest.ErrConflict
	}

	if db.Can(ctx, session, isDomain, database, collection, id, backd.PermissionUpdate) {

		// first get the old entry
		err = db.GetOneByID(ctx, database, collection, id, &oldData)
		if err != nil {
			return nil, err
		}

		// updated metadata
		var metadata structs.Metadata
		err = metadata.FromInterface(oldData["meta"].(map[string]interface{}))
		metadata.SetUpdate(session.GetDomainId(), session.GetUserId())
		this["meta"] = metadata
		return this, db.UpdateByID(ctx, database, collection, id, this)
	}

	return nil, rest.ErrUnauthorized

}

// UpdateByIDRBACInterface updates the data passed as interface and updates the database if ok
func (db *Mongo) UpdateByIDRBACInterface(ctx context.Context, session *pbsessions.Session, isDomain bool, database, collection, id string, this interface{}) error {

	if db.Can(ctx, session, isDomain, database, collection, id, backd.PermissionUpdate) {
		return db.UpdateByID(ctx, database, collection, id, this)
	}

	return rest.ErrUnauthorized

}

// DeleteByIDRBAC deletes from the item from the collection if user has permission for it
func (db *Mongo) DeleteByIDRBAC(ctx context.Context, session *pbsessions.Session, isDomain bool, database, collection, id string) (int64, error) {

	if db.Can(ctx, session, isDomain, database, collection, id, backd.PermissionDelete) {
		return db.DeleteByID(ctx, database, collection, id)
	}

	return 0, rest.ErrUnauthorized
}

// DeleteByQueryRBAC deletes from the item from the collection if user has permission for it
func (db *Mongo) DeleteByQueryRBAC(ctx context.Context, session *pbsessions.Session, isDomain bool, database, collection string, query map[string]interface{}) (int64, error) {

	if db.Can(ctx, session, isDomain, database, collection, "*", backd.PermissionDelete) {
		return db.Delete(ctx, database, collection, query)
	}

	return 0, rest.ErrUnauthorized
}
