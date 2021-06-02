package models

import (
	"budget-tracker-api/observability"
	"budget-tracker-api/services"
	"context"
	"errors"
	"time"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"go.opentelemetry.io/otel/attribute"
)

// CreateCard creates a card for a given owner_id
func CreateCard(parentCtx context.Context, c CreditCard) (id string, err error) {
	spanTags := []attribute.KeyValue{
		attribute.Key("card.owner.id").String(c.OwnerID.String()),
	}

	ctx, span := observability.Span(parentCtx, "mongodb", "CreateCard", spanTags)
	defer span.End()

	dbClient, err := services.InitMongoDB()
	if err != nil {
		return "", err
	}

	col := dbClient.Database(mongodbDatabase).Collection(mongodbCardsCollection)
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	_, err = col.Indexes().CreateOne(
		ctx,
		mongo.IndexModel{
			Keys:    bsonx.Doc{{Key: "owner_id", Value: bsonx.Int32(1)}, {Key: "last_digits", Value: bsonx.Int32(1)}},
			Options: options.Index().SetUnique(true),
		},
	)

	// adding timestamp to creationDate
	t := time.Now()
	c.CreatedAt = primitive.NewDateTimeFromTime(t)

	r, err := col.InsertOne(ctx, c)
	if err != nil {
		cancel()
		return "", err
	}

	span.SetAttributes(attribute.Key("card.id").String(r.InsertedID.(primitive.ObjectID).Hex()))
	defer cancel()

	observability.Metrics.Cards.CardsCreated.Inc()
	log.Infoln("created card", c.Alias)
	return r.InsertedID.(primitive.ObjectID).Hex(), nil
}

// GetAllCards will return a list of all cards from the database
func GetAllCards(parentCtx context.Context) (cards []CreditCard, err error) {
	ctx, span := observability.Span(parentCtx, "mongodb", "GetAllCards", []attribute.KeyValue{})
	defer span.End()

	dbClient, err := services.InitMongoDB()
	if err != nil {
		return []CreditCard{}, err
	}

	col := dbClient.Database(mongodbDatabase).Collection(mongodbCardsCollection)
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	cursor, err := col.Find(ctx, bson.M{})
	if err != nil {
		cancel()
		return []CreditCard{}, err
	}

	defer cursor.Close(ctx)
	defer cancel()

	for cursor.Next(ctx) {
		var card CreditCard
		cursor.Decode(&card)
		cards = append(cards, card)
		defer cancel()
	}

	if err := cursor.Err(); err != nil {
		cancel()
		return []CreditCard{}, err
	}

	return cards, nil
}

// GetCards will return a list of cards from a owner_id
func GetCards(parentCtx context.Context, ownerID string) (cards []CreditCard, err error) {
	spanTags := []attribute.KeyValue{
		attribute.Key("card.owner.id").String(ownerID),
	}

	ctx, span := observability.Span(parentCtx, "mongodb", "GetUserCards", spanTags)
	defer span.End()

	dbClient, err := services.InitMongoDB()
	if err != nil {
		return []CreditCard{}, err
	}

	pid, err := primitive.ObjectIDFromHex(ownerID)
	if err != nil {
		return []CreditCard{}, err
	}

	col := dbClient.Database(mongodbDatabase).Collection(mongodbCardsCollection)
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	cursor, err := col.Find(ctx, bson.M{"owner_id": pid})
	if err != nil {
		cancel()
		return []CreditCard{}, err
	}

	defer cursor.Close(ctx)
	defer cancel()

	for cursor.Next(ctx) {
		var card CreditCard
		cursor.Decode(&card)
		cards = append(cards, card)
		defer cancel()
	}

	if err := cursor.Err(); err != nil {
		cancel()
		return []CreditCard{}, err
	}

	return cards, nil
}

// DeleteCard creates an user based on request body payload
func DeleteCard(parentCtx context.Context, id string) (err error) {
	spanTags := []attribute.KeyValue{
		attribute.Key("card.id").String(id),
	}

	ctx, span := observability.Span(parentCtx, "mongodb", "DeleteUserCard", spanTags)
	defer span.End()

	pid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	dbClient, err := services.InitMongoDB()
	if err != nil {
		return err
	}

	col := dbClient.Database(mongodbDatabase).Collection(mongodbCardsCollection)
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	log.Infoln("deleting card", id)
	result, err := col.DeleteOne(ctx, bson.M{"_id": pid})
	if err != nil {
		cancel()
		return err
	}

	log.Infoln("number of cards deleted:", result.DeletedCount)

	if result.DeletedCount == 0 {
		cancel()
		return errors.New("non existent card")
	}

	defer cancel()

	log.Infoln("deleted card", id)
	return nil
}
