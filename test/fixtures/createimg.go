package main

import (
"image"
"image/color"
"image/jpeg"
"image/png"
"os"
)

func main() {
img := image.NewRGBA(image.Rect(0, 0, 400, 300))
for y := 0; y < 300; y++ {
for x := 0; x < 400; x++ {
img.Set(x, y, color.RGBA{uint8(x % 256), uint8(y % 256), 150, 255})
}
}

f, _ := os.Create("sample.jpg")
jpeg.Encode(f, img, &jpeg.Options{Quality: 90})
f.Close()

f2, _ := os.Create("sample.png")
png.Encode(f2, img)
f2.Close()
}
