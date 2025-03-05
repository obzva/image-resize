package main

import (
	"flag"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"log"
	"math"
	"os"
	"regexp"
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

func nearestNeighbor(src *image.RGBA, w, h int, output string) {
	srcRect := src.Bounds()
	srcW, srcH := srcRect.Dx(), srcRect.Dy()

	res := image.NewRGBA(image.Rect(0, 0, w, h))

	for y := range h {
		for x := range w {
			// interpolated x and y
			// we can get the value we want by just doing integer division (without subtracting 0.5)
			iX, iY := x*srcW/w, y*srcH/h
			res.Set(x, y, src.At(iX, iY))
		}
	}

	if output == "" {
		output = "nearest-neighbor.jpg"
	}
	f, err := os.Create(output)
	checkError(err)
	defer f.Close()
	if err := jpeg.Encode(f, res, &jpeg.Options{Quality: 100}); err != nil {
		log.Fatal("encoding failed")
	}
}

func bilinear(src *image.RGBA, w, h int, output string) {
	srcRect := src.Bounds()
	srcW, srcH := srcRect.Dx(), srcRect.Dy()

	res := image.NewRGBA(image.Rect(0, 0, w, h))

	for y := range h {
		for x := range w {
			// interpolated x and y
			// subtract 0.5 for convinience
			iX := float64(x)*float64(srcW)/float64(w) - 0.5
			iY := float64(y)*float64(srcH)/float64(h) - 0.5

			// interpolated color that will be calculated
			iC := color.RGBA{}

			// boundary check
			outX := iX < 0 || iX > float64(w-1)
			outY := iY < 0 || iY > float64(h-1)

			// use weighted mean method(https://en.wikipedia.org/wiki/Bilinear_interpolation#Weighted_mean)
			// n: nearest
			// l: left
			// r: right
			// t: top
			// b: bottom
			// W: weight

			// use just one nearest surrounding point
			if outX && outY {
				var nX int
				var nY int

				if iX < 0 {
					nX = 0
				} else {
					nX = w - 1
				}

				if iY < 0 {
					nY = 0
				} else {
					nY = h - 1
				}

				iC = src.RGBAAt(nX, nY)
			} else if outX { // use two surrounding points (only y-axis)
				var nX float64

				if iX < 0 {
					nX = 0
				} else {
					nX = float64(w - 1)
				}

				tY := math.Floor(iY)
				tR, tB, tG, tA := src.At(int(nX), int(tY)).RGBA()

				bY := tY + 1
				bR, bB, bG, bA := src.At(int(nX), int(bY)).RGBA()

				iC.R = uint8(float64(tR)*(iY-tY) + float64(bR)*(bY-iY))
				iC.G = uint8(float64(tG)*(iY-tY) + float64(bG)*(bY-iY))
				iC.B = uint8(float64(tB)*(iY-tY) + float64(bB)*(bY-iY))
				iC.A = uint8(float64(tA)*(iY-tY) + float64(bA)*(bY-iY))
			} else if outY { // use two surrounding points (only x-axis)
				var nY float64

				if iY < 0 {
					nY = 0
				} else {
					nY = float64(h - 1)
				}

				lX := math.Floor(iX)
				lR, lG, lB, lA := src.At(int(lX), int(nY)).RGBA()

				rX := lX + 1
				rR, rG, rB, rA := src.At(int(rX), int(nY)).RGBA()

				iC.R = uint8(float64(lR)*(iX-lX) + float64(rR)*(rX-iX))
				iC.G = uint8(float64(lG)*(iX-lX) + float64(rG)*(rX-iX))
				iC.B = uint8(float64(lB)*(iX-lX) + float64(rB)*(rX-iX))
				iC.A = uint8(float64(lA)*(iX-lX) + float64(rA)*(rX-iX))
			} else { // use four surrounding points
				ltX := math.Floor(iX)
				ltY := math.Floor(iY)
				ltW := (iX - ltX) * (iY - ltY)
				ltR, ltG, ltB, ltA := src.At(int(ltX), int(ltY)).RGBA()

				rtX := ltX + 1
				rtY := ltY
				rtW := (rtX - iX) * (iY - rtY)
				rtR, rtG, rtB, rtA := src.At(int(rtX), int(rtY)).RGBA()

				lbX := ltX
				lbY := ltY + 1
				lbW := (iX - lbX) * (lbY - iY)
				lbR, lbG, lbB, lbA := src.At(int(lbX), int(lbY)).RGBA()

				rbX := ltX + 1
				rbY := ltY + 1
				rbW := (rbX - iX) * (rbY - iY)
				rbR, rbG, rbB, rbA := src.At(int(rbX), int(rbY)).RGBA()

				iC.R = uint8(float64(ltR)*ltW + float64(rtR)*rtW + float64(lbR)*lbW + float64(rbR)*rbW)
				iC.G = uint8(float64(ltG)*ltW + float64(rtG)*rtW + float64(lbG)*lbW + float64(rbG)*rbW)
				iC.B = uint8(float64(ltB)*ltW + float64(rtB)*rtW + float64(lbB)*lbW + float64(rbB)*rbW)
				iC.A = uint8(float64(ltA)*ltW + float64(rtA)*rtW + float64(lbA)*lbW + float64(rbA)*rbW)
			}
			res.Set(x, y, iC)
		}
	}

	if output == "" {
		output = "bilinear.jpg"
	}
	f, err := os.Create(output)
	checkError(err)
	defer f.Close()
	if err := jpeg.Encode(f, res, &jpeg.Options{Quality: 100}); err != nil {
		log.Fatal("encoding failed")
	}
}

func main() {
	pathPtr := flag.String("p", "", "input image path")
	wPtr := flag.Int("w", 0, "desired width of an output image (defaults to the original width when omitted)")
	hPtr := flag.Int("h", 0, "desired height of an output image (defaults to the original height when omitted)")
	methodPtr := flag.String("m", "nearest-neighbor", "desired interpolation method (defaults to nearest-neighbor when omitted)")
	outputPtr := flag.String("o", "", "desired output filename (defaults to the method name when omitted)")

	flag.Parse()

	if *pathPtr == "" {
		log.Fatal("input image path is required")
	}
	matched, err := regexp.MatchString(`\.jpe?g$`, *pathPtr)
	checkError(err)
	if !matched {
		log.Fatal("input image only available for jpg and jpeg")
	}

	src := getRGBAFromPath(*pathPtr)

	if *wPtr == 0 {
		*wPtr = src.Bounds().Dx()
	}
	if *hPtr == 0 {
		*hPtr = src.Bounds().Dy()
	}

	switch *methodPtr {
	case "nearest-neighbor":
		nearestNeighbor(src, *wPtr, *hPtr, *outputPtr)
	case "bilinear":
		bilinear(src, *wPtr, *hPtr, *outputPtr)
	}
}
