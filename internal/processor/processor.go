package processor

import (
	"fmt"
	"image"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/disintegration/imaging"
	"github.com/miracleio/imgconvert/internal/converter"
	"github.com/miracleio/imgconvert/internal/scanner"
	"github.com/schollz/progressbar/v3"
)

// Options for processing
type Options struct {
	Quality         int
	Lossless        bool
	Effort          int
	OutputDir       string
	Suffix          string
	Overwrite       bool
	DeleteOriginal  bool
	Workers         int
	Quiet           bool
	DryRun          bool
	
	// Resize options
	Width          int
	Height         int
	MaxDimension   int
	Scale          float64
}

// Result of a conversion
type Result struct {
	SourcePath   string
	OutputPath   string
	OriginalSize int64
	NewSize      int64
	Success      bool
	Error        error
}

// Stats tracks overall conversion statistics
type Stats struct {
	TotalFiles      int
	Converted       int
	Failed          int
	Skipped         int
	TotalOrigSize   int64
	TotalNewSize    int64
}

// Processor handles parallel image conversion
type Processor struct {
	converter converter.Converter
	opts      Options
}

// New creates a new processor
func New(conv converter.Converter, opts Options) *Processor {
	return &Processor{
		converter: conv,
		opts:      opts,
	}
}

// Process converts all files
func (p *Processor) Process(files []scanner.FileInfo) (*Stats, []Result) {
	stats := &Stats{TotalFiles: len(files)}
	results := make([]Result, 0, len(files))
	resultsChan := make(chan Result, len(files))
	
	// Create progress bar
	var bar *progressbar.ProgressBar
	if !p.opts.Quiet && !p.opts.DryRun {
		bar = progressbar.Default(int64(len(files)), "Converting")
	}
	
	// Worker pool
	var wg sync.WaitGroup
	sem := make(chan struct{}, p.opts.Workers)
	
	var converted, failed, skipped int32
	var totalOrigSize, totalNewSize int64
	
	for _, file := range files {
		wg.Add(1)
		go func(f scanner.FileInfo) {
			defer wg.Done()
			sem <- struct{}{}        // Acquire
			defer func() { <-sem }() // Release
			
			result := p.processFile(f)
			resultsChan <- result
			
			if result.Success {
				atomic.AddInt32(&converted, 1)
				atomic.AddInt64(&totalOrigSize, result.OriginalSize)
				atomic.AddInt64(&totalNewSize, result.NewSize)
			} else if result.Error == nil {
				atomic.AddInt32(&skipped, 1)
			} else {
				atomic.AddInt32(&failed, 1)
			}
			
			if bar != nil {
				bar.Add(1)
			}
		}(file)
	}
	
	// Wait for all workers
	go func() {
		wg.Wait()
		close(resultsChan)
	}()
	
	// Collect results
	for result := range resultsChan {
		results = append(results, result)
	}
	
	if bar != nil {
		bar.Finish()
	}
	
	stats.Converted = int(converted)
	stats.Failed = int(failed)
	stats.Skipped = int(skipped)
	stats.TotalOrigSize = totalOrigSize
	stats.TotalNewSize = totalNewSize
	
	return stats, results
}

// processFile converts a single file
func (p *Processor) processFile(file scanner.FileInfo) Result {
	result := Result{
		SourcePath:   file.Path,
		OriginalSize: file.Size,
	}
	
	// Determine output path
	outputPath := p.getOutputPath(file.Path)
	result.OutputPath = outputPath
	
	// Check if output exists
	if !p.opts.Overwrite {
		if _, err := os.Stat(outputPath); err == nil {
			return result // Skip, file exists
		}
	}
	
	// Dry run mode
	if p.opts.DryRun {
		// Estimate output size (rough estimate: 30% of original for quality 85)
		estimatedSize := int64(float64(file.Size) * 0.3)
		result.NewSize = estimatedSize
		result.Success = true
		return result
	}
	
	// Open source file
	srcFile, err := os.Open(file.Path)
	if err != nil {
		result.Error = fmt.Errorf("cannot open: %w", err)
		return result
	}
	defer srcFile.Close()
	
	// Decode image
	img, err := p.converter.Decode(srcFile)
	if err != nil {
		result.Error = fmt.Errorf("cannot decode: %w", err)
		return result
	}
	
	// Apply resizing if needed
	img = p.resizeImage(img)
	
	// Create output directory if needed
	outDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(outDir, 0755); err != nil {
		result.Error = fmt.Errorf("cannot create output dir: %w", err)
		return result
	}
	
	// Create output file
	outFile, err := os.Create(outputPath)
	if err != nil {
		result.Error = fmt.Errorf("cannot create output: %w", err)
		return result
	}
	defer outFile.Close()
	
	// Encode to WebP
	encOpts := &converter.EncodeOptions{
		Quality:  p.opts.Quality,
		Lossless: p.opts.Lossless,
		Effort:   p.opts.Effort,
	}
	
	if err := p.converter.Encode(outFile, img, encOpts); err != nil {
		result.Error = fmt.Errorf("cannot encode: %w", err)
		os.Remove(outputPath) // Clean up failed output
		return result
	}
	
	// Get output file size
	if stat, err := os.Stat(outputPath); err == nil {
		result.NewSize = stat.Size()
	}
	
	result.Success = true
	
	// Delete original if requested
	if p.opts.DeleteOriginal {
		os.Remove(file.Path)
	}
	
	return result
}

// getOutputPath determines the output file path
func (p *Processor) getOutputPath(sourcePath string) string {
	dir := filepath.Dir(sourcePath)
	base := filepath.Base(sourcePath)
	ext := filepath.Ext(base)
	nameWithoutExt := strings.TrimSuffix(base, ext)
	
	// Add suffix if specified
	if p.opts.Suffix != "" {
		nameWithoutExt += p.opts.Suffix
	}
	
	// Determine output directory
	if p.opts.OutputDir != "" {
		dir = p.opts.OutputDir
	}
	
	return filepath.Join(dir, nameWithoutExt+".webp")
}

// resizeImage applies resizing transformations
func (p *Processor) resizeImage(img image.Image) image.Image {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	
	// Apply scale percentage
	if p.opts.Scale > 0 && p.opts.Scale != 100 {
		newWidth := int(float64(width) * p.opts.Scale / 100.0)
		newHeight := int(float64(height) * p.opts.Scale / 100.0)
		return imaging.Resize(img, newWidth, newHeight, imaging.Lanczos)
	}
	
	// Apply max dimension
	if p.opts.MaxDimension > 0 {
		maxDim := max(width, height)
		if maxDim > p.opts.MaxDimension {
			scale := float64(p.opts.MaxDimension) / float64(maxDim)
			newWidth := int(float64(width) * scale)
			newHeight := int(float64(height) * scale)
			return imaging.Resize(img, newWidth, newHeight, imaging.Lanczos)
		}
	}
	
	// Apply width constraint
	if p.opts.Width > 0 && width > p.opts.Width {
		return imaging.Resize(img, p.opts.Width, 0, imaging.Lanczos)
	}
	
	// Apply height constraint
	if p.opts.Height > 0 && height > p.opts.Height {
		return imaging.Resize(img, 0, p.opts.Height, imaging.Lanczos)
	}
	
	return img
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
