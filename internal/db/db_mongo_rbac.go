package db

import (
	"github.com/backd-io/backd/backd"
	"github.com/backd-io/backd/internal/constants"
	"github.com/backd-io/backd/internal/pbsessions"
	"github.com/backd-io/backd/internal/rest"
	"github.com/backd-io/backd/internal/structs"
	"github.com/globalsign/mgo/bson"
)

// Can validates the ability to make something by an user
func (db *Mongo) Can(session *pbsessions.Session, isDomain bool, database, collection, id string, perm backd.Permission) bool {

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
		"did": session.GetDomainId(),
		"iid": bson.M{
			"$in": db.getIdentities(session),
		},
		"c": bson.M{
			"$in": []string{collection, "*"},
		},
		"cid": bson.M{
			"$in": cid,
		},
		"p": bson.M{
			"$in": []backd.Permission{perm, backd.PermissionAdmin},
		},
	}

	err = db.GetOne(database, constants.ColRBAC, query, &rbac)
	if err != nil {
		return false
	}

	return true
}

// VisibleID returns only those IDs that the user is able to see from a collection
func (db *Mongo) VisibleID(session *pbsessions.Session, isDomain bool, database, collection string, perm backd.Permission) (all bool, ids []string, err error) {

	if isDomain {
		collection = constants.ColDomains
	}

	var (
		query map[string]interface{}
	)

	query = map[string]interface{}{
		"did": session.GetDomainId(),
		"iid": bson.M{
			"$in": db.getIdentities(session),
		},
		"c": bson.M{
			"$in": []string{collection, "*"},
		},
		"p": bson.M{
			"$in": []backd.Permission{perm, backd.PermissionAdmin},
		},
	}

	err = db.session.DB(database).C(constants.ColRBAC).Find(query).Distinct("cid", &ids)

	// see if can see all the items to simplify the query
	for _, id := range ids {
		if id == "*" {
			all = true
			break
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
func (db *Mongo) GetOneRBAC(session *pbsessions.Session, isDomain bool, perm backd.Permission, database, collection string, query map[string]interface{}) (map[string]interface{}, error) {

	var (
		data map[string]interface{}
		err  error
	)

	err = db.GetOne(database, collection, query, &data)
	if err != nil {
		return data, err
	}

	if _, ok := data["_id"]; ok {
		return data, constants.ErrItemWithoutID
	}

	// allowed
	if db.Can(session, isDomain, database, collection, data["_id"].(string), perm) {
		return data, nil
	}

	return nil, rest.ErrUnauthorized
}

// GetOneByIDRBAC returns one object by ID
func (db *Mongo) GetOneByIDRBAC(session *pbsessions.Session, isDomain bool, perm backd.Permission, database, collection, id string) (map[string]interface{}, error) {

	var (
		data map[string]interface{}
		err  error
	)

	// allowed
	if db.Can(session, isDomain, database, collection, id, perm) {
		err = db.GetOneByID(database, collection, id, &data)
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
func (db *Mongo) GetOneByIDRBACInterface(session *pbsessions.Session, isDomain bool, perm backd.Permission, database, collection, id string, this interface{}) error {

	// allowed
	if db.Can(session, isDomain, database, collection, id, perm) {
		return db.GetOneByID(database, collection, id, this)
	}

	return rest.ErrUnauthorized
}

// GetManyRBAC returns all records that meets RBAC and the desired filter,
//   skip and limit must be passed to limit the number of results
func (db *Mongo) GetManyRBAC(session *pbsessions.Session, isDomain bool, perm backd.Permission, database, collection string, query map[string]interface{}, sort []string, this interface{}, skip, limit int) error {

	var (
		all          bool
		accesibleIDs []string
		err          error
	)

	if query == nil {
		query = make(map[string]interface{})
	}

	all, accesibleIDs, err = db.VisibleID(session, isDomain, database, collection, perm)
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

	if len(sort) > 0 {
		return db.session.DB(database).C(collection).Find(query).Sort(sort...).Skip(skip).Limit(limit).All(this)
	}
	return db.session.DB(database).C(collection).Find(query).Skip(skip).Limit(limit).All(this)
}

// InsertRBAC a new entry on the collection, adding metadata, returns errors if any
func (db *Mongo) InsertRBAC(session *pbsessions.Session, isDomain bool, database, collection string, this map[string]interface{}) (map[string]interface{}, error) {

	if db.Can(session, isDomain, database, collection, "*", backd.PermissionCreate) {
		// set metadata
		var metadata structs.Metadata
		metadata.SetCreate(session.GetDomainId(), session.GetUserId())
		this["_id"] = NewXID().String()
		this["_meta"] = metadata
		return this, db.Insert(database, collection, this)
	}

	return map[string]interface{}{}, rest.ErrUnauthorized
}

// InsertRBACInterface a new entry on the collection, metadata and ID must be written in advance
func (db *Mongo) InsertRBACInterface(session *pbsessions.Session, isDomain bool, database, collection string, this interface{}) error {

	if db.Can(session, isDomain, database, collection, "*", backd.PermissionCreate) {
		return db.Insert(database, collection, this)
	}

	return rest.ErrUnauthorized

}

// UpdateByIDRBAC updates the data and metadata on the collections, returning errors if any
func (db *Mongo) UpdateByIDRBAC(session *pbsessions.Session, isDomain bool, database, collection, id string, this map[string]interface{}) (map[string]interface{}, error) {

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

	if db.Can(session, isDomain, database, collection, id, backd.PermissionUpdate) {

		// first get the old entry
		err = db.GetOneByID(database, collection, id, &oldData)
		if err != nil {
			return nil, err
		}

		// updated metadata
		var metadata structs.Metadata
		err = metadata.FromInterface(oldData["_meta"].(map[string]interface{}))
		metadata.SetUpdate(session.GetDomainId(), session.GetUserId())
		this["_meta"] = metadata
		return this, db.UpdateByID(database, collection, id, this)
	}

	return nil, rest.ErrUnauthorized

}

// UpdateByIDRBACInterface updates the data passed as interface and updates the database if ok
func (db *Mongo) UpdateByIDRBACInterface(session *pbsessions.Session, isDomain bool, database, collection, id string, this interface{}) error {

	if db.Can(session, isDomain, database, collection, id, backd.PermissionUpdate) {
		return db.UpdateByID(database, collection, id, this)
	}

	return rest.ErrUnauthorized

}

// DeleteByIDRBAC deletes from the item from the collection if user has permission for it
func (db *Mongo) DeleteByIDRBAC(session *pbsessions.Session, isDomain bool, database, collection, id string) error {

	if db.Can(session, isDomain, database, collection, id, backd.PermissionDelete) {
		return db.session.DB(database).C(collection).RemoveId(id)
	}

	return rest.ErrUnauthorized
}

// DeleteByQueryRBAC deletes from the item from the collection if user has permission for it
func (db *Mongo) DeleteByQueryRBAC(session *pbsessions.Session, isDomain bool, database, collection string, query map[string]interface{}) error {

	if db.Can(session, isDomain, database, collection, "*", backd.PermissionDelete) {
		return db.session.DB(database).C(collection).Remove(query)
	}

	return rest.ErrUnauthorized
}
