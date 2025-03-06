package imageio

import (
	"image"
	"image/draw"
	"image/jpeg"
	"log"
	"os"
)

func ReadImage(path string) *image.RGBA {
	// read a jpeg image from file path
	reader, err := os.Open(path)
	if err != nil {
		log.Fatal("imageio: error occurred when opening the image")
	}
	defer reader.Close()

	// decode in-memory image into image.Image interface
	src, err := jpeg.Decode(reader)
	if err != nil {
		log.Fatal("imageio: error occurred when decoding the image")
	}

	// get the size of original source image
	srcRect := src.Bounds()
	srcW, srcH := srcRect.Size().X, srcRect.Size().Y

	// convert image into more useful form (RGBA interface)
	// so that we can pass it to draw.Draw
	res := image.NewRGBA(image.Rect(0, 0, srcW, srcH))
	draw.Draw(res, srcRect, src, srcRect.Min, draw.Src)

	return res
}

func CreateImageFile(output string, img *image.RGBA) {
	f, err := os.Create(output)
	if err != nil {
		log.Fatal("error occurred when creating output file")
	}
	defer f.Close()
	if err := jpeg.Encode(f, img, &jpeg.Options{Quality: 100}); err != nil {
		log.Fatal("jpeg encoding failed")
	}
}
