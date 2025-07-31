package repositories

import (
	"context"

	"github.com/InkForge/Blog_Website/domain"
	"github.com/InkForge/Blog_Website/repositories/mongo/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BlogMongoRepository struct {
	blogCollection *mongo.Collection
}

func NewBlogMongoRepository(db *mongo.Database) *BlogMongoRepository {
	return &BlogMongoRepository{
		blogCollection: db.Collection("blogs"),
	}
}

func (b *BlogMongoRepository) Create(ctx context.Context, blog domain.Blog) (string, error) {
	blogMongo := models.FromDomain(&blog)
	result, err := b.blogCollection.InsertOne(ctx, blogMongo)
	if err != nil {
		return "", domain.ErrInsertingDocuments
	}
	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (b *BlogMongoRepository) GetAll(ctx context.Context) ([]domain.Blog, error) {
	var blogs []domain.Blog
	filter := bson.M{}
	cursor, err := b.blogCollection.Find(ctx, filter)
	if err != nil {
		return nil, domain.ErrRetrievingDocuments
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var blogMongo models.BlogMongo
		err := cursor.Decode(&blogMongo)
		if err != nil {
			return nil, domain.ErrDecodingDocument
		}
		blogs = append(blogs, *blogMongo.ToDomain())
	}

	if err := cursor.Err(); err != nil {
		return nil, domain.ErrCursorIteration
	}
	return blogs, nil
}

func (b *BlogMongoRepository) GetByID(ctx context.Context, blogID string) (domain.Blog, error) {
	objID, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return domain.Blog{}, domain.ErrInvalidBlogID
	}

	filter := bson.M{"_id": objID}
	var blogModel models.BlogMongo

	err = b.blogCollection.FindOne(ctx, filter).Decode(&blogModel)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return domain.Blog{}, domain.ErrBlogNotFound
		}
		return domain.Blog{}, domain.ErrRetrievingDocuments
	}

	return *blogModel.ToDomain(), nil
}

func (b *BlogMongoRepository) Update(ctx context.Context, blog domain.Blog) error {
	objID, err := primitive.ObjectIDFromHex(blog.Blog_id)
	if err != nil {
		return domain.ErrInvalidBlogID
	}

	filter := bson.M{"_id": objID}
	blogMongo := models.FromDomain(&blog)

	update := bson.M{
		"$set": bson.M{
			"title":      blogMongo.Title,
			"content":    blogMongo.Content,
			"user_id":    blogMongo.UserID,
			"images":     blogMongo.Images,
			"tag_ids":    blogMongo.TagIDs,
			"posted_at":  blogMongo.PostedAt,
			"updated_at": blogMongo.UpdatedAt,
		},
	}

	result, err := b.blogCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return domain.ErrUpdatingDocument
	}
	if result.MatchedCount == 0 {
		return domain.ErrBlogNotFound
	}
	return nil
}

func (b *BlogMongoRepository) Delete(ctx context.Context, blogID string) error {
	objID, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return domain.ErrInvalidBlogID
	}

	filter := bson.M{"_id": objID}

	result, err := b.blogCollection.DeleteOne(ctx, filter)
	if err != nil {
		return domain.ErrDeletingDocument
	}

	if result.DeletedCount == 0 {
		return domain.ErrBlogNotFound
	}

	return nil
}
