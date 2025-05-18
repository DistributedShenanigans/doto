package tasks

import (
	"context"
	"errors"
	"fmt"

	dotoapi "github.com/DistributedShenanigans/doto/api"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository struct {
	collection *mongo.Collection
}

type mongoTask struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	ChatID      int64              `bson:"chatID"`
	Description string             `bson:"description"`
	Status      string             `bson:"status"`
}

// WARN: Порнография
var (
	ErrNotFound = errors.New("task not found")
)

func New(db *mongo.Database, collectionName string) *Repository {
	return &Repository{
		collection: db.Collection(collectionName),
	}
}

func (r *Repository) Get(ctx context.Context, chatID int64) ([]dotoapi.Task, error) {
	filter := bson.M{"chatID": chatID}
	opts := options.Find().SetSort(bson.D{{Key: "createdAt", Value: 1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tasks []dotoapi.Task
	for cursor.Next(ctx) {
		var mt mongoTask
		if err := cursor.Decode(&mt); err != nil {
			return nil, err
		}

		tasks = append(tasks, dotoapi.Task{
			Id:          mt.ID.Hex(),
			Status:      mt.Status,
			Description: mt.Description,
		})
	}

	return tasks, nil
}

func (r *Repository) Add(ctx context.Context, chatID int64, task dotoapi.TaskCreation) error {
	mt := mongoTask{
		ChatID:      chatID,
		Status:      task.Status,
		Description: task.Description,
	}

	if _, err := r.collection.InsertOne(ctx, mt); err != nil {
		return err
	}

	return nil
}

func (r *Repository) UpdateStatus(
	ctx context.Context,
	chatID int64,
	taskID string,
	update dotoapi.TaskStatusUpdate,
) (dotoapi.Task, error) {
	objID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		return dotoapi.Task{}, fmt.Errorf("failed to update status: %w", ErrNotFound)
	}

	filter := bson.M{
		"_id":    objID,
		"chatID": chatID,
	}

	updateDoc := bson.M{
		"$set": bson.M{
			"status": update.Status,
		},
	}

	opts := options.FindOneAndUpdate().
		SetReturnDocument(options.After)

	var updated mongoTask
	if err := r.collection.FindOneAndUpdate(
		ctx,
		filter,
		updateDoc,
		opts,
	).Decode(&updated); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return dotoapi.Task{}, fmt.Errorf("failed to update status: %w", ErrNotFound)
		}

		return dotoapi.Task{}, fmt.Errorf("failed to update status: %w", err)
	}

	return dotoapi.Task{
		Id:          updated.ID.Hex(),
		Status:      updated.Status,
		Description: updated.Description,
	}, nil
}

func (r *Repository) Delete(ctx context.Context, chatID int64, taskID string) error {
	objID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		return fmt.Errorf("failed to update status: %w", ErrNotFound)
	}

	filter := bson.M{
		"_id":    objID,
		"chatID": chatID,
	}

	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return ErrNotFound
	}

	return nil
}
