package interpolation

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
)

func getOffset(k float64) float64 {
	return (k - 1) / (2 * k)
}

// transform result rectangle's coordinate
// scale it down and
// if subtractOffset is true then
// subtract offset from it for convinience (its effect is as same as linearly movin original coordinates)
func transformCoords(x, y int, scaleW, scaleH float64, subtractOffset bool) (float64, float64) {
	offsetX := getOffset(scaleW)
	offsetY := getOffset(scaleH)

	transX := float64(x) / scaleW
	transY := float64(y) / scaleH

	if subtractOffset {
		transX -= offsetX
		transY -= offsetY
	}

	return transX, transY
}

func NearestNeighbor(src *image.RGBA, w, h int) *image.RGBA {
	srcRect := src.Bounds()
	srcW, srcH := srcRect.Dx(), srcRect.Dy()

	scaleW := float64(w) / float64(srcW)
	scaleH := float64(h) / float64(srcH)

	res := image.NewRGBA(image.Rect(0, 0, w, h))

	for y := range h {
		for x := range w {
			// transformed x and y
			transX, transY := transformCoords(x, y, scaleW, scaleH, false)

			res.Set(x, y, src.At(int(transX), int(transY)))
		}
	}

	return res
}

func internalDivision(r1, g1, b1, a1, r2, g2, b2, a2, v1, v2, v float64) (r float64, g float64, b float64, a float64) {
	if !(v1 <= v && v <= v2) {
		log.Fatalf("it should be like: v1 <= v <= v2 but got v1: %f, v2: %f, and v:%f\n", v1, v2, v)
	}

	r = (v2-v)/(v2-v1)*float64(r1) + (v-v1)/(v2-v1)*float64(r2)
	g = (v2-v)/(v2-v1)*float64(g1) + (v-v1)/(v2-v1)*float64(g2)
	b = (v2-v)/(v2-v1)*float64(b1) + (v-v1)/(v2-v1)*float64(b2)
	a = (v2-v)/(v2-v1)*float64(a1) + (v-v1)/(v2-v1)*float64(a2)

	return r, g, b, a
}

func Bilinear(src *image.RGBA, w, h int) *image.RGBA {
	srcRect := src.Bounds()
	srcW := srcRect.Dx()
	srcH := srcRect.Dy()

	scaleW := float64(w) / float64(srcW)
	scaleH := float64(h) / float64(srcH)

	res := image.NewRGBA(image.Rect(0, 0, w, h))

	for y := range h {
		for x := range w {
			// transformed x and y
			transX, transY := transformCoords(x, y, scaleW, scaleH, true)

			// boundary check
			outX := transX < 0 || transX > float64(srcW-1)
			outY := transY < 0 || transY > float64(srcH-1)

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

				if transX < 0 {
					nX = 0
				} else {
					nX = srcW - 1
				}

				if transY < 0 {
					nY = 0
				} else {
					nY = srcH - 1
				}

				iC = src.RGBAAt(nX, nY)
			} else if outX { // use two surrounding points (only y-axis)
				var nX float64
				if transX < 0 {
					nX = 0
				} else {
					nX = float64(srcW - 1)
				}

				tY := math.Floor(transY)
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

				iR, iG, iB, iA := internalDivision(tR, tG, tB, tA, bR, bG, bB, bA, tY, bY, transY)
				// use round to dodge edge cases like:
				// uint8(1.99999999999) -> 1
				iC = color.RGBA{uint8(math.Round(iR)), uint8(math.Round(iG)), uint8(math.Round(iB)), uint8(math.Round(iA))}
			} else if outY { // use two surrounding points (only x-axis)
				var nY float64

				if transY < 0 {
					nY = 0
				} else {
					nY = float64(srcH - 1)
				}

				lX := math.Floor(transX)
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

				iR, iG, iB, iA := internalDivision(lR, lG, lB, lA, rR, rG, rB, rA, lX, rX, transX)
				// use round to dodge edge cases like:
				// uint8(1.99999999999) -> 1
				iC = color.RGBA{uint8(math.Round(iR)), uint8(math.Round(iG)), uint8(math.Round(iB)), uint8(math.Round(iA))}
			} else { // use four surrounding points
				ltX := math.Floor(transX)
				ltY := math.Floor(transY)
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

				tmp1R, tmp1G, tmp1B, tmp1A := internalDivision(ltR, ltG, ltB, ltA, rtR, rtG, rtB, rtA, ltX, rtX, transX)
				tmp2R, tmp2G, tmp2B, tmp2A := internalDivision(lbR, lbG, lbB, lbA, rbR, rbG, rbB, rbA, lbX, rbX, transX)

				iR, iG, iB, iA := internalDivision(tmp1R, tmp1G, tmp1B, tmp1A, tmp2R, tmp2G, tmp2B, tmp2A, ltY, lbY, transY)
				// use round to dodge edge cases like:
				// uint8(1.99999999999) -> 1
				iC = color.RGBA{uint8(math.Round(iR)), uint8(math.Round(iG)), uint8(math.Round(iB)), uint8(math.Round(iA))}
			}
			res.Set(x, y, iC)
		}
	}
	return res
}

