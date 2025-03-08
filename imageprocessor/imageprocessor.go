package imageprocessor

import (
	"errors"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"log"
	"math"
	"os"
	"regexp"

	"gthub.com/obzva/image-resize/interpolator"
)

type ImageProcessor struct {
	path         string      // path to the input file
	iExt         string      // "jpeg" | "png" extension of the input file, only jpeg(jpg), png are available
	src          *image.NRGBA // in-memory input image converted to *image.NRGBA
	w, h         int         // width and height of output image file
	name         string      // name of output image file
	oExt         string      // "jpeg" | "png" extension of the output file, only jpeg(jpg), png are available
	concurrency  bool
	interpolator interpolator.Interpolator
}

// readImageFile the input image and then convert it into *image.NRGBA
// after that, set that into ip.src
func (ip *ImageProcessor) readImageFile() error {
	// read the image from file path
	r, err := os.Open(ip.path)
	if err != nil {
		return errors.New("error occurred when opening the image file")
	}
	defer r.Close()

	var i image.Image
	// decode in-memory image into image.Image interface
	if ip.iExt == "jpeg" {
		d, err := jpeg.Decode(r)
		if err != nil {
			return errors.New("error occurred when decoding the image")
		}
		i = d

	} else {
		d, err := png.Decode(r)
		if err != nil {
			return errors.New("error occurred when decoding the image")
		}
		i = d
	}

	// get the size of the input image
	iRect := i.Bounds()
	iW, iH := iRect.Size().X, iRect.Size().Y

	// convert image into more useful form *image.NRGBA
	// so that we can pass it to draw.Draw
	rgba := image.NewNRGBA(image.Rect(0, 0, iW, iH))
	draw.Draw(rgba, iRect, i, iRect.Min, draw.Src)

	ip.src = rgba

	return nil
}

func (ip *ImageProcessor) CreateImageFile() error {
	f, err := os.Create(ip.name)
	if err != nil {
		return errors.New("error occurred when creating output file")
	}
	defer f.Close()

	p := ip.interpolator.Interpolate(ip.concurrency)

	if ip.oExt == "jpeg" {
		if err := jpeg.Encode(f, p, &jpeg.Options{Quality: 100}); err != nil {
			return err
		}
	} else {
		if err := png.Encode(f, p); err != nil {
			return err
		}
	}
	return nil
}

func New(path string, w, h int, method string, concurrency bool, name string) *ImageProcessor {
	// check path
	if path == "" {
		log.Fatal("input image path is required")
	}

	// check input file extension
	iExt := extCheck(path)

	// set path, extension, and concurrency
	ip := &ImageProcessor{
		path:        path,
		iExt:        iExt,
		concurrency: concurrency,
	}

	// read input and set src
	err := ip.readImageFile()
	if err != nil {
		log.Fatal(err)
	}

	// set w and h
	if w == 0 && h == 0 {
		log.Fatal("at least one dimension, w or h, is required")
	} else if w == 0 {
		iH := ip.src.Bounds().Dy()
		scale := float64(h) / float64(iH)
		w = int(math.Round(float64(ip.src.Bounds().Dx()) * scale))
	} else if h == 0 {
		iW := ip.src.Bounds().Dx()
		scale := float64(w) / float64(iW)
		h = int(math.Round(float64(ip.src.Bounds().Dy()) * scale))
	}
	ip.w = w
	ip.h = h

	// set name (output filename)
	if name == "" {
		name = method + "." + ip.iExt
	}
	ip.name = name

	// set extension of output file
	oExt := extCheck(name)
	ip.oExt = oExt

	// set interpolator
	ip.interpolator = interpolator.New(ip.src, ip.w, ip.h, method)

	return ip
}

func extCheck(s string) string {
	re, err := regexp.Compile(`\.(jpe?g|png)$`)
	if err != nil {
		log.Fatal("error occurred while compiling regexp")
	}
	matches := re.FindStringSubmatch(s)
	if matches == nil {
		log.Fatal("input image only available for jpg/jpeg and png")
	}
	extension := matches[1]
	if extension == "jpg" {
		extension = "jpeg"
	}
	return extension
}
