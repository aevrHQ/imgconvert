# Image Converter - Detailed Specification Document

**Project Goal**: Build a simple, high-performance, portable image converter for bulk conversion to WebP format with reduced file sizes.

---

## DECISIONS LOG

### 1. Core Requirements
- **Bulk Operation**: Primary use case (covers single operations)
- **Source Formats**: JPG, PNG, etc.
- **Target Format**: WebP
- **Performance**: High performance required
- **Portability**: Must be portable
- **File Size**: Must reduce file sizes

---

## QUESTIONS & DECISIONS

### Q1: Implementation Language
**Decision**: Go (Golang)
**Rationale**: Maximum portability, performance, and simple distribution
**Trade-offs**: 
- ✅ Compiles to single binary with no runtime dependencies
- ✅ Excellent performance and concurrency
- ✅ Cross-platform compilation built-in
- ✅ Small executable size (~5-15MB)
- ⚠️ May require learning Go if unfamiliar

### Q2: Distribution & Installation
**Decision**: go install + GitHub Releases + install script
**Binary Name**: `imgconvert` (future-proof for other formats)
**Rationale**: Start simple, MVP for personal use first, can expand later
**Details**:
- Users can do: `go install github.com/USER/imgconvert@latest`
- GitHub Releases with pre-compiled binaries (Windows, macOS Intel/ARM, Linux)
- Simple install script to detect OS and download correct binary
- Defer package managers (Homebrew, apt, etc.) until later

### Q3: Image Processing Library
**Decision**: Use existing high-performance libraries (no custom implementation)
**Rationale**: Leverage battle-tested, optimized libraries
**Options to evaluate**:
- libwebp (official WebP library with Go bindings)
- bimg (uses libvips - extremely fast)
- imaging/resize libraries for Go

### Q4: Input File Discovery
**Decision**: Support multiple input methods (file globs, directories, individual files)
**Recursive Scanning**: Opt-in with `-r` or `--recursive` flag (safer default)
**Symlinks**: Do NOT follow symlinks (prevent infinite loops/security issues)
**Hidden Files**: Exclude hidden files/folders by default
**Examples**:
- `imgconvert *.jpg` - Current directory only, glob pattern
- `imgconvert -r .` - Recursive from current directory
- `imgconvert -d /path/to/images` - Specific directory, non-recursive
- `imgconvert -r -d /path/to/images` - Recursive directory scan
- `imgconvert file1.jpg file2.png` - Specific files
**Rationale**: Maximum flexibility while maintaining safety

### Q5: Source Format Support
**Decision**: Tier 1 formats with enhanced features
**Supported Input Formats**:
- JPEG/JPG (most common)
- PNG (with transparency)
- GIF (including animated)
- BMP (Windows standard)
**Special Features**:
- **Animated GIFs**: Convert to animated WebP (preserve all frames)
- **EXIF Metadata**: Preserve by default, add `--strip-metadata` flag to remove
- **Transparency**: Always preserve alpha channel (PNG, GIF → WebP with alpha)
**Rationale**: Cover 95%+ of common use cases without over-complicating

### Q6: Output Strategy
**Decision**: Flexible output with safe defaults
**Default Behavior**: Same directory with `.webp` extension
- Example: `photo.jpg` → `photo.webp`
**Optional Flags**:
- `-o <dir>` or `--output <dir>`: Output to specific directory
- `--preserve-structure`: When using `-o`, mirror source directory structure
- `--suffix <string>`: Add suffix before extension (e.g., `--suffix _compressed` → `photo_compressed.webp`)
- `--overwrite`: Force overwrite existing output files
- `--delete-original`: Delete source files after successful conversion (requires confirmation or `--force`)
**Existing File Handling**:
- Default: Skip if output file exists (safe), warn user
- Require `--overwrite` flag to replace existing files
**Rationale**: Intuitive defaults with power-user flexibility, prioritize safety

