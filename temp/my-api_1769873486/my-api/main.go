package main

import (
	"log"
	"github.com/gin-gonic/gin"
	"my-api/config"
	"github.com/OkanUysal/go-logger"
)

func main() {
	cfg := config.Load()

	logger.Init(logger.Config{
		Level: cfg.LogLevel,
	})
	defer logger.Sync()

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := router.Group("/api/v1")
	{
		api.GET("/users", func(c *gin.Context) {
			c.JSON(200, gin.H{"users": []string{}})
		})
	}

	port := cfg.Port
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
