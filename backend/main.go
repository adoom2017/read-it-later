package main

import (
	"log"
	"os"
	"path/filepath"
	"read-it-later/backend/handler"
	"read-it-later/backend/middleware"
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
		// 认证相关路由（公开访问）
		auth := api.Group("/auth")
		{
			auth.POST("/register", handler.Register)
			auth.POST("/login", handler.Login)
		}

		// 用户相关路由（需要认证）
		user := api.Group("/user")
		user.Use(middleware.AuthMiddleware())
		{
			user.GET("/profile", handler.GetProfile)
		}

		// 需要认证的文章相关路由
		articles := api.Group("/articles")
		articles.Use(middleware.AuthMiddleware())
		{
			articles.GET("", handler.GetArticles)
			articles.GET("/search", handler.SearchArticles)
			articles.POST("", handler.AddArticle)
			articles.GET("/:id", handler.GetArticle)
			articles.POST("/:id/tags", handler.AddTagToArticle)
			articles.DELETE("/:id/tags/:tagId", handler.RemoveTagFromArticle)
			articles.DELETE("/:id", handler.DeleteArticle)
		}

		// Image proxy to handle anti-hotlinking (公开访问)
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