### Q7: Quality & Compression Settings
**Decision**: Smart defaults with comprehensive control options
**Default Quality**: 85 (excellent balance between size and visual quality)
**Quality Control Flags**:
- `-q <0-100>` or `--quality <0-100>`: User-specified quality level
- `--lossless`: Use lossless WebP compression (perfect quality, larger files)
**Convenience Presets**:
- `--preset low` (quality 60) - Maximum compression, smallest files
- `--preset medium` (quality 80) - Good balance
- `--preset high` (quality 90) - Minimal compression, best quality
- `--preset lossless` - No quality loss
**Advanced Options**:
- `--effort <0-6>`: Encoding effort (0=fast, 6=slowest/best) - Default: 4
- `--target-size <KB>`: Attempt to achieve specific file size (best effort)
**Override Behavior**: Manual quality settings override presets (last flag wins)
**Statistics**: Show size reduction after conversion
- Example: "Converted 10 files: 45.2MB → 12.3MB (72.8% reduction)"
**Rationale**: Easy for beginners (good defaults), powerful for advanced users

### Q8: Concurrency & Performance
**Decision**: Intelligent parallel processing with user control
**Default Concurrency**: Use all available CPU cores (`runtime.NumCPU()`)
**Concurrency Control**:
- `-j <N>` or `--jobs <N>`: Limit concurrent conversions to N
- `-j 1`: Force sequential processing (debugging, low memory)
**Memory Management**:
- `--low-memory`: Process one file at a time, release immediately
- Useful for very large images or constrained systems
**Progress Indication**:
- Default: Show progress bar for bulk operations
  - Format: `Converting: [=====>    ] 45/100 (45%) | Current: photo_045.jpg`
- `-q` or `--quiet`: Suppress progress output (only show errors/summary)
- `-v` or `--verbose`: Show detailed info for each file
**Error Handling**:
- Continue processing remaining files if one fails
- Collect and report all errors at the end
- `--fail-fast`: Stop on first error (optional flag)
**Rationale**: Maximum performance by default, flexible control for edge cases

### Q9: Image Resizing & Optimization
**Decision**: Optional resizing with aspect ratio preservation (NO cropping in v1)
**Resizing Flags**:
- `--width <px>`: Set maximum width (maintain aspect ratio)
- `--height <px>`: Set maximum height (maintain aspect ratio)
- `--max-dimension <px>`: Set maximum for longest side (maintain aspect ratio)
- `--scale <percent>`: Scale by percentage (e.g., `--scale 50` = 50% of original)
- `--fit <WxH>`: Fit within dimensions maintaining aspect ratio (e.g., `--fit 1920x1080`)
**Percentage Presets** (common use cases):
- `--preset thumb` or `--scale 33.33` - Thumbnail size (1/3 original)
- `--preset half` or `--scale 50` - Half size
- `--preset web` or `--scale 80` - Web optimized (80% of original)
**Scaling Behavior**:
- Only downscale by default (never upscale - prevents quality loss)
- Add `--allow-upscale` flag to permit enlarging
- **Always preserve aspect ratio** (no distortion)
- **No cropping in v1** (future feature)
**Resampling Algorithm**:
- Default: Lanczos (best quality for downscaling)
- `--filter <type>`: Choose filter (lanczos, bilinear, nearest) for advanced users
**Rationale**: Major value-add for web optimization while keeping implementation simple

### Q10: Dry Run & Preview Mode
**Decision**: Yes - Add comprehensive dry run mode for safety and planning
**Dry Run Flag**:
- `--dry-run` or `-n`: Show what would be converted without actually converting
**Output Format**:
```
[DRY RUN] Would convert 47 files:

photo_001.jpg (2.4 MB) → photo_001.webp (~0.7 MB estimated)
photo_002.png (1.8 MB) → photo_002.webp (~0.5 MB estimated)
...

Estimated total: 112.3 MB → ~32.1 MB (71.4% reduction)

Run without --dry-run to perform conversion
```
**Additional Preview Features**:
- `--list`: Just list discovered files (no size estimates, faster)
- Show warnings for:
  - Files that would be skipped (output already exists)
  - Unsupported formats found
  - Permission issues
**Behavior**:
- Include estimated size reductions (based on typical compression ratios)
- Verbose by default in dry run mode (users want details)
**Rationale**: Safety first - let users verify before bulk operations

