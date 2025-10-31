package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// Bookmark represents a single bookmark entry
type Bookmark struct {
	Name        string `yaml:"name" json:"name"`
	URL         string `yaml:"url" json:"url"`
	Description string `yaml:"description" json:"description"`
	Shortcode   string `yaml:"shortcode" json:"shortcode"`
}

// Group represents a collection of related bookmarks
type Group struct {
	Name      string     `yaml:"name" json:"name"`
	Bookmarks []Bookmark `yaml:"bookmarks" json:"bookmarks"`
}

// Config represents the bookmark configuration
type Config struct {
	Groups []Group `yaml:"groups" json:"groups"`
}

var (
	config        Config
	shortcodeMap  map[string]string
	configFile    = flag.String("config", "bookmarks.yaml", "Path to bookmarks YAML file")
	staticDir     = flag.String("static", "static", "Path to static files directory")
	port          = flag.String("port", "8080", "Port to serve on")
)

func main() {
	flag.Parse()

	// Load configuration
	if err := loadConfig(*configFile); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	log.Printf("Loaded %d groups with %d total bookmarks", len(config.Groups), countBookmarks())

	// Setup routes
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/api/bookmarks", handleAPIBookmarks)
	http.HandleFunc("/s/", handleShortcode)

	// Start server
	addr := ":" + *port
	log.Printf("Starting server on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

// loadConfig reads and parses the YAML configuration file
func loadConfig(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("reading config file: %w", err)
	}

	if err := yaml.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("parsing YAML: %w", err)
	}

	// Build shortcode map
	shortcodeMap = make(map[string]string)
	for _, group := range config.Groups {
		for _, bookmark := range group.Bookmarks {
			if bookmark.Shortcode != "" {
				shortcodeMap[bookmark.Shortcode] = bookmark.URL
			}
		}
	}

	return nil
}

// countBookmarks returns total number of bookmarks across all groups
func countBookmarks() int {
	count := 0
	for _, group := range config.Groups {
		count += len(group.Bookmarks)
	}
	return count
}

// handleHome serves the main HTML page
func handleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	indexPath := *staticDir + "/index.html"
	data, err := os.ReadFile(indexPath)
	if err != nil {
		http.Error(w, "Failed to load page", http.StatusInternalServerError)
		log.Printf("Error reading index.html: %v", err)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(data)
}

// handleAPIBookmarks serves bookmark data as JSON
func handleAPIBookmarks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(config); err != nil {
		http.Error(w, "Failed to encode bookmarks", http.StatusInternalServerError)
		log.Printf("Error encoding bookmarks: %v", err)
	}
}

// handleShortcode redirects shortcodes to their target URLs
func handleShortcode(w http.ResponseWriter, r *http.Request) {
	code := strings.TrimPrefix(r.URL.Path, "/s/")

	if url, ok := shortcodeMap[code]; ok {
		http.Redirect(w, r, url, http.StatusFound)
		return
	}

	http.NotFound(w, r)
}
