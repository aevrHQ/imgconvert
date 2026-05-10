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

### Quick Install (Recommended)
```bash
git clone https://github.com/aevrHQ/imgconvert.git
cd imgconvert
./install.sh
```

The installer will:
1. Set up Go if needed (installs to `~/gosdk/bin`)
2. Install the tool to `~/go/bin/imgconvert` (wrapper script)
3. Add `~/go/bin` to your PATH in `~/.zshrc` or `~/.bashrc`

### ⚠️ Important: First-Time Setup
After installation, **reload your shell configuration**:
```bash
source ~/.zshrc  # or source ~/.bashrc
```

You only need to do this once in each open terminal. New terminals will work automatically.

### Verify Installation
```bash
imgconvert --version
```

### Alternative Methods

#### Via go install
```bash
go install github.com/aevrHQ/imgconvert/cmd/imgconvert@latest
```

#### Build from Source
```bash
git clone https://github.com/aevrHQ/imgconvert.git
cd imgconvert
go build -o imgconvert ./cmd/imgconvert
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

#### "command not found: imgconvert"
**Solution:** Reload your shell configuration:
```bash
source ~/.zshrc  # or source ~/.bashrc
```

#### "killed" or "zsh: killed" errors (macOS)
This happens with compiled binaries on macOS due to code signing. The installer uses a **wrapper script** that avoids this issue by running the tool via `go run` instead of a compiled binary.

**How it works:**
- The wrapper at `~/go/bin/imgconvert` automatically converts relative paths to absolute paths
- Then runs the tool from the project directory
- You don't need to think about it - just use `imgconvert` normally!

#### Cannot access file errors
If you see errors like `cannot access file.jpg: no such file or directory`:
1. Make sure you've sourced your shell config (`source ~/.zshrc`)
2. Verify the files exist: `ls *.jpg`
3. Try using absolute paths: `imgconvert /full/path/to/file.jpg`

#### Permission denied
Make sure the wrapper script is executable:
```bash
chmod +x ~/go/bin/imgconvert
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