### Q11: Error Handling & Validation
**Decision**: Robust error handling with clear, actionable messages
**Error Categories**:
1. **Input Validation**: Invalid quality values, bad paths, conflicting flags, no files found
2. **File Operations**: Permission denied, disk full, file locked
3. **Conversion Errors**: Corrupted images, decode/encode failures, out of memory
**Error Reporting**:
- Clear, actionable messages with file names and specific issues
- Suggest fixes when possible
- Example: "Error: Cannot read 'photo.jpg': Permission denied. Try running with appropriate permissions."
**Exit Codes**:
- `0`: Success (all files converted)
- `1`: Partial success (some files failed)
- `2`: Total failure (no files converted)
- `3`: Invalid arguments/usage
**Error Handling Strategy**:
- Continue processing after individual failures
- Collect all errors and show summary at end:
```
Conversion complete:
✓ 95 files converted successfully
✗ 5 files failed:
  - photo_032.jpg: Corrupted image data
  - photo_054.png: Permission denied
  ...
```
**Validation**:
- Validate all inputs before starting conversion (fail fast on argument errors)
- Check available disk space and warn if estimated output exceeds capacity
- `--log <file>`: Write detailed error logs (useful for large batches)
**Detail Level**:
- Detailed error messages by default
- Concise in `--quiet` mode (only show counts)
**Rationale**: Prevent data loss, provide clear feedback, enable debugging

### Q12: CLI Interface Design
**Decision**: Intuitive, discoverable CLI with no subcommands (keep v1 simple)
**Command Structure Examples**:
```bash
imgconvert *.jpg                    # Convert all JPGs in current dir
imgconvert -r .                     # Convert all images recursively
imgconvert -q 90 photo.jpg          # Convert with specific quality
```
**Complete Flag Specification**:
```
Input/Output:
  -d, --dir <path>           Directory to scan
  -r, --recursive            Recursive directory scanning
  -o, --output <dir>         Output directory
      --preserve-structure   Mirror source directory structure
      --suffix <string>      Add suffix to output filenames

Quality/Compression:
  -q, --quality <0-100>      Quality level (default: 85)
      --preset <name>        Quality preset (low/medium/high/lossless/thumb/half/web)
      --lossless             Use lossless compression
      --effort <0-6>         Encoding effort (default: 4)
      --strip-metadata       Remove EXIF/metadata

Resizing:
      --width <px>           Maximum width
      --height <px>          Maximum height
      --max-dimension <px>   Maximum for longest side
      --scale <percent>      Scale by percentage
      --fit <WxH>            Fit within dimensions
      --allow-upscale        Allow upscaling images

Concurrency:
  -j, --jobs <N>             Concurrent jobs (default: all CPUs)
      --low-memory           Process one file at a time

Output Control:
      --overwrite            Overwrite existing files
      --delete-original      Delete source after conversion
      --force                Skip confirmations

Modes:
  -n, --dry-run              Preview without converting
      --list                 List files only
      --fail-fast            Stop on first error

Verbosity:
      --quiet                Suppress progress output
      --verbose              Show detailed info
      --log <file>           Write errors to log file

Info:
  -h, --help                 Show help
      --version              Show version
```
**Design Principles**:
- Support both short and long forms for common flags
- Group related flags in help output
- Use consistent kebab-case for long flags
- No subcommands in v1 (single clear purpose)
**Rationale**: Balance simplicity and power, discoverable through --help

### Q13: Configuration File Support
**Decision**: Defer to v2 - Focus on perfecting CLI first
**Rationale**: Config files are nice-to-have, not essential for MVP. Gather user feedback on v1 before designing config system.
**Deferred to v2** (see V2 Features section below)

