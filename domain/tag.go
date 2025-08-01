package domain

import "context"

type Tag struct {
	Tag_id  string
	TagName string
}

type ITagRepository interface {
	FindByNames(ctx context.Context, names []string) ([]Tag, error)
	CreateMany(ctx context.Context, names []string) ([]Tag, error)
}
