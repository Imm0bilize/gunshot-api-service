package repository

import (
	"context"
	"github.com/Imm0bilize/gunshot-api-service/internal/entities"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type ClientRepo struct {
	collection *mongo.Collection
	tracer     trace.Tracer
}

// Create save a new user with uuid and information about him
func (c ClientRepo) Create(ctx context.Context, client *entities.Client) (string, error) {
	ctx, span := c.tracer.Start(ctx, "ClientRepo.Create")
	defer span.End()

	client.ID = primitive.NewObjectID()

	_, err := c.collection.InsertOne(ctx, client)
	if err != nil {
		return "", errors.Wrap(err, "error during create client")
	}

	return client.ID.Hex(), nil
}

func (c ClientRepo) Get(ctx context.Context, id string) (entities.Client, error) {
	ctx, span := c.tracer.Start(ctx, "ClientRepo.Get")
	defer span.End()

	castedID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return entities.Client{}, errors.Wrap(err, "invalid client id")
	}

	filter := bson.M{
		"_id": castedID,
	}

	var client entities.Client
	if err := c.collection.FindOne(ctx, filter).Decode(&client); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entities.Client{}, ErrClientNotFound
		}

		span.RecordError(err)
		return entities.Client{}, errors.Wrap(err, "error during get client from db")
	}

	return client, nil
}

func (c ClientRepo) Update(ctx context.Context, id string, client *entities.Client) error {
	ctx, span := c.tracer.Start(ctx, "ClientRepo.Update")
	defer span.End()

	castedID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.Wrap(err, "invalid client id")
	}

	filter := bson.M{
		"_id": castedID,
	}

	update := bson.M{
		"$set": bson.M{
			"locationName":       client.LocationName,
			"fullName":           client.FullName,
			"latitude":           client.Latitude,
			"longitude":          client.Longitude,
			"notificationMethod": client.NotificationMethods,
		},
	}

	res := c.collection.FindOneAndUpdate(ctx, filter, update)

	if res.Err() != nil {
		if errors.Is(res.Err(), mongo.ErrNoDocuments) {
			return ErrClientNotFound
		}

		span.RecordError(err)
		return errors.Wrap(err, "error during update information")
	}

	return nil
}

// Delete remove the user from database
func (c ClientRepo) Delete(ctx context.Context, id string) error {
	ctx, span := c.tracer.Start(ctx, "ClientRepo.Delete")
	defer span.End()

	castedID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.Wrap(err, "invalid client id")
	}

	filter := bson.M{
		"_id": castedID,
	}

	res := c.collection.FindOneAndDelete(ctx, filter)

	if res.Err() != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrClientNotFound
		}

		span.RecordError(err)
		return errors.Wrap(err, "error during delete client")
	}

	return nil
}

func NewClientRepo(database *mongo.Database) *ClientRepo {
	tracer := otel.Tracer("ClientRepo")

	return &ClientRepo{
		collection: database.Collection(_clientsCollection),
		tracer:     tracer,
	}
}
