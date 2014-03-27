/*
	HTTP Server produces a PNG given a URL with 9-tuples in the path:

		x1,y1,x2,y2,x3,y3,r,g,b

	(x1, y2), (x2, y2), (x3, y3) are corners of triangle.
	Visible screen is real 0..100 on both x and y.
	r, g, & b are in 0..255.
	Background is black.
*/
package main

import (
	"flag"
	"fmt"
	// "html"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

import "github.com/strickyak/canvas"

var PORT = flag.Int("p", 8080, "port to listen on")
var WIDTH = flag.Int("w", 640, "Width of PNG in pixels")
var HEIGHT = flag.Int("h", 360, "Height of PNG in pixels")

var notnum = regexp.MustCompile("[^0-9]+")

/*
var colors = []canvas.Color{
	canvas.RGB(0, 0, 0),
	canvas.RGB(255, 255, 255),
	canvas.RGB(255, 0, 0),
	canvas.RGB(0, 255, 0),
	canvas.RGB(0, 0, 255),
}
*/

type H complex128

func (h H) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	strs := notnum.Split(r.URL.Path, -1)
	var nums []float64
	for _, s := range strs {
		if s == "" {
			continue
		}
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			panic(err)
		}
		nums = append(nums, f)
	}

	can := canvas.NewCanvas(*WIDTH, *HEIGHT)
	for i := 0; i < len(nums)-8; i += 9 {
		x1 := int(nums[i+0] / 100.0 * float64(*WIDTH))
		y1 := int(nums[i+1] / 100.0 * float64(*HEIGHT))
		x2 := int(nums[i+2] / 100.0 * float64(*WIDTH))
		y2 := int(nums[i+3] / 100.0 * float64(*HEIGHT))
		x3 := int(nums[i+4] / 100.0 * float64(*WIDTH))
		y3 := int(nums[i+5] / 100.0 * float64(*HEIGHT))

		red := byte(nums[i+6])
		green := byte(nums[i+7])
		blue := byte(nums[i+8])

		can.PaintTriangle(x1, y1, x2, y2, x3, y3, canvas.RGB(red, green, blue))
	}

	w.Header().Set("Content-Type", "image/png")
	can.WritePng(w)
}

func main() {
	var myHandler H

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", *PORT),
		Handler:        myHandler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(s.ListenAndServe())
}
