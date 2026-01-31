package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/OkanUysal/go-logger"
)

// ProjectConfig holds project configuration
type ProjectConfig struct {
	Name       string
	ModulePath string
	Structure  string // "simple" or "standard"
	Database   string // "postgres", "mysql", "mongodb", "none"
	Libraries  []string
	Deployment string // "railway", "local", "docker"
	OutputDir  string
}

// GenerateProject generates a complete project
func GenerateProject(config *ProjectConfig) error {
	logger.Info("Starting project generation",
		logger.String("name", config.Name),
		logger.String("structure", config.Structure),
		logger.Int("libraries", len(config.Libraries)),
	)

	// Create output directory
	if err := os.MkdirAll(config.OutputDir, 0755); err != nil {
		logger.Error("Failed to create output directory", logger.Err(err))
		return err
	}

	// Set defaults
	if config.Structure == "" {
		config.Structure = "simple"
	}
	if config.Deployment == "" {
		config.Deployment = "railway"
	}
	if config.Database == "" {
		config.Database = "none"
	}

	logger.Debug("Creating directory structure")
	// Create directory structure
	if err := createDirectoryStructure(config); err != nil {
		logger.Error("Failed to create directory structure", logger.Err(err))
		return err
	}

	// Generate files
	logger.Debug("Generating go.mod")
	if err := generateGoMod(config); err != nil {
		logger.Error("Failed to generate go.mod", logger.Err(err))
		return err
	}

	logger.Debug("Generating main.go")
	if err := generateMain(config); err != nil {
		logger.Error("Failed to generate main.go", logger.Err(err))
		return err
	}

	logger.Debug("Generating config")
	if err := generateConfig(config); err != nil {
		logger.Error("Failed to generate config", logger.Err(err))
		return err
	}

	logger.Debug("Generating handlers")
	if err := generateHandlers(config); err != nil {
		logger.Error("Failed to generate handlers", logger.Err(err))
		return err
	}

	if config.hasLibrary("go-auth") {
		logger.Debug("Generating auth middleware")
		if err := generateMiddleware(config); err != nil {
			logger.Error("Failed to generate middleware", logger.Err(err))
			return err
		}
	}

	logger.Debug("Generating env files")
	if err := generateEnvFiles(config); err != nil {
		logger.Error("Failed to generate env files", logger.Err(err))
		return err
	}

	logger.Debug("Generating .gitignore")
	if err := generateGitignore(config); err != nil {
		logger.Error("Failed to generate .gitignore", logger.Err(err))
		return err
	}

	if config.Deployment == "railway" {
		logger.Debug("Generating Railway config")
		if err := generateRailwayConfig(config); err != nil {
			logger.Error("Failed to generate Railway config", logger.Err(err))
			return err
		}
	}

	logger.Debug("Generating README")
	if err := generateReadme(config); err != nil {
		logger.Error("Failed to generate README", logger.Err(err))
		return err
	}

	logger.Info("Project generation completed successfully", logger.String("output", config.OutputDir))
	return nil
}

// createDirectoryStructure creates project directories
func createDirectoryStructure(config *ProjectConfig) error {
	dirs := []string{"config"}

	if config.Structure == "standard" {
		dirs = append(dirs,
			"cmd/server",
			"internal/handlers",
			"internal/middleware",
			"internal/models",
		)
	}

	if config.hasLibrary("go-migration") {
		dirs = append(dirs, "migrations")
	}

	for _, dir := range dirs {
		path := filepath.Join(config.OutputDir, dir)
		if err := os.MkdirAll(path, 0755); err != nil {
			return err
		}
	}

	return nil
}

// hasLibrary checks if a library is selected
func (c *ProjectConfig) hasLibrary(lib string) bool {
	for _, l := range c.Libraries {
		if l == lib {
			return true
		}
	}
	return false
}

// writeFile writes content to a file
func writeFile(path, content string) error {
	return os.WriteFile(path, []byte(content), 0644)
}

// generateGoMod creates go.mod
func generateGoMod(config *ProjectConfig) error {
	content := fmt.Sprintf(`module %s

go 1.21

require (
	github.com/gin-gonic/gin v1.9.1
`, config.ModulePath)

	// Add library dependencies
	for _, lib := range config.Libraries {
		content += fmt.Sprintf("\tgithub.com/OkanUysal/%s v1.0.0\n", lib)
	}

	// Add database drivers
	switch config.Database {
	case "postgres":
		content += "\tgithub.com/lib/pq v1.10.9\n"
		content += "\tgorm.io/driver/postgres v1.5.7\n"
		content += "\tgorm.io/gorm v1.25.9\n"
	case "mysql":
		content += "\tgorm.io/driver/mysql v1.5.6\n"
		content += "\tgorm.io/gorm v1.25.9\n"
	case "mongodb":
		content += "\tgo.mongodb.org/mongo-driver v1.14.0\n"
	}

	content += ")\n"

	return writeFile(filepath.Join(config.OutputDir, "go.mod"), content)
}

