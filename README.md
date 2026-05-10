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

### Pre-compiled Binaries
Download from [releases](https://github.com/aevrHQ/imgconvert/releases)

### Via go install
```bash
go install github.com/aevrHQ/imgconvert/cmd/imgconvert@latest
```

### Build from Source
```bash
git clone https://github.com/aevrHQ/imgconvert.git
cd imgconvert
go build -o imgconvert ./cmd/imgconvert
```

## Documentation

See [SPEC.md](SPEC.md) for complete technical specification and design decisions.

## License

MIT License - See [LICENSE](LICENSE) for details

## Status

🚧 **Under Active Development** - v0.1.0 coming soon!
