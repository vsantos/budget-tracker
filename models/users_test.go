package models

import (
	"context"
	"testing"
	"time"

	"github.com/mongo-go/testdb"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
)

var testDb *testdb.TestDB

func setup(t *testing.T) *mongo.Collection {
	if testDb == nil {
		testDb = testdb.NewTestDB("mongodb://localhost", "your_db", time.Duration(2)*time.Second)

		err := testDb.Connect()
		if err != nil {
			t.Fatal(err)
		}
	}

	coll, err := testDb.CreateRandomCollection(testdb.NoIndexes)
	if err != nil {
		t.Fatal(err)
	}

	return coll // random *mongo.Collection in "your_db"
}

func GetUsersTest(t *testing.T) {
	setup(t)
	coll := setup(t)
	defer coll.Drop(context.Background())
	assert.Equal(t, 123, 123, "they should be equal")
}
