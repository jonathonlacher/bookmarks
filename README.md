# Bookmarks

A simple bookmark manager in Go. Just a single-page app that reads from a YAML file.

## Quick Start

```bash
go build
./bookmarks
```

Open `http://localhost:8080`

## Setup

Edit `bookmarks.yaml`:

```yaml
groups:
  - name: Dev Tools
    bookmarks:
      - name: GitHub
        url: https://github.com
        description: Code stuff
        shortcode: gh
```

## What it does

- Displays bookmarks in groups
- Search box to filter them
- URL shortcodes: `http://localhost:8080/s/gh` → redirects to GitHub
- That's it

## Docker

```bash
docker build -t bookmarks .
docker run -p 8080:8080 bookmarks
```

Mount your own bookmarks file:
```bash
docker run -p 8080:8080 -v $(pwd)/my-bookmarks.yaml:/app/bookmarks.yaml bookmarks
```

## Flags

```
-config string    bookmarks file (default "bookmarks.yaml")
-port string      port (default "8080")
```

## How it works

- `main.go` - HTTP server, YAML parser
- `static/` - HTML/CSS/JS (embedded in binary at build time)
- `bookmarks.yaml` - Your bookmarks

Binary is standalone (~9MB). Just needs the YAML file to run.

## Notes

Only dependency is `gopkg.in/yaml.v3`. Everything else is standard library.

The code is intentionally simple. It's meant to be easy to read and modify.
