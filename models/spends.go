package models

import (
	"budget-tracker-api/observability"
	"budget-tracker-api/services"
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.opentelemetry.io/otel/attribute"
)

// CreateSpend creates a card for a given owner_id
func CreateSpend(parentCtx context.Context, s Spend) (id string, err error) {
	spanTags := []attribute.KeyValue{
		attribute.Key("spend.owner.id").String(s.OwnerID.String()),
	}

	ctx, span := observability.Span(parentCtx, "mongodb", "models.CreateSpend", spanTags)
	defer span.End()

	dbClient, err := services.InitMongoDB()
	if err != nil {
		return "", err
	}

	col := dbClient.Database(mongodbDatabase).Collection(mongodbSpendsCollection)
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	// adding timestamp to creationDate
	t := time.Now()
	s.CreatedAt = primitive.NewDateTimeFromTime(t)

	r, err := col.InsertOne(ctx, s)
	if err != nil {
		cancel()
		return "", err
	}

	span.SetAttributes(attribute.Key("spend.id").String(r.InsertedID.(primitive.ObjectID).Hex()))
	defer cancel()

	observability.Metrics.Spends.SpendsCreated.Inc()
	log.Infoln("created spend", r.InsertedID.(primitive.ObjectID).Hex())
	return r.InsertedID.(primitive.ObjectID).Hex(), nil
}

// GetSpends will return all spends from a specific owner_id
func GetSpends(parentCtx context.Context, ownerID string) (spends []Spend, err error) {
	spanTags := []attribute.KeyValue{
		attribute.Key("spend.owner.id").String(ownerID),
	}

	ctx, span := observability.Span(parentCtx, "mongodb", "models.GetSpends", spanTags)
	defer span.End()

	dbClient, err := services.InitMongoDB()
	if err != nil {
		return []Spend{}, err
	}

	pid, err := primitive.ObjectIDFromHex(ownerID)
	if err != nil {
		return []Spend{}, err
	}

	col := dbClient.Database(mongodbDatabase).Collection(mongodbSpendsCollection)
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	cursor, err := col.Find(ctx, bson.M{"owner_id": pid})
	if err != nil {
		cancel()
		return []Spend{}, err
	}

	defer cursor.Close(ctx)
	defer cancel()

	for cursor.Next(ctx) {
		var spend Spend
		cursor.Decode(&spend)
		spends = append(spends, spend)
		defer cancel()
	}

	if err := cursor.Err(); err != nil {
		cancel()
		return []Spend{}, err
	}

	return spends, nil
}
