package mandelbrot

import (
	"math/cmplx"
)

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