### Q14: Project Structure & Code Organization
**Decision**: Clean, maintainable architecture following Go best practices
**Structure**:
```
imgconvert/
├── cmd/
│   └── imgconvert/
│       └── main.go              # Entry point, CLI setup
├── internal/
│   ├── converter/
│   │   ├── converter.go         # Core conversion logic
│   │   ├── webp.go              # WebP-specific encoding
│   │   └── resize.go            # Image resizing logic
│   ├── scanner/
│   │   └── scanner.go           # File discovery & filtering
│   ├── processor/
│   │   └── processor.go         # Parallel processing coordinator
│   └── config/
│       └── config.go            # Configuration & validation
├── pkg/
│   └── stats/
│       └── stats.go             # Statistics tracking
├── test/
│   ├── integration/             # Integration tests
│   └── fixtures/                # Sample test images
├── docs/                         # Additional documentation
├── go.mod
├── go.sum
├── README.md
├── SPEC.md                       # This spec document
├── LICENSE
└── .goreleaser.yml              # For automated releases
```
**Architecture Principles**:
- `cmd/`: Application entry points
- `internal/`: Private application code (not importable by other projects)
- `pkg/`: Public reusable packages
- `test/`: Integration tests with sample images
- `docs/`: Additional documentation
- Clear separation of concerns, easy to test
**Rationale**: Follows Go community standards, maintainable, testable

### Q15: Testing Strategy
**Decision**: Comprehensive testing with pragmatic coverage targets
**Testing Levels**:
1. **Unit Tests**: Test individual functions in isolation, mock I/O, target 70-80% coverage
2. **Integration Tests**: Real image conversions with test fixtures
3. **End-to-End Tests**: CLI flags, exit codes, error scenarios
**Test Organization**:
```
internal/converter/converter_test.go    # Unit tests
internal/scanner/scanner_test.go        # Unit tests
test/integration/conversion_test.go     # Integration tests
test/fixtures/                          # Sample images
  ├── test.jpg
  ├── test.png
  ├── test.gif (animated)
  └── test.bmp
```
**Test Coverage**:
- Happy path conversions (all formats)
- Quality settings and presets
- Resizing with aspect ratio preservation
- Animated GIF → animated WebP
- Transparency preservation (PNG, GIF)
- EXIF metadata handling
- Error handling (corrupted files, permissions)
- Concurrent processing
- Dry run mode accuracy
- File overwrite behavior
**Benchmarks**: Add performance benchmarks for conversion operations
**CI/CD**:
- Test on: Linux, macOS, Windows
- Go version: 1.21+ (or latest stable)
- GitHub Actions for automation
**Targets**:
- Minimum coverage: 70% (focus on critical paths)
- Include benchmarks for performance tracking
- Defer fuzzing to v2
**Rationale**: Balance thorough testing with development velocity

### Q16: Plugin Architecture & Dependencies
**Decision**: Plugin-based converter architecture with format-specific implementations
**Architecture Design**:
```go
// Converter interface - all format converters implement this
type Converter interface {
    // Decode reads and decodes the source image
    Decode(src io.Reader) (image.Image, error)
    
    // Encode writes the image in target format
    Encode(dst io.Writer, img image.Image, opts *EncodeOptions) error
    
    // SupportsFormat checks if this converter handles the format
    SupportsFormat(format string) bool
    
    // Info returns converter metadata
    Info() ConverterInfo
}

// Registry manages available converters
type Registry struct {
    converters []Converter
}

// Select chooses best converter based on format and conditions
func (r *Registry) Select(format string, conditions Conditions) Converter
```

**Plugin Structure**:
```
internal/
├── converter/
│   ├── converter.go           # Interface definition
│   ├── registry.go            # Converter registry & selection logic
│   ├── options.go             # Shared encoding options
│   └── plugins/
│       ├── webp_kolesa/       # kolesa-team/go-webp implementation
│       │   └── webp.go
│       ├── webp_chai2010/     # chai2010/webp (pure Go fallback)
│       │   └── webp.go
│       ├── libvips/           # h2non/bimg (fastest, optional)
│       │   └── vips.go
│       └── stdlib/            # Standard library decoders
│           ├── jpeg.go
│           ├── png.go
│           ├── gif.go
│           └── bmp.go
```

**Selection Strategy**:
```
Conditions to consider:
- Input format (JPEG, PNG, GIF, BMP)
- Output format (WebP for v1)
- File size (large files may need different handling)
- Animation (GIF → animated WebP)
- Transparency (alpha channel)
- Available libraries (CGO vs pure Go)
- Performance requirements
```

