package database

import (
	"context"
	"fmt"
	"time"

	// "movieapi/app/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"movieapi/app/config"
)

type MovieDB interface {
	Open() error
	Close() error
	InsertOne(dataBase, col string, doc interface{}) (*mongo.InsertOneResult, error)
	UpdateOne(dataBase, col string, filter, update interface{}) (result *mongo.UpdateResult, err error)
	Query(dataBase, col string, query interface{}, opts *options.FindOptions) (result *mongo.Cursor, err error)
}

// mongo.Client will be used for further database operation.
// context.Context will be used set deadlines for process.
// context.CancelFunc will be used to cancel context and
// resource associated with it.
type DB struct {
	client *mongo.Client
	ctx    context.Context
	cancel context.CancelFunc
}

// Method that set mongo.Client,
// context.Context, context.CancelFunc and error.
func (db *DB) Open() error {
	// ctx will be used to set deadline for process, here
	// deadline will of 30 seconds.
	ctx, cancel := context.WithTimeout(context.Background(),
		30*time.Second)

	// mongo.Connect return mongo.Client method
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MgConnectionString))
	if err != nil {
		return err
	}
	db.client = client
	db.ctx = ctx
	db.cancel = cancel

	err = db.ping(client, ctx)
	if err != nil {
		return err
	}
	return nil
}

// This is a user defined method to close resources.
// This method closes mongoDB connection and cancel context.
func (db *DB) Close() error {
	client := db.client
	ctx := db.ctx
	cancel := db.cancel
	// CancelFunc to cancel to context
	defer cancel()

	// client provides a method to close
	// a mongoDB connection.
	defer func() {
		// client.Disconnect method also has deadline.
		// returns error if any,
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	return nil
}

// This is a user defined method that accepts
// mongo.Client and context.Context
// This method used to ping the mongoDB, return error if any.
func (db *DB) ping(client *mongo.Client, ctx context.Context) error {

	// mongo.Client has Ping to ping mongoDB, deadline of
	// the Ping method will be determined by cxt
	// Ping method return error if any occurred, then
	// the error can be handled.
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}
	fmt.Println("connected successfully")
	return nil
}