// generateMain creates main.go
func generateMain(config *ProjectConfig) error {
	mainPath := "main.go"
	if config.Structure == "standard" {
		mainPath = "cmd/server/main.go"
	}

	imports := []string{
		`"log"`,
		`"github.com/gin-gonic/gin"`,
		fmt.Sprintf(`"%s/config"`, config.ModulePath),
	}

	for _, lib := range config.Libraries {
		switch lib {
		case "go-logger":
			imports = append(imports, `"github.com/OkanUysal/go-logger"`)
		case "go-metrics":
			imports = append(imports, `"github.com/OkanUysal/go-metrics"`)
		case "go-migration":
			imports = append(imports, `"github.com/OkanUysal/go-migration"`)
		case "go-swagger":
			imports = append(imports, `"github.com/OkanUysal/go-swagger"`)
		}
	}

	content := "package main\n\nimport (\n"
	for _, imp := range imports {
		content += fmt.Sprintf("\t%s\n", imp)
	}
	content += ")\n\n"

	content += "func main() {\n"
	content += "\tcfg := config.Load()\n\n"

	if config.hasLibrary("go-logger") {
		content += "\tlogger.Init(logger.Config{\n"
		content += "\t\tLevel: cfg.LogLevel,\n"
		content += "\t})\n"
		content += "\tdefer logger.Sync()\n\n"
	}

	if config.hasLibrary("go-migration") && config.Database != "none" {
		content += "\tif err := migration.Up(cfg.DatabaseURL, \"./migrations\"); err != nil {\n"
		content += "\t\tlog.Fatal(err)\n"
		content += "\t}\n\n"
	}

	if config.hasLibrary("go-metrics") {
		content += "\tmetricsCollector := metrics.NewMetrics(metrics.Config{\n"
		content += "\t\tNamespace: cfg.AppName,\n"
		content += "\t})\n\n"
	}

	content += "\trouter := gin.Default()\n\n"

	if config.hasLibrary("go-metrics") {
		content += "\trouter.Use(metricsCollector.Middleware())\n\n"
	}

	content += "\trouter.GET(\"/health\", func(c *gin.Context) {\n"
	content += "\t\tc.JSON(200, gin.H{\"status\": \"ok\"})\n"
	content += "\t})\n\n"

	if config.hasLibrary("go-metrics") {
		content += "\trouter.GET(\"/metrics\", metricsCollector.Handler())\n\n"
	}

	content += "\tapi := router.Group(\"/api/v1\")\n"
	content += "\t{\n"
	content += "\t\tapi.GET(\"/users\", func(c *gin.Context) {\n"
	content += "\t\t\tc.JSON(200, gin.H{\"users\": []string{}})\n"
	content += "\t\t})\n"
	content += "\t}\n\n"

	content += "\tport := cfg.Port\n"
	content += "\tif port == \"\" {\n"
	content += "\t\tport = \"8080\"\n"
	content += "\t}\n\n"

	content += "\tlog.Printf(\"Starting server on port %s\", port)\n"
	content += "\tif err := router.Run(\":\" + port); err != nil {\n"
	content += "\t\tlog.Fatal(err)\n"
	content += "\t}\n"
	content += "}\n"

	return writeFile(filepath.Join(config.OutputDir, mainPath), content)
}

// generateConfig creates config/config.go
func generateConfig(config *ProjectConfig) error {
	content := `package config

import "os"

type Config struct {
	AppName string
	Port    string
`

	if config.Database != "none" {
		content += "\tDatabaseURL string\n"
	}

	if config.hasLibrary("go-auth") {
		content += "\tJWTSecret string\n"
	}

	if config.hasLibrary("go-logger") {
		content += "\tLogLevel string\n"
	}

	content += "}\n\n"
	content += "func Load() *Config {\n"
	content += "\treturn &Config{\n"
	content += fmt.Sprintf("\t\tAppName: getEnv(\"APP_NAME\", \"%s\"),\n", config.Name)
	content += "\t\tPort: getEnv(\"PORT\", \"8080\"),\n"

	if config.Database != "none" {
		content += "\t\tDatabaseURL: getEnv(\"DATABASE_URL\", \"\"),\n"
	}

	if config.hasLibrary("go-auth") {
		content += "\t\tJWTSecret: getEnv(\"JWT_SECRET\", \"\"),\n"
	}

	if config.hasLibrary("go-logger") {
		content += "\t\tLogLevel: getEnv(\"LOG_LEVEL\", \"info\"),\n"
	}

	content += "\t}\n}\n\n"
	content += "func getEnv(key, def string) string {\n"
	content += "\tif v := os.Getenv(key); v != \"\" {\n"
	content += "\t\treturn v\n"
	content += "\t}\n"
	content += "\treturn def\n"
	content += "}\n"

	return writeFile(filepath.Join(config.OutputDir, "config", "config.go"), content)
}