**Proposed Dependencies**:
```go
require (
    // CLI & Core
    github.com/spf13/cobra v1.8.0                   // CLI framework
    
    // WebP Encoders (plugins)
    github.com/kolesa-team/go-webp v1.0.4          // Primary WebP (CGO, best quality)
    github.com/chai2010/webp v1.1.1                 // Fallback WebP (pure Go)
    
    // Optional: Ultra-fast processing
    github.com/h2non/bimg v1.1.9                    // libvips (optional, for speed)
    
    // Image Processing
    github.com/disintegration/imaging v1.6.2        // Resizing & manipulation
    golang.org/x/image v0.15.0                      // BMP + extended formats
    
    // UI
    github.com/schollz/progressbar/v3 v3.14.1       // Progress bars
    
    // Optional: Metadata
    github.com/rwcarlsen/goexif v0.0.0-20190401172101-9e8deecbddbd  // EXIF handling
)
```

**Plugin Selection Logic (v1)**:
1. **Default**: kolesa-team/go-webp (best quality)
2. **Fallback**: chai2010/webp if CGO unavailable
3. **Optional**: bimg if installed and file size > 10MB (performance)
4. **Future**: Allow user to specify via flag `--engine <name>`

**Rationale**: 
- Flexibility to swap implementations without changing core code
- Easy to add new formats (AVIF, JPEG XL) in future
- Can benchmark and choose best library per use case
- Graceful fallbacks if dependencies unavailable
- Testable in isolation

### Q17: Plugin Selection Criteria & Strategy
**Decision**: Smart automatic selection with user override options
**Automatic Selection Priority**:

**Priority 1 - File Format Requirements:**
- Animated GIF → Use converter that supports animated WebP
- PNG with transparency → Ensure alpha channel support
- Large files (>10MB) → Prefer faster converter (libvips if available)

**Priority 2 - Library Availability:**
- Check if CGO-based libraries are available
- Fall back to pure Go if CGO unavailable
- Detect libvips installation for bimg

**Priority 3 - Quality vs Speed:**
- Default: Prioritize quality (kolesa-team/go-webp)
- With `--fast` flag: Prefer speed (bimg/libvips) - **KEY FEATURE**
- With `--lossless`: Use converter with best lossless support

**Selection Flow**:
```
1. Check file format and features needed
2. If animated GIF → Select animated-capable converter
3. If --fast flag → Try bimg (libvips) first, fall back to fastest available
4. If --lossless → Select converter with best lossless support
5. Default → kolesa-team/go-webp (best quality)
6. Fallback → chai2010/webp (pure Go)
7. Error if no suitable converter found
```

**User Control Flags (v1)**:
- `--fast`: Prioritize speed over quality (use libvips/bimg) - **IMPORTANT**
- `--prefer-pure-go`: Avoid CGO dependencies
- `--engine <name>`: Force specific converter (kolesa, chai2010, vips)

**Verbose Mode**:
- In verbose mode, show which converter is used:
  - `[DEBUG] Using kolesa-team/go-webp for JPEG → WebP`
  - `[DEBUG] Using libvips (bimg) for fast processing`

**Rationale**: 
- Default to quality (most users want best output)
- `--fast` flag enables high-performance workflows (important for large batches)
- Transparent about which library is used (debugging, trust)
- Graceful fallbacks ensure tool works even with missing dependencies

### Q18: Release & Versioning Strategy
**Decision**: Automated multi-platform releases with SemVer
**Versioning**: Semantic Versioning (SemVer)
- Format: `vMAJOR.MINOR.PATCH` (e.g., v0.1.0, v1.0.0, v1.2.3)
- Major: Breaking changes
- Minor: New features, backward compatible
- Patch: Bug fixes
- Start with v0.1.0 (signals "in development"), v1.0.0 when stable and battle-tested

**Release Automation**: GoReleaser
- Triggered by Git tags: `git tag v0.1.0 && git push --tags`
- Automatically builds binaries for all platforms
- Generates checksums (SHA256)
- Creates GitHub Release with changelog

**Target Platforms**:
```yaml
# .goreleaser.yml builds for:
- imgconvert-linux-amd64
- imgconvert-linux-arm64
- imgconvert-darwin-amd64 (Intel Mac)
- imgconvert-darwin-arm64 (Apple Silicon)
- imgconvert-windows-amd64.exe
- imgconvert-windows-arm64.exe
```

