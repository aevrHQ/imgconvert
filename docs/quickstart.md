# Quick Start Guide

## Installation

```bash
cd /Users/miracleio/Documents/devprojects/image-converter
export PATH="$HOME/gosdk/bin:$PATH"
go build -o imgconvert ./cmd/imgconvert
```

Or use directly with `go run`:
```bash
go run cmd/imgconvert/main.go [flags] [files...]
```

## Common Usage Examples

### Basic Conversion
```bash
# Convert single file
./imgconvert photo.jpg

# Convert multiple files
./imgconvert photo1.jpg photo2.png photo3.gif

# Convert with glob pattern
./imgconvert *.jpg

# Recursive directory conversion
./imgconvert -r /path/to/photos
```

### Quality Control
```bash
# High quality (90)
./imgconvert -q 90 photo.jpg

# Low quality for maximum compression (60)
./imgconvert -q 60 photo.jpg

# Use preset
./imgconvert --preset high photo.jpg
./imgconvert --preset low photo.jpg

# Lossless compression
./imgconvert --lossless photo.png
```

### Resizing
```bash
# Scale to 50%
./imgconvert --scale 50 photo.jpg

# Thumbnail (33.33%)
./imgconvert --preset thumb photo.jpg

# Max width 800px (maintains aspect ratio)
./imgconvert --width 800 photo.jpg

# Max dimension 1920px
./imgconvert --max-dimension 1920 photo.jpg
```

### Output Control
```bash
# Output to specific directory
./imgconvert -o ./webp-output *.jpg

# Add suffix to filenames
./imgconvert --suffix _compressed *.jpg

# Overwrite existing files
./imgconvert --overwrite *.jpg

# Preview before converting (dry-run)
./imgconvert -n -r .
```

### Advanced
```bash
# Control parallel jobs
./imgconvert -j 4 -r /large/directory

# Quiet mode (no progress bar)
./imgconvert --quiet *.jpg

# Verbose output
./imgconvert --verbose *.jpg

# Combine multiple options
./imgconvert -r -q 90 --scale 80 --suffix _web -o ./output /photos
```

## Presets Reference

- `--preset low` - Quality 60, maximum compression
- `--preset medium` - Quality 80, balanced
- `--preset high` - Quality 90, best quality
- `--preset lossless` - No quality loss
- `--preset thumb` - Scale to 33.33% (thumbnail)
- `--preset half` - Scale to 50%
- `--preset web` - Scale to 80% (web optimized)

## Tips

1. **Dry run first**: Use `-n` to preview conversions
2. **Keep originals**: Don't use `--delete-original` until you're confident
3. **Use presets**: Faster than remembering quality numbers
4. **Parallel processing**: Adjust `-j` based on your CPU cores
5. **Batch processing**: Use recursive mode for directories

## Current Version: v0.1.0

**Features**:
- ✅ JPEG, PNG, GIF, BMP → WebP conversion
- ✅ Parallel processing
- ✅ Quality presets and manual control
- ✅ Image resizing
- ✅ Progress bars
- ✅ Dry-run mode

**Coming in v0.2.0**:
- Animated GIF support
- Additional converter plugins (libvips for speed)
- EXIF metadata preservation
- Configuration file support

See SPEC.md for complete technical specification.
