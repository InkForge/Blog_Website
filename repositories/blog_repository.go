package repositories

import (
	"context"
	"errors"
	"reflect"
	"time"

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
	mongoBlog, err := models.FromDomain(&blog)
	if err != nil {
		return "", domain.ErrInvalidBlogIdFormat
	}
	// if they came in empty
	if mongoBlog.Created_at.IsZero() {
		mongoBlog.Created_at = time.Now()
	}
	if mongoBlog.Updated_at.IsZero() {
		mongoBlog.Updated_at = time.Now()
	}

	result, err := b.blogCollection.InsertOne(ctx, mongoBlog)
	if err != nil {
		return "", domain.ErrInsertingDocuments
	}
	// type assertion
	objectID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", domain.ErrInvalidBlogIdFormat
	}

	return objectID.Hex(), nil
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
		var mongoBlog models.MongoBlog
		if err := cursor.Decode(&mongoBlog); err != nil {
			return nil, domain.ErrDecodingDocument
		}
		blogs = append(blogs, *mongoBlog.ToDomain())
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
	var mongoBlog models.MongoBlog

	err = b.blogCollection.FindOne(ctx, filter).Decode(&mongoBlog)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return domain.Blog{}, domain.ErrBlogNotFound
		}
		return domain.Blog{}, domain.ErrRetrievingDocuments
	}

	return *mongoBlog.ToDomain(), nil
}

func (b *BlogMongoRepository) Update(ctx context.Context, blog domain.Blog) error {
	objID, err := primitive.ObjectIDFromHex(blog.Blog_id)
	if err != nil {
		return domain.ErrInvalidBlogID
	}

	var existingMongoBlog models.MongoBlog
	err = b.blogCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&existingMongoBlog)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.ErrBlogNotFound
		}
		return domain.ErrRetrievingDocuments
	}

	if blog.Title == existingMongoBlog.Title &&
		blog.Content == existingMongoBlog.Content &&
		reflect.DeepEqual(blog.Images, existingMongoBlog.Images) &&
		reflect.DeepEqual(blog.Tag_ids, existingMongoBlog.Tag_ids) {
		return domain.ErrNoBlogChangesMade
	}

	mongoBlog, err := models.FromDomain(&blog)
	if err != nil {
		return domain.ErrInvalidBlogIdFormat
	}
	mongoBlog.Updated_at = time.Now()

	update := bson.M{
		"$set": bson.M{
			"title":      mongoBlog.Title,
			"content":    mongoBlog.Content,
			"images":     mongoBlog.Images,
			"tag_ids":    mongoBlog.Tag_ids,
			"updated_at": mongoBlog.Updated_at,
		},
	}

	result, err := b.blogCollection.UpdateOne(ctx, bson.M{"_id": objID}, update)
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

func (b *BlogMongoRepository) Search(ctx context.Context, title string, user_ids []string) ([]domain.Blog, error) {
	filter := bson.M{
		"title": title,
		"user_id": bson.M{
			"$in": user_ids,
		},
	}
	var blogs []domain.Blog
	cursor, err := b.blogCollection.Find(ctx, filter)
	if err != nil {
		return nil, domain.ErrQueryFailed
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var mongoBlog models.MongoBlog
		if err := cursor.Decode(&mongoBlog); err != nil {
			return nil, domain.ErrDocumentDecoding
		}
		blogs = append(blogs, *mongoBlog.ToDomain())
	}

	if err := cursor.Err(); err != nil {
		return nil, domain.ErrCursorFailed
	}
	return blogs, nil
}

// related to Blog Reactions

func (r *BlogMongoRepository) IncrementLike(ctx context.Context, blogID string) error {
	objID, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return domain.ErrInvalidBlogIdFormat
	}

	_, err = r.blogCollection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{
		"$inc": bson.M{"like_count": 1},
	})
	if err != nil {
		return domain.ErrUpdatingDocument
	}
	return nil
}

func (r *BlogMongoRepository) DecrementLike(ctx context.Context, blogID string) error {
	objID, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return domain.ErrInvalidBlogIdFormat
	}

	_, err = r.blogCollection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{
		"$inc": bson.M{"like_count": -1},
	})
	if err != nil {
		return domain.ErrUpdatingDocument
	}
	return nil
}

func (r *BlogMongoRepository) IncrementDisLike(ctx context.Context, blogID string) error {
	objID, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return domain.ErrInvalidBlogIdFormat
	}

	_, err = r.blogCollection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{
		"$inc": bson.M{"dislike_count": 1},
	})
	if err != nil {
		return domain.ErrUpdatingDocument
	}
	return nil
}

func (r *BlogMongoRepository) DecrementDisLike(ctx context.Context, blogID string) error {
	objID, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return domain.ErrInvalidBlogIdFormat
	}

	_, err = r.blogCollection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{
		"$inc": bson.M{"dislike_count": -1},
	})
	if err != nil {
		return domain.ErrUpdatingDocument
	}
	return nil
}

func (r *BlogMongoRepository) ToggleLikeDislikeCounts(ctx context.Context, blogID string, to_like, to_dislike int) error {
	objID, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return domain.ErrInvalidBlogIdFormat
	}

	update := bson.M{}
	if to_like != 0 {
		update["like_count"] = to_like
	}
	if to_dislike != 0 {
		update["dislike_count"] = to_dislike
	}

	if len(update) == 0 {
		return nil
	}

	_, err = r.blogCollection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{
		"$inc": update,
	})
	if err != nil {
		return domain.ErrUpdatingDocument
	}
	return nil
}

// related to blog_view
func (r *BlogMongoRepository) IncrementView(ctx context.Context, blogID string) error {
	objID, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return domain.ErrInvalidBlogIdFormat
	}

	_, err = r.blogCollection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{
		"$inc": bson.M{"view_count": 1},
	})
	if err != nil {
		return domain.ErrUpdatingDocument
	}
	return nil
}
