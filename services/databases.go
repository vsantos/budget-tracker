package services

import "go.mongodb.org/mongo-driver/mongo"

var (
	// NoSQLClient will define a global client
	NoSQLClient *mongo.Client
)

// Storage will define
type Storage struct {
	NoSQLClient *mongo.Client
}
