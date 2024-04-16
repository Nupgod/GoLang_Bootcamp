package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
)

func main() {
	// Define the dimensions of the image
	width := 300
	height := 300

	// Create a new image with the specified dimensions
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Define the colors forlogo
	Red := color.RGBA{255, 0, 128, 255}
	Blue := color.RGBA{0, 255, 255, 255}

	// Draw base
	draw.Draw(img, image.Rect(0, 0, width, height), &image.Uniform{Blue}, image.Point{}, draw.Src)
	draw.Draw(img, image.Rect(0, 0, width, height/2), &image.Uniform{Red}, image.Point{}, draw.Src)

	// Define data for drawing circle
	centerX := width / 2
	centerY := height / 2
	radius := width / 4
	drawCircle(img, centerX, centerY, radius)

	// Create a file to save the image
	file, err := os.Create("amazing_logo.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Encode the image to PNG and write it to the file
	err = png.Encode(file, img)
	if err != nil {
		panic(err)
	}
}
func drawCircle(img *image.RGBA, cx, cy, r int) {
	// Draw the perimeter in green
	drawCirclePerimeter(img, cx, cy, r, color.RGBA{0, 128, 0, 255}) // Green

	// Draw the circle with the top half purple and the bottom half yellow
	for y := -r; y <= r; y++ {
		for x := -r; x <= r; x++ {
			if x*x+y*y <= r*r {
				setPixelColor(img, cx+x, cy+y, cy)
			}
		}
	}
}

func drawCirclePerimeter(img *image.RGBA, cx, cy, r int, col color.Color) {
	x, y, dx, dy := r, 0, 1, 1
	err := dx - (r * 2)

	for x >= y {
		img.Set(cx+x, cy+y, col)
		img.Set(cx+y, cy+x, col)
		img.Set(cx-y, cy+x, col)
		img.Set(cx-x, cy+y, col)
		img.Set(cx-x, cy-y, col)
		img.Set(cx-y, cy-x, col)
		img.Set(cx+y, cy-x, col)
		img.Set(cx+x, cy-y, col)

		img.Set(cx+x+1, cy+y, col)
		img.Set(cx+y+1, cy+x, col)
		img.Set(cx-y-1, cy+x, col)
		img.Set(cx-x-1, cy+y, col)
		img.Set(cx-x-1, cy-y, col)
		img.Set(cx-y-1, cy-x, col)
		img.Set(cx+y+1, cy-x, col)
		img.Set(cx+x+1, cy-y, col)

		if err <= 0 {
			y++
			err += dy
			dy += 2
		}
		if err > 0 {
			x--
			dx += 2
			err += dx - (r * 2)
		}
	}
}

func setPixelColor(img *image.RGBA, x, y, cy int) {
	var col color.RGBA
	if y < cy {
		// Top half of the circle is purple
		col = color.RGBA{255, 255, 0, 255} // Purple
	} else if y > cy {
		// Bottom half of the circle is yellow
		col = color.RGBA{128, 0, 128, 255} // Yellow
	}
	img.Set(x, y, col)
}
