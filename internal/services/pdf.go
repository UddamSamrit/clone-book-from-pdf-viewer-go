package services

import (
	"fmt"
	"image"
	"log"
	"os"
	"path/filepath"

	"github.com/jung-kurt/gofpdf/v2"
)

// PDFService handles PDF creation operations
type PDFService struct {
	bookInfo *BookInfo
}

// NewPDFService creates a new PDF service
func NewPDFService() *PDFService {
	return &PDFService{}
}

// SetBookInfo sets the book information for PDF creation
func (pdf *PDFService) SetBookInfo(bookInfo *BookInfo) {
	pdf.bookInfo = bookInfo
}

// CreateFromImages creates a PDF from all processed images
func (pdf *PDFService) CreateFromImages() error {
	if pdf.bookInfo == nil {
		return fmt.Errorf("book info not set")
	}

	doc := gofpdf.New("P", "mm", "A4", "")
	doc.SetAutoPageBreak(false, 0)

	startPage := pdf.bookInfo.StartPage
	endPage := pdf.bookInfo.EndPage
	totalPages := endPage - startPage + 1

	addedCount := 0
	skippedCount := 0

	for i := startPage; i <= endPage; i++ {
		imagePath, err := pdf.FindImagePath(i)
		if err != nil {
			fmt.Printf("Warning: Image %d not found, skipping\n", i)
			skippedCount++
			continue
		}

		width, height, err := pdf.GetImageDimensions(imagePath)
		if err != nil {
			log.Printf("Warning: Failed to get dimensions for %s: %v", imagePath, err)
			skippedCount++
			continue
		}

		x, y, displayWidth, displayHeight := pdf.CalculatePageLayout(width, height)

		doc.AddPage()
		doc.ImageOptions(imagePath, x, y, displayWidth, displayHeight, false, gofpdf.ImageOptions{ImageType: "JPG", ReadDpi: true}, 0, "")
		addedCount++

		// Show progress every 50 pages or at completion
		if addedCount%50 == 0 || i == endPage {
			fmt.Printf("PDF Progress: %d/%d pages added (%.1f%%)\n", addedCount, totalPages, float64(addedCount)/float64(totalPages)*100)
		}
	}

	if skippedCount > 0 {
		fmt.Printf("\n⚠️  Warning: %d pages were skipped\n", skippedCount)
	}

	// Use book name for PDF filename
	pdfFilename := pdf.bookInfo.BookName + ".pdf"

	err := doc.OutputFileAndClose(pdfFilename)
	if err != nil {
		return err
	}

	fmt.Printf("\n✅ PDF creation complete! %d pages added to PDF.\n", addedCount)
	return nil
}

// FindImagePath finds the path to an original image
func (pdf *PDFService) FindImagePath(pageNum int) (string, error) {
	originalPath := filepath.Join(pdf.bookInfo.GetImageDir(), fmt.Sprintf("%d.jpg", pageNum))

	// Use original image only
	if _, err := os.Stat(originalPath); err == nil {
		return originalPath, nil
	}

	return "", fmt.Errorf("image not found for page %d", pageNum)
}

// GetImageDimensions gets the dimensions of an image
func (pdf *PDFService) GetImageDimensions(imagePath string) (float64, float64, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return 0, 0, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return 0, 0, err
	}

	bounds := img.Bounds()
	return float64(bounds.Dx()), float64(bounds.Dy()), nil
}

// CalculatePageLayout calculates the layout for an image on an A4 page
func (pdf *PDFService) CalculatePageLayout(imgWidth, imgHeight float64) (float64, float64, float64, float64) {
	pageWidth := 210.0  // A4 width in mm
	pageHeight := 297.0 // A4 height in mm

	widthRatio := pageWidth / imgWidth
	heightRatio := pageHeight / imgHeight
	ratio := widthRatio
	if heightRatio < widthRatio {
		ratio = heightRatio
	}

	displayWidth := imgWidth * ratio
	displayHeight := imgHeight * ratio

	// Center image on page
	x := (pageWidth - displayWidth) / 2
	y := (pageHeight - displayHeight) / 2

	return x, y, displayWidth, displayHeight
}
