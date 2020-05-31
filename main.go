package main

import (
	"github.com/secretworry/crossing_shader/scenes/rock_garden"
	"image/png"
	"os"
)

func main() {
	scene := rock_garden.New(8)
	img := scene.Render()
	f, err := os.OpenFile("test.png", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}
	defer func() { _ = f.Close() }()
	if err := png.Encode(f, img); err != nil {
		panic(err)
	}
}