// https://en.wikipedia.org/wiki/Cubic_Hermite_spline#Interpolation_on_the_unit_interval_with_matched_derivatives_at_endpoints
// p1: p_n-1
// p2: p_n
// p3: p_n+1
// p4: p_n+2
func catmullRomSpline(u, p1, p2, p3, p4 float64) float64 {
	u2 := u * u
	u3 := u2 * u

	term1 := (-p1 + 3*p2 - 3*p3 + p4) * u3
	term2 := (2*p1 - 5*p2 + 4*p3 - p4) * u2
	term3 := (-p1 + p3) * u
	term4 := 2 * p2

	return 0.5 * (term1 + term2 + term3 + term4)
}

func Bicubic(src *image.RGBA, w, h int) (*image.RGBA, error) {
	srcRect := src.Bounds()
	srcW := srcRect.Dx()
	srcH := srcRect.Dy()

	// we need at least 4 points to do bicubic interpolation
	if srcW < 4 || srcH < 4 {
		return nil, fmt.Errorf("src image should be larger than or equal to 4x4 but passed src has width: %d, height: %d", srcW, srcH)
	}

	scaleW := float64(w) / float64(srcW)
	scaleH := float64(h) / float64(srcH)

	res := image.NewRGBA(image.Rect(0, 0, w, h))

	for y := range h {
		for x := range w {
			// transformed x and y
			transX, transY := transformCoords(x, y, scaleW, scaleH, true)

			// boundary check
			outX := transX < 1 || transX > float64(srcW-2)
			outY := transY < 1 || transY > float64(srcH-2)

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

				if transX < 0.5 {
					nX = 0
				} else if transX < 1 {
					nX = 1
				} else if transX <= float64(srcW)-1.5 {
					nX = srcW - 2
				} else {
					nX = srcW - 1
				}

				if transY < 0.5 {
					nY = 0
				} else if transY < 1 {
					nY = 1
				} else if transY <= float64(srcH)-1.5 {
					nY = srcH - 2
				} else {
					nY = srcH - 1
				}

				iC = src.RGBAAt(nX, nY)
			} else if outX { // use only y-axis
				var nX int

				if transX < 0.5 {
					nX = 0
				} else if transX < 1 {
					nX = 1
				} else if transX <= float64(srcW)-1.5 {
					nX = srcW - 2
				} else {
					nX = srcW - 1
				}

				floorY := math.Floor(transY)
				fractionY := transY - floorY

				intY := int(floorY)

				var pR, pG, pB, pA [4]float64

				for i := range 4 {
					pR[i] = float64(src.RGBAAt(nX, intY-1+i).R)
					pG[i] = float64(src.RGBAAt(nX, intY-1+i).G)
					pB[i] = float64(src.RGBAAt(nX, intY-1+i).B)
					pA[i] = float64(src.RGBAAt(nX, intY-1+i).A)
				}

				iR := uint8(math.Round(catmullRomSpline(fractionY, pR[0], pR[1], pR[2], pR[3])))
				iG := uint8(math.Round(catmullRomSpline(fractionY, pG[0], pG[1], pG[2], pG[3])))
				iB := uint8(math.Round(catmullRomSpline(fractionY, pB[0], pB[1], pB[2], pB[3])))
				iA := uint8(math.Round(catmullRomSpline(fractionY, pA[0], pA[1], pA[2], pA[3])))

				iC = color.RGBA{iR, iG, iB, iA}
			} else if outY { // use only x-axis
				var nY int

				if transY < 0.5 {
					nY = 0
				} else if transY < 1 {
					nY = 1
				} else if transY <= float64(srcH)-1.5 {
					nY = srcH - 2
				} else {
					nY = srcH - 1
				}

				floorX := math.Floor(transX)
				fractionX := transX - floorX

				intX := int(floorX)

				var pR, pG, pB, pA [4]float64

				for i := range 4 {
					pR[i] = float64(src.RGBAAt(intX-1+i, nY).R)
					pG[i] = float64(src.RGBAAt(intX-1+i, nY).G)
					pB[i] = float64(src.RGBAAt(intX-1+i, nY).B)
					pA[i] = float64(src.RGBAAt(intX-1+i, nY).A)
				}

				iR := uint8(math.Round(catmullRomSpline(fractionX, pR[0], pR[1], pR[2], pR[3])))
				iG := uint8(math.Round(catmullRomSpline(fractionX, pG[0], pG[1], pG[2], pG[3])))
				iB := uint8(math.Round(catmullRomSpline(fractionX, pB[0], pB[1], pB[2], pB[3])))
				iA := uint8(math.Round(catmullRomSpline(fractionX, pA[0], pA[1], pA[2], pA[3])))

				iC = color.RGBA{iR, iG, iB, iA}
			} else { // use both two axes, x first y later
				floorX := math.Floor(transX)
				fractionX := transX - floorX

				intX := int(floorX)

				floorY := math.Floor(transY)
				fractionY := transY - floorY

				intY := int(floorY)

				var tmpR, tmpG, tmpB, tmpA [4]float64

				for i := range 4 {
					var pR, pG, pB, pA [4]float64

					for j := range 4 {
						pR[j] = float64(src.RGBAAt(intX-1+j, intY-1+i).R)
						pG[j] = float64(src.RGBAAt(intX-1+j, intY-1+i).G)
						pB[j] = float64(src.RGBAAt(intX-1+j, intY-1+i).B)
						pA[j] = float64(src.RGBAAt(intX-1+j, intY-1+i).A)
					}

					tmpR[i] = catmullRomSpline(fractionX, pR[0], pR[1], pR[2], pR[3])
					tmpG[i] = catmullRomSpline(fractionX, pG[0], pG[1], pG[2], pG[3])
					tmpB[i] = catmullRomSpline(fractionX, pB[0], pB[1], pB[2], pB[3])
					tmpA[i] = catmullRomSpline(fractionX, pA[0], pA[1], pA[2], pA[3])
				}

				iR := uint8(math.Round(catmullRomSpline(fractionY, tmpR[0], tmpR[1], tmpR[2], tmpR[3])))
				iG := uint8(math.Round(catmullRomSpline(fractionY, tmpG[0], tmpG[1], tmpG[2], tmpG[3])))
				iB := uint8(math.Round(catmullRomSpline(fractionY, tmpB[0], tmpB[1], tmpB[2], tmpB[3])))
				iA := uint8(math.Round(catmullRomSpline(fractionY, tmpA[0], tmpA[1], tmpA[2], tmpA[3])))

				iC = color.RGBA{iR, iG, iB, iA}
			}
			res.Set(x, y, iC)
		}
	}
	return res, nil
}
