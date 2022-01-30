package database

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// insertOne is a user defined method, used to insert
// documents into collection returns result of InsertOne
// and error if any.
func (db *DB) InsertOne(dataBase, col string, doc interface{}) (*mongo.InsertOneResult, error) {

	// select database and collection ith Client.Database method
	// and Database.Collection method
	collection := db.client.Database(dataBase).Collection(col)

	// InsertOne accept two argument of type Context
	// and of empty interface
	result, err := collection.InsertOne(context.TODO(), doc)
	return result, err
}

// query is user defined method used to query MongoDB,
// that accepts mongo.client,context, database name,
// collection name, a query and field.

//  database name and collection name is of type
// string. query is of type interface.
// field is of type interface, which limts
// the field being returned.

// query method returns a cursor and error.
func (db *DB) Query(dataBase, col string, query interface{}, opts *options.FindOptions) (result *mongo.Cursor, err error) {

	// select database and collection.
	collection := db.client.Database(dataBase).Collection(col)

	// collection has an method Find,
	// that returns a mongo.cursor
	// based on query and field.
	result, err = collection.Find(context.TODO(), query, opts)
	return
}

// UpdateOne is a user defined method, that update
// a single document matching the filter.
// This methods accepts client, context, database,
// collection, filter and update filter and update
// is of type interface this method returns
// UpdateResult and an error if any.
func (db *DB) UpdateOne(dataBase, col string, filter, update interface{}) (result *mongo.UpdateResult, err error) {

	// select the database and the collection
	collection := db.client.Database(dataBase).Collection(col)

	// A single document that match with the
	// filter will get updated.
	// update contains the filed which should get updated.
	result, err = collection.UpdateOne(context.TODO(), filter, update)
	return
}
