package models

import (
	"fmt"

	"github.com/InkForge/Blog_Website/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MongoTag struct {
	Tag_id  primitive.ObjectID `bson:"_id,omitempty"`
	TagName string             `bson:"tag_name"`
}

func TagFromDomain(tag *domain.Tag) (*MongoTag, error) {
	var objID primitive.ObjectID
	if tag.Tag_id != "" {
		var err error
		objID, err = primitive.ObjectIDFromHex(tag.Tag_id)
		if err != nil {
			return nil, fmt.Errorf("invalid tag id format: %v", err)
		}
	}
	return &MongoTag{
		Tag_id:  objID,
		TagName: tag.TagName,
	}, nil
}

func (m *MongoTag) ToDomain() *domain.Tag {
	return &domain.Tag{
		Tag_id:  m.Tag_id.Hex(),
		TagName: m.TagName,
	}
}
