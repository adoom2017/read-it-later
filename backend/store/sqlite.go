package store

import (
	"database/sql"
	"log"
	"read-it-later/backend/model"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

// InitDB initializes the SQLite database and creates tables if they don't exist.
func InitDB(dataSourceName string) {
	var err error
	DB, err = sql.Open("sqlite", dataSourceName)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	createTables()
}

// createTables creates the necessary tables in the database.
func createTables() {
	articlesTable := `
	CREATE TABLE IF NOT EXISTS articles (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		url TEXT NOT NULL UNIQUE,
		title TEXT NOT NULL,
		content TEXT,
		excerpt TEXT,
		image_url TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	tagsTable := `
	CREATE TABLE IF NOT EXISTS tags (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE
	);`

	articleTagsTable := `
	CREATE TABLE IF NOT EXISTS article_tags (
		article_id INTEGER,
		tag_id INTEGER,
		PRIMARY KEY (article_id, tag_id),
		FOREIGN KEY (article_id) REFERENCES articles(id) ON DELETE CASCADE,
		FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
	);`

	_, err := DB.Exec(articlesTable)
	if err != nil {
		log.Fatalf("Error creating articles table: %v", err)
	}

	_, err = DB.Exec(tagsTable)
	if err != nil {
		log.Fatalf("Error creating tags table: %v", err)
	}

	_, err = DB.Exec(articleTagsTable)
	if err != nil {
		log.Fatalf("Error creating article_tags table: %v", err)
	}
}

// SaveArticle inserts a new article into the database and returns it with the new ID.
func SaveArticle(article model.Article) (model.Article, error) {
	stmt, err := DB.Prepare("INSERT INTO articles(url, title, content, excerpt, image_url) VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		return model.Article{}, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(article.URL, article.Title, article.Content, article.Excerpt, article.ImageURL)
	if err != nil {
		return model.Article{}, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return model.Article{}, err
	}

	article.ID = int(id)

	return article, nil
}

// GetAllArticles retrieves all articles from the database.
func GetAllArticles() ([]model.Article, error) {
	rows, err := DB.Query("SELECT id, url, title, excerpt, image_url, created_at FROM articles ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []model.Article
	for rows.Next() {
		var article model.Article
		if err := rows.Scan(&article.ID, &article.URL, &article.Title, &article.Excerpt, &article.ImageURL, &article.CreatedAt); err != nil {
			return nil, err
		}

		// Get tags for this article
		tags, err := GetTagsForArticle(article.ID)
		if err != nil {
			// Log the error but don't fail the entire request
			log.Printf("Error getting tags for article %d: %v", article.ID, err)
			tags = []model.Tag{} // Empty slice instead of nil
		}
		article.Tags = tags

		articles = append(articles, article)
	}

	return articles, nil
}

// GetArticleByID retrieves a single article by its ID.
func GetArticleByID(id int) (model.Article, error) {
	var article model.Article
	err := DB.QueryRow("SELECT id, url, title, content, excerpt, image_url, created_at FROM articles WHERE id = ?", id).
		Scan(&article.ID, &article.URL, &article.Title, &article.Content, &article.Excerpt, &article.ImageURL, &article.CreatedAt)

	if err != nil {
		return model.Article{}, err
	}

	// Get tags for this article
	tags, err := GetTagsForArticle(id)
	if err != nil {
		return model.Article{}, err
	}
	article.Tags = tags

	return article, nil
}

// DeleteArticleByID deletes an article by its ID.
func DeleteArticleByID(id int) error {
	stmt, err := DB.Prepare("DELETE FROM articles WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// GetOrCreateTag gets an existing tag or creates a new one.
func GetOrCreateTag(tagName string) (model.Tag, error) {
	var tag model.Tag

	// Try to get existing tag
	err := DB.QueryRow("SELECT id, name FROM tags WHERE name = ?", tagName).Scan(&tag.ID, &tag.Name)
	if err == nil {
		return tag, nil
	}

	if err != sql.ErrNoRows {
		return model.Tag{}, err
	}

	// Create new tag
	stmt, err := DB.Prepare("INSERT INTO tags(name) VALUES(?)")
	if err != nil {
		return model.Tag{}, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(tagName)
	if err != nil {
		return model.Tag{}, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return model.Tag{}, err
	}

	tag.ID = int(id)
	tag.Name = tagName
	return tag, nil
}

// AddTagToArticleByID adds a tag to an article.
func AddTagToArticleByID(articleID int, tagName string) error {
	// Check if article exists
	var exists bool
	err := DB.QueryRow("SELECT EXISTS(SELECT 1 FROM articles WHERE id = ?)", articleID).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return sql.ErrNoRows
	}

	// Get or create tag
	tag, err := GetOrCreateTag(tagName)
	if err != nil {
		return err
	}

	// Add relationship
	stmt, err := DB.Prepare("INSERT OR IGNORE INTO article_tags(article_id, tag_id) VALUES(?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(articleID, tag.ID)
	return err
}

// RemoveTagFromArticle removes a tag from an article.
func RemoveTagFromArticle(articleID int, tagID int) error {
	stmt, err := DB.Prepare("DELETE FROM article_tags WHERE article_id = ? AND tag_id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(articleID, tagID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// GetTagsForArticle retrieves all tags for a specific article.
func GetTagsForArticle(articleID int) ([]model.Tag, error) {
	rows, err := DB.Query(`
		SELECT t.id, t.name
		FROM tags t
		JOIN article_tags at ON t.id = at.tag_id
		WHERE at.article_id = ?`, articleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []model.Tag
	for rows.Next() {
		var tag model.Tag
		if err := rows.Scan(&tag.ID, &tag.Name); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

// SearchArticlesByTitle searches articles by title using LIKE query.
func SearchArticlesByTitle(query string) ([]model.Article, error) {
	rows, err := DB.Query("SELECT id, url, title, excerpt, image_url, created_at FROM articles WHERE title LIKE ? ORDER BY created_at DESC", "%"+query+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []model.Article
	for rows.Next() {
		var article model.Article
		if err := rows.Scan(&article.ID, &article.URL, &article.Title, &article.Excerpt, &article.ImageURL, &article.CreatedAt); err != nil {
			return nil, err
		}

		// Get tags for this article
		tags, err := GetTagsForArticle(article.ID)
		if err != nil {
			log.Printf("Error getting tags for article %d: %v", article.ID, err)
			tags = []model.Tag{}
		}
		article.Tags = tags

		articles = append(articles, article)
	}

	return articles, nil
}

// SearchArticlesByTag searches articles by tag name.
func SearchArticlesByTag(tagName string) ([]model.Article, error) {
	rows, err := DB.Query(`
		SELECT a.id, a.url, a.title, a.excerpt, a.image_url, a.created_at
		FROM articles a
		JOIN article_tags at ON a.id = at.article_id
		JOIN tags t ON at.tag_id = t.id
		WHERE t.name LIKE ?
		ORDER BY a.created_at DESC`, "%"+tagName+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []model.Article
	for rows.Next() {
		var article model.Article
		if err := rows.Scan(&article.ID, &article.URL, &article.Title, &article.Excerpt, &article.ImageURL, &article.CreatedAt); err != nil {
			return nil, err
		}

		// Get tags for this article
		tags, err := GetTagsForArticle(article.ID)
		if err != nil {
			log.Printf("Error getting tags for article %d: %v", article.ID, err)
			tags = []model.Tag{}
		}
		article.Tags = tags

		articles = append(articles, article)
	}

	return articles, nil
}
