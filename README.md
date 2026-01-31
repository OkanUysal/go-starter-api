# go-starter-api

Backend API for go-starter project generator.

## 🚀 Features

- REST API for project generation
- ZIP download support
- CORS enabled for Flutter web
- 10 production libraries support
- Simple & Standard project structures

## 📦 API Endpoints

### GET /health
Health check endpoint

**Response:**
```json
{
  "status": "ok",
  "service": "go-starter-api"
}
```

### GET /api/libraries
Get all available libraries

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "name": "go-auth",
      "displayName": "Authentication",
      "description": "JWT authentication & authorization",
      "version": "v1.0.0",
      "repoURL": "https://github.com/OkanUysal/go-auth",
      "category": "Security",
      "requiresDB": false
    }
  ],
  "count": 10
}
```

### POST /api/generate
Generate a new project and download as ZIP

**Request:**
```json
{
  "name": "my-api",
  "modulePath": "github.com/user/my-api",
  "structure": "simple",
  "database": {
    "type": "postgres"
  },
  "libraries": [
    "go-auth",
    "go-logger",
    "go-response",
    "go-metrics"
  ],
  "deployment": "railway"
}
```

**Response:**
- Content-Type: `application/zip`
- File download: `my-api.zip`

## 🛠️ Development

### Install dependencies
```bash
go mod tidy
```

### Run server
```bash
go run main.go
```

Server starts on `http://localhost:8080`

### Test endpoints

```bash
# Health check
curl http://localhost:8080/health

# Get libraries
curl http://localhost:8080/api/libraries

# Generate project
curl -X POST http://localhost:8080/api/generate \
  -H "Content-Type: application/json" \
  -d '{
    "name": "demo-api",
    "modulePath": "github.com/demo/demo-api",
    "structure": "simple",
    "database": {"type": "postgres"},
    "libraries": ["go-auth", "go-logger", "go-response"],
    "deployment": "railway"
  }' \
  --output demo-api.zip
```

## 📁 Project Structure

```
go-starter-api/
├── main.go              # Server entry point
├── handlers/
│   ├── libraries.go     # GET /api/libraries
│   └── generate.go      # POST /api/generate (with ZIP)
├── generator/
│   └── generator.go     # Project generation logic
├── types/
│   └── types.go         # Type definitions
├── temp/                # Temporary ZIP files (auto-cleanup)
└── README.md
```

## 🌐 CORS

CORS is enabled for all origins to support Flutter web frontend.

## 🔧 Configuration

- **Port**: 8080 (configurable)
- **Body Limit**: 10MB
- **ZIP Cleanup**: 10 minutes after generation

## 📦 Dependencies

- **Fiber v2**: Web framework
- **Archive/zip**: ZIP file creation

## 🚀 Deployment

### Railway
```bash
# railway.json included in generated projects
git push origin main
```

### Docker
```bash
docker build -t go-starter-api .
docker run -p 8080:8080 go-starter-api
```

## 🤝 Integration with Flutter Web

```dart
// Flutter Web
final response = await http.post(
  Uri.parse('http://localhost:8080/api/generate'),
  headers: {'Content-Type': 'application/json'},
  body: jsonEncode({
    'name': 'my-api',
    'modulePath': 'github.com/user/my-api',
    'structure': 'simple',
    'database': {'type': 'postgres'},
    'libraries': ['go-auth', 'go-logger'],
    'deployment': 'railway',
  }),
);

// Download ZIP
final bytes = response.bodyBytes;
final blob = html.Blob([bytes]);
final url = html.Url.createObjectUrlFromBlob(blob);
html.AnchorElement(href: url)
  ..setAttribute('download', 'my-api.zip')
  ..click();
```

## 📝 License

MIT
