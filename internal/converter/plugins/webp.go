package plugins

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"

	"github.com/chai2010/webp"
	"github.com/miracleio/imgconvert/internal/converter"
	"golang.org/x/image/bmp"
)

// WebPConverter implements WebP encoding using chai2010/webp
type WebPConverter struct{}

// NewWebPConverter creates a new WebP converter
func NewWebPConverter() *WebPConverter {
	return &WebPConverter{}
}

// Decode reads and decodes an image
func (c *WebPConverter) Decode(src io.Reader) (image.Image, error) {
	img, format, err := image.Decode(src)
	if err != nil {
		// Try BMP if standard decode fails
		if format == "" {
			return bmp.Decode(src)
		}
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}
	return img, nil
}

// Encode writes the image as WebP
func (c *WebPConverter) Encode(dst io.Writer, img image.Image, opts *converter.EncodeOptions) error {
	quality := float32(opts.Quality)
	if opts.Lossless {
		// Use lossless encoding
		return webp.Encode(dst, img, &webp.Options{Lossless: true})
	}
	
	// Use lossy encoding with specified quality
	return webp.Encode(dst, img, &webp.Options{
		Lossless: false,
		Quality:  quality,
	})
}

// SupportsFormat checks if format is supported for input
func (c *WebPConverter) SupportsFormat(format string) bool {
	supported := map[string]bool{
		"jpeg": true,
		"jpg":  true,
		"png":  true,
		"gif":  true,
		"bmp":  true,
		"webp": true,
	}
	return supported[format]
}

// Info returns converter information
func (c *WebPConverter) Info() converter.ConverterInfo {
	return converter.ConverterInfo{
		Name:          "chai2010/webp",
		Version:       "1.4.0",
		SupportsAnim:  false, // Basic version doesn't support animated GIF yet
		SupportsAlpha: true,
	}
}
