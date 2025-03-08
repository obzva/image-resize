package main

import (
	"flag"
	"log"

	"gthub.com/obzva/image-resize/imageprocessor"
)

func main() {
	// flags
	pathPtr := flag.String("p", "", "input image path")
	wPtr := flag.Int("w", 0, "desired width of output image, defaults to keep the ratio of the original image when omitted (at least one of two, width or height, is required)")
	hPtr := flag.Int("h", 0, "desired height of output image, defaults to keep the ratio of the original image when omitted (at least one of two, width or height, is required)")
	methodPtr := flag.String("m", "nearestneighbor", "desired interpolation method, defaults to nearestneighbor (options: nearestneighbor, bilinear, and bicubic)")
	outputPtr := flag.String("o", "", "desired output filename, defaults to the method name when omitted")
	concurrencyPtr := flag.Bool("c", true, "concurrency mode, defaults to true when omitted")

	flag.Parse()

	ip := imageprocessor.New(*pathPtr, *wPtr, *hPtr, *methodPtr, *concurrencyPtr, *outputPtr)

	err := ip.CreateImageFile()
	if err != nil {
		log.Fatal(err)
	}
}
