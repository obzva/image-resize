package interpolation

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"runtime"
	"time"
)

func timeTrack(start time.Time, funcName string) {
	elapsed := time.Since(start)
	fmt.Printf("%s interpolation took %v to run\n", funcName, elapsed)
}

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

func NearestNeighbor(src *image.RGBA, w, h int, concurrency bool) *image.RGBA {
	defer timeTrack(time.Now(), "NearestNeighbor")

	srcRect := src.Bounds()
	srcW, srcH := srcRect.Dx(), srcRect.Dy()

	scaleW := float64(w) / float64(srcW)
	scaleH := float64(h) / float64(srcH)

	res := image.NewRGBA(image.Rect(0, 0, w, h))

	operate := func(start, end int) {
		for ; start < end; start++ {
			x := start % w
			y := start / w
			// transformed x and y
			transX, transY := transformCoords(x, y, scaleW, scaleH, false)

			res.Set(x, y, src.At(int(transX), int(transY)))
		}
	}
	
	if concurrency {
		numCPU := runtime.NumCPU()
		c := make(chan int, numCPU)
		for i := range numCPU {
			go func() {
				operate(i*w*h/numCPU, (i+1)*w*h/numCPU)
				c <- 1
			}()
		}
		for i := 0; i < numCPU; i++ {
			<-c
		}
	} else {
		operate(0, w*h)
	}

	return res
}

// p*: two surrounding points' * values
// nv: integer part of v
func internalDivision(pR, pG, pB, pA *[2]float64, nv, v float64) (r float64, g float64, b float64, a float64) {
	r = (nv+1-v)*float64(pR[0]) + (v-nv)*float64(pR[1])
	g = (nv+1-v)*float64(pG[0]) + (v-nv)*float64(pG[1])
	b = (nv+1-v)*float64(pB[0]) + (v-nv)*float64(pB[1])
	a = (nv+1-v)*float64(pA[0]) + (v-nv)*float64(pA[1])

	return r, g, b, a
}

func clamp(v float64) uint8 {
	if v > 255 { // overshoot
		return 255
	} else if v < 0 { // undershoot
		return 0
	} else {
		return uint8(math.Round(v))
	}
}

func Bilinear(src *image.RGBA, w, h int, concurrency bool) *image.RGBA {
	defer timeTrack(time.Now(), "Bilinear")

	srcRect := src.Bounds()
	srcW := srcRect.Dx()
	srcH := srcRect.Dy()

	scaleW := float64(w) / float64(srcW)
	scaleH := float64(h) / float64(srcH)

	res := image.NewRGBA(image.Rect(0, 0, w, h))

	operate := func(start, end int) {
		for ; start < end; start++ {
			x := start % w
			y := start / w

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

				nY := math.Floor(transY)

				// surrounding points' r/g/b/a values
				// p*[0] <- * value of [nX, nY] at src
				// p*[1] <- * value of [nX, nY + 1] at src
				var pR, pG, pB, pA [2]float64

				for i := range 2 {
					pRGBA := src.RGBAAt(int(nX), int(nY)+i)
					pR[i] = float64(pRGBA.R)
					pG[i] = float64(pRGBA.G)
					pB[i] = float64(pRGBA.B)
					pA[i] = float64(pRGBA.A)
				}

				iR, iG, iB, iA := internalDivision(&pR, &pG, &pB, &pA, nY, transY)
				// use round to dodge edge cases like:
				// uint8(1.99999999999) -> 1
				iC = color.RGBA{clamp(iR), clamp(iG), clamp(iB), clamp(iA)}
			} else if outY { // use two surrounding points (only x-axis)
				var nY float64

				if transY < 0 {
					nY = 0
				} else {
					nY = float64(srcH - 1)
				}

				nX := math.Floor(transX)

				// surrounding points' r/g/b/a values
				// p*[0] <- * value of [nX, nY] at src
				// p*[1] <- * value of [nX + 1, nY] at src
				var pR, pG, pB, pA [2]float64

				for i := range 2 {
					pRGBA := src.RGBAAt(int(nX)+i, int(nY))
					pR[i] = float64(pRGBA.R)
					pG[i] = float64(pRGBA.G)
					pB[i] = float64(pRGBA.B)
					pA[i] = float64(pRGBA.A)
				}

				iR, iG, iB, iA := internalDivision(&pR, &pG, &pB, &pA, nX, transX)
				// use round to dodge edge cases like:
				// uint8(1.99999999999) -> 1
				iC = color.RGBA{clamp(iR), clamp(iG), clamp(iB), clamp(iA)}
			} else { // use four surrounding points
				nX := math.Floor(transX)
				nY := math.Floor(transY)

				// surrounding points' r/g/b/a values
				// p*[0][0] <- * value of [nX, nY] at src
				// p*[0][1] <- * value of [nX + 1, nY] at src
				// p*[1][0] <- * value of [nX, nY + 1] at src
				// p*[1][1] <- * value of [nX + 1, nY + 1] at src
				var pR, pG, pB, pA [2][2]float64

				// temporarily saved value got from internalDivision in x-axis
				// tmp*[0] <- * value got from internalDivision in y = nY
				// tmp*[1] <- * value got from internalDivision in y = nY + 1
				var tmpR, tmpG, tmpB, tmpA [2]float64

				for i := range 2 {
					for j := range 2 {
						pRGBA := src.RGBAAt(int(nX)+j, int(nY)+i)
						pR[i][j] = float64(pRGBA.R)
						pG[i][j] = float64(pRGBA.G)
						pB[i][j] = float64(pRGBA.B)
						pA[i][j] = float64(pRGBA.A)
					}
					tmpR[i], tmpG[i], tmpB[i], tmpA[i] = internalDivision(&pR[i], &pG[i], &pB[i], &pA[i], nX, transX)
				}

				iR, iG, iB, iA := internalDivision(&tmpR, &tmpG, &tmpB, &tmpA, nY, transY)
				// use round to dodge edge cases like:
				// uint8(1.99999999999) -> 1
				iC = color.RGBA{clamp(iR), clamp(iG), clamp(iB), clamp(iA)}
			}
			res.Set(x, y, iC)
		}
	}

	if concurrency {
		numCPU := runtime.NumCPU()
		c := make(chan int, numCPU)
		for i := range numCPU {
			go func() {
				operate(i*w*h/numCPU, (i+1)*w*h/numCPU)
				c <- 1
			}()
		}
		for i := 0; i < numCPU; i++ {
			<-c
		}
	} else {
		operate(0, w*h)
	}

	return res
}

