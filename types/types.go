package types

// GenerateRequest represents the project generation request
type GenerateRequest struct {
	Name       string         `json:"name"`
	ModulePath string         `json:"modulePath"`
	Structure  string         `json:"structure"` // "simple" or "standard"
	Database   DatabaseConfig `json:"database"`
	Libraries  []string       `json:"libraries"`
	Deployment string         `json:"deployment"` // "railway", "local", "docker"
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Type string `json:"type"` // "postgres", "mysql", "mongodb", "none"
}

// Library represents an available library
type Library struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	Description string `json:"description"`
	Version     string `json:"version"`
	RepoURL     string `json:"repoURL"`
	Category    string `json:"category"`
	RequiresDB  bool   `json:"requiresDB"`
}

// GenerateResponse represents the generation response
type GenerateResponse struct {
	Success     bool   `json:"success"`
	DownloadURL string `json:"downloadUrl"`
	FileName    string `json:"fileName"`
	Message     string `json:"message,omitempty"`
	Error       string `json:"error,omitempty"`
}

// GetAvailableLibraries returns all available libraries with latest versions from GitHub
func GetAvailableLibraries() []Library {
	// Fetch latest versions from GitHub
	versions := FetchLatestVersions()

	// Default versions (fallback if GitHub API fails)
	defaultVersions := map[string]string{
		"go-auth":       "v1.0.0",
		"go-migration":  "v1.0.0",
		"go-logger":     "v1.0.1",
		"go-cache":      "v1.0.0",
		"go-swagger":    "v1.1.1",
		"go-response":   "v1.0.0",
		"go-validator":  "v1.0.0",
		"go-pagination": "v1.0.0",
		"go-websocket":  "v1.0.0",
		"go-metrics":    "v1.0.5",
	}

	// Helper to get version (GitHub or fallback)
	getVersion := func(name string) string {
		if v, ok := versions[name]; ok && v != "" {
			return v
		}
		return defaultVersions[name]
	}

	return []Library{
		{
			Name:        "go-auth",
			DisplayName: "Authentication",
			Description: "JWT authentication & authorization with middleware",
			Version:     getVersion("go-auth"),
			RepoURL:     "https://github.com/OkanUysal/go-auth",
			Category:    "Security",
			RequiresDB:  false,
		},
		{
			Name:        "go-migration",
			DisplayName: "Database Migration",
			Description: "Database migrations for PostgreSQL, MySQL, MongoDB",
			Version:     getVersion("go-migration"),
			RepoURL:     "https://github.com/OkanUysal/go-migration",
			Category:    "Database",
			RequiresDB:  true,
		},
		{
			Name:        "go-logger",
			DisplayName: "Structured Logger",
			Description: "High-performance structured logging with Zap",
			Version:     getVersion("go-logger"),
			RepoURL:     "https://github.com/OkanUysal/go-logger",
			Category:    "Observability",
			RequiresDB:  false,
		},
		{
			Name:        "go-cache",
			DisplayName: "Caching",
			Description: "Multi-backend caching (Redis, In-Memory)",
			Version:     getVersion("go-cache"),
			RepoURL:     "https://github.com/OkanUysal/go-cache",
			Category:    "Performance",
			RequiresDB:  false,
		},
		{
			Name:        "go-swagger",
			DisplayName: "API Documentation",
			Description: "Automatic Swagger/OpenAPI documentation",
			Version:     getVersion("go-swagger"),
			RepoURL:     "https://github.com/OkanUysal/go-swagger",
			Category:    "Documentation",
			RequiresDB:  false,
		},
		{
			Name:        "go-response",
			DisplayName: "API Response",
			Description: "Standardized API response format",
			Version:     getVersion("go-response"),
			RepoURL:     "https://github.com/OkanUysal/go-response",
			Category:    "API",
			RequiresDB:  false,
		},
		{
			Name:        "go-validator",
			DisplayName: "Request Validator",
			Description: "Request validation with custom rules",
			Version:     getVersion("go-validator"),
			RepoURL:     "https://github.com/OkanUysal/go-validator",
			Category:    "API",
			RequiresDB:  false,
		},
		{
			Name:        "go-pagination",
			DisplayName: "Pagination",
			Description: "Offset & cursor-based pagination",
			Version:     getVersion("go-pagination"),
			RepoURL:     "https://github.com/OkanUysal/go-pagination",
			Category:    "API",
			RequiresDB:  false,
		},
		{
			Name:        "go-websocket",
			DisplayName: "WebSocket",
			Description: "Real-time WebSocket with room management",
			Version:     getVersion("go-websocket"),
			RepoURL:     "https://github.com/OkanUysal/go-websocket",
			Category:    "Real-time",
			RequiresDB:  false,
		},
		{
			Name:        "go-metrics",
			DisplayName: "Metrics & Monitoring",
			Description: "Prometheus metrics with Grafana integration",
			Version:     getVersion("go-metrics"),
			RepoURL:     "https://github.com/OkanUysal/go-metrics",
			Category:    "Observability",
			RequiresDB:  false,
		},
	}
}