**Distribution Channels**:
1. **GitHub Releases** (primary): Binaries + checksums + changelog
2. **go install**: `go install github.com/USER/imgconvert/cmd/imgconvert@latest`
3. **Install Script**: One-liner that detects OS/arch and installs correct binary
   ```bash
   curl -sSL https://raw.githubusercontent.com/USER/imgconvert/main/install.sh | bash
   ```

**Changelog**:
- Auto-generate from Git commits (conventional commits format recommended)
- Include in GitHub Release notes
- Highlight breaking changes, new features, bug fixes

**Release Cadence**:
- Feature-based (release when features complete and tested)
- Not time-based (avoid rushing incomplete features)

**Security**:
- Always include SHA256 checksums (GoReleaser automatic)
- Consider code signing for binaries (v2 feature)

**Rationale**: Industry-standard approach, automated, multi-platform, secure

### Q19: Documentation Strategy
**Decision**: Comprehensive, user-friendly documentation with multiple layers
**Documentation Components**:

**1. README.md** (Essential - Primary entry point)
- Project overview and features
- Quick install instructions (all methods)
- Quick start examples
- Basic benchmarks (build trust, show performance)
- Links to detailed documentation
- Badge indicators (build status, version, license)

**2. docs/ Directory** (Detailed guides)
- `installation.md` - Detailed installation, dependency setup, troubleshooting
- `usage.md` - Complete flag reference, detailed examples, workflows
- `performance.md` - Benchmarks, optimization tips, comparison with other tools
- `troubleshooting.md` - Common issues, solutions, FAQ
- `contributing.md` - Contribution guidelines, development setup, code standards
- `architecture.md` - System design, plugin architecture explanation

**3. Built-in Help** (CLI)
- `imgconvert --help` - Comprehensive CLI help with grouped flags
- Clear descriptions for each flag
- Examples in help text

**4. SPEC.md** (Technical specification)
- This document - Architecture decisions
- Rationale for technical choices
- Feature roadmap (v1 vs v2)
- Decision log from this session

**5. Code Documentation** (For developers)
- GoDoc comments for all public APIs
- Package-level documentation
- Inline comments for complex logic only
- Code examples in GoDoc

**6. Examples & Tutorials** (Practical guides)
- Sample images in `test/fixtures/` for testing
- Tutorial: "Converting your photo library"
- Tutorial: "Optimizing images for the web"
- Tutorial: "Batch processing large directories"

**FAQ Section**: Include in troubleshooting.md
- "Why is my conversion slow?" → Check --fast flag, CPU cores
- "Why are my WebP files larger?" → Quality settings, lossless mode
- "How do I preserve EXIF data?" → Default behavior explanation
- "Can I undo conversions?" → No, always keep backups or don't use --delete-original

**Deferred to v2**:
- GitHub Pages website (README sufficient initially)
- Video tutorials
- Interactive examples

**Rationale**: Progressive documentation - quick start for beginners, deep dive for power users

### Q20: Performance Benchmarking & Metrics
**Decision**: Track key performance indicators with comparative benchmarks
**Metrics to Track**:

**1. Conversion Speed:**
- Images per second (throughput)
- Megabytes per second
- Time per image (average, min, max)
- Separate metrics for different source formats

**2. Compression Effectiveness:**
- Size reduction percentage
- Output file size vs input file size
- Quality-to-size ratio

**3. Resource Usage:**
- CPU utilization (per core and total)
- Memory consumption (peak and average)
- Disk I/O patterns

**4. Quality Metrics (v1: file size only, v2: add PSNR/SSIM):**
- File size comparison (primary metric)
- PSNR/SSIM deferred to v2 (for quality-focused users)

**Benchmark Suite**:
```go
// Go benchmarks
BenchmarkConvertJPEG_Quality85
BenchmarkConvertPNG_WithAlpha
BenchmarkConvertGIF_Animated
BenchmarkConvertBulk_100Files
BenchmarkResize_2000x2000_To_800x800
BenchmarkParallel_vs_Sequential
BenchmarkKolesaConverter
BenchmarkChai2010Converter
BenchmarkLibvipsConverter

// Include memory profiling
```

