package repositories

import (
	"context"
	"time"

	"github.com/InkForge/Blog_Website/domain"
	"github.com/InkForge/Blog_Website/repositories/mongo/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CommentMongoRepository struct {
	commentCollection *mongo.Collection
}

func NewCommentMongoRepository(db *mongo.Database) *CommentMongoRepository {
	return &CommentMongoRepository{
		commentCollection: db.Collection("comments"),
	}
}

func (c CommentMongoRepository) Create(ctx context.Context, comment domain.Comment) (string, error) {
	commentMongo := models.FromDomainComment(&comment)
	result, err := c.commentCollection.InsertOne(ctx, commentMongo)
	if err != nil {
		return "", domain.ErrInsertingDocuments
	}
	
	// Safe type assertion with error handling
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		return oid.Hex(), nil
	}
	return "", domain.ErrInsertingDocuments
}

func (c CommentMongoRepository) GetByID(ctx context.Context, commentID string) (domain.Comment, error) {
	filter := bson.M{"comment_id": commentID}
	var commentModel models.CommentMongo

	err := c.commentCollection.FindOne(ctx, filter).Decode(&commentModel)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return domain.Comment{}, domain.ErrCommentNotFound
		}
		return domain.Comment{}, domain.ErrRetrievingDocuments
	}

	return *commentModel.ToDomain(), nil
}

func (c CommentMongoRepository) GetByBlogID(ctx context.Context, blogID string) ([]domain.Comment, error) {
	var comments []domain.Comment
	filter := bson.M{"blog_id": blogID}
	cursor, err := c.commentCollection.Find(ctx, filter)
	if err != nil {
		return nil, domain.ErrRetrievingDocuments
	}
	defer cursor.Close(ctx)
	
	for cursor.Next(ctx) {
		var commentMongo models.CommentMongo
		err := cursor.Decode(&commentMongo)
		if err != nil {
			return nil, domain.ErrDecodingDocument
		}
		comments = append(comments, *commentMongo.ToDomain())
	}
	if err := cursor.Err(); err != nil {
		return nil, domain.ErrCursorIteration
	}
	return comments, nil
}

func (c CommentMongoRepository) Update(ctx context.Context, comment domain.Comment) error {
	filter := bson.M{"comment_id": comment.Comment_id}
	commentMongo := models.FromDomainComment(&comment)
	
	update := bson.M{
		"$set": bson.M{
			"content":    commentMongo.Content,
			"updated_at": commentMongo.UpdatedAt,
		},
	}

	result, err := c.commentCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return domain.ErrUpdatingDocument
	}
	if result.MatchedCount == 0 {
		return domain.ErrCommentNotFound
	}
	return nil
}

func (c CommentMongoRepository) Delete(ctx context.Context, commentID string) error {
	filter := bson.M{"comment_id": commentID}

	result, err := c.commentCollection.DeleteOne(ctx, filter)
	if err != nil {
		return domain.ErrDeletingDocument
	}

	if result.DeletedCount == 0 {
		return domain.ErrCommentNotFound
	}

	return nil
}

func (c CommentMongoRepository) UpdateReactionCounts(ctx context.Context, commentID string, likeCount, dislikeCount int) error {
	filter := bson.M{"comment_id": commentID}
	
	update := bson.M{
		"$set": bson.M{
			"like":       likeCount,
			"dislike":    dislikeCount,
			"updated_at": time.Now(),
		},
	}

	result, err := c.commentCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return domain.ErrUpdatingDocument
	}
	if result.MatchedCount == 0 {
		return domain.ErrCommentNotFound
	}
	return nil
} 