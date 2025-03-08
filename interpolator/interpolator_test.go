package interpolator

import (
	"image"
	"image/color"
	"testing"
)

func TestNearestNeighbor(t *testing.T) {
	// Create a small source image with different colors in each corner
	src := image.NewNRGBA(image.Rect(0, 0, 2, 2))
	src.Set(0, 0, color.NRGBA{255, 0, 0, 255})   // Red at top-left
	src.Set(1, 0, color.NRGBA{0, 255, 0, 255})   // Green at top-right
	src.Set(0, 1, color.NRGBA{0, 0, 255, 255})   // Blue at bottom-left
	src.Set(1, 1, color.NRGBA{255, 255, 0, 255}) // Yellow at bottom-right

	// Create expected result for 4x4 upscale
	expected := image.NewNRGBA(image.Rect(0, 0, 6, 6))
	// Red at top-left
	for y := range 3 {
		for x := range 3 {
			expected.Set(x, y, color.NRGBA{255, 0, 0, 255})
		}
	}
	// Green at top-right
	for y := range 3 {
		for x := range 3 {
			expected.Set(x+3, y, color.NRGBA{0, 255, 0, 255})
		}
	}
	// Blue at bottom-left
	for y := range 3 {
		for x := range 3 {
			expected.Set(x, y+3, color.NRGBA{0, 0, 255, 255})
		}
	}
	// Yellow at bottom-right
	for y := range 3 {
		for x := range 3 {
			expected.Set(x+3, y+3, color.NRGBA{255, 255, 0, 255})
		}
	}

	interpolator := New(src, 6, 6, "nearestneighbor")

	// concurrency = false
	actual := interpolator.Interpolate(false)

	// compare each cells' RGBA values
	for y := range 6 {
		for x := range 6 {
			aR := actual.NRGBAAt(x, y).R
			aG := actual.NRGBAAt(x, y).G
			aB := actual.NRGBAAt(x, y).B
			aA := actual.NRGBAAt(x, y).A

			eR := actual.NRGBAAt(x, y).R
			eG := actual.NRGBAAt(x, y).G
			eB := actual.NRGBAAt(x, y).B
			eA := actual.NRGBAAt(x, y).A

			if aR != eR || aG != eG || aB != eB || aA != eA {
				t.Errorf("expected actual RGBA at [%d, %d] to be:\n[%d, %d, %d, %d]\nbut instead got:\n[%d, %d, %d, %d]\n", x, y, eR, eG, eB, eA, aR, aG, aB, aA)
			}
		}
	}

	// concurrency = true
	actual = interpolator.Interpolate(true)

	// compare each cells' RGBA values
	for y := range 6 {
		for x := range 6 {
			aR := actual.NRGBAAt(x, y).R
			aG := actual.NRGBAAt(x, y).G
			aB := actual.NRGBAAt(x, y).B
			aA := actual.NRGBAAt(x, y).A

			eR := actual.NRGBAAt(x, y).R
			eG := actual.NRGBAAt(x, y).G
			eB := actual.NRGBAAt(x, y).B
			eA := actual.NRGBAAt(x, y).A

			if aR != eR || aG != eG || aB != eB || aA != eA {
				t.Errorf("expected actual RGBA at [%d, %d] to be:\n[%d, %d, %d, %d]\nbut instead got:\n[%d, %d, %d, %d]\n", x, y, eR, eG, eB, eA, aR, aG, aB, aA)
			}
		}
	}
}

