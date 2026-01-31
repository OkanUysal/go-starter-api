package main

import (
	"log"

	"github.com/OkanUysal/go-logger"
	"github.com/OkanUysal/go-metrics"
	"github.com/OkanUysal/go-swagger"
	"github.com/gin-gonic/gin"

	_ "github.com/OkanUysal/go-starter-api/docs" // Import generated docs
	"github.com/OkanUysal/go-starter-api/handlers"

	docs "github.com/OkanUysal/go-starter-api/docs" // Explicit import for SwaggerInfo
)

var metricsInstance *metrics.Metrics

// @title           Go Starter API
// @version         1.0.0
// @description     REST API for Go project generator with 10 production libraries
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api

// @schemes http https
func main() {
	// Initialize logger
	loggerConfig := &logger.Config{
		Level:  logger.LevelInfo,
		Format: logger.FormatJSON,
	}
	logger.SetDefault(logger.New(loggerConfig))

	// Initialize metrics
	metricsInstance = metrics.NewMetrics(&metrics.Config{
		ServiceName: "go-starter-api",
	})

	// Initialize Gin
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// Setup metrics endpoints (/metrics and /health)
	metricsInstance.Setup(r)

	// Metrics middleware (automatic HTTP metrics collection)
	r.Use(metricsInstance.GinMiddleware())

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Set metrics instance for handlers
	handlers.SetMetrics(metricsInstance)

	// Swagger documentation with auto host detection
	swagSpec, err := swagger.LoadSwagDocs(docs.SwaggerInfo.ReadDoc())
	if err != nil {
		logger.Error("Failed to load swagger docs", logger.Err(err))
	} else {
		swagger.SetupWithSwag(r, swagSpec, swagger.DefaultConfig())
		logger.Info("Swagger UI enabled", logger.String("path", "/swagger/index.html"))
	}

	// API routes
	api := r.Group("/api")
	{
		api.GET("/libraries", handlers.GetLibraries)
		api.POST("/generate", handlers.GenerateProject)
	}

	// Start server
	port := ":8080"
	logger.Info("Server starting", logger.String("port", port))
	if err := r.Run(port); err != nil {
		log.Fatal(err)
	}
}
