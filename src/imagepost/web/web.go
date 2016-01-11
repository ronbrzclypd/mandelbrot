package web

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	"image"
	"image/jpeg"
	"imagepost/mandelbrot"
	"net/http"
	"net/url"
	"strconv"
)

// url parameters will be of the form
// "/?centerX=...&centerY=...&sizeY=...&sizeX=...&height=...&width=..."
func getURLparams(u *url.URL) (r mandelbrot.Range, f mandelbrot.Frame) {
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

	r = mandelbrot.Range{Max: max, Min: min, Iter: 1000}
	f = mandelbrot.Frame{Width: width, Height: height}

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
	m := mandelbrot.MandelImage(ran, f)
	err := writeImageWithTemplate(w, &m)
	if err != nil {
		return
	}
}

func Server() {
	http.HandleFunc("/", handler)
}