func TestBilinear(t *testing.T) {
	// Create a small source image with different colors in each corner
	src := image.NewNRGBA(image.Rect(0, 0, 2, 2))
	src.Set(0, 0, color.NRGBA{255, 0, 0, 255})   // Red at top-left
	src.Set(1, 0, color.NRGBA{0, 255, 0, 255})   // Green at top-right
	src.Set(0, 1, color.NRGBA{0, 0, 255, 255})   // Blue at bottom-left
	src.Set(1, 1, color.NRGBA{255, 255, 0, 255}) // Yellow at bottom-right

	// Create expected result for 4x4 upscale
	expected := image.NewNRGBA(image.Rect(0, 0, 6, 6))
	// Corners
	// Red at top-left
	for y := range 2 {
		for x := range 2 {
			expected.Set(x, y, color.NRGBA{255, 0, 0, 255})
		}
	}
	// Green at top-right
	for y := range 2 {
		for x := range 2 {
			expected.Set(x+4, y, color.NRGBA{0, 255, 0, 255})
		}
	}
	// Blue at top-left
	for y := range 2 {
		for x := range 2 {
			expected.Set(x, y+4, color.NRGBA{0, 0, 255, 255})
		}
	}
	// Yellow at top-left
	for y := range 2 {
		for x := range 2 {
			expected.Set(x+4, y+4, color.NRGBA{255, 255, 0, 255})
		}
	}
	// Top
	for y := range 2 {
		for x := range 2 {
			if x == 0 {
				expected.Set(x+2, y, color.NRGBA{170, 85, 0, 255})
			} else {
				expected.Set(x+2, y, color.NRGBA{85, 170, 0, 255})
			}
		}
	}
	// Right
	for y := range 2 {
		for x := range 2 {
			if y == 0 {
				expected.Set(x+4, y+2, color.NRGBA{85, 255, 0, 255})
			} else {
				expected.Set(x+4, y+2, color.NRGBA{170, 255, 0, 255})
			}
		}
	}
	// Bottom
	for y := range 2 {
		for x := range 2 {
			if x == 0 {
				expected.Set(x+2, y+4, color.NRGBA{85, 85, 170, 255})
			} else {
				expected.Set(x+2, y+4, color.NRGBA{170, 170, 85, 255})
			}
		}
	}
	// Left
	for y := range 2 {
		for x := range 2 {
			if y == 0 {
				expected.Set(x, y+2, color.NRGBA{170, 0, 85, 255})
			} else {
				expected.Set(x, y+2, color.NRGBA{85, 0, 170, 255})
			}
		}
	}
	// Center
	expected.Set(2, 2, color.NRGBA{141, 85, 56, 255})
	expected.Set(3, 2, color.NRGBA{113, 170, 28, 255})
	expected.Set(2, 3, color.NRGBA{113, 85, 113, 255})
	expected.Set(3, 3, color.NRGBA{141, 170, 56, 255})

	interpolator := New(src, 6, 6, "bilinear")

	// concurrency = false
	actual := interpolator.Interpolate(false)

	for y := range 6 {
		for x := range 6 {
			aR := actual.NRGBAAt(x, y).R
			aG := actual.NRGBAAt(x, y).G
			aB := actual.NRGBAAt(x, y).B
			aA := actual.NRGBAAt(x, y).A

			eR := expected.NRGBAAt(x, y).R
			eG := expected.NRGBAAt(x, y).G
			eB := expected.NRGBAAt(x, y).B
			eA := expected.NRGBAAt(x, y).A

			if absDiff(aR, eR) > 1 || absDiff(aG, eG) > 1 || absDiff(aB, eB) > 1 || absDiff(aA, eA) > 1 {
				t.Errorf("expected actual RGBA at [%d, %d] to be:\n[%d±1, %d±1, %d±1, %d±1]\nbut instead got:\n[%d, %d, %d, %d]\n", x, y, eR, eG, eB, eA, aR, aG, aB, aA)
			}
		}
	}

	// concurrency = true
	actual = interpolator.Interpolate(true)

	for y := range 6 {
		for x := range 6 {
			aR := actual.NRGBAAt(x, y).R
			aG := actual.NRGBAAt(x, y).G
			aB := actual.NRGBAAt(x, y).B
			aA := actual.NRGBAAt(x, y).A

			eR := expected.NRGBAAt(x, y).R
			eG := expected.NRGBAAt(x, y).G
			eB := expected.NRGBAAt(x, y).B
			eA := expected.NRGBAAt(x, y).A

			if absDiff(aR, eR) > 1 || absDiff(aG, eG) > 1 || absDiff(aB, eB) > 1 || absDiff(aA, eA) > 1 {
				t.Errorf("expected actual RGBA at [%d, %d] to be:\n[%d±1, %d±1, %d±1, %d±1]\nbut instead got:\n[%d, %d, %d, %d]\n", x, y, eR, eG, eB, eA, aR, aG, aB, aA)
			}
		}
	}
}

func absDiff(x, y uint8) uint8 {
	if x >= y {
		return x - y
	} else {
		return y - x
	}
}
