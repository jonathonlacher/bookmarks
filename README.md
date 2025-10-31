# Bookmarks

A simple, minimal, single-page bookmark manager written in Go. Perfect for internal company use with easy Docker deployment.

## Features

- **Simple YAML Configuration**: Define bookmarks in a clear YAML structure
- **Organized Groups**: Categorize bookmarks into logical groups
- **URL Shortcodes**: Create short URLs for quick access (e.g., `/s/gh` → GitHub)
- **Live Search**: Filter bookmarks with fuzzy matching
- **Minimal Dependencies**: Uses only `gopkg.in/yaml.v3` outside the standard library
- **Clean UI**: Attractive, responsive single-page interface
- **Docker Ready**: Easy containerized deployment
- **Lightweight**: Small binary, minimal resource usage

## Project Structure

```
.
├── main.go              # Go application (API, routing, YAML parsing)
├── bookmarks.yaml       # Bookmark configuration
├── static/
│   └── index.html      # Frontend (HTML/CSS/JS in a single file)
├── Dockerfile           # Docker build configuration
├── go.mod              # Go dependencies
└── README.md           # Documentation
```

## Quick Start

### Running Locally

```bash
# Build and run
go build -o bookmarks
./bookmarks

# Custom config and port
./bookmarks -config my-bookmarks.yaml -port 3000
```

Visit `http://localhost:8080` in your browser.

### Docker Deployment

```bash
# Build the Docker image
docker build -t bookmarks .

# Run with default bookmarks.yaml
docker run -p 8080:8080 bookmarks

# Run with custom bookmarks file
docker run -p 8080:8080 -v $(pwd)/my-bookmarks.yaml:/app/bookmarks.yaml bookmarks

# Run with custom port
docker run -p 3000:3000 bookmarks ./bookmarks -port 3000
```

## Configuration

Create a `bookmarks.yaml` file with your bookmarks:

```yaml
groups:
  - name: Development Tools
    bookmarks:
      - name: GitHub
        url: https://github.com
        description: Code repository and collaboration platform
        shortcode: gh

      - name: Stack Overflow
        url: https://stackoverflow.com
        description: Programming Q&A community
        shortcode: so

  - name: Internal Resources
    bookmarks:
      - name: Company Wiki
        url: https://wiki.company.internal
        description: Internal documentation
        shortcode: wiki
```

### YAML Structure

- **groups**: Array of bookmark groups
  - **name**: Group display name
  - **bookmarks**: Array of bookmarks in this group
    - **name**: Bookmark display name (required)
    - **url**: Target URL (required)
    - **description**: Brief description (optional)
    - **shortcode**: Short URL code (optional)

## Usage

### Main Page

Access the main bookmark page at `http://localhost:8080/`

Features:
- All bookmarks displayed in organized groups
- Search bar for filtering bookmarks
- Click any bookmark to visit the URL
- Hover effects for better interaction

### URL Shortcodes

If a bookmark has a shortcode, you can access it via:
```
http://localhost:8080/s/<shortcode>
```

Example: `http://localhost:8080/s/gh` redirects to GitHub

### API Endpoint

Get all bookmarks as JSON:
```
http://localhost:8080/api/bookmarks
```

## Command Line Options

```
-config string
    Path to bookmarks YAML file (default "bookmarks.yaml")
-static string
    Path to static files directory (default "static")
-port string
    Port to serve on (default "8080")
```

## Building from Source

```bash
# Clone the repository
git clone <repository-url>
cd bookmarks

# Download dependencies
go mod download

# Build
go build -o bookmarks

# Run
./bookmarks
```

## Security Considerations

- **Minimal Dependencies**: Only uses `gopkg.in/yaml.v3` for YAML parsing
- **No External Resources**: All HTML/CSS/JS is embedded in the binary
- **Read-Only**: Application only reads the config file, no write operations
- **Standard Library**: Maximizes use of Go standard library for security
- **No Database**: No database dependencies or connections
- **Stateless**: No session management or authentication complexity

## Docker Image Details

The Docker image uses multi-stage builds for minimal size:
- **Build Stage**: Uses `golang:1.21-alpine` to compile the application
- **Final Stage**: Uses `alpine:latest` with only the compiled binary
- **Size**: Final image is approximately 10-15 MB
- **Security**: Includes CA certificates for HTTPS support

## Customization

### Changing the UI

The HTML/CSS/JS is located in `static/index.html`. Edit this file to customize the appearance and functionality.

### Adding Features

The code is intentionally simple and easy to modify. Some ideas:
- Add categories/tags
- Export bookmarks
- Import from browser bookmarks
- Add bookmark icons/favicons
- User authentication

## Troubleshooting

**Port already in use:**
```bash
./bookmarks -port 3000
```

**Config file not found:**
```bash
./bookmarks -config /path/to/bookmarks.yaml
```

**Docker build fails:**
Ensure you have `go.mod` and `go.sum` files:
```bash
go mod tidy
docker build -t bookmarks .
```

## License

This is internal company software. Adjust licensing as needed.

## Contributing

Keep it simple! The goal is clarity and maintainability over cleverness.
