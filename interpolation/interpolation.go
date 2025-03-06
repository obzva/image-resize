package interpolation

import (
	"image"
	"image/color"
	"log"
	"math"
)

func NearestNeighbor(src *image.RGBA, w, h int) *image.RGBA {
	srcRect := src.Bounds()
	srcW, srcH := srcRect.Dx(), srcRect.Dy()

	res := image.NewRGBA(image.Rect(0, 0, w, h))

	for y := range h {
		for x := range w {
			// interpolated x and y
			// we can get the value we want by just doing integer division (without subtracting offset)
			iX, iY := x*srcW/w, y*srcH/h
			res.Set(x, y, src.At(iX, iY))
		}
	}

	return res
}

func Bilinear(src *image.RGBA, w, h int) *image.RGBA {
	srcRect := src.Bounds()
	srcW := srcRect.Dx()
	srcH := srcRect.Dy()
	
	scaleW := float64(w) / float64(srcW)
	scaleH := float64(h) / float64(srcH)

	offsetX := getOffset(scaleW)
	offsetY := getOffset(scaleH)

	res := image.NewRGBA(image.Rect(0, 0, w, h))

	for y := range h {
		for x := range w {
			// interpolated x and y
			// subtract offset for convinience (its effect is as same as linearly movin original coordinates)
			iX := float64(x)/scaleW - offsetX
			iY := float64(y)/scaleH - offsetY

			// boundary check
			outX := iX < 0 || iX > float64(srcW-1)
			outY := iY < 0 || iY > float64(srcH-1)

			var iC color.RGBA

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
					nX = srcW - 1
				}

				if iY < 0 {
					nY = 0
				} else {
					nY = srcH - 1
				}

				iC = src.RGBAAt(nX, nY)
			} else if outX { // use two surrounding points (only y-axis)
				var nX float64
				if iX < 0 {
					nX = 0
				} else {
					nX = float64(srcW - 1)
				}

				tY := math.Floor(iY)
				tRGBA := src.RGBAAt(int(nX), int(tY))
				tR := float64(tRGBA.R)
				tG := float64(tRGBA.G)
				tB := float64(tRGBA.B)
				tA := float64(tRGBA.A)

				bY := tY + 1
				bRGBA := src.RGBAAt(int(nX), int(bY))
				bR := float64(bRGBA.R)
				bG := float64(bRGBA.G)
				bB := float64(bRGBA.B)
				bA := float64(bRGBA.A)

				iR, iG, iB, iA := internalDivisionY(tR, tG, tB, tA, bR, bG, bB, bA, tY, bY, iY)
				// use round to dodge edge cases like:
				// uint8(1.99999999999) -> 1
				iC = color.RGBA{uint8(math.Round(iR)), uint8(math.Round(iG)), uint8(math.Round(iB)), uint8(math.Round(iA))}
			} else if outY { // use two surrounding points (only x-axis)
				var nY float64

				if iY < 0 {
					nY = 0
				} else {
					nY = float64(srcH - 1)
				}

				lX := math.Floor(iX)
				lRGBA := src.RGBAAt(int(lX), int(nY))
				lR := float64(lRGBA.R)
				lG := float64(lRGBA.G)
				lB := float64(lRGBA.B)
				lA := float64(lRGBA.A)

				rX := lX + 1
				rRGBA := src.RGBAAt(int(rX), int(nY))
				rR := float64(rRGBA.R)
				rG := float64(rRGBA.G)
				rB := float64(rRGBA.B)
				rA := float64(rRGBA.A)

				iR, iG, iB, iA := internalDivisionX(lR, lG, lB, lA, rR, rG, rB, rA, lX, rX, iX)
				// use round to dodge edge cases like:
				// uint8(1.99999999999) -> 1
				iC = color.RGBA{uint8(math.Round(iR)), uint8(math.Round(iG)), uint8(math.Round(iB)), uint8(math.Round(iA))}
			} else { // use four surrounding points
				ltX := math.Floor(iX)
				ltY := math.Floor(iY)
				ltRGBA := src.RGBAAt(int(ltX), int(ltY))
				ltR := float64(ltRGBA.R)
				ltG := float64(ltRGBA.G)
				ltB := float64(ltRGBA.B)
				ltA := float64(ltRGBA.A)

				rtX := ltX + 1
				rtY := ltY
				rtRGBA := src.RGBAAt(int(rtX), int(rtY))
				rtR := float64(rtRGBA.R)
				rtG := float64(rtRGBA.G)
				rtB := float64(rtRGBA.B)
				rtA := float64(rtRGBA.A)

				lbX := ltX
				lbY := ltY + 1
				lbRGBA := src.RGBAAt(int(lbX), int(lbY))
				lbR := float64(lbRGBA.R)
				lbG := float64(lbRGBA.G)
				lbB := float64(lbRGBA.B)
				lbA := float64(lbRGBA.A)

				rbX := ltX + 1
				rbY := ltY + 1
				rbRGBA := src.RGBAAt(int(rbX), int(rbY))
				rbR := float64(rbRGBA.R)
				rbG := float64(rbRGBA.G)
				rbB := float64(rbRGBA.B)
				rbA := float64(rbRGBA.A)

				tmp1R, tmp1G, tmp1B, tmp1A := internalDivisionX(ltR, ltG, ltB, ltA, rtR, rtG, rtB, rtA, ltX, rtX, iX)
				tmp2R, tmp2G, tmp2B, tmp2A := internalDivisionX(lbR, lbG, lbB, lbA, rbR, rbG, rbB, rbA, lbX, rbX, iX)

				iR, iG, iB, iA := internalDivisionY(tmp1R, tmp1G, tmp1B, tmp1A, tmp2R, tmp2G, tmp2B, tmp2A, ltY, lbY, iY)
				// use round to dodge edge cases like:
				// uint8(1.99999999999) -> 1
				iC = color.RGBA{uint8(math.Round(iR)), uint8(math.Round(iG)), uint8(math.Round(iB)), uint8(math.Round(iA))}
			}
			res.Set(x, y, iC)
		}
	}
	return res
}

