# imgconvert - High-Performance Image Converter

Convert images to WebP format with optimal compression and blazing-fast performance.

## Features

- 🚀 **Bulk Conversion** - Convert multiple files/directories in parallel
- 📁 **Multiple Input Formats** - JPEG, PNG, GIF (animated), BMP
- 🎯 **Smart Quality Control** - Presets and manual quality settings
- ⚡ **High Performance** - Parallel processing using all CPU cores
- 🔄 **Image Resizing** - Scale images with aspect ratio preservation
- 🎨 **Transparency Support** - Preserves alpha channel from PNG/GIF
- 📊 **Progress Tracking** - Real-time progress bars and statistics
- 🔍 **Dry Run Mode** - Preview conversions before executing
- 🔌 **Plugin Architecture** - Multiple conversion engines for optimal performance

## Quick Start

```bash
# Convert all JPEGs in current directory
imgconvert *.jpg

# Recursive conversion with custom quality
imgconvert -r -q 90 .

# Fast mode for large batches
imgconvert --fast -r /path/to/photos

# Resize and convert
imgconvert --scale 50 *.jpg

# Dry run to preview
imgconvert -n -r .
```

## Installation

### Prerequisites
- Go 1.20 or later (installer can set this up for you)

### macOS

#### Quick Install (Recommended)
```bash
git clone https://github.com/aevrHQ/imgconvert.git
cd imgconvert
./install.sh
```

The installer will:
1. Set up Go if needed (installs to `~/gosdk/bin`)
2. Install the tool to `~/go/bin/imgconvert` (wrapper script)
3. Add `~/go/bin` to your PATH

**⚠️ Important: First-Time Setup**
After installation, reload your shell configuration:
```bash
source ~/.zshrc  # if using zsh (default on macOS)
source ~/.bashrc # if using bash
```

You only need to do this once in each open terminal. New terminals will work automatically.

**Why a wrapper script?**  
macOS Gatekeeper blocks unsigned binaries. The installer creates a wrapper that uses `go run` to avoid code signing issues.

#### Manual macOS Installation
If the installer doesn't work:
```bash
git clone https://github.com/aevrHQ/imgconvert.git
cd imgconvert

# Create wrapper script
cat > ~/go/bin/imgconvert << 'EOF'
#!/bin/bash
CURRENT_DIR="$(pwd)"
export PATH="$HOME/gosdk/bin:$PATH"
ARGS=()
for arg in "$@"; do
    if [[ -e "$arg" ]] || [[ -e "$CURRENT_DIR/$arg" ]]; then
        if [[ "$arg" = /* ]]; then
            ARGS+=("$arg")
        else
            ARGS+=("$CURRENT_DIR/$arg")
        fi
    else
        ARGS+=("$arg")
    fi
done
cd /Users/$(whoami)/Documents/devprojects/image-converter
exec go run cmd/imgconvert/main.go "${ARGS[@]}"
EOF

chmod +x ~/go/bin/imgconvert

# Add to PATH
echo 'export PATH="$HOME/go/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc
```

### Linux

#### Via go install (Recommended)
```bash
go install github.com/aevrHQ/imgconvert/cmd/imgconvert@latest
```

This installs a compiled binary to `~/go/bin/imgconvert`.

**Add to PATH (if not already):**
```bash
# For bash
echo 'export PATH="$HOME/go/bin:$PATH"' >> ~/.bashrc
source ~/.bashrc

# For zsh
echo 'export PATH="$HOME/go/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc
```

#### Build from Source
```bash
git clone https://github.com/aevrHQ/imgconvert.git
cd imgconvert
go build -o imgconvert ./cmd/imgconvert

# Install globally
sudo mv imgconvert /usr/local/bin/
```

#### Using package script
```bash
git clone https://github.com/aevrHQ/imgconvert.git
cd imgconvert
./build.sh
sudo mv imgconvert /usr/local/bin/
```

### Windows

#### Via go install
```powershell
go install github.com/aevrHQ/imgconvert/cmd/imgconvert@latest
```

This installs to `%USERPROFILE%\go\bin\imgconvert.exe`.

**Add to PATH:**
1. Search for "Environment Variables" in Windows
2. Edit "Path" for your user
3. Add: `%USERPROFILE%\go\bin`
4. Restart your terminal

#### Build from Source
```powershell
git clone https://github.com/aevrHQ/imgconvert.git
cd imgconvert
go build -o imgconvert.exe .\cmd\imgconvert
```

Then move `imgconvert.exe` to a directory in your PATH, or add the current directory to PATH.

#### Using PowerShell
```powershell
# Convert images
.\imgconvert.exe *.jpg

# Or if in PATH
imgconvert *.jpg
```

#### Using Command Prompt
```cmd
imgconvert *.jpg
```