// generateHandlers creates handlers
func generateHandlers(config *ProjectConfig) error {
	handlerPath := "handlers.go"
	if config.Structure == "standard" {
		handlerPath = "internal/handlers/handlers.go"
	}

	pkg := "main"
	if config.Structure == "standard" {
		pkg = "handlers"
	}

	content := fmt.Sprintf("package %s\n\n", pkg)
	content += "import \"github.com/gin-gonic/gin\"\n\n"

	if config.hasLibrary("go-response") {
		content += "import \"github.com/OkanUysal/go-response\"\n\n"
	}

	content += "func GetUsers(c *gin.Context) {\n"
	if config.hasLibrary("go-response") {
		content += "\tresponse.Success(c, []string{})\n"
	} else {
		content += "\tc.JSON(200, gin.H{\"users\": []string{}})\n"
	}
	content += "}\n"

	return writeFile(filepath.Join(config.OutputDir, handlerPath), content)
}

// generateMiddleware creates middleware
func generateMiddleware(config *ProjectConfig) error {
	middlewarePath := "middleware.go"
	if config.Structure == "standard" {
		middlewarePath = "internal/middleware/auth.go"
	}

	pkg := "main"
	if config.Structure == "standard" {
		pkg = "middleware"
	}

	content := fmt.Sprintf("package %s\n\n", pkg)
	content += "import (\n"
	content += "\t\"github.com/gin-gonic/gin\"\n"
	content += "\t\"github.com/OkanUysal/go-auth\"\n"
	content += ")\n\n"

	content += "func AuthMiddleware() gin.HandlerFunc {\n"
	content += "\treturn func(c *gin.Context) {\n"
	content += "\t\ttoken := c.GetHeader(\"Authorization\")\n"
	content += "\t\tif token == \"\" {\n"
	content += "\t\t\tc.JSON(401, gin.H{\"error\": \"unauthorized\"})\n"
	content += "\t\t\tc.Abort()\n"
	content += "\t\t\treturn\n"
	content += "\t\t}\n\n"
	content += "\t\tif _, err := auth.ValidateToken(token); err != nil {\n"
	content += "\t\t\tc.JSON(401, gin.H{\"error\": \"invalid token\"})\n"
	content += "\t\t\tc.Abort()\n"
	content += "\t\t\treturn\n"
	content += "\t\t}\n\n"
	content += "\t\tc.Next()\n"
	content += "\t}\n"
	content += "}\n"

	return writeFile(filepath.Join(config.OutputDir, middlewarePath), content)
}

// generateEnvFiles creates .env files
func generateEnvFiles(config *ProjectConfig) error {
	env := fmt.Sprintf("APP_NAME=%s\n", config.Name)
	env += "PORT=8080\n"

	if config.Database != "none" {
		env += "DATABASE_URL=\n"
	}

	if config.hasLibrary("go-auth") {
		env += "JWT_SECRET=your-secret-key\n"
	}

	if config.hasLibrary("go-logger") {
		env += "LOG_LEVEL=info\n"
	}

	writeFile(filepath.Join(config.OutputDir, ".env"), env)
	writeFile(filepath.Join(config.OutputDir, ".env.example"), env)

	return nil
}

// generateGitignore creates .gitignore
func generateGitignore(config *ProjectConfig) error {
	content := `*.exe
*.dll
*.so
*.dylib
*.test
*.out
vendor/
.env
.DS_Store
tmp/
temp/
`

	return writeFile(filepath.Join(config.OutputDir, ".gitignore"), content)
}

// generateRailwayConfig creates railway.json
func generateRailwayConfig(config *ProjectConfig) error {
	startCmd := "go run main.go"
	if config.Structure == "standard" {
		startCmd = "go run cmd/server/main.go"
	}

	content := fmt.Sprintf(`{
  "$schema": "https://railway.app/railway.schema.json",
  "build": {
    "builder": "NIXPACKS"
  },
  "deploy": {
    "startCommand": "%s",
    "restartPolicyType": "ON_FAILURE",
    "restartPolicyMaxRetries": 10
  }
}
`, startCmd)

	return writeFile(filepath.Join(config.OutputDir, "railway.json"), content)
}

// generateReadme creates README.md
func generateReadme(config *ProjectConfig) error {
	content := fmt.Sprintf("# %s\n\n", config.Name)
	content += "A Go API generated with go-starter.\n\n"
	content += "## Features\n\n"

	for _, lib := range config.Libraries {
		libName := strings.TrimPrefix(lib, "go-")
		content += fmt.Sprintf("- %s\n", libName)
	}

	content += "\n## Getting Started\n\n"
	content += "```bash\n"
	content += "go mod tidy\n"
	content += "cp .env.example .env\n"
	content += "go run main.go\n"
	content += "```\n\n"

	content += "## API Endpoints\n\n"
	content += "- `GET /health` - Health check\n"
	if config.hasLibrary("go-metrics") {
		content += "- `GET /metrics` - Prometheus metrics\n"
	}
	content += "- `GET /api/v1/users` - Get users\n\n"

	if config.Deployment == "railway" {
		content += "## Deployment\n\n"
		content += "This project is ready for Railway deployment.\n\n"
		content += "1. Push to GitHub\n"
		content += "2. Connect to Railway\n"
		content += "3. Deploy!\n"
	}

	return writeFile(filepath.Join(config.OutputDir, "README.md"), content)
}
