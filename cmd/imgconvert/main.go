package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/miracleio/imgconvert/internal/converter/plugins"
	"github.com/miracleio/imgconvert/internal/processor"
	"github.com/miracleio/imgconvert/internal/scanner"
	"github.com/spf13/cobra"
)

var (
	version = "0.1.0"
	
	// Global flags
	recursive      bool
	outputDir      string
	quality        int
	lossless       bool
	jobs           int
	quiet          bool
	verbose        bool
	dryRun         bool
	overwrite      bool
	deleteOriginal bool
	suffix         string
	preset         string
	
	// Resize flags
	width        int
	height       int
	maxDimension int
	scale        float64
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "imgconvert [files/directories...]",
		Short: "High-performance image converter to WebP",
		Long: `imgconvert - Convert images to WebP format with optimal compression

Supports: JPEG, PNG, GIF, BMP → WebP
Features: Bulk conversion, parallel processing, resizing, quality control`,
		Version: version,
		Args:    cobra.MinimumNArgs(1),
		RunE:    run,
	}
	
	// Input/Output flags
	rootCmd.Flags().BoolVarP(&recursive, "recursive", "r", false, "Recursively scan directories")
	rootCmd.Flags().StringVarP(&outputDir, "output", "o", "", "Output directory")
	rootCmd.Flags().StringVar(&suffix, "suffix", "", "Add suffix to output filenames")
	rootCmd.Flags().BoolVar(&overwrite, "overwrite", false, "Overwrite existing files")
	rootCmd.Flags().BoolVar(&deleteOriginal, "delete-original", false, "Delete original files after conversion")
	
	// Quality flags
	rootCmd.Flags().IntVarP(&quality, "quality", "q", 85, "Quality level (0-100)")
	rootCmd.Flags().StringVar(&preset, "preset", "", "Quality preset (low/medium/high/lossless/thumb/half/web)")
	rootCmd.Flags().BoolVar(&lossless, "lossless", false, "Use lossless compression")
	
	// Resize flags
	rootCmd.Flags().IntVar(&width, "width", 0, "Maximum width (maintains aspect ratio)")
	rootCmd.Flags().IntVar(&height, "height", 0, "Maximum height (maintains aspect ratio)")
	rootCmd.Flags().IntVar(&maxDimension, "max-dimension", 0, "Maximum dimension for longest side")
	rootCmd.Flags().Float64Var(&scale, "scale", 0, "Scale by percentage (e.g., 50 for 50%)")
	
	// Processing flags
	rootCmd.Flags().IntVarP(&jobs, "jobs", "j", runtime.NumCPU(), "Number of parallel jobs")
	rootCmd.Flags().BoolVar(&quiet, "quiet", false, "Suppress progress output")
	rootCmd.Flags().BoolVar(&verbose, "verbose", false, "Verbose output")
	rootCmd.Flags().BoolVarP(&dryRun, "dry-run", "n", false, "Preview without converting")
	
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) error {
	// Apply presets
	applyPreset()
	
	// Validate quality
	if quality < 0 || quality > 100 {
		return fmt.Errorf("quality must be between 0 and 100")
	}
	
	// Create scanner
	scanOpts := scanner.Options{
		Recursive:     recursive,
		FollowLinks:   false,
		IncludeHidden: false,
	}
	scan := scanner.New(scanOpts)
	
	// Discover files
	if verbose {
		fmt.Println("Scanning for images...")
	}
	
	files, err := scan.Scan(args)
	if err != nil {
		return fmt.Errorf("scanning failed: %w", err)
	}
	
	if len(files) == 0 {
		return fmt.Errorf("no image files found")
	}
	
	if verbose || dryRun {
		fmt.Printf("Found %d image(s)\n", len(files))
	}
	
	// Create converter
	conv := plugins.NewWebPConverter()
	
	// Create processor
	procOpts := processor.Options{
		Quality:        quality,
		Lossless:       lossless,
		Effort:         4,
		OutputDir:      outputDir,
		Suffix:         suffix,
		Overwrite:      overwrite,
		DeleteOriginal: deleteOriginal,
		Workers:        jobs,
		Quiet:          quiet,
		DryRun:         dryRun,
		Width:          width,
		Height:         height,
		MaxDimension:   maxDimension,
		Scale:          scale,
	}
	proc := processor.New(conv, procOpts)
	
	// Process files
	stats, results := proc.Process(files)
	
	// Display results
	printResults(stats, results)
	
	// Exit with appropriate code
	if stats.Failed > 0 {
		return fmt.Errorf("%d file(s) failed to convert", stats.Failed)
	}
	
	return nil
}

func applyPreset() {
	switch preset {
	case "low":
		quality = 60
	case "medium":
		quality = 80
	case "high":
		quality = 90
	case "lossless":
		lossless = true
	case "thumb":
		scale = 33.33
	case "half":
		scale = 50
	case "web":
		scale = 80
	}
}

func printResults(stats *processor.Stats, results []processor.Result) {
	fmt.Println()
	
	if dryRun {
		fmt.Println("[DRY RUN] Would convert", stats.TotalFiles, "file(s)")
		fmt.Println()
	}
	
	// Show summary
	fmt.Printf("✓ %d file(s) converted successfully\n", stats.Converted)
	
	if stats.Skipped > 0 {
		fmt.Printf("⊘ %d file(s) skipped (already exist)\n", stats.Skipped)
	}
	
	if stats.Failed > 0 {
		fmt.Printf("✗ %d file(s) failed:\n", stats.Failed)
		for _, r := range results {
			if !r.Success && r.Error != nil {
				fmt.Printf("  - %s: %v\n", r.SourcePath, r.Error)
			}
		}
	}
	
	// Size reduction stats
	if stats.Converted > 0 && !dryRun {
		reduction := float64(stats.TotalOrigSize-stats.TotalNewSize) / float64(stats.TotalOrigSize) * 100
		fmt.Printf("\nSize: %s → %s (%.1f%% reduction)\n",
			formatSize(stats.TotalOrigSize),
			formatSize(stats.TotalNewSize),
			reduction)
	} else if dryRun && stats.Converted > 0 {
		fmt.Printf("\nEstimated size: %s → %s (~70%% reduction)\n",
			formatSize(stats.TotalOrigSize),
			formatSize(stats.TotalNewSize))
	}
	
	if dryRun {
		fmt.Println("\nRun without --dry-run to perform conversion")
	}
}

func formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
