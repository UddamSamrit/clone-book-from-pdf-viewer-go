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

## Git Workflow

**⚠️ Important: Do not push directly to `main` branch! (Repository owners can bypass this)**

### Working with branches

1. **Create a feature branch:**
```bash
git checkout -b feature/your-feature-name
```

2. **Make your changes and commit:**
```bash
git add .
git commit -m "Your commit message"
```

3. **Push to your feature branch:**
```bash
git push origin feature/your-feature-name
```

4. **Create a Pull Request** on GitHub to merge into `main`

### Prevent accidental pushes to main

#### For Contributors (Non-Owners)

**Option 1: Git config (recommended)**
```bash
git config branch.main.pushRemote no_push
```

**Option 2: Use pre-push hook**
```bash
# Install the pre-push hook
git config core.hooksPath .githooks
chmod +x .githooks/pre-push
```

#### For Repository Owner

**To allow owner to push directly to main:**

1. **Local Git Hook (if using pre-push hook):**
   ```bash
   git config user.isOwner true
   ```

2. **GitHub Actions:** Automatically allows repository owner to push (no configuration needed)

3. **GitHub Branch Protection Settings:**
   - Go to Settings → Branches → Add rule for `main`
   - Enable: "Require a pull request before merging"
   - Enable: "Require approvals" (set to 1)
   - **Important:** Leave "Do not allow bypassing" **UNCHECKED** to allow owner bypass
   - Or add yourself to "Restrict who can push to matching branches" as an exception

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
