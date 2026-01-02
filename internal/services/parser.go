package services

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// ParserService handles URL parsing and auto-detection
type ParserService struct{}

// NewParserService creates a new parser service
func NewParserService() *ParserService {
	return &ParserService{}
}

// BookInfo contains parsed book information
type BookInfo struct {
	BaseURL   string
	ImageURL  string
	StartPage int
	EndPage   int
	Referer   string
	BookName  string
}

// ParseURL parses a book URL and auto-detects the structure
func (p *ParserService) ParseURL(bookURL string) (*BookInfo, error) {
	parsedURL, err := url.Parse(bookURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %v", err)
	}

	// Remove hash fragment (#p=1, etc.) - it's just for navigation
	parsedURL.Fragment = ""

	// Extract base URL (without hash fragment)
	baseURL := fmt.Sprintf("%s://%s%s", parsedURL.Scheme, parsedURL.Host, parsedURL.Path)

	// Clean up the path - remove trailing slash
	baseURL = strings.TrimSuffix(baseURL, "/")

	// Extract book name from URL path (last segment)
	pathParts := strings.Split(strings.Trim(parsedURL.Path, "/"), "/")
	bookName := "book"
	if len(pathParts) > 0 {
		bookName = pathParts[len(pathParts)-1]
	}

	// Try to detect image URL pattern
	// Common patterns:
	// 1. /ebooks/.../files/mobile/{page}.jpg
	// 2. /ebooks/.../files/{page}.jpg
	// 3. /ebooks/.../pages/{page}.jpg
	// 4. /ebooks/.../images/{page}.jpg
	// 5. /ebooks/.../mobile/{page}.jpg

	imageURL := p.detectImageURLPattern(baseURL)
	if imageURL == "" {
		// Default pattern - try /files/mobile first
		imageURL = baseURL + "/files/mobile"
		// Verify it exists, if not try /files
		if !p.checkURLExists(imageURL + "/1.jpg") {
			imageURL = baseURL + "/files"
			if !p.checkURLExists(imageURL + "/1.jpg") {
				// Last resort - try direct pattern
				imageURL = baseURL
			}
		}
	}

	// Auto-detect page range
	startPage, endPage, err := p.detectPageRange(imageURL)
	if err != nil {
		return nil, fmt.Errorf("failed to detect page range: %v", err)
	}

	// Generate referer from base URL
	referer := baseURL
	if !strings.HasSuffix(referer, "/") {
		referer += "/"
	}

	return &BookInfo{
		BaseURL:   baseURL,
		ImageURL:  imageURL,
		StartPage: startPage,
		EndPage:   endPage,
		Referer:   referer,
		BookName:  bookName,
	}, nil
}

// GetImageDir returns the image directory path for this book
func (b *BookInfo) GetImageDir() string {
	return "images/" + b.BookName
}

// detectImageURLPattern tries to detect the image URL pattern
func (p *ParserService) detectImageURLPattern(baseURL string) string {
	// Try common patterns in order of likelihood
	patterns := []string{
		"/files/mobile",
		"/files",
		"/pages",
		"/images",
		"/mobile",
		"/page",
	}

	base := strings.TrimSuffix(baseURL, "/")
	for _, pattern := range patterns {
		testURL := base + pattern + "/1.jpg"
		if p.checkURLExists(testURL) {
			fmt.Printf("Detected image pattern: %s\n", pattern)
			return base + pattern
		}
	}

	return ""
}

// detectPageRange auto-detects the page range by trying pages
func (p *ParserService) detectPageRange(imageURL string) (int, int, error) {
	fmt.Println("Auto-detecting page range...")

	// Start from page 1
	startPage := 1

	// Find the last page by binary search
	// First, find an upper bound
	maxPage := 1
	for p.checkURLExists(fmt.Sprintf("%s/%d.jpg", imageURL, maxPage)) {
		maxPage *= 2
		if maxPage > 10000 {
			return 0, 0, fmt.Errorf("page range too large or detection failed")
		}
	}

	// Binary search for the last page
	low := 1
	high := maxPage
	lastPage := 1

	for low <= high {
		mid := (low + high) / 2
		if p.checkURLExists(fmt.Sprintf("%s/%d.jpg", imageURL, mid)) {
			lastPage = mid
			low = mid + 1
		} else {
			high = mid - 1
		}
	}

	if lastPage < 1 {
		return 0, 0, fmt.Errorf("no pages found")
	}

	fmt.Printf("Detected pages: %d to %d\n", startPage, lastPage)
	return startPage, lastPage, nil
}

// checkURLExists checks if a URL exists (returns 200 OK)
func (p *ParserService) checkURLExists(url string) bool {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return false
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}
