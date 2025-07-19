package handler

import (
	"database/sql"
	"net/http"
	"read-it-later/backend/extractor"
	"read-it-later/backend/model"
	"read-it-later/backend/store"
	"strconv"

	"github.com/gin-gonic/gin"
)

// AddArticle handles the creation of a new article from a URL.
func AddArticle(c *gin.Context) {
	var json struct {
		URL string `json:"url" binding:"required"`
	}

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Extract content from the URL
	article, err := extractor.Extract(json.URL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to extract article: " + err.Error()})
		return
	}

	// Save the article to the database
	savedArticle, err := store.SaveArticle(article)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save article: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, savedArticle)
}

// GetArticles handles listing all articles.
func GetArticles(c *gin.Context) {
	articles, err := store.GetAllArticles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve articles"})
		return
	}
	c.JSON(http.StatusOK, articles)
}

// SearchArticles handles searching articles by title or tags.
func SearchArticles(c *gin.Context) {
	query := c.Query("q")
	tag := c.Query("tag")

	if query == "" && tag == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query or tag is required"})
		return
	}

	var articles []model.Article
	var err error

	if tag != "" {
		// Search by tag
		articles, err = store.SearchArticlesByTag(tag)
	} else if query != "" {
		// Search by title
		articles, err = store.SearchArticlesByTitle(query)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search articles"})
		return
	}

	c.JSON(http.StatusOK, articles)
}

// GetArticle handles retrieving a single article by its ID.
func GetArticle(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid article ID"})
		return
	}

	article, err := store.GetArticleByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve article"})
		}
		return
	}

	c.JSON(http.StatusOK, article)
}

// AddTagToArticle handles adding a tag to an article.
func AddTagToArticle(c *gin.Context) {
	idParam := c.Param("id")
	articleID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid article ID"})
		return
	}

	var json struct {
		TagName string `json:"tag_name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = store.AddTagToArticleByID(articleID, json.TagName)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add tag to article"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tag added successfully"})
}

// RemoveTagFromArticle handles removing a tag from an article.
func RemoveTagFromArticle(c *gin.Context) {
	idParam := c.Param("id")
	articleID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid article ID"})
		return
	}

	tagIdParam := c.Param("tagId")
	tagID, err := strconv.Atoi(tagIdParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tag ID"})
		return
	}

	err = store.RemoveTagFromArticle(articleID, tagID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Article or tag not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove tag from article"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tag removed successfully"})
}

// DeleteArticle handles deleting an article.
func DeleteArticle(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid article ID"})
		return
	}

	err = store.DeleteArticleByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete article"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Article deleted successfully"})
}
