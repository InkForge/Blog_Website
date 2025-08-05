package repositories

import (
	"context"
	"time"

	"github.com/InkForge/Blog_Website/domain"
	"github.com/InkForge/Blog_Website/repositories/mongo/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BlogViewRepository struct {
	collection *mongo.Collection
}

func NewBlogViewRepository(db *mongo.Database) domain.IBlogViewRepository {
	collection := db.Collection("blog_views")
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "blog_id", Value: 1}, {Key: "user_id", Value: 1}},
		Options: (&options.IndexOptions{}).SetUnique(true),
	}
	_, err := collection.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
	}

	return &BlogViewRepository{
		collection: collection,
	}
}

func (r *BlogViewRepository) CreateViewRecord(ctx context.Context, blogID string, userID string) error {
	domainView := &domain.BlogView{
		Blog_id:  blogID,
		User_id:  userID,
		ViewedAt: time.Now(),
	}
	mongoBlogView, err := models.FromDomainBlogView(domainView)
	if err != nil {
		return domain.ErrInvalidBlogIdFormat
	}
	_, err = r.collection.InsertOne(ctx, mongoBlogView)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return domain.ErrViewRecordAlreadyExists
		}
		return domain.ErrInsertingDocuments
	}
	return nil
}
