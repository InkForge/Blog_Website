package domain

import (
	"context"
	"time"
)

type Blog struct {
	Blog_id        string   
	User_id        string    
	Title          string    
	Images         []string  
	Content        string    
	Tag_ids        []string  
	Comment_ids    []string  
	Posted_at      time.Time 
	Like_counts    int       
	Dislike_counts int       
	Share_count    int       
	Created_at     time.Time 
	Updated_at     time.Time 
}

type IBlogRepository interface {
	Create(ctx context.Context, blog Blog) (string, error)
	GetAll(ctx context.Context) ([]Blog, error)
	GetByID(ctx context.Context, blogID string) (Blog, error)
	Update(ctx context.Context, blog Blog) error
	Delete(ctx context.Context, blogID string) error
}

