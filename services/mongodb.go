package services

import (
	"context"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// MongoOperationsInterface defines methods to interact with mongoDB
type MongoOperationsInterface interface {
	GetOne(parentCtx context.Context, database string, collection string, filter interface{}) (*mongo.SingleResult, error)
	Get(parentCtx context.Context, database string, collection string, filter interface{}) (*mongo.Cursor, error)
	CreateIndex(parentCtx context.Context, database string, collection string, keys interface{}) (err error)
	CreateOne(parentCtx context.Context, database string, collection string, document interface{}) (*mongo.InsertOneResult, error)
	DeleteOne()
	Ping(parentCtx context.Context) (err error)
}

// InitMongoDB will return a database client for usage
func InitMongoDB() (c *mongo.Client, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	uri := "mongodb://mongodb:27017/"

	clientOptions := options.Client().ApplyURI(uri)
	clientOptions.Monitor = otelmongo.NewMonitor("mongodb")
	dbClient, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		cancel()
		return &mongo.Client{}, err
	}

	defer cancel()

	return dbClient, err
}

// Get will return a cursor to be posterior iterated
func (s Storage) Get(parentCtx context.Context, database string, collection string, filter interface{}) (*mongo.Cursor, error) {
	col := s.NoSQLClient.Database(database).Collection(collection)
	ctx, cancel := context.WithTimeout(parentCtx, 5*time.Second)

	cursor, err := col.Find(ctx, filter)
	if err != nil {
		cancel()
		return nil, err
	}

	defer cancel()
	return cursor, nil

}

// GetOne will return a single mongodb result to be posterior decoded
func (s Storage) GetOne(parentCtx context.Context, database string, collection string, filter interface{}) (*mongo.SingleResult, error) {
	col := s.NoSQLClient.Database(database).Collection(collection)
	ctx, cancel := context.WithTimeout(parentCtx, 5*time.Second)

	r := col.FindOne(ctx, filter)
	if r.Err() != nil {
		cancel()
		return nil, r.Err()
	}

	defer cancel()
	return r, nil

}

// CreateIndex will create a index based on document
func (s Storage) CreateIndex(parentCtx context.Context, database string, collection string, keys interface{}) (err error) {
	col := s.NoSQLClient.Database(database).Collection(collection)

	_, err = col.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    keys,
			Options: options.Index().SetUnique(true),
		},
	)

	if err != nil {
		return err
	}

	return nil
}

// CreateOne will create a single registry from mongodb
func (s Storage) CreateOne(parentCtx context.Context, database string, collection string, document interface{}) (*mongo.InsertOneResult, error) {
	col := s.NoSQLClient.Database(database).Collection(collection)
	ctx, cancel := context.WithTimeout(parentCtx, 5*time.Second)

	r, err := col.InsertOne(ctx, document)
	if err != nil {
		cancel()
		return nil, err
	}

	defer cancel()
	return r, nil
}

// DeleteOne will delete a single registry from mongodb
func (s Storage) DeleteOne() {

}

// Ping will execute a test connection to MongoDB server
func (s Storage) Ping(parentCtx context.Context) (err error) {
	ctx, cancel := context.WithTimeout(parentCtx, 2*time.Second)
	err = s.NoSQLClient.Ping(ctx, &readpref.ReadPref{})
	if err != nil {
		cancel()
		return err
	}

	defer cancel()
	return nil
}
