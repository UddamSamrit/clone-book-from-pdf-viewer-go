package services

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"clone-book/internal/config"
)

// DownloaderService handles all image downloading operations
type DownloaderService struct {
	bookInfo *BookInfo
}

// NewDownloaderService creates a new downloader service
func NewDownloaderService() *DownloaderService {
	return &DownloaderService{}
}

// SetBookInfo sets the book information for downloading
func (d *DownloaderService) SetBookInfo(bookInfo *BookInfo) {
	d.bookInfo = bookInfo
}

// DownloadAll downloads only missing images concurrently
// Skips images that already exist
func (d *DownloaderService) DownloadAll() error {
	if d.bookInfo == nil {
		return fmt.Errorf("book info not set")
	}

	startPage := d.bookInfo.StartPage
	endPage := d.bookInfo.EndPage

	// Use book-specific image directory
	imageDir := d.bookInfo.GetImageDir()

	// Create directory if it doesn't exist
	if err := os.MkdirAll(imageDir, 0755); err != nil {
		return fmt.Errorf("failed to create image directory: %v", err)
	}

	// First, check which images need to be downloaded
	missingImages := []int{}
	for i := startPage; i <= endPage; i++ {
		filePath := filepath.Join(imageDir, fmt.Sprintf("%d.jpg", i))
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			missingImages = append(missingImages, i)
		}
	}

	if len(missingImages) == 0 {
		fmt.Println("All images already exist. Skipping download.")
		return nil
	}

	fmt.Printf("Found %d missing images. Downloading...\n", len(missingImages))
	fmt.Printf("Total pages: %d (from %d to %d)\n\n", endPage-startPage+1, startPage, endPage)

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, config.MaxWorkers)
	errors := make(chan error, len(missingImages))

	// Progress tracking
	var mu sync.Mutex
	downloadedCount := 0
	totalToDownload := len(missingImages)

	for _, pageNum := range missingImages {
		wg.Add(1)
		go func(pageNum int) {
			defer wg.Done()
			semaphore <- struct{}{}        // Acquire semaphore
			defer func() { <-semaphore }() // Release semaphore

			url := fmt.Sprintf("%s/%d.jpg", d.bookInfo.ImageURL, pageNum)
			filePath := filepath.Join(d.bookInfo.GetImageDir(), fmt.Sprintf("%d.jpg", pageNum))

			if err := d.DownloadSingle(url, filePath, d.bookInfo.Referer); err != nil {
				errors <- fmt.Errorf("failed to download page %d: %v", pageNum, err)
				return
			}

			// Update progress
			mu.Lock()
			downloadedCount++
			current := downloadedCount
			mu.Unlock()

			// Show progress every 10 images or at completion
			if current%10 == 0 || current == totalToDownload {
				fmt.Printf("Progress: %d/%d images downloaded (%.1f%%)\n", current, totalToDownload, float64(current)/float64(totalToDownload)*100)
			}
		}(pageNum)
	}

	wg.Wait()
	close(errors)

	// Check for errors
	errorCount := 0
	for err := range errors {
		log.Printf("Error: %v", err)
		errorCount++
	}

	if errorCount > 0 {
		fmt.Printf("\n⚠️  Warning: %d images failed to download\n", errorCount)
	} else {
		fmt.Printf("\n✅ Download complete! All %d images downloaded successfully.\n", totalToDownload)
	}

	return nil
}

// DownloadSingle downloads a single image from URL to file path
func (d *DownloaderService) DownloadSingle(url, filePath, referer string) error {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	// Set headers to mimic browser request
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36")
	if referer != "" {
		req.Header.Set("Referer", referer)
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	return err
}
