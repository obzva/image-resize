package main

import (
	"image"
	"image/draw"
	"image/jpeg"
	"log"
	"os"
	"strconv"
)

func checkError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func getRGBAFromPath(path string) *image.RGBA {
	// read a jpeg image from file path
	reader, err := os.Open(path)
	checkError(err)
	defer reader.Close()

	// decode in-memory image into image.Image interface
	src, err := jpeg.Decode(reader)
	checkError(err)

	// get the size of original source image
	srcRect := src.Bounds()
	srcW, srcH := srcRect.Size().X, srcRect.Size().Y

	// convert image into more useful form (RGBA interface)
	// so that we can pass it to draw.Draw
	res := image.NewRGBA(image.Rect(0, 0, srcW, srcH))
	draw.Draw(res, srcRect, src, srcRect.Min, draw.Src)

	return res
}

func nearestNeighbor(src *image.RGBA, w, h int) {
	srcRect := src.Bounds()
	srcW, srcH := srcRect.Dx(), srcRect.Dy()

	res := image.NewRGBA(image.Rect(0, 0, w, h))

	for y := range h {
		for x := range w {
			srcX, srcY := x*srcW/w, y*srcH/h
			res.Set(x, y, src.At(srcX, srcY))
		}
	}

	f, err := os.Create("nearest-neighbor.jpg")
	checkError(err)
	defer f.Close()
	if err := jpeg.Encode(f, res, &jpeg.Options{Quality: 100}); err != nil {
		log.Fatal("encoding failed")
	}
}

func main() {
	src := getRGBAFromPath("original.jpg")

	w, err := strconv.Atoi(os.Args[1])
	checkError(err)

	h, err := strconv.Atoi(os.Args[2])
	checkError(err)

	nearestNeighbor(src, w, h)
}
