package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"clone-book/internal/services"
)

func main() {
	// Get URL from command line argument or prompt
	var bookURL string
	if len(os.Args) > 1 {
		bookURL = os.Args[1]
	} else {
		bookURL = promptURL()
	}

	if bookURL == "" {
		log.Fatal("No URL provided")
	}

	fmt.Println("Starting book cloning process...")
	fmt.Printf("Parsing URL: %s\n", bookURL)

	// Initialize parser
	parser := services.NewParserService()
	bookInfo, err := parser.ParseURL(bookURL)
	if err != nil {
		log.Fatalf("Failed to parse URL: %v", err)
	}

	fmt.Printf("\nBook Information:")
	fmt.Printf("\n  Book Name: %s", bookInfo.BookName)
	fmt.Printf("\n  Base URL: %s", bookInfo.BaseURL)
	fmt.Printf("\n  Image URL: %s", bookInfo.ImageURL)
	fmt.Printf("\n  Pages: %d to %d", bookInfo.StartPage, bookInfo.EndPage)
	fmt.Println()

	// Image directory will be created per book in downloader

	// Initialize services
	downloader := services.NewDownloaderService()
	downloader.SetBookInfo(bookInfo)

	pdfService := services.NewPDFService()
	pdfService.SetBookInfo(bookInfo)

	// Download only missing images (skips if already exist)
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("STEP 1: Downloading Images")
	fmt.Println(strings.Repeat("=", 60))
	if err := downloader.DownloadAll(); err != nil {
		log.Fatalf("Failed to download images: %v", err)
	}

	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("STEP 2: Creating PDF")
	fmt.Println(strings.Repeat("=", 60))
	if err := pdfService.CreateFromImages(); err != nil {
		log.Fatalf("Failed to create PDF: %v", err)
	}

	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("✅ PROCESS COMPLETE!")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf("📚 Book Name: %s\n", bookInfo.BookName)
	fmt.Printf("📄 PDF file created: %s.pdf\n", bookInfo.BookName)
	fmt.Printf("📁 Images stored in: %s/\n", bookInfo.GetImageDir())
	fmt.Printf("📊 Total pages: %d\n", bookInfo.EndPage-bookInfo.StartPage+1)
	fmt.Println(strings.Repeat("=", 60))
}

// promptURL prompts the user to enter a URL
func promptURL() string {
	fmt.Print("Enter book URL: ")
	reader := bufio.NewReader(os.Stdin)
	url, err := reader.ReadString('\n')
	if err != nil {
		return ""
	}
	return strings.TrimSpace(url)
}