**Comparison Documentation** (in docs/performance.md):
```markdown
## Performance Comparison

Test set: 100 JPEG images (average 3MB each)
Hardware: MacBook Pro M2, 8 cores

| Tool              | Time  | Total Size | Reduction | Quality |
|-------------------|-------|------------|-----------|---------|
| imgconvert        | 8.2s  | 42MB       | 86%       | 85      |
| imgconvert --fast | 3.1s  | 44MB       | 85%       | 85      |
| cwebp             | 12.4s | 41MB       | 86.3%     | 85      |
| ImageMagick       | 45.2s | 48MB       | 84%       | 85      |
```

**README Performance Summary**: Include highlights, link to detailed benchmarks

**Continuous Performance Monitoring**:
- GitHub Actions workflow for performance regression testing
- Compare against baseline on each PR
- Alert if performance degrades >10%
- Include memory profiling in benchmarks (prevent memory leaks)

**Deferred to v2**:
- Built-in benchmark mode: `imgconvert --benchmark photo.jpg`
- PSNR/SSIM quality metrics
- Real-time performance dashboard

**Rationale**: Data-driven optimization, build user trust, prevent regressions

### Q21: CI/CD Pipeline
**Decision**: Comprehensive GitHub Actions workflows for quality and automation
**CI Workflows**:

**1. Test & Build** (on every push/PR):
```yaml
# .github/workflows/test.yml
- Platforms: Linux, macOS, Windows
- Go versions: 1.21, 1.22 (latest stable)
- Steps: Checkout → Setup Go → Install deps (libwebp) → go vet → go test -race → benchmarks → coverage
- Upload coverage to Codecov
```

**2. Release** (on Git tag v*.*.*):
```yaml
# .github/workflows/release.yml
- Uses GoReleaser
- Builds all platform binaries
- Creates GitHub Release with auto-generated changelog (editable)
- Uploads binaries, checksums, assets
```

**3. Performance Regression** (on PR):
```yaml
# .github/workflows/benchmark.yml
- Run benchmarks, compare vs main
- Comment on PR if performance degrades >10% (informational, doesn't block)
```

**4. Linting & Code Quality** (on push/PR):
```yaml
# .github/workflows/lint.yml
- golangci-lint (comprehensive)
- gofmt check
- go vet
```

**Branch Protection**:
- Require all checks to pass before merge (except benchmark - informational)
- Full test suite on PRs, quick tests on pushes (save CI time)
- Require 1 approval for PRs (when team grows)

**Pre-commit Hooks** (optional local):
- Run tests, format code, catch issues early

**Rationale**: Automated quality gates, fast feedback, reliable releases

### Q22: License
**Decision**: MIT License
**Rationale**: 
- Permissive, widely adopted
- Allows commercial use
- Simple, well-understood
- Encourages contributions and adoption
- Compatible with most dependencies

### Q23: Security Considerations
**Decision**: Defense-in-depth for untrusted input
**Security Measures**:
- **Input Validation**: Check file magic bytes (not just extension)
- **Size Limits**: Default max file size (e.g., 100MB), configurable via flag
- **Memory Limits**: Prevent OOM with --low-memory flag
- **Dependency Scanning**: Dependabot enabled, regular updates
- **Sandboxing**: Consider resource limits (future)
- **Error Messages**: Don't leak system paths in production
**Vulnerability Management**:
- GitHub Security Advisories enabled
- Regular dependency updates
- Clear security policy (SECURITY.md)
**Rationale**: Handle untrusted images safely, prevent DoS attacks

### Q24: Development Workflow
**Decision**: Trunk-based development with conventional commits
**Git Strategy**:
- Main branch: `main` (always stable)
- Feature branches: `feature/description` or `fix/description`
- Merge via PRs with squash commits
**Commit Convention**: Conventional Commits
```
feat: add animated GIF support
fix: preserve EXIF orientation
docs: update installation guide
perf: optimize parallel processing
test: add integration tests for PNG
```
**Development Process**:
1. Create feature branch from main
2. Develop with tests
3. Open PR with description
4. Pass CI checks
5. Squash merge to main
**Rationale**: Clean history, automated changelogs, clear intent

