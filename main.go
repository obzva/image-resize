package main

import (
	"flag"
	"log"
	"math"
	"regexp"

	"gthub.com/obzva/image-resize/imageio"
	"gthub.com/obzva/image-resize/interpolation"
)

func main() {
	// flags
	pathPtr := flag.String("p", "", "input image path")
	wPtr := flag.Int("w", 0, "desired width of an output image (defaults to the original width when omitted)")
	hPtr := flag.Int("h", 0, "desired height of an output image (defaults to the original height when omitted)")
	methodPtr := flag.String("m", "nearestneighbor", "desired interpolation method (defaults to nearest-neighbor when omitted)")
	outputPtr := flag.String("o", "", "desired output filename (defaults to the method name when omitted)")
	concurrencyPtr := flag.Bool("c", true, "concurrency mode")

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

	if *wPtr == 0 && *hPtr == 0 {
		log.Fatal("at least one dimension, w or h, is required")
	} else if *wPtr == 0 {
		iH := src.Bounds().Dy()
		scale := float64(*hPtr) /  float64(iH)
		*wPtr = int(math.Round(float64(src.Bounds().Dx()) * scale))
	} else if *hPtr == 0 {
		iW := src.Bounds().Dx()
		scale := float64(*wPtr) /  float64(iW)
		*hPtr = int(math.Round(float64(src.Bounds().Dx()) * scale))
	}

	if *outputPtr == "" {
		*outputPtr = *methodPtr + ".jpg"
	}

	it := interpolation.Init(src, *wPtr, *hPtr, *methodPtr)

	res := it.Interpolate(*concurrencyPtr)

	imageio.CreateImageFile(*outputPtr, res)
}
