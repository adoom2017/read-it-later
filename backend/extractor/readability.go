package extractor

import (
	"net/http"
	"net/url"
	"time"

	"github.com/go-shiori/go-readability"
	"read-it-later/backend/model"
)

// Extract fetches the content from a URL and uses go-readability to parse it
// into a structured Article object.
func Extract(urlString string) (model.Article, error) {
	// Parse the URL string
	parsedURL, err := url.ParseRequestURI(urlString)
	if err != nil {
		return model.Article{}, err
	}

	// Make HTTP request
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Get(urlString)
	if err != nil {
		return model.Article{}, err
	}
	defer resp.Body.Close()

	// Use go-readability to parse the document
	article, err := readability.FromReader(resp.Body, parsedURL)
	if err != nil {
		return model.Article{}, err
	}

	// Create our own Article model from the parsed data
	result := model.Article{
		URL:      urlString,
		Title:    article.Title,
		Content:  article.TextContent, // Using TextContent for a cleaner reading view
		Excerpt:  article.Excerpt,
		ImageURL: article.Image,
	}

	return result, nil
}
