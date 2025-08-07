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
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BlogMongoRepository struct {
	blogCollection *mongo.Collection
}

func NewBlogMongoRepository(db *mongo.Database) domain.IBlogRepository {
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

func (b *BlogMongoRepository) GetAll(ctx context.Context, page, limit int) ([]domain.Blog, int, error) {
	var blogs []domain.Blog
	filter := bson.M{}

	skip := int64((page - 1) * limit)

	findOptions := options.Find()
	findOptions.SetSkip(skip)
	findOptions.SetLimit(int64(limit))

	total, err := b.blogCollection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, domain.ErrRetrievingDocuments
	}

	cursor, err := b.blogCollection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, domain.ErrRetrievingDocuments
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var mongoBlog models.MongoBlog
		if err := cursor.Decode(&mongoBlog); err != nil {
			return nil, 0, domain.ErrDecodingDocument
		}
		blogs = append(blogs, *mongoBlog.ToDomain())
	}

	if err := cursor.Err(); err != nil {
		return nil, 0, domain.ErrCursorIteration
	}

	return blogs, int(total), nil
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

func (b *BlogMongoRepository) Search(ctx context.Context, title string, user_ids []string, page, limit int) ([]domain.Blog, int, error) {
	filter := bson.M{}
	if title != "" {
		filter["title"] = bson.M{"$regex": title, "$options": "i"}
	}
	if len(user_ids) > 0 {
		filter["user_id"] = bson.M{"$in": user_ids}
	}

	skip := int64((page - 1) * limit)
	findOptions := options.Find()
	findOptions.SetSkip(skip)
	findOptions.SetLimit(int64(limit))

	total, err := b.blogCollection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, domain.ErrRetrievingDocuments
	}

	cursor, err := b.blogCollection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, domain.ErrQueryFailed
	}
	defer cursor.Close(ctx)

	var blogs []domain.Blog
	for cursor.Next(ctx) {
		var mongoBlog models.MongoBlog
		if err := cursor.Decode(&mongoBlog); err != nil {
			return nil, 0, domain.ErrDocumentDecoding
		}
		blogs = append(blogs, *mongoBlog.ToDomain())
	}

	if err := cursor.Err(); err != nil {
		return nil, 0, domain.ErrCursorFailed
	}
	return blogs, int(total), nil
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

// at the bottom of BlogMongoRepository:

// Filter implements filtering by tag, date, and popularity
func (b *BlogMongoRepository) Filter(ctx context.Context, params domain.FilterParams) ([]domain.Blog, int, error) {
	filter := bson.M{}
	if len(params.TagIDs) > 0 {
		filter["tag_ids"] = bson.M{"$in": params.TagIDs}
	}

	findOptions := options.Find()
	skip := int64((params.Page - 1) * params.Limit)
	findOptions.SetSkip(skip)
	findOptions.SetLimit(int64(params.Limit))

	// Sort by popularity
	switch params.Popularity {
	case "views":
		findOptions.SetSort(bson.D{{Key: "view_count", Value: -1}})
	case "comments":
		findOptions.SetSort(bson.D{{Key: "comment_count", Value: -1}})
	case "likes":
		findOptions.SetSort(bson.D{{Key: "like_count", Value: -1}})
	case "dislikes":
		findOptions.SetSort(bson.D{{Key: "dislike_count", Value: -1}})
	}

	total, err := b.blogCollection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, domain.ErrRetrievingDocuments
	}

	cursor, err := b.blogCollection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, domain.ErrRetrievingDocuments
	}
	defer cursor.Close(ctx)

	var blogs []domain.Blog
	for cursor.Next(ctx) {
		var mongoBlog models.MongoBlog
		if err := cursor.Decode(&mongoBlog); err != nil {
			return nil, 0, domain.ErrDecodingDocument
		}
		blogs = append(blogs, *mongoBlog.ToDomain())
	}
	if err := cursor.Err(); err != nil {
		return nil, 0, domain.ErrCursorIteration
	}
	return blogs, int(total), nil
}

// Operations related to comments
func (r *BlogMongoRepository) AddCommentID(ctx context.Context, blogID, commentID string) error {
	objID, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return domain.ErrInvalidBlogIdFormat
	}

	_, err = r.blogCollection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{
		"$inc":  bson.M{"comment_count": 1},
	})
	if err != nil {
		return domain.ErrUpdatingDocument
	}
	return nil
}

func (r *BlogMongoRepository) RemoveCommentID(ctx context.Context, blogID, commentID string) error {
	objID, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return domain.ErrInvalidBlogIdFormat
	}

	_, err = r.blogCollection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{
		"$inc":  bson.M{"comment_count": -1},
	})
	if err != nil {
		return domain.ErrUpdatingDocument
	}
	return nil
}
