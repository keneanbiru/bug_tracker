package repository

import (
	"context"
	"time"

	"bug-tracker/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepository struct {
	db *mongo.Database
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	collection := r.db.Collection("users")

	// Set timestamps
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	result, err := collection.InsertOne(ctx, user)
	if err != nil {
		return err
	}

	user.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	collection := r.db.Collection("users")

	var user models.User
	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.User, error) {
	collection := r.db.Collection("users")

	var user models.User
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	collection := r.db.Collection("users")

	user.UpdatedAt = time.Now()

	_, err := collection.ReplaceOne(
		ctx,
		bson.M{"_id": user.ID},
		user,
		options.Replace().SetUpsert(true),
	)

	return err
}

func (r *UserRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	collection := r.db.Collection("users")

	_, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *UserRepository) FindByRole(ctx context.Context, role string) ([]*models.User, error) {
	collection := r.db.Collection("users")

	cursor, err := collection.Find(ctx, bson.M{"role": role})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []*models.User
	if err = cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}
