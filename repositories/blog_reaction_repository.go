package repositories

import (
	"context"

	"github.com/InkForge/Blog_Website/domain"
	"github.com/InkForge/Blog_Website/repositories/mongo/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type BlogReactionRepository struct {
	collection *mongo.Collection
}

func NewBlogReactionRepository(db *mongo.Database) domain.IBlogReactionRepository {
	return &BlogReactionRepository{
		collection: db.Collection("blog_reactions"),
	}
}

func (r *BlogReactionRepository) CreateReaction(ctx context.Context, blogReaction domain.BlogReaction) error {
	mongoBlogReaction, err := models.FromDomainBlogReaction(&blogReaction)
	if err != nil {
		return domain.ErrInvalidBlogIdFormat
	}
	_, err = r.collection.InsertOne(ctx, mongoBlogReaction)
	if err != nil {
		return domain.ErrInsertingDocuments
	}
	return nil
}

func (r *BlogReactionRepository) GetReactionByUserAndBlog(ctx context.Context, blogID, userID string) (domain.BlogReaction, error) {
	filter := bson.M{
		"user_id": userID,
		"blog_id": blogID,
	}

	var mongoBlogReaction models.MongoBlogReaction
	err := r.collection.FindOne(ctx, filter).Decode(&mongoBlogReaction)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return domain.BlogReaction{}, domain.ErrBlogReactionNotFound
		}
		return domain.BlogReaction{}, domain.ErrDecodingDocument
	}

	return *mongoBlogReaction.ToDomainBlogReaction(), nil
}

func (r *BlogReactionRepository) UpdateReaction(ctx context.Context, blogReaction domain.BlogReaction) error {
	mongoReaction, err := models.FromDomainBlogReaction(&blogReaction)
	if err != nil {
		return domain.ErrInvalidBlogIdFormat
	}

	filter := bson.M{"_id": mongoReaction.ID}
	update := bson.M{
		"$set": bson.M{
			"reaction_type": mongoReaction.Reaction_type,
		},
	}
	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return domain.ErrUpdateBlogReactionFailed
	}
	if res.MatchedCount == 0 {
		return domain.ErrBlogReactionNotFound
	}
	return nil
}

func (r *BlogReactionRepository) DeleteReaction(ctx context.Context, blogID, userID string) error {
	filter := bson.M{"blog_id": blogID, "user_id": userID}

	res, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return domain.ErrDeletingBlogReaction
	}
	if res.DeletedCount == 0 {
		return domain.ErrBlogReactionNotFound
	}
	return nil
}
