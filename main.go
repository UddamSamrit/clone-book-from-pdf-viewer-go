package main

import (
	"fmt"
	"log"
	"os"

	"clone-book/internal/config"
	"clone-book/internal/services"
)

func main() {
	fmt.Println("Starting book cloning process...")
	fmt.Printf("Processing images from %d.jpg to %d.jpg\n", config.StartPage, config.EndPage)

	// Create download directory
	if err := os.MkdirAll(config.DownloadDir, 0755); err != nil {
		log.Fatalf("Failed to create download directory: %v", err)
	}

	// Initialize services
	downloader := services.NewDownloaderService()
	pdfService := services.NewPDFService()

	// Download only missing images (skips if already exist)
	fmt.Println("\nChecking for missing images...")
	if err := downloader.DownloadAll(); err != nil {
		log.Fatalf("Failed to download images: %v", err)
	}

	fmt.Println("\nCreating PDF from original images...")
	if err := pdfService.CreateFromImages(); err != nil {
		log.Fatalf("Failed to create PDF: %v", err)
	}

	fmt.Printf("\n✓ Successfully created %s\n", config.OutputPDF)
}
