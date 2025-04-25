package repository

import (
	"bug-tracker/models"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BugRepository struct {
	db *mongo.Database
}

func NewBugRepository(db *mongo.Database) *BugRepository {
	return &BugRepository{db: db}
}

func (r *BugRepository) Create(ctx context.Context, bug *models.Bug) error {
	collection := r.db.Collection("bugs")

	bug.CreatedAt = time.Now()
	bug.UpdatedAt = time.Now()
	bug.Status = "open"

	result, err := collection.InsertOne(ctx, bug)
	if err != nil {
		return err
	}

	bug.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *BugRepository) FindAll(ctx context.Context) ([]*models.Bug, error) {
	collection := r.db.Collection("bugs")

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var bugs []*models.Bug
	if err = cursor.All(ctx, &bugs); err != nil {
		return nil, err
	}

	return bugs, nil
}

func (r *BugRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.Bug, error) {
	collection := r.db.Collection("bugs")

	var bug models.Bug
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&bug)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &bug, nil
}

func (r *BugRepository) FindByAssignee(ctx context.Context, developerID primitive.ObjectID) ([]*models.Bug, error) {
	collection := r.db.Collection("bugs")

	cursor, err := collection.Find(ctx, bson.M{"assigned_to": developerID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var bugs []*models.Bug
	if err = cursor.All(ctx, &bugs); err != nil {
		return nil, err
	}

	return bugs, nil
}

func (r *BugRepository) UpdateStatus(ctx context.Context, id primitive.ObjectID, status string) error {
	collection := r.db.Collection("bugs")

	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": time.Now(),
		},
	}

	_, err := collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

func (r *BugRepository) AssignToDeveloper(ctx context.Context, bugID, developerID primitive.ObjectID) error {
	collection := r.db.Collection("bugs")

	update := bson.M{
		"$set": bson.M{
			"assigned_to": developerID,
			"updated_at":  time.Now(),
		},
	}

	_, err := collection.UpdateOne(ctx, bson.M{"_id": bugID}, update)
	return err
}

func (r *BugRepository) Update(ctx context.Context, bug *models.Bug) error {
	collection := r.db.Collection("bugs")

	bug.UpdatedAt = time.Now()

	_, err := collection.ReplaceOne(
		ctx,
		bson.M{"_id": bug.ID},
		bug,
	)

	return err
}

func (r *BugRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	collection := r.db.Collection("bugs")

	_, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
