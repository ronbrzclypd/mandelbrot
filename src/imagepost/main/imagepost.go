package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	"image"
	"image/color"
	"image/jpeg"
	"math/cmplx"
	"net/http"
	"net/url"
	"strconv"
)

// View is what you are looking at, the rectangle
// within these two complex numbers in the plane
type Range struct {
	max  complex128
	min  complex128
	iter int
}

// width and height of picture
type Frame struct {
	width  int
	height int
}

// numIterate is the logic of the mandelbrot set. This checks if a complex
// number is in the mandelbrot set or not, and if it bails returns how many
// iterations it took to do so
func numIterate(n complex128, maxIteration int) (numIter int, bailp bool) {
	var start complex128
	i := 0
	for i < maxIteration && cmplx.Abs(start) < 2 {
		start = start*start + n
		i++
	}
	if i == maxIteration {
		return 0, true
	}
	return i, false

}

// Converts pixel at x,y coordinate to corresponding complex number
func pixelToCmplx(r Range, f Frame, x, y int) complex128 {
	minX := real(r.min)
	minY := imag(r.min)

	maxX := real(r.max)
	maxY := imag(r.max)

	pixelWidth := (maxX - minX) / float64(f.width)
	pixelHeight := (maxY - minY) / float64(f.height)

	resultX := float64(x)*pixelWidth + minX
	resultY := maxY - float64(y)*pixelHeight //how pixel coordinates are

	return complex(resultX, resultY)
}

// Creates a JPEG mandelbrot image
func mandelImage(r Range, f Frame) image.Image {
	m := image.NewRGBA(image.Rect(0, 0, f.width, f.height))
	ch := make(chan struct{}, 4)
	for x := 0; x < f.width; x++ {
		go setMandelImageCol(x, m, ch, r, f)
	}
	for count := 0; count < f.width; count++ {
		<-ch
	}
	return m
}

// helper function to mandelImage function. Sets all coordinates in
// given image in column x
func setMandelImageCol(x int, m *image.RGBA, ch chan struct{}, r Range, f Frame) {
	maxIteration := r.iter
	for y := 0; y < f.height; y++ {
		num := pixelToCmplx(r, f, x, y)
		i, inSet := numIterate(num, maxIteration)
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

// url parameters will be of the form
// "/?centerX=...&centerY=...&sizeY=...&sizeX=...&height=...&width=..."
func getURLparams(u *url.URL) (r Range, f Frame) {
	m := u.Query()
	var maxX, maxY, minX, minY float64
	var width, height int
	var max, min complex128

	allhere := len(m) > 5
	if !allhere {
		max = complex(1, 1)
		min = complex(-2, -1)
		width = 600
		height = 400
	} else { // Is there a better way to do this?
		mFloat := map[string]float64{}
		mInt := map[string]int{}
		fmt.Println(m)
		mFloat["centerX"], _ = strconv.ParseFloat(m["centerX"][0], 64)
		mFloat["centerY"], _ = strconv.ParseFloat(m["centerY"][0], 64)
		mFloat["sizeX"], _ = strconv.ParseFloat(m["sizeX"][0], 64)
		mFloat["sizeY"], _ = strconv.ParseFloat(m["sizeY"][0], 64)

		gg, _ := strconv.ParseInt(m["width"][0], 0, 0)
		gg2, _ := strconv.ParseInt(m["height"][0], 10, 0)
		mInt["width"] = int(gg)
		mInt["height"] = int(gg2)

		maxX = mFloat["centerX"] + mFloat["sizeX"]/2
		maxY = mFloat["centerY"] + mFloat["sizeY"]/2

		minX = mFloat["centerX"] - mFloat["sizeY"]/2
		minY = mFloat["centerY"] - mFloat["sizeY"]/2

		max = complex(maxX, maxY)
		min = complex(minX, minY)

		width = mInt["width"]
		height = mInt["height"]
	}

	r = Range{max: max, min: min, iter: 1000}
	f = Frame{width: width, height: height}

	return r, f

}

// Writes image to template, then writes that to the responseWriter
func writeImageWithTemplate(w http.ResponseWriter, img *image.Image) (err error) {
	buffer := new(bytes.Buffer)
	jpeg.Encode(buffer, *img, nil)
	str := base64.StdEncoding.EncodeToString(buffer.Bytes())
	t, _ := template.ParseFiles("static/page.html")
	data := map[string]interface{}{"Image": str}
	err = t.Execute(w, data)
	return
}

func handler(w http.ResponseWriter, r *http.Request) {
	u := r.URL
	ran, f := getURLparams(u)
	m := mandelImage(ran, f)

	err := writeImageWithTemplate(w, &m)
	if err != nil {
		return
	}
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
