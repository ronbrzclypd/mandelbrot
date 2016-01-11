package mandelbrot

import "testing"

type testpair struct {
	point complex128
	inset bool
}

var tests = []testpair{
	{0, true},
	{10, false},
	{complex(1, 1), false},
	{1, false},
	{-1, true},
	{complex(1, -1), false},
}

func TestinSet(t *testing.T) {
	for _, y := range tests {
		x, _ := inMandelbrotSet(y.point, 1000)
		if x != y.inset {
			t.Error("For", y.point,
				"expected", y.inset,
				"got", x)
		}
	}

}
