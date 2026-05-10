package scanner

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Options for file scanning
type Options struct {
	Recursive   bool
	FollowLinks bool
	IncludeHidden bool
}

// FileInfo represents a discovered file
type FileInfo struct {
	Path      string
	Size      int64
	Extension string
}

// Scanner discovers image files
type Scanner struct {
	opts Options
}

// New creates a new file scanner
func New(opts Options) *Scanner {
	return &Scanner{opts: opts}
}

// Scan discovers files from the given paths
func (s *Scanner) Scan(paths []string) ([]FileInfo, error) {
	var files []FileInfo
	seen := make(map[string]bool) // Avoid duplicates

	for _, path := range paths {
		info, err := os.Stat(path)
		if err != nil {
			return nil, fmt.Errorf("cannot access %s: %w", path, err)
		}

		if info.IsDir() {
			// Scan directory
			dirFiles, err := s.scanDirectory(path)
			if err != nil {
				return nil, err
			}
			for _, f := range dirFiles {
				if !seen[f.Path] {
					files = append(files, f)
					seen[f.Path] = true
				}
			}
		} else {
			// Single file
			if s.isImageFile(path) {
				absPath, _ := filepath.Abs(path)
				if !seen[absPath] {
					files = append(files, FileInfo{
						Path:      absPath,
						Size:      info.Size(),
						Extension: strings.ToLower(filepath.Ext(path)),
					})
					seen[absPath] = true
				}
			}
		}
	}

	return files, nil
}

// scanDirectory recursively scans a directory
func (s *Scanner) scanDirectory(dir string) ([]FileInfo, error) {
	var files []FileInfo

	walkFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip hidden files/directories if needed
		if !s.opts.IncludeHidden && strings.HasPrefix(filepath.Base(path), ".") {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Handle directories
		if info.IsDir() {
			if path != dir && !s.opts.Recursive {
				return filepath.SkipDir
			}
			return nil
		}

		// Check if it's an image file
		if s.isImageFile(path) {
			absPath, _ := filepath.Abs(path)
			files = append(files, FileInfo{
				Path:      absPath,
				Size:      info.Size(),
				Extension: strings.ToLower(filepath.Ext(path)),
			})
		}

		return nil
	}

	err := filepath.Walk(dir, walkFn)
	return files, err
}

// isImageFile checks if the file is a supported image format
func (s *Scanner) isImageFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	supportedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".bmp":  true,
		".webp": true,
	}
	return supportedExts[ext]
}
