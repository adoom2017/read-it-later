package main

import (
	"log"
	"os"
	"path/filepath"
	"read-it-later/backend/handler"
	"read-it-later/backend/store"

	"github.com/gin-gonic/gin"
)

func main() {
	// 设置数据库文件路径
	dataDir := os.Getenv("DATA_DIR")
	if dataDir == "" {
		dataDir = "/app/data"
	}

	// 确保数据目录存在
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		log.Fatalf("Failed to create data directory: %v", err)
	}

	dbFileName := filepath.Join(dataDir, "read-it-later.db")

	// Initialize the database
	store.InitDB(dbFileName)
	log.Println("Database initialized successfully at:", dbFileName)

	// Set up the Gin router
	router := gin.Default()

	// 添加 CORS 中间件
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// API routes
	api := router.Group("/api")
	{
		api.GET("/articles", handler.GetArticles)
		api.GET("/articles/search", handler.SearchArticles)
		api.POST("/articles", handler.AddArticle)
		api.GET("/articles/:id", handler.GetArticle)
		api.POST("/articles/:id/tags", handler.AddTagToArticle)
		api.DELETE("/articles/:id/tags/:tagId", handler.RemoveTagFromArticle)
		api.DELETE("/articles/:id", handler.DeleteArticle)

		// Image proxy to handle anti-hotlinking
		api.GET("/proxy/image", handler.ProxyImage)
	}

	// Simple health check route
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Server is running",
			"status":  "healthy",
		})
	})

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "healthy",
		})
	})

	// 获取服务器端口
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start the server
	log.Printf("Starting server on :%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
