# Clone Book

A Go application to download book images and create a PDF from them.

## Prerequisites

- Go 1.21 or higher
- Internet connection

## Installation

### 1. Clone the repository

```bash
git clone <repository-url>
cd clone-book
```

### 2. Install Go (if not already installed)

#### macOS
```bash
brew install go
```

#### Linux
```bash
sudo apt-get update
sudo apt-get install golang-go
```

#### Windows
Download and install from [https://golang.org/dl/](https://golang.org/dl/)

### 3. Verify Go installation

```bash
go version
```

## Usage

### Run the application

```bash
go run main.go
```

Or build and run:

```bash
go build -o clone-book
./clone-book
```

## What it does

1. **Downloads images**: Downloads missing images from the configured URL (pages 1-337)
2. **Creates PDF**: Generates a PDF book from all downloaded images

## Configuration

Edit `internal/config/config.go` to change:
- `BaseURL`: The base URL for downloading images
- `StartPage`: Starting page number
- `EndPage`: Ending page number
- `DownloadDir`: Directory to store downloaded images
- `OutputPDF`: Output PDF filename
- `MaxWorkers`: Number of concurrent downloads

## Project Structure

```
clone-book/
├── main.go                 # Main entry point
├── go.mod                  # Go module file
├── internal/
│   ├── config/
│   │   └── config.go      # Configuration constants
│   └── services/
│       ├── downloader.go  # Image downloading service
│       └── pdf.go         # PDF creation service
└── README.md              # This file
```

## Output

- Downloaded images are stored in the `images/` directory
- Generated PDF is saved as `book.pdf` in the project root

## Notes

- The application skips downloading images that already exist
- Images are downloaded concurrently for faster processing
- The PDF is created from original downloaded images
