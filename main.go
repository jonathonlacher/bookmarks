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

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, htmlPage)
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

// htmlPage is the embedded HTML content
const htmlPage = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Bookmarks</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
            line-height: 1.6;
            color: #333;
            background: #f5f5f5;
            padding: 20px;
        }

        .container {
            max-width: 1200px;
            margin: 0 auto;
            background: white;
            padding: 30px;
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }

        h1 {
            font-size: 28px;
            margin-bottom: 10px;
            color: #2c3e50;
        }

        .subtitle {
            color: #7f8c8d;
            margin-bottom: 25px;
            font-size: 14px;
        }

        .search-box {
            width: 100%;
            padding: 12px 16px;
            font-size: 16px;
            border: 2px solid #e0e0e0;
            border-radius: 6px;
            margin-bottom: 30px;
            transition: border-color 0.2s;
        }

        .search-box:focus {
            outline: none;
            border-color: #3498db;
        }

        .group {
            margin-bottom: 35px;
        }

        .group-name {
            font-size: 20px;
            font-weight: 600;
            color: #2c3e50;
            margin-bottom: 15px;
            padding-bottom: 8px;
            border-bottom: 2px solid #3498db;
        }

        .bookmarks {
            display: grid;
            grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
            gap: 15px;
        }

        .bookmark {
            padding: 15px;
            border: 1px solid #e0e0e0;
            border-radius: 6px;
            transition: all 0.2s;
            background: #fafafa;
        }

        .bookmark:hover {
            border-color: #3498db;
            box-shadow: 0 2px 8px rgba(52, 152, 219, 0.15);
            transform: translateY(-2px);
        }

        .bookmark-name {
            font-size: 16px;
            font-weight: 600;
            margin-bottom: 6px;
        }

        .bookmark-name a {
            color: #3498db;
            text-decoration: none;
        }

        .bookmark-name a:hover {
            text-decoration: underline;
        }

        .bookmark-url {
            font-size: 12px;
            color: #7f8c8d;
            margin-bottom: 8px;
            word-break: break-all;
        }

        .bookmark-description {
            font-size: 14px;
            color: #555;
            margin-bottom: 8px;
        }

        .bookmark-shortcode {
            display: inline-block;
            font-size: 11px;
            background: #ecf0f1;
            color: #2c3e50;
            padding: 3px 8px;
            border-radius: 3px;
            font-family: 'Courier New', monospace;
        }

        .no-results {
            text-align: center;
            padding: 40px;
            color: #7f8c8d;
            font-size: 16px;
        }

        .stats {
            font-size: 13px;
            color: #95a5a6;
            margin-bottom: 20px;
        }

        .hidden {
            display: none;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>📚 Bookmarks</h1>
        <p class="subtitle">Your centralized bookmark collection</p>

        <input
            type="text"
            id="searchInput"
            class="search-box"
            placeholder="Search bookmarks... (name, description, or URL)"
            autocomplete="off"
        >

        <div class="stats" id="stats"></div>

        <div id="bookmarksContainer"></div>
        <div id="noResults" class="no-results hidden">No bookmarks found</div>
    </div>

    <script>
        let allBookmarks = [];

        // Fetch bookmarks on page load
        async function loadBookmarks() {
            try {
                const response = await fetch('/api/bookmarks');
                const data = await response.json();
                allBookmarks = data.groups || [];
                displayBookmarks(allBookmarks);
                updateStats();
            } catch (error) {
                console.error('Failed to load bookmarks:', error);
            }
        }

        // Display bookmarks
        function displayBookmarks(groups) {
            const container = document.getElementById('bookmarksContainer');
            const noResults = document.getElementById('noResults');

            container.innerHTML = '';

            if (!groups || groups.length === 0) {
                noResults.classList.remove('hidden');
                return;
            }

            let hasVisibleBookmarks = false;

            groups.forEach(group => {
                if (!group.bookmarks || group.bookmarks.length === 0) return;

                const groupDiv = document.createElement('div');
                groupDiv.className = 'group';

                const groupTitle = document.createElement('div');
                groupTitle.className = 'group-name';
                groupTitle.textContent = group.name;
                groupDiv.appendChild(groupTitle);

                const bookmarksDiv = document.createElement('div');
                bookmarksDiv.className = 'bookmarks';

                group.bookmarks.forEach(bookmark => {
                    hasVisibleBookmarks = true;
                    const bookmarkDiv = createBookmarkElement(bookmark);
                    bookmarksDiv.appendChild(bookmarkDiv);
                });

                groupDiv.appendChild(bookmarksDiv);
                container.appendChild(groupDiv);
            });

            noResults.classList.toggle('hidden', hasVisibleBookmarks);
        }

        // Create bookmark HTML element
        function createBookmarkElement(bookmark) {
            const div = document.createElement('div');
            div.className = 'bookmark';

            const name = document.createElement('div');
            name.className = 'bookmark-name';
            const link = document.createElement('a');
            link.href = bookmark.url;
            link.textContent = bookmark.name;
            link.target = '_blank';
            name.appendChild(link);
            div.appendChild(name);

            const url = document.createElement('div');
            url.className = 'bookmark-url';
            url.textContent = bookmark.url;
            div.appendChild(url);

            if (bookmark.description) {
                const desc = document.createElement('div');
                desc.className = 'bookmark-description';
                desc.textContent = bookmark.description;
                div.appendChild(desc);
            }

            if (bookmark.shortcode) {
                const shortcode = document.createElement('div');
                shortcode.className = 'bookmark-shortcode';
                shortcode.textContent = '/s/' + bookmark.shortcode;
                div.appendChild(shortcode);
            }

            return div;
        }

        // Simple fuzzy search implementation
        function fuzzyMatch(text, query) {
            if (!query) return true;

            text = text.toLowerCase();
            query = query.toLowerCase();

            // Simple substring match for performance
            return text.includes(query);
        }

        // Filter bookmarks based on search query
        function filterBookmarks(query) {
            if (!query.trim()) {
                displayBookmarks(allBookmarks);
                updateStats();
                return;
            }

            const filtered = allBookmarks.map(group => {
                const matchingBookmarks = group.bookmarks.filter(bookmark => {
                    const searchText = [
                        bookmark.name,
                        bookmark.description || '',
                        bookmark.url,
                        bookmark.shortcode || ''
                    ].join(' ');

                    return fuzzyMatch(searchText, query);
                });

                return matchingBookmarks.length > 0
                    ? { ...group, bookmarks: matchingBookmarks }
                    : null;
            }).filter(Boolean);

            displayBookmarks(filtered);
            updateStats(query, filtered);
        }

        // Update statistics display
        function updateStats(query, filtered) {
            const statsDiv = document.getElementById('stats');
            const data = filtered || allBookmarks;

            let totalBookmarks = 0;
            data.forEach(group => {
                totalBookmarks += (group.bookmarks || []).length;
            });

            if (query) {
                statsDiv.textContent = 'Showing ' + totalBookmarks + ' bookmark(s) matching "' + query + '"';
            } else {
                statsDiv.textContent = totalBookmarks + ' bookmark(s) across ' + data.length + ' group(s)';
            }
        }

        // Setup search input handler
        document.getElementById('searchInput').addEventListener('input', (e) => {
            filterBookmarks(e.target.value);
        });

        // Load bookmarks on page load
        loadBookmarks();
    </script>
</body>
</html>
`
