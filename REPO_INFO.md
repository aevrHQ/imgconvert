# Repository Information

## GitHub Repository
**Official URL**: https://github.com/aevrHQ/imgconvert

## Go Module Path
```
github.com/aevrHQ/imgconvert
```

## Installation from GitHub

### Via go install
```bash
go install github.com/aevrHQ/imgconvert/cmd/imgconvert@latest
```

### Clone and Build
```bash
git clone https://github.com/aevrHQ/imgconvert.git
cd imgconvert
./build.sh
```

### Download Release Binary
Visit: https://github.com/aevrHQ/imgconvert/releases

## Import in Go Code

```go
import (
    "github.com/aevrHQ/imgconvert/internal/converter"
    "github.com/aevrHQ/imgconvert/internal/processor"
    "github.com/aevrHQ/imgconvert/internal/scanner"
)
```

## Publishing to GitHub

1. Create repository: https://github.com/aevrHQ/imgconvert

2. Push code:
```bash
cd /Users/miracleio/Documents/devprojects/image-converter
git init
git add .
git commit -m "Initial commit: imgconvert v0.1.0"
git branch -M main
git remote add origin https://github.com/aevrHQ/imgconvert.git
git push -u origin main
```

3. Create first release:
   - Go to https://github.com/aevrHQ/imgconvert/releases
   - Click "Create a new release"
   - Tag: v0.1.0
   - Title: "imgconvert v0.1.0 - Initial Release"
   - Upload binaries (see build instructions below)

## Building Release Binaries

```bash
cd /Users/miracleio/Documents/devprojects/image-converter
export PATH="$HOME/gosdk/bin:$PATH"

# macOS Apple Silicon (M1/M2/M3)
go build -o imgconvert-v0.1.0-darwin-arm64 -ldflags "-s -w" ./cmd/imgconvert

# macOS Intel
GOOS=darwin GOARCH=amd64 go build -o imgconvert-v0.1.0-darwin-amd64 -ldflags "-s -w" ./cmd/imgconvert

# Linux
GOOS=linux GOARCH=amd64 go build -o imgconvert-v0.1.0-linux-amd64 -ldflags "-s -w" ./cmd/imgconvert

# Windows
GOOS=windows GOARCH=amd64 go build -o imgconvert-v0.1.0-windows-amd64.exe -ldflags "-s -w" ./cmd/imgconvert
```

## Users Can Install With

Once published on GitHub:

```bash
# Via go install
go install github.com/aevrHQ/imgconvert/cmd/imgconvert@latest

# Or download binary from releases
# https://github.com/aevrHQ/imgconvert/releases
```
