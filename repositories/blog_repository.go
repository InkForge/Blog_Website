package repositories

import (
	"context"
	"fmt"

	"github.com/InkForge/Blog_Website/domain"
	"github.com/InkForge/Blog_Website/infrastructures/mongo/models"
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
	blog_mongo := models.FromDomain(&blog)
	result, err := b.blogCollection.InsertOne(ctx, blog_mongo)
	if err != nil {
		return "", err
	}
	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (b *BlogMongoRepository) GetAll(ctx context.Context) ([]domain.Blog, error) {
	var blogs []domain.Blog
	filter := bson.M{}
	cursor, err := b.blogCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var blogMongo models.BlogMongo
		err := cursor.Decode(&blogMongo)
		if err != nil {
			return nil, err
		}
		blogs = append(blogs, *blogMongo.ToDomain())
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return blogs, nil
}

func (b *BlogMongoRepository) GetByID(ctx context.Context, blogID string) (domain.Blog, error) {
	filter := bson.M{"blog_id": blogID}
	var blogModel models.BlogMongo

	err := b.blogCollection.FindOne(ctx, filter).Decode(&blogModel)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return domain.Blog{}, fmt.Errorf("blog with ID '%s' not found", blogID)
		}
		return domain.Blog{}, fmt.Errorf("failed to get blog by ID '%s' from MongoDB: %w", blogID, err)
	}

	return *blogModel.ToDomain(), nil
}

func (b *BlogMongoRepository) Update(ctx context.Context, blog domain.Blog) error {
	filter := bson.M{"blog_id": blog.Blog_id}
	blogMongo := models.FromDomain(&blog)

	// Only update user-editable fields (system-managed fields are excluded)
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
		return fmt.Errorf("failed to update blog: %w", err)
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("blog with ID '%s' not found", blog.Blog_id)
	}
	return nil
}

func (b *BlogMongoRepository) Delete(ctx context.Context, blogID string) error {
	filter := bson.M{"blog_id": blogID}

	result, err := b.blogCollection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete blog: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("blog with ID '%s' not found", blogID)
	}

	return nil
}
