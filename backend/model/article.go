package model

import "time"

// Article represents a saved article.
type Article struct {
	ID          int       `json:"id"`
	URL         string    `json:"url"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	Excerpt     string    `json:"excerpt"`
	ImageURL    string    `json:"image_url"`
	CreatedAt   time.Time `json:"created_at"`
	Tags        []Tag     `json:"tags"`
}

// Tag represents a tag for an article.
type Tag struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
