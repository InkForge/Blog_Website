package repositories

import (
	"context"

	"github.com/InkForge/Blog_Website/domain"
	"github.com/InkForge/Blog_Website/repositories/mongo/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TagMongoRepository struct {
	tagCollection *mongo.Collection
}

func NewTagMongoRepository(db *mongo.Database) *TagMongoRepository {
	return &TagMongoRepository{
		tagCollection: db.Collection("tags"),
	}
}

func (t *TagMongoRepository) FindByNames(ctx context.Context, names []string) ([]domain.Tag, error) {
	filter := bson.M{"tag_name": bson.M{"$in": names}}

	cursor, err := t.tagCollection.Find(ctx, filter)
	if err != nil {
		return nil, domain.ErrQueryFailed
	}
	defer cursor.Close(ctx)

	var mongoTags []models.MongoTag
	if err := cursor.All(ctx, &mongoTags); err != nil {
		return nil, domain.ErrDocumentDecoding
	}

	var tags []domain.Tag
	for _, mt := range mongoTags {
		tags = append(tags, *mt.ToDomain())
	}
	return tags, nil
}

func (t *TagMongoRepository) CreateMany(ctx context.Context, names []string) ([]domain.Tag, error) {
	var toInsert []interface{}
	for _, name := range names {
		mongoTag := models.MongoTag{TagName: name}
		toInsert = append(toInsert, mongoTag)
	}

	result, err := t.tagCollection.InsertMany(ctx, toInsert)
	if err != nil {
		return nil, domain.ErrInsertingDocuments
	}

	var tags []domain.Tag
	for i, id := range result.InsertedIDs {
		objID, ok := id.(primitive.ObjectID)
		if !ok {
			continue
		}
		tags = append(tags, domain.Tag{
			Tag_id:  objID.Hex(),
			TagName: names[i],
		})
	}
	return tags, nil
}