### Verify Installation

All platforms:
```bash
imgconvert --version
imgconvert --help
```

## Getting Started

### Basic Usage

**Convert a single file:**
```bash
imgconvert photo.jpg
# Creates: photo.webp
```

**Convert multiple files:**
```bash
imgconvert *.jpg *.png
# Processes all JPEG and PNG files
```

**Convert from any directory:**
```bash
cd ~/Pictures/vacation
imgconvert *.jpg
# Works from anywhere after installation
```

### Common Patterns

**Bulk convert with quality preset:**
```bash
imgconvert --preset high *.jpg
```

**Recursive directory conversion:**
```bash
imgconvert -r ~/Pictures/
```

**Resize while converting:**
```bash
imgconvert --scale 50 *.jpg    # 50% size
imgconvert --width 1920 *.jpg  # Max 1920px wide
```

**Preview before converting:**
```bash
imgconvert -n -r .
# Shows what will be converted without doing it
```

**Custom output directory:**
```bash
imgconvert -o ./webp-output *.jpg
```

**Delete originals after conversion:**
```bash
imgconvert --delete-original *.jpg
```

### Quality Presets

| Preset | Quality | Use Case |
|--------|---------|----------|
| `lossless` | 100 | Perfect quality, larger files |
| `high` | 90 | Photography, high-quality images |
| `medium` | 85 | **Default** - Great balance |
| `web` | 80 | Web images, good quality |
| `low` | 60 | Previews, thumbnails |
| `thumb` | 50 | Small thumbnails |
| `half` | 85 + 50% scale | Half size with good quality |

### Troubleshooting

#### macOS Issues

**"killed" or "zsh: killed" errors**  
This happens with compiled binaries on macOS due to code signing. Use the installer script which creates a wrapper that avoids this issue.

**"command not found: imgconvert"**  
Solution: Reload your shell configuration:
```bash
source ~/.zshrc  # or source ~/.bashrc
```

**Cannot access file errors**  
If you see `cannot access file.jpg: no such file or directory`:
1. Make sure you've sourced your shell config (`source ~/.zshrc`)
2. Verify the files exist: `ls *.jpg`
3. Try using absolute paths: `imgconvert /full/path/to/file.jpg`

**Permission denied**  
Make sure the wrapper script is executable:
```bash
chmod +x ~/go/bin/imgconvert
```

#### Linux Issues

**"command not found: imgconvert"**  
Make sure `~/go/bin` is in your PATH:
```bash
export PATH="$HOME/go/bin:$PATH"
# Add to ~/.bashrc or ~/.zshrc to make permanent
```

**Permission denied**  
Make the binary executable:
```bash
chmod +x ~/go/bin/imgconvert
# Or if installed globally:
sudo chmod +x /usr/local/bin/imgconvert
```

#### Windows Issues

**"imgconvert is not recognized"**  
Add Go's bin directory to your PATH:
1. Open "Edit environment variables for your account"
2. Edit "Path" variable
3. Add: `%USERPROFILE%\go\bin`
4. Restart terminal/PowerShell

**Wildcard not working**  
PowerShell and Command Prompt handle wildcards differently. Use quotes:
```powershell
imgconvert "*.jpg"
```

Or use `Get-ChildItem`:
```powershell
Get-ChildItem *.jpg | ForEach-Object { imgconvert $_.FullName }
```

#### General Issues

**"go: command not found"**  
Install Go from [golang.org/dl](https://golang.org/dl/)

**Out of memory errors**  
Reduce parallel jobs:
```bash
imgconvert -j 2 *.jpg  # Use only 2 workers
```

**Slow performance**  
Ensure you're not limiting CPU:
```bash
imgconvert -j 0 *.jpg  # Use all available cores (default)
```

### Performance Tips

**Use parallel processing (default):**
- Tool automatically uses all CPU cores
- Adjust with `-j` flag: `imgconvert -j 4 *.jpg`

**Fast mode for large batches:**
```bash
imgconvert --fast -r /large/photo/library
# Uses quality 75 with optimized settings
```

**Process in batches:**
```bash
# For thousands of files, process subdirectories separately
for dir in */; do imgconvert -r "$dir"; done
```

### Expected Results

Typical file size reductions:
- **High quality (q=90):** 60-70% smaller
- **Default (q=85):** 70-80% smaller  
- **Thumbnail preset:** 90%+ smaller
- **50% scaling:** 85%+ smaller

Example:
```
Input:  photo.jpg (23.6 MB)
Output: photo.webp (7.0 MB)
Reduction: 70.3%
```

## Documentation

See [SPEC.md](SPEC.md) for complete technical specification and design decisions.

## License

MIT License - See [LICENSE](LICENSE) for details

## Status

🚧 **Under Active Development** - v0.1.0 coming soon!
