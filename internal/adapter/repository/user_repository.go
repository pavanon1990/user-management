package repository

import (
	"context"
	"errors"
	"user-management-api/internal/core/domain"
	"user-management-api/internal/core/entity"
	"user-management-api/internal/core/port"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type userRepository struct {
	collection *mongo.Collection
}

const USER_COLLECTION = "user"

func NewUserRepository(db *mongo.Database) port.UserRepository {
	return &userRepository{collection: db.Collection(USER_COLLECTION)}
}

func (r *userRepository) Insert(ctx context.Context, user entity.User) error {
	_, err := r.collection.InsertOne(ctx, user)
	return err
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (entity.User, error) {
	var user entity.User

	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entity.User{}, entity.ErrUserNotFound
		}
		return entity.User{}, err
	}

	return user, nil
}

func (r *userRepository) GetByID(ctx context.Context, id string) (entity.User, error) {
	var user entity.User

	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entity.User{}, entity.ErrUserNotFound
		}
		return entity.User{}, err
	}

	return user, nil
}

func (r *userRepository) List(ctx context.Context, filter domain.ListUsersFilter) ([]entity.User, error) {
	query := bson.M{
		"created_at": bson.M{
			"$gte": filter.From,
			"$lte": filter.To,
		},
	}

	opts := options.Find().SetSkip(filter.Offset).SetLimit(filter.Limit).SetSort(bson.M{"created_at": -1})

	cursor, err := r.collection.Find(ctx, query, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []entity.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *userRepository) Delete(ctx context.Context, id string) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return entity.ErrUserNotFound
	}
	return nil
}

func (r *userRepository) ExistsByEmail(ctx context.Context, email, excludeID string) (bool, error) {

	filter := bson.M{"email": email}
	if excludeID != "" {
		filter["_id"] = bson.M{"$ne": excludeID}
	}

	err := r.collection.FindOne(ctx, filter).Err()
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (r *userRepository) Update(ctx context.Context, id string, input domain.UpdateUserInput) (entity.User, error) {
	setFields := bson.M{}
	if input.Name != nil {
		setFields["name"] = *input.Name
	}
	if input.Email != nil {
		setFields["email"] = *input.Email
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updated entity.User
	err := r.collection.FindOneAndUpdate(ctx, bson.M{"_id": id}, bson.M{"$set": setFields}, opts).Decode(&updated)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entity.User{}, entity.ErrUserNotFound
		}
		return entity.User{}, err
	}

	return updated, nil
}

func (r *userRepository) Count(ctx context.Context) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{})
}
