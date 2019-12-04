package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Mongo is the struct used to interact with a MongoDB database from backd apis
type Mongo struct {
	client *mongo.Client
}

// NewMongo returns a DB struct to interact with MongoDB
func NewMongo(ctx context.Context, mongoURL string) (*Mongo, error) {
	// DB connection
	sess, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURL))
	if err != nil {
		return nil, err
	}

	return &Mongo{
		client: sess,
	}, nil
}

// Client is a direct access to the Client struct
func (db *Mongo) Client() *mongo.Client {
	return db.client
}

func isID(keys map[string]interface{}) bool {

	if len(keys) == 1 {
		_, ok := keys["_id"]
		if ok {
			return true
		}
	}
	return false

}

// CreateIndex creates required indexes with some default settings that seems to
// be enough for our needs
func (db *Mongo) CreateIndex(ctx context.Context, database, collection string, keys map[string]interface{}, unique bool) (err error) {

	var index mongo.IndexModel

	index.Keys = bson.M(keys)

	if unique && !isID(keys) { // mongo does not allow to set unique to _id
		index.Options = options.Index().SetUnique(unique)
	}

	_, err = db.client.Database(database).Collection(collection).Indexes().CreateOne(ctx, index)

	return

}

// IsInitialized returns an error if there is no connection with the DB
func (db *Mongo) IsInitialized(ctx context.Context, database string) error {

	var (
		collections []string
		err         error
	)

	collections, err = db.client.Database(database).ListCollectionNames(ctx, bson.D{})
	if err != nil {
		return err
	}

	if len(collections) == 0 {
		return ErrDatabaseNotInitialized
	}

	return nil

}

// Insert a new entry on the collection, returns errors if any
func (db *Mongo) Insert(ctx context.Context, database, collection string, this interface{}) (*mongo.InsertOneResult, error) {
	return db.client.Database(database).Collection(collection).InsertOne(ctx, this)
}

// Count returns the number of ocurrencies returnd from the database using that query
func (db *Mongo) Count(ctx context.Context, database, collection string, query map[string]interface{}) (int64, error) {
	return db.client.Database(database).Collection(collection).CountDocuments(ctx, query)
}

// GetMany returns all records that meets the desired filter,
//   skip and limit must be passed to limit the number of results
func (db *Mongo) GetMany(ctx context.Context, database, collection string, query interface{}, sort map[string]interface{}, this interface{}, skip, limit int64) (err error) {

	var (
		cursor *mongo.Cursor
		opts   []*options.FindOptions
	)

	opts = append(opts, options.Find().SetSkip(skip))
	opts = append(opts, options.Find().SetLimit(limit))

	if len(sort) > 0 {
		opts = append(opts, options.Find().SetSort(bson.M(sort)))
	}

	cursor, err = db.client.Database(database).Collection(collection).Find(ctx, query, opts...)
	if err != nil {
		return
	}

	err = cursor.All(ctx, this)
	return

}

// // GetAll returns all records that meets the desired filter
// func (db *Mongo) GetAll(database, collection string, query interface{}, sort []string, this interface{}) error {

// 	var err error

// 	if len(sort) > 0 {
// 		err = db.client.DB(database).C(collection).Find(query).Sort(sort...).All(this)
// 		if err == mgo.ErrNotFound {
// 			return nil // do no return error, return an empty array
// 		}
// 		return err
// 	}
// 	err = db.client.DB(database).C(collection).Find(query).All(this)
// 	if err == mgo.ErrNotFound {
// 		return nil // do no return error, return an empty array
// 	}
// 	return err

// }

// GetOne returns one object by query
func (db *Mongo) GetOne(ctx context.Context, database, collection string, query, this interface{}) error {
	return db.client.Database(database).Collection(collection).FindOne(ctx, query).Decode(this)
}

// GetOneByID returns one object by ID
func (db *Mongo) GetOneByID(ctx context.Context, database, collection, id string, this interface{}) error {
	return db.client.Database(database).Collection(collection).FindOne(ctx, bson.M{"_id": id}).Decode(this)
}

// Update updates the database by using a selector and an object
func (db *Mongo) Update(ctx context.Context, database, collection string, selector, to interface{}) error {
	return db.client.Database(database).Collection(collection).FindOneAndUpdate(ctx, selector, to).Decode(to)
}

// UpdateByID updates the database when object used ObjectID as unique ID
func (db *Mongo) UpdateByID(ctx context.Context, database, collection, id string, to interface{}) error {
	return db.client.Database(database).Collection(collection).FindOneAndUpdate(ctx, bson.M{"_id": id}, to).Decode(to)
}

// Delete deletes from the collection the referenced object
func (db *Mongo) Delete(ctx context.Context, database, collection string, selector map[string]interface{}) (count int64, err error) {
	var res *mongo.DeleteResult

	res, err = db.client.Database(database).Collection(collection).DeleteMany(ctx, selector)
	if err != nil {
		return
	}

	count = res.DeletedCount
	return
}

// DeleteByID deletes from the collection the referenced object using an ObjectID passed as string
func (db *Mongo) DeleteByID(ctx context.Context, database, collection, id string) (count int64, err error) {
	var res *mongo.DeleteResult

	res, err = db.client.Database(database).Collection(collection).DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return
	}

	count = res.DeletedCount
	return

}
