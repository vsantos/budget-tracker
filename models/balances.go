package models

import (
	"budget-tracker-api/observability"
	"budget-tracker-api/services"
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"go.opentelemetry.io/otel/attribute"
)

// CreateBalance creates a balance for a given owner_id
func CreateBalance(parentCtx context.Context, b Balance) (id string, err error) {
	spanTags := []attribute.KeyValue{
		attribute.Key("balance.owner.id").String(b.OwnerID.String()),
	}

	ctx, span := observability.Span(parentCtx, "mongodb", "models.CreateBalance", spanTags)
	defer span.End()

	dbClient, err := services.InitMongoDB()
	if err != nil {
		return "", err
	}

	col := dbClient.Database(mongodbDatabase).Collection(mongodbBalanceCollection)
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	_, err = col.Indexes().CreateOne(
		ctx,
		mongo.IndexModel{
			Keys: bsonx.Doc{
				{Key: "owner_id", Value: bsonx.Int32(1)},
				{Key: "month", Value: bsonx.Int32(1)},
				{Key: "year", Value: bsonx.Int32(1)},
			},
			Options: options.Index().SetUnique(true),
		},
	)

	// adding timestamp to creationDate
	t := time.Now()
	b.CreatedAt = primitive.NewDateTimeFromTime(t)
	b.UpdatedAt = primitive.NewDateTimeFromTime(t)
	b.SpendableAmount = b.Income.NetIncome
	b.Historic = []Spend{}

	r, err := col.InsertOne(ctx, b)
	if err != nil {
		cancel()
		return "", err
	}

	span.SetAttributes(attribute.Key("balance.id").String(r.InsertedID.(primitive.ObjectID).Hex()))
	defer cancel()

	observability.Metrics.Balances.BalancesCreated.Inc()
	log.Infoln("created balance", r.InsertedID.(primitive.ObjectID).Hex())
	return r.InsertedID.(primitive.ObjectID).Hex(), nil
}

// GetBalance will return a balance from an owner_id based on month and year
func GetBalance(parentCtx context.Context, ownerID string, month int64, year int64) (balance *Balance, err error) {
	spanTags := []attribute.KeyValue{
		attribute.Key("balance.owner.id").String(ownerID),
	}

	ctx, span := observability.Span(parentCtx, "mongodb", "models.GetBalance", spanTags)
	defer span.End()

	oid, err := primitive.ObjectIDFromHex(ownerID)
	if err != nil {
		return &Balance{}, err
	}

	dbClient, err := services.InitMongoDB()
	if err != nil {
		return &Balance{}, err
	}

	col := dbClient.Database(mongodbDatabase).Collection(mongodbBalanceCollection)
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	err = col.FindOne(ctx, bson.M{
		"owner_id": oid,
		"month":    month,
		"year":     year,
	}).Decode(&balance)
	if err != nil {
		cancel()
		return &Balance{}, err
	}

	defer cancel()
	return balance, nil
}

// GetAllBalances will return all balances from an owner_id
func GetAllBalances(parentCtx context.Context, ownerID string) (balances []Balance, err error) {
	spanTags := []attribute.KeyValue{
		attribute.Key("balance.owner.id").String(ownerID),
	}

	ctx, span := observability.Span(parentCtx, "mongodb", "models.GetAllBalances", spanTags)
	defer span.End()

	dbClient, err := services.InitMongoDB()
	if err != nil {
		return []Balance{}, err
	}

	oid, err := primitive.ObjectIDFromHex(ownerID)
	if err != nil {
		return []Balance{}, err
	}

	col := dbClient.Database(mongodbDatabase).Collection(mongodbBalanceCollection)
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	cursor, err := col.Find(ctx, bson.M{"owner_id": oid})
	if err != nil {
		cancel()
		return []Balance{}, err
	}

	defer cursor.Close(ctx)
	defer cancel()

	for cursor.Next(ctx) {
		var balance Balance
		cursor.Decode(&balance)
		balances = append(balances, balance)
		defer cancel()
	}

	if err := cursor.Err(); err != nil {
		cancel()
		return []Balance{}, err
	}

	return balances, nil
}
