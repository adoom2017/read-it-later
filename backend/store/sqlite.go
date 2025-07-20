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
	// 用户表
	usersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	// 修改文章表，添加用户ID
	articlesTable := `
	CREATE TABLE IF NOT EXISTS articles (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		url TEXT NOT NULL,
		title TEXT NOT NULL,
		content TEXT,
		excerpt TEXT,
		image_url TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		UNIQUE(user_id, url)
	);`

	tagsTable := `
	CREATE TABLE IF NOT EXISTS tags (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		name TEXT NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		UNIQUE(user_id, name)
	);`

	articleTagsTable := `
	CREATE TABLE IF NOT EXISTS article_tags (
		article_id INTEGER,
		tag_id INTEGER,
		PRIMARY KEY (article_id, tag_id),
		FOREIGN KEY (article_id) REFERENCES articles(id) ON DELETE CASCADE,
		FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
	);`

	// 执行表创建
	_, err := DB.Exec(usersTable)
	if err != nil {
		log.Fatalf("Error creating users table: %v", err)
	}

	_, err = DB.Exec(articlesTable)
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
	stmt, err := DB.Prepare("INSERT INTO articles(user_id, url, title, content, excerpt, image_url) VALUES(?, ?, ?, ?, ?, ?)")
	if err != nil {
		return model.Article{}, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(article.UserID, article.URL, article.Title, article.Content, article.Excerpt, article.ImageURL)
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

// GetAllArticles retrieves all articles from the database for a specific user.
func GetAllArticles(userID int) ([]model.Article, error) {
	rows, err := DB.Query("SELECT id, user_id, url, title, excerpt, image_url, created_at FROM articles WHERE user_id = ? ORDER BY created_at DESC", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []model.Article
	for rows.Next() {
		var article model.Article
		if err := rows.Scan(&article.ID, &article.UserID, &article.URL, &article.Title, &article.Excerpt, &article.ImageURL, &article.CreatedAt); err != nil {
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

// GetArticleByID retrieves a single article by its ID and user ID.
func GetArticleByID(id int, userID int) (model.Article, error) {
	var article model.Article
	err := DB.QueryRow("SELECT id, user_id, url, title, content, excerpt, image_url, created_at FROM articles WHERE id = ? AND user_id = ?", id, userID).
		Scan(&article.ID, &article.UserID, &article.URL, &article.Title, &article.Content, &article.Excerpt, &article.ImageURL, &article.CreatedAt)

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

// DeleteArticleByID deletes an article by its ID and user ID.
func DeleteArticleByID(id int, userID int) error {
	stmt, err := DB.Prepare("DELETE FROM articles WHERE id = ? AND user_id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(id, userID)
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

// GetOrCreateTag gets an existing tag or creates a new one for a specific user.
func GetOrCreateTag(tagName string, userID int) (model.Tag, error) {
	var tag model.Tag

	// Try to get existing tag for this user
	err := DB.QueryRow("SELECT id, name FROM tags WHERE name = ? AND user_id = ?", tagName, userID).Scan(&tag.ID, &tag.Name)
	if err == nil {
		return tag, nil
	}

	if err != sql.ErrNoRows {
		return model.Tag{}, err
	}

	// Create new tag for this user
	stmt, err := DB.Prepare("INSERT INTO tags(user_id, name) VALUES(?, ?)")
	if err != nil {
		return model.Tag{}, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(userID, tagName)
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
func AddTagToArticleByID(articleID int, tagName string, userID int) error {
	// Check if article exists and belongs to the user
	var exists bool
	err := DB.QueryRow("SELECT EXISTS(SELECT 1 FROM articles WHERE id = ? AND user_id = ?)", articleID, userID).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return sql.ErrNoRows
	}

	// Get or create tag for this user
	tag, err := GetOrCreateTag(tagName, userID)
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

// SearchArticlesByTitle searches articles by title using LIKE query for a specific user.
func SearchArticlesByTitle(query string, userID int) ([]model.Article, error) {
	rows, err := DB.Query("SELECT id, user_id, url, title, excerpt, image_url, created_at FROM articles WHERE title LIKE ? AND user_id = ? ORDER BY created_at DESC", "%"+query+"%", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []model.Article
	for rows.Next() {
		var article model.Article
		if err := rows.Scan(&article.ID, &article.UserID, &article.URL, &article.Title, &article.Excerpt, &article.ImageURL, &article.CreatedAt); err != nil {
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

// SearchArticlesByTag searches articles by tag name for a specific user.
func SearchArticlesByTag(tagName string, userID int) ([]model.Article, error) {
	rows, err := DB.Query(`
		SELECT a.id, a.user_id, a.url, a.title, a.excerpt, a.image_url, a.created_at
		FROM articles a
		JOIN article_tags at ON a.id = at.article_id
		JOIN tags t ON at.tag_id = t.id
		WHERE t.name LIKE ? AND a.user_id = ?
		ORDER BY a.created_at DESC`, "%"+tagName+"%", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []model.Article
	for rows.Next() {
		var article model.Article
		if err := rows.Scan(&article.ID, &article.UserID, &article.URL, &article.Title, &article.Excerpt, &article.ImageURL, &article.CreatedAt); err != nil {
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

// ===== 用户相关数据库操作 =====

// UserExists 检查用户名或邮箱是否已存在
func UserExists(username, email string) bool {
	query := "SELECT COUNT(*) FROM users WHERE username = ? OR email = ?"
	var count int
	err := DB.QueryRow(query, username, email).Scan(&count)
	if err != nil {
		log.Printf("Error checking user existence: %v", err)
		return false
	}
	return count > 0
}

// CreateUser 创建新用户
func CreateUser(user model.User) (int, error) {
	stmt, err := DB.Prepare("INSERT INTO users(username, email, password) VALUES(?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(user.Username, user.Email, user.Password)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// GetUserByUsername 根据用户名获取用户
func GetUserByUsername(username string) (*model.User, error) {
	query := "SELECT id, username, email, password, created_at FROM users WHERE username = ?"
	var user model.User
	err := DB.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByID 根据ID获取用户
func GetUserByID(userID int) (*model.User, error) {
	query := "SELECT id, username, email, password, created_at FROM users WHERE id = ?"
	var user model.User
	err := DB.QueryRow(query, userID).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
