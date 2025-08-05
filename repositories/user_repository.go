package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/InkForge/Blog_Website/domain"
	"github.com/InkForge/Blog_Website/repositories/mongo/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	userCollection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) domain.IUserRepository {
	return &UserRepository{
		userCollection: db.Collection("users"),
	}
}

// consolidateUserError extracts the type of error and returns it
// it is a helper function to avoid code repetition
func consolidateUserError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, mongo.ErrNoDocuments) {
		return domain.ErrUserNotFound
	}
	return err
}

// IsEmailVerified query and check if the specified id is verified or not
// returns ErrUserNotFound if user is not found
func (ur *UserRepository) IsEmailVerified(ctx context.Context, id string) (bool, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, domain.ErrInvalidUserID
	}
	filter := bson.D{{Key: "_id", Value: objID}}
	result := ur.userCollection.FindOne(ctx, filter)

	if err := consolidateUserError(result.Err()); err != nil {
		return false, err
	}

	var userModel models.User
	if err := result.Decode(&userModel); err != nil {
		return false, domain.ErrDecodingDocument
	}
	return userModel.IsVerified, nil
}

// SetEmailVerified sets the specified id user is_verified entry to true
// returns ErrUserNotFound if user is not found
func (ur *UserRepository) SetEmailVerified(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.ErrInvalidUserID
	}

	// prepare the update data
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "is_verified", Value: true},
			{Key: "updated_at", Value: time.Now().Format("2006-01-02")},
		}},
	}

	result, err := ur.userCollection.UpdateByID(ctx, objID, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

func (ur *UserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, domain.ErrInvalidUserID
	}

	// prepate the filter
	filter := bson.D{{Key: "_id", Value: objID}}
	result := ur.userCollection.FindOne(ctx, filter)

	if err := consolidateUserError(result.Err()); err != nil {
		return nil, err
	}

	var userModel models.User
	if err := result.Decode(&userModel); err != nil {
		return nil, domain.ErrDecodingDocument
	}
	user := userModel.ToDomain()

	return &user, nil
}

// DeleteByID deletes a user by the id specified
// returns ErrUserNotFound if user is not found
func (ur *UserRepository) DeleteByID(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.ErrInvalidUserID
	}
	filter := bson.D{{Key: "_id", Value: objID}}
	result, err := ur.userCollection.DeleteOne(ctx, filter)
	if err != nil {
		return domain.ErrDeletingDocument
	}
	if result.DeletedCount == 0 {
		return domain.ErrUserNotFound
	}
	return nil
}

// Updates a user completely, it returns ErrInvalidUserID or ErrUserNotFound
func (ur *UserRepository) UpdateUser(ctx context.Context, user *domain.User) error {
	userData, err := models.UserFromDomain(*user)
	if err != nil {
		return err
	}
	filter := bson.D{{Key: "_id", Value: userData.UserID}}
	result, err := ur.userCollection.ReplaceOne(ctx, filter, userData)
	if err != nil {
		return err
	}
	// check if a user is found
	if result.MatchedCount == 0 {
		return domain.ErrUserNotFound
	}
	return nil
}

// CreateUser inserts the specified data into user collection
func (ur *UserRepository) CreateUser(ctx context.Context, user *domain.User) error {
	userModel, err := models.UserFromDomain(*user)
	if err != nil {
		return err
	}
	result, err := ur.userCollection.InsertOne(ctx, userModel)
	if err != nil {
		return domain.ErrUserCreationFailed
	}

	objID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return err
	}
	user.UserID = objID.Hex()
	return nil
}

// FindByEmail query the user collection based on specified email
func (ur *UserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	filter := bson.D{{Key: "email", Value: email}}
	result := ur.userCollection.FindOne(ctx, filter)

	if err := consolidateUserError(result.Err()); err != nil {
		return nil, err
	}

	var userModel models.User
	if err := result.Decode(&userModel); err != nil {
		return nil, domain.ErrDecodingDocument
	}

	user := userModel.ToDomain()

	return &user, nil
}

// FindByUsername query the user collection based on specified username
func (ur *UserRepository) FindByUserName(ctx context.Context, username string) (*domain.User, error) {
	filter := bson.D{{Key: "username", Value: username}}
	result := ur.userCollection.FindOne(ctx, filter)

	if err := consolidateUserError(result.Err()); err != nil {
		return nil, err
	}

	var userModel models.User
	if err := result.Decode(&userModel); err != nil {
		return nil, domain.ErrDecodingDocument
	}

	user := userModel.ToDomain()

	return &user, nil
}

// CountByEmail counts how many entry exist in the collection with the specified email
func (ur *UserRepository) CountByEmail(ctx context.Context, email string) (int64, error) {
	filter := bson.D{{Key: "email", Value: email}}
	return ur.userCollection.CountDocuments(ctx, filter)
}

// CountAll counts the number of documents inside user collection
func (ur *UserRepository) CountAll(ctx context.Context) (int64, error) {
	return ur.userCollection.CountDocuments(ctx, bson.D{})
}

// FindUsersByName searches users by first or last name (case-insensitive, partial match)
func (ur *UserRepository) FindUsersByName(ctx context.Context, name string) ([]*domain.User, error) {
	filter := bson.D{
		{Key: "$or", Value: bson.A{
			bson.D{{Key: "first_name", Value: bson.D{{Key: "$regex", Value: name}, {Key: "$options", Value: "i"}}}},
			bson.D{{Key: "last_name", Value: bson.D{{Key: "$regex", Value: name}, {Key: "$options", Value: "i"}}}},
		}},
	}
	cursor, err := ur.userCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []*domain.User
	for cursor.Next(ctx) {
		var userModel models.User
		if err := cursor.Decode(&userModel); err != nil {
			return nil, err
		}
		user := userModel.ToDomain()
		users = append(users, &user)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return users, nil
}
