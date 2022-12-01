package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/commedesvlados/go-app/internal/user"
	"github.com/commedesvlados/go-app/pkg/logging"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type db struct {
	collection *mongo.Collection
	logger     *logging.Logger
}

func (d *db) Create(ctx context.Context, user user.User) (string, error) {
	d.logger.Debug("create user")
	result, err := d.collection.InsertOne(ctx, user)
	if err != nil {
		return "", fmt.Errorf("failed to create user due to error: %v", err)
	}

	d.logger.Debug("convert InsertId to ObjectId")
	oid, ok := result.InsertedID.(primitive.ObjectID)
	if ok {
		return oid.Hex(), nil
	}
	d.logger.Trace(user)
	return "", fmt.Errorf("failed to convert objectid to hex")
}

func (d *db) FindOne(ctx context.Context, id string) (u user.User, err error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return u, fmt.Errorf("failed to convert hex: %s to objectid due to error: %v", id, err)
	}
	//mongo.getDatabsase("test").getCollection("docs").find({})
	filter := bson.M{"_id": oid}

	result := d.collection.FindOne(ctx, filter)
	if err = result.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			// TODO ErrEntityNotfound
			return u, fmt.Errorf("not found")
		}
		return u, fmt.Errorf("failed to find one user by id: %s due to error: %v", id, err)
	}

	if err = result.Decode(&u); err != nil {
		return u, fmt.Errorf("failed to decode user by id: %s due to error: %v", id, err)
	}

	return u, nil
}

func (d *db) Update(ctx context.Context, user user.User) error {
	objectId, err := primitive.ObjectIDFromHex(user.Id)
	if err != nil {
		return fmt.Errorf("failed to convert userId: %s to objectId due to error: %v", user.Id, err)
	}

	filter := bson.M{"_id": objectId}

	userBytes, err := bson.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshal user due to error: %v", err)
	}

	var updateUserObj bson.M
	if err = bson.Unmarshal(userBytes, &updateUserObj); err != nil {
		return fmt.Errorf("failed to unmarshal user due to error: %v", err)
	}

	delete(updateUserObj, "_id")

	update := bson.M{
		"$set": updateUserObj,
	}

	result, err := d.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to execute update user query due to error: %v", err)
	}

	if result.MatchedCount == 0 {
		// TODO ErrEntituNotFound
		return fmt.Errorf("not found ")
	}

	d.logger.Tracef("Matched %d documents and Modified %d documents", result.MatchedCount, result.ModifiedCount)

	return nil
}

func (d *db) Delete(ctx context.Context, id string) error {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("failed to convert userId: %s to objectId due to error: %v", id, err)
	}

	filter := bson.M{"_id": objectId}
	result, err := d.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to execute filter due to error: %v", err)
	}
	if result.DeletedCount == 0 {
		// TODO ErrEntituNotFound
		return fmt.Errorf("not found")
	}
	d.logger.Tracef("Deleted %d documents", result.DeletedCount)

	return nil
}

func NewStorage(database *mongo.Database, collection string, logger *logging.Logger) user.Storage {
	return &db{
		collection: database.Collection(collection),
		logger:     logger,
	}
}
