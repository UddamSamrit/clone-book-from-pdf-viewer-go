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
type DownloaderService struct{}

// NewDownloaderService creates a new downloader service
func NewDownloaderService() *DownloaderService {
	return &DownloaderService{}
}

// DownloadAll downloads only missing images concurrently
// Skips images that already exist
func (d *DownloaderService) DownloadAll() error {
	// First, check which images need to be downloaded
	missingImages := []int{}
	for i := config.StartPage; i <= config.EndPage; i++ {
		filePath := filepath.Join(config.DownloadDir, fmt.Sprintf("%d.jpg", i))
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			missingImages = append(missingImages, i)
		}
	}

	if len(missingImages) == 0 {
		fmt.Println("All images already exist. Skipping download.")
		return nil
	}

	fmt.Printf("Found %d missing images. Downloading...\n", len(missingImages))

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, config.MaxWorkers)
	errors := make(chan error, len(missingImages))

	for _, pageNum := range missingImages {
		wg.Add(1)
		go func(pageNum int) {
			defer wg.Done()
			semaphore <- struct{}{}        // Acquire semaphore
			defer func() { <-semaphore }() // Release semaphore

			url := fmt.Sprintf("%s/%d.jpg", config.BaseURL, pageNum)
			filePath := filepath.Join(config.DownloadDir, fmt.Sprintf("%d.jpg", pageNum))

			if err := d.DownloadSingle(url, filePath); err != nil {
				errors <- fmt.Errorf("failed to download page %d: %v", pageNum, err)
				return
			}

			if pageNum%10 == 0 {
				fmt.Printf("Downloaded %d/%d images...\n", pageNum, config.EndPage)
			}
		}(pageNum)
	}

	wg.Wait()
	close(errors)

	// Check for errors
	hasErrors := false
	for err := range errors {
		log.Printf("Error: %v", err)
		hasErrors = true
	}

	if hasErrors {
		return fmt.Errorf("some images failed to download")
	}

	return nil
}

// DownloadSingle downloads a single image from URL to file path
func (d *DownloaderService) DownloadSingle(url, filePath string) error {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	// Set headers to mimic browser request
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36")
	req.Header.Set("Referer", "https://www.elibraryofcambodia.org/ebooks/2019/06/sk-23-06-19-agneung/")

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
