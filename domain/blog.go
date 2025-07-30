package domain

import (
	"time"
)

type Blog struct {
	Blog_id        string    `json:"blog_id"`
	User_id        string    `json:"user_id"`
	Title          string    `json:"title"`
	Images         []string  `json:"images"`
	Content        string    `json:"content"`
	Tag_ids        []string  `json:"tag_ids"`
	Comment_ids    []string  `json:"comment_ids"`
	Posted_at      time.Time `json:"posted_at"`
	Like_counts    int       `json:"like_counts"`
	Dislike_counts int       `json:"dislike_counts"`
	Share_count    int       `json:"share_count"`
	Created_at     time.Time `json:"created_at"`
	Updated_at     time.Time `json:"updated_at"`
}
