package main

import (
	"flag"
	"gthub.com/obzva/image-resize/imageio"
	"gthub.com/obzva/image-resize/interpolation"
	"image"
	"log"
	"regexp"
)

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
	if err != nil {
		log.Fatal("error occurred while compiling regexp")
	}
	if !matched {
		log.Fatal("input image only available for jpg and jpeg")
	}

	src := imageio.ReadImage(*pathPtr)

	if *wPtr == 0 {
		*wPtr = src.Bounds().Dx()
	}
	if *hPtr == 0 {
		*hPtr = src.Bounds().Dy()
	}

	if *outputPtr == "" {
		*outputPtr = *methodPtr + ".jpg"
	}

	var res *image.RGBA
	switch *methodPtr {
	case "nearest-neighbor":
		res = interpolation.NearestNeighbor(src, *wPtr, *hPtr)
	case "bilinear":
		res = interpolation.Bilinear(src, *wPtr, *hPtr)
	}

	imageio.CreateImageFile(*outputPtr, res)
}
