package main

import (
"image"
"image/color"
"image/jpeg"
"image/png"
"os"
)

func main() {
// Create a simple test JPEG
img := image.NewRGBA(image.Rect(0, 0, 200, 150))
for y := 0; y < 150; y++ {
for x := 0; x < 200; x++ {
img.Set(x, y, color.RGBA{uint8(x % 256), uint8(y % 256), 100, 255})
}
}

// Save as JPEG
f, _ := os.Create("test.jpg")
jpeg.Encode(f, img, &jpeg.Options{Quality: 90})
f.Close()

// Save as PNG
f2, _ := os.Create("test.png")
png.Encode(f2, img)
f2.Close()
}
