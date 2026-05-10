package converter

import (
	"image"
	"io"
)

// EncodeOptions contains settings for image encoding
type EncodeOptions struct {
	Quality      int
	Lossless     bool
	Effort       int
	StripMetadata bool
}

// ConverterInfo provides metadata about a converter
type ConverterInfo struct {
	Name        string
	Version     string
	SupportsAnim bool
	SupportsAlpha bool
}

// Converter defines the interface for image format converters
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