func getOffset(k float64) float64 {
	return (k-1)/(2*k)
}

func internalDivisionX(r1, g1, b1, a1, r2, g2, b2, a2, x1, x2, x float64) (r float64, g float64, b float64, a float64) {
	if !(x1 <= x && x <= x2) {
		log.Fatalf("it should be like: x1 <= x <= x2 but got x1: %f, x2: %f, and x:%f\n", x1, x2, x)
	}

	r = (x2 - x) / (x2 - x1) * float64(r1) + (x - x1) / (x2 - x1) * float64(r2)
	g = (x2 - x) / (x2 - x1) * float64(g1) + (x - x1) / (x2 - x1) * float64(g2)
	b = (x2 - x) / (x2 - x1) * float64(b1) + (x - x1) / (x2 - x1) * float64(b2)
	a = (x2 - x) / (x2 - x1) * float64(a1) + (x - x1) / (x2 - x1) * float64(a2)

	return r, g, b, a
}

func internalDivisionY(r1, g1, b1, a1, r2, g2, b2, a2, y1, y2, y float64) (r float64, g float64, b float64, a float64) {
	if !(y1 <= y && y <= y2) {
		log.Fatalf("it should be like: y1 <= y <= y2but got y1: %f, y2: %f, and y:%f", y1, y2, y)
	}

	r = (y2 - y) / (y2 - y1) * float64(r1) + (y - y1) / (y2 - y1) * float64(r2)
	g = (y2 - y) / (y2 - y1) * float64(g1) + (y - y1) / (y2 - y1) * float64(g2)
	b = (y2 - y) / (y2 - y1) * float64(b1) + (y - y1) / (y2 - y1) * float64(b2)
	a = (y2 - y) / (y2 - y1) * float64(a1) + (y - y1) / (y2 - y1) * float64(a2)

	return r, g, b, a
}