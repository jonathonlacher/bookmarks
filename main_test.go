package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// TestLoadConfig tests YAML configuration loading
func TestLoadConfig(t *testing.T) {
	// Create a temporary test config
	testYAML := `
groups:
  - name: Test Group
    bookmarks:
      - name: Test Bookmark
        url: https://example.com
        description: Test description
        shortcode: test
`
	tmpfile, err := os.CreateTemp("", "bookmarks-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(testYAML)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Test loading the config
	if err := loadConfig(tmpfile.Name()); err != nil {
		t.Fatalf("Failed to load valid config: %v", err)
	}

	// Verify the config was loaded correctly
	if len(config.Groups) != 1 {
		t.Errorf("Expected 1 group, got %d", len(config.Groups))
	}

	if config.Groups[0].Name != "Test Group" {
		t.Errorf("Expected group name 'Test Group', got '%s'", config.Groups[0].Name)
	}

	if len(config.Groups[0].Bookmarks) != 1 {
		t.Errorf("Expected 1 bookmark, got %d", len(config.Groups[0].Bookmarks))
	}

	// Verify shortcode map was built
	if url, ok := shortcodeMap["test"]; !ok || url != "https://example.com" {
		t.Errorf("Shortcode map not built correctly")
	}
}

// TestLoadConfigInvalidFile tests error handling for missing files
func TestLoadConfigInvalidFile(t *testing.T) {
	err := loadConfig("/nonexistent/file.yaml")
	if err == nil {
		t.Error("Expected error for nonexistent file, got nil")
	}
}

// TestHandleHome tests the home page handler
func TestHandleHome(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	handleHome(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "text/html") {
		t.Errorf("Expected HTML content type, got %s", contentType)
	}
}

// TestHandleHomeNotFound tests 404 for non-root paths
func TestHandleHomeNotFound(t *testing.T) {
	req := httptest.NewRequest("GET", "/nonexistent", nil)
	w := httptest.NewRecorder()

	handleHome(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", resp.StatusCode)
	}
}

// TestHandleAPIBookmarks tests the API endpoint
func TestHandleAPIBookmarks(t *testing.T) {
	// Set up test config
	config = Config{
		Groups: []Group{
			{
				Name: "Test",
				Bookmarks: []Bookmark{
					{
						Name:        "Example",
						URL:         "https://example.com",
						Description: "Test site",
						Shortcode:   "ex",
					},
				},
			},
		},
	}

	req := httptest.NewRequest("GET", "/api/bookmarks", nil)
	w := httptest.NewRecorder()

	handleAPIBookmarks(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		t.Errorf("Expected JSON content type, got %s", contentType)
	}

	// Verify JSON is valid
	var result Config
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Errorf("Failed to decode JSON response: %v", err)
	}

	if len(result.Groups) != 1 {
		t.Errorf("Expected 1 group in response, got %d", len(result.Groups))
	}
}

// TestHandleShortcode tests shortcode redirects
func TestHandleShortcode(t *testing.T) {
	// Set up test shortcode map
	shortcodeMap = map[string]string{
		"test": "https://example.com",
	}

	req := httptest.NewRequest("GET", "/s/test", nil)
	w := httptest.NewRecorder()

	handleShortcode(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusFound {
		t.Errorf("Expected status 302, got %d", resp.StatusCode)
	}

	location := resp.Header.Get("Location")
	if location != "https://example.com" {
		t.Errorf("Expected redirect to https://example.com, got %s", location)
	}
}

// TestHandleShortcodeNotFound tests 404 for invalid shortcodes
func TestHandleShortcodeNotFound(t *testing.T) {
	shortcodeMap = map[string]string{
		"test": "https://example.com",
	}

	req := httptest.NewRequest("GET", "/s/nonexistent", nil)
	w := httptest.NewRecorder()

	handleShortcode(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", resp.StatusCode)
	}
}

// TestHandleStatic tests serving static files
func TestHandleStatic(t *testing.T) {
	req := httptest.NewRequest("GET", "/static/style.css", nil)
	w := httptest.NewRecorder()

	handleStatic(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType != "text/css" {
		t.Errorf("Expected CSS content type, got %s", contentType)
	}
}

// TestHandleStaticNotFound tests 404 for missing static files
func TestHandleStaticNotFound(t *testing.T) {
	req := httptest.NewRequest("GET", "/static/nonexistent.css", nil)
	w := httptest.NewRecorder()

	handleStatic(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", resp.StatusCode)
	}
}

// TestCountBookmarks tests the bookmark counter
func TestCountBookmarks(t *testing.T) {
	config = Config{
		Groups: []Group{
			{Bookmarks: []Bookmark{{}, {}}},
			{Bookmarks: []Bookmark{{}, {}, {}}},
		},
	}

	count := countBookmarks()
	if count != 5 {
		t.Errorf("Expected 5 bookmarks, got %d", count)
	}
}