### Q25: Community & Contribution
**Decision**: Welcoming, low-friction contribution process
**Templates**:
- Issue templates: Bug report, Feature request, Question
- PR template: Checklist (tests, docs, changelog)
- CONTRIBUTING.md: Setup guide, coding standards, PR process
**Code of Conduct**: Contributor Covenant
**Recognition**: CONTRIBUTORS.md listing all contributors
**Rationale**: Encourage community growth, make contributing easy

---

## IMPLEMENTATION SUMMARY

### V1 Scope (MVP - Build This First)
✅ **Core Features**:
- Bulk image conversion (JPEG, PNG, GIF, BMP → WebP)
- Multiple input methods (globs, directories, files)
- Parallel processing (all CPU cores)
- Quality control (presets and manual)
- Image resizing with aspect ratio preservation
- Animated GIF → animated WebP
- Transparency preservation
- EXIF metadata handling
- Dry-run mode
- Progress bars and statistics
- Plugin architecture for converters
- `--fast` flag for high-performance mode

✅ **Implementation Language**: Go
✅ **Distribution**: Binary releases, go install, install script
✅ **Documentation**: README, docs/, built-in help
✅ **Testing**: 70%+ coverage, CI/CD
✅ **License**: MIT

### V2 Features (Future Enhancements)
- Configuration file support (YAML)
- Profile system
- More output formats (AVIF, JPEG XL)
- More input formats (HEIC, TIFF, RAW)
- Image cropping
- Watermarking
- Advanced filters
- Built-in benchmark mode
- Quality metrics (PSNR/SSIM)
- Watch mode
- GUI option
- Cloud integration
- Plugin system for external converters

---

## NEXT STEPS

1. **Initialize Go Project**
   ```bash
   go mod init github.com/USER/imgconvert
   ```

2. **Create Project Structure**
   - Set up directories (cmd, internal, test, docs)
   
3. **Implement Core Components** (Priority Order):
   - Plugin architecture & converter interface
   - File scanner
   - WebP converter (kolesa + chai2010 fallback)
   - Parallel processor
   - CLI with cobra
   - Progress bars
   - Resize functionality
   - Dry-run mode
   
4. **Testing & Documentation**
   - Unit tests
   - Integration tests
   - README and docs
   
5. **Release v0.1.0**
   - GitHub Actions setup
   - GoReleaser config
   - First release

---

## V2 FEATURES (Future Enhancements)

### Configuration File Support
**Format**: YAML (human-readable, comments supported)
**Example Config**:
```yaml
# .imgconvert.yaml
quality: 85
effort: 4
recursive: true
overwrite: false
jobs: 8
preserve_structure: true

# Profiles for different use cases
profiles:
  web:
    quality: 80
    scale: 80
    strip_metadata: true
  
  thumbnails:
    preset: thumb
    max_dimension: 400
```
**Config File Locations** (priority order):
1. `./.imgconvert.yaml` (current directory - highest priority)
2. `~/.config/imgconvert/config.yaml` (user config)
3. `~/.imgconvert.yaml` (user home - fallback)
**Behavior**:
- CLI flags override config file settings
- `--no-config`: Ignore all config files
- `--profile <name>`: Use specific profile from config
**Usage**: `imgconvert --profile web *.jpg`

### Additional V2 Features to Consider
- **More Output Formats**: Support converting to AVIF, JPEG XL
- **More Input Formats**: HEIC/HEIF, TIFF, AVIF, RAW formats
- **Image Cropping**: Smart crop, manual crop coordinates
- **Watermarking**: Add text/image watermarks
- **Batch Naming**: Advanced rename patterns (e.g., `--rename "img_{index}_{date}"`)
- **Image Filters**: Sharpen, blur, brightness, contrast adjustments
- **Smart Optimization**: Analyze image and auto-select best quality
- **Resume Support**: Resume interrupted bulk operations
- **Watch Mode**: Monitor directory and auto-convert new images
- **GUI Mode**: Optional simple GUI for non-CLI users
- **Cloud Integration**: Direct upload to S3, Cloudinary, etc.
- **Plugin System**: Allow custom processing plugins

---

