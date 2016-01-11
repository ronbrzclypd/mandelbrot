package mandelbrot

import (
	"image"
	"image/color"
	"math/cmplx"
)

// View is what you are looking at, the rectangle
// within these two complex numbers in the plane
type Range struct {
	Max  complex128
	Min  complex128
	Iter int
}

// width and height of picture
type Frame struct {
	Width  int
	Height int
}

// inMandelbrotSet is the logic of the mandelbrot set. This checks if a complex
// number is in the mandelbrot set or not, and if it bails returns how many
// iterations it took to do so
func inMandelbrotSet(n complex128, maxIteration int) (bool, int) {
	var start complex128
	i := 0
	for i < maxIteration && cmplx.Abs(start) < 2 {
		start = start*start + n
		i++
	}
	if i == maxIteration {
		return true, 0
	}
	return false, i

}

// Converts pixel at x,y coordinate to corresponding complex number
func pixelToCmplx(x, y int, r Range, f Frame) complex128 {
	minX := real(r.Min)
	minY := imag(r.Min)

	maxX := real(r.Max)
	maxY := imag(r.Max)

	pixelWidth := (maxX - minX) / float64(f.Width)
	pixelHeight := (maxY - minY) / float64(f.Height)

	resultX := float64(x)*pixelWidth + minX
	resultY := maxY - float64(y)*pixelHeight //how pixel coordinates are

	return complex(resultX, resultY)
}

// Creates a JPEG mandelbrot image
func MandelImage(r Range, f Frame) image.Image {
	m := image.NewRGBA(image.Rect(0, 0, f.Width, f.Height))
	ch := make(chan struct{}, 4)
	for x := 0; x < f.Width; x++ {
		go setMandelImageCol(x, m, ch, r, f)
	}
	for count := 0; count < f.Width; count++ {
		<-ch
	}
	return m
}

// helper function to mandelImage function. Sets all coordinates in
// given image in column x
func setMandelImageCol(x int, m *image.RGBA, ch chan struct{}, r Range, f Frame) {
	maxIteration := r.Iter
	for y := 0; y < f.Height; y++ {
		num := pixelToCmplx(x, y, r, f)
		inSet, i := inMandelbrotSet(num, maxIteration)
		if inSet {
			m.Set(x, y, color.Black)
		} else {
			bComp := i % 255
			rComp := 255 - bComp
			colr := color.RGBA{R: uint8(rComp), G: 0, B: uint8(bComp), A: 0}
			m.Set(x, y, colr)
		}
	}
	ch <- struct{}{}

}
