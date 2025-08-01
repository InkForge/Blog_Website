package repositories

import (
	"context"

	"github.com/InkForge/Blog_Website/domain"
	"github.com/InkForge/Blog_Website/repositories/mongo/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CommentReactionMongoRepository struct {
	reactionCollection *mongo.Collection
}

func NewCommentReactionMongoRepository(db *mongo.Database) *CommentReactionMongoRepository {
	return &CommentReactionMongoRepository{
		reactionCollection: db.Collection("comment_reactions"),
	}
}

func (c CommentReactionMongoRepository) Create(ctx context.Context, reaction domain.CommentReaction) error {
	reactionMongo := models.FromDomainCommentReaction(&reaction)
	result, err := c.reactionCollection.InsertOne(ctx, reactionMongo)
	if err != nil {
		return domain.ErrInsertingDocuments
	}
	
	// Safe type assertion with error handling
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		_ = oid // Use the ObjectID if needed in the future
	}
	return nil
}

func (c CommentReactionMongoRepository) GetByCommentAndUser(ctx context.Context, commentID, userID string) (domain.CommentReaction, error) {
	filter := bson.M{
		"comment_id": commentID,
		"user_id":    userID,
	}
	var reactionModel models.CommentReactionMongo

	err := c.reactionCollection.FindOne(ctx, filter).Decode(&reactionModel)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return domain.CommentReaction{}, domain.ErrCommentReactionNotFound
		}
		return domain.CommentReaction{}, domain.ErrRetrievingDocuments
	}

	return *reactionModel.ToDomain(), nil
}

func (c CommentReactionMongoRepository) Update(ctx context.Context, reaction domain.CommentReaction) error {
	filter := bson.M{
		"comment_id": reaction.Comment_id,
		"user_id":    reaction.User_id,
	}
	reactionMongo := models.FromDomainCommentReaction(&reaction)
	
	update := bson.M{
		"$set": bson.M{
			"action":     reactionMongo.Action,
			"updated_at": reactionMongo.UpdatedAt,
		},
	}

	result, err := c.reactionCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return domain.ErrUpdatingDocument
	}
	if result.MatchedCount == 0 {
		return domain.ErrCommentReactionNotFound
	}
	return nil
}

func (c CommentReactionMongoRepository) Delete(ctx context.Context, commentID, userID string) error {
	filter := bson.M{
		"comment_id": commentID,
		"user_id":    userID,
	}

	result, err := c.reactionCollection.DeleteOne(ctx, filter)
	if err != nil {
		return domain.ErrDeletingDocument
	}

	if result.DeletedCount == 0 {
		return domain.ErrCommentReactionNotFound
	}

	return nil
}

func (c CommentReactionMongoRepository) GetReactionCounts(ctx context.Context, commentID string) (int, int, error) {
	// Count likes (action = 1)
	likeFilter := bson.M{
		"comment_id": commentID,
		"action":     1,
	}
	likeCount, err := c.reactionCollection.CountDocuments(ctx, likeFilter)
	if err != nil {
		return 0, 0, domain.ErrRetrievingDocuments
	}

	// Count dislikes (action = -1)
	dislikeFilter := bson.M{
		"comment_id": commentID,
		"action":     -1,
	}
	dislikeCount, err := c.reactionCollection.CountDocuments(ctx, dislikeFilter)
	if err != nil {
		return 0, 0, domain.ErrRetrievingDocuments
	}

	return int(likeCount), int(dislikeCount), nil
} 