// https://en.wikipedia.org/wiki/Cubic_Hermite_spline#Interpolation_on_the_unit_interval_with_matched_derivatives_at_endpoints
// p[0]: p_n-1
// p[1]: p_n
// p[2]: p_n+1
// p[3]: p_n+2
func catmullRomSpline(u float64, p *[4]float64) float64 {
	u2 := u * u
	u3 := u2 * u

	term1 := (-p[0] + 3*p[1] - 3*p[2] + p[3]) * u3
	term2 := (2*p[0] - 5*p[1] + 4*p[2] - p[3]) * u2
	term3 := (-p[0] + p[2]) * u
	term4 := 2 * p[1]

	return 0.5 * (term1 + term2 + term3 + term4)
}

func Bicubic(src *image.RGBA, w, h int, concurrency bool) (*image.RGBA, error) {
	defer timeTrack(time.Now(), "Bicubic")

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

	operate := func(start, end int) {
		for ; start < end; start++ {
			x := start % w
			y := start / w

			// transformed x and y
			transX, transY := transformCoords(x, y, scaleW, scaleH, true)

			// boundary check
			outX := transX < 1 || transX > float64(srcW-2)
			outY := transY < 1 || transY > float64(srcH-2)

			var iC color.RGBA

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

				iR := clamp(catmullRomSpline(fractionY, &pR))
				iG := clamp(catmullRomSpline(fractionY, &pG))
				iB := clamp(catmullRomSpline(fractionY, &pB))
				iA := clamp(catmullRomSpline(fractionY, &pA))

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

				iR := clamp(catmullRomSpline(fractionX, &pR))
				iG := clamp(catmullRomSpline(fractionX, &pG))
				iB := clamp(catmullRomSpline(fractionX, &pB))
				iA := clamp(catmullRomSpline(fractionX, &pA))

				iC = color.RGBA{iR, iG, iB, iA}
			} else { // use both two axes, x first y later
				floorX := math.Floor(transX)
				fractionX := transX - floorX

				intX := int(floorX)

				floorY := math.Floor(transY)
				fractionY := transY - floorY

				intY := int(floorY)

				var tmpR, tmpG, tmpB, tmpA [4]float64
				var pR, pG, pB, pA [4][4]float64

				for i := range 4 {
					for j := range 4 {
						pR[i][j] = float64(src.RGBAAt(intX-1+j, intY-1+i).R)
						pG[i][j] = float64(src.RGBAAt(intX-1+j, intY-1+i).G)
						pB[i][j] = float64(src.RGBAAt(intX-1+j, intY-1+i).B)
						pA[i][j] = float64(src.RGBAAt(intX-1+j, intY-1+i).A)
					}

					tmpR[i] = catmullRomSpline(fractionX, &pR[i])
					tmpG[i] = catmullRomSpline(fractionX, &pG[i])
					tmpB[i] = catmullRomSpline(fractionX, &pB[i])
					tmpA[i] = catmullRomSpline(fractionX, &pA[i])
				}

				iR := clamp(catmullRomSpline(fractionY, &tmpR))
				iG := clamp(catmullRomSpline(fractionY, &tmpG))
				iB := clamp(catmullRomSpline(fractionY, &tmpB))
				iA := clamp(catmullRomSpline(fractionY, &tmpA))

				iC = color.RGBA{iR, iG, iB, iA}
			}
			res.Set(x, y, iC)
		}
	}

	if concurrency {
		numCPU := runtime.NumCPU()
		c := make(chan int, numCPU)
		for i := range numCPU {
			go func() {
				operate(i*w*h/numCPU, (i+1)*w*h/numCPU)
				c <- 1
			}()
		}
		for i := 0; i < numCPU; i++ {
			<-c
		}
	} else {
		operate(0, w*h)
	}

	return res, nil
}
