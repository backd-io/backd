package db

import (
	"github.com/backd-io/backd/backd"
	"github.com/backd-io/backd/internal/constants"
	"github.com/backd-io/backd/internal/pbsessions"
	"github.com/backd-io/backd/internal/structs"
	mgo "github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// Mongo is the struct used to interact with a MongoDB database from backd apis
type Mongo struct {
	session *mgo.Session
}

// NewMongo returns a DB struct to interact with MongoDB
func NewMongo(mongoURL string) (*Mongo, error) {
	// DB connection
	sess, err := mgo.Dial(mongoURL)
	if err != nil {
		return nil, err
	}

	sess.SetMode(mgo.Monotonic, true)

	return &Mongo{
		session: sess,
	}, nil
}

// Session is a direct access to the Session struct
func (db *Mongo) Session() *mgo.Session {
	return db.session
}

// CreateDefaultDomainIndexes cretes the required indexes for a domain
//  - entities
//  - memberships
//  - sessions
func (db *Mongo) CreateDefaultDomainIndexes(database string) error {

	var (
		err error
	)

	// // index to help find user / group
	// if err = db.CreateIndex(database, constants.ColEntities, []string{"_type"}, false); err != nil {
	// 	return err
	// }

	// // index to help find user (_type=u) / group (_type=g)
	// if err = db.CreateIndex(database, constants.ColEntities, []string{"_type", "name"}, false); err != nil {
	// 	return err
	// }

	// index to help find relations of user (_type=u) & group (_type=g)
	if err = db.CreateIndex(database, constants.ColMembership, []string{"u", "g"}, true); err != nil {
		return err
	}

	// index to help find relations of user (_type=u)
	if err = db.CreateIndex(database, constants.ColMembership, []string{"u"}, false); err != nil {
		return err
	}

	// index to help find relations of  group (_type=g)
	if err = db.CreateIndex(database, constants.ColMembership, []string{"g"}, false); err != nil {
		return err
	}

	return nil
}

// CreateDefaultAppIndexes creates the required indexes for the basic API services:
//  - data relationship
func (db *Mongo) CreateDefaultAppIndexes(database string) error {

	var (
		err error
	)

	// relation -> src, sid (source)
	if err = db.CreateIndex(database, constants.ColRelation, []string{"src", "sid"}, false); err != nil {
		return err
	}

	// relation -> dst, did (destination)
	if err = db.CreateIndex(database, constants.ColRelation, []string{"dst", "did"}, false); err != nil {
		return err
	}

	// relation -> src, sid, rel (source + relation)
	if err = db.CreateIndex(database, constants.ColRelation, []string{"src", "sid", "rel"}, false); err != nil {
		return err
	}

	// relation -> dst, did, rel (destination + relation)
	if err = db.CreateIndex(database, constants.ColRelation, []string{"dst", "did", "rel"}, false); err != nil {
		return err
	}

	// relation -> src, sid, rel, dst (source + relation + destinationType)
	if err = db.CreateIndex(database, constants.ColRelation, []string{"src", "sid", "rel", "dst"}, false); err != nil {
		return err
	}

	// relation -> dst, did, rel (destination + relation + sourceType)
	if err = db.CreateIndex(database, constants.ColRelation, []string{"dst", "did", "rel", "src"}, false); err != nil {
		return err
	}

	// relation -> src, sid, rel, dst (source + relation + destination) - must be unique
	if err = db.CreateIndex(database, constants.ColRelation, []string{"src", "sid", "rel", "dst", "did"}, true); err != nil {
		return err
	}

	return nil

}

// CreateIndex creates required indexes with some default settings that seems to
// be enough for our needs
func (db *Mongo) CreateIndex(database, collection string, fields []string, unique bool) error {

	index := mgo.Index{
		Key:        fields,
		Unique:     unique,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}

	return db.session.DB(database).C(collection).EnsureIndex(index)

}

// IsInitialized returns an error if there is no connection with the DB
func (db *Mongo) IsInitialized(database string) error {

	var (
		collections []string
		err         error
	)

	collections, err = db.session.DB(database).CollectionNames()
	if err != nil {
		return err
	}

	if len(collections) == 0 {
		return ErrDatabaseNotInitialized
	}

	return nil

}

// Insert a new entry on the collection, returns errors if any
func (db *Mongo) Insert(database, collection string, this interface{}) error {
	return db.session.DB(database).C(collection).Insert(this)
}

// Count returns the number of ocurrencies returnd from the database using that query
func (db *Mongo) Count(database, collection string, query interface{}) (int, error) {
	return db.session.DB(database).C(collection).Find(query).Count()
}

// GetManyRBAC returns all records that meets RBAC and the desired filter,
//   skip and limit must be passed to limit the number of results
func (db *Mongo) GetManyRBAC(session *pbsessions.Session, perm backd.Permission, database, collection string, query map[string]interface{}, sort []string, this interface{}, skip, limit int) error {

	var (
		all          bool
		accesibleIDs []string
		err          error
	)

	if query == nil {
		query = make(map[string]interface{})
	}

	all, accesibleIDs, err = db.VisibleID(session, database, collection, perm)
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

// GetMany returns all records that meets the desired filter,
//   skip and limit must be passed to limit the number of results
func (db *Mongo) GetMany(database, collection string, query interface{}, sort []string, this interface{}, skip, limit int) error {
	if len(sort) > 0 {
		return db.session.DB(database).C(collection).Find(query).Sort(sort...).Skip(skip).Limit(limit).All(this)
	}
	return db.session.DB(database).C(collection).Find(query).Skip(skip).Limit(limit).All(this)
}

// GetOne returns one object by ID
func (db *Mongo) GetOne(database, collection string, query, this interface{}) error {
	return db.session.DB(database).C(collection).Find(query).One(this)
}

// GetOneByID returns one object by ID
func (db *Mongo) GetOneByID(database, collection, id string, this interface{}) error {
	return db.session.DB(database).C(collection).FindId(id).One(this)
}

// Update updates the database by using a selector and an object
func (db *Mongo) Update(database, collection string, selector, to interface{}) error {
	return db.session.DB(database).C(collection).Update(selector, to)
}

// UpdateByID updates the database when object used ObjectID as unique ID
func (db *Mongo) UpdateByID(database, collection, id string, to interface{}) error {
	return db.session.DB(database).C(collection).UpdateId(id, to)
}

// Delete deletes from the collection the referenced object
func (db *Mongo) Delete(database, collection string, selector interface{}) error {
	return db.session.DB(database).C(collection).Remove(selector)
}

// DeleteByID deletes from the collection the referenced object using an ObjectID passed as string
func (db *Mongo) DeleteByID(database, collection, id string) error {
	return db.session.DB(database).C(collection).RemoveId(id)
}

// Can validates the ability to make something by an user
func (db *Mongo) Can(session *pbsessions.Session, database, collection, id string, perm backd.Permission) bool {

	var (
		rbac  structs.RBAC
		query map[string]interface{}
		err   error
	)

	query = map[string]interface{}{
		"did": session.GetDomainId(),
		"uid": db.getIdentities(session),
		"c":   collection,
		"cid": id,
		"p":   perm,
	}

	err = db.GetOne(database, constants.ColRBAC, query, &rbac)
	if err != nil {
		return false
	}

	return true
}

// VisibleID returns only those IDs that the user is able to see from a collection
func (db *Mongo) VisibleID(session *pbsessions.Session, database, collection string, perm backd.Permission) (all bool, ids []string, err error) {

	var (
		query map[string]interface{}
	)

	query = map[string]interface{}{
		"did": session.GetDomainId(),
		"uid": db.getIdentities(session),
		"c":   collection,
		"p":   perm,
	}

	err = db.session.DB(database).C(constants.ColRBAC).Find(query).Distinct("cid", ids)

	// see if can see all the items to simplify the query
	for _, id := range ids {
		if id == "*" {
			all = true
			break
		}
	}
	return

}

// getIdentities returns the
func (db *Mongo) getIdentities(session *pbsessions.Session) (identities []string) {

	identities = append(identities, session.GetUserId())
	for _, identity := range session.GetGroups() {
		identities = append(identities, identity)
	}
	return

}
