package rock_garden

import (
	"github.com/golang/geo/r2"
	"github.com/golang/geo/r3"
	"image"
	"image/color"
	"math"
)

type Light struct {
	Pos   r3.Vector
	Color float64
}

type Scene struct {
	Interval float64
	Lights   []Light
}

func New(interval float64) *Scene {
	return &Scene{
		Interval: interval,
		Lights: []Light{
			{
				Pos: r3.Vector{
					X: 3 * segmentWidth,
					Y: 0,
					Z: 3 * segmentWidth,
				},
				Color: 2.0,
			},
		},
	}
}

const ambient float64 = 0.5

const segmentWidth float64 = 32

var dimen = r2.Point{
	X: segmentWidth * 3.0,
	Y: segmentWidth * 3.0,
}

func (s Scene) Norm(p r2.Point) r3.Vector {
	var n r3.Vector
	var sections = struct {
		X int
		Y int
	}{
		X: int(p.X / segmentWidth),
		Y: int(p.Y / segmentWidth),
	}
	if sections.X == 1 && sections.Y == 1 {
		n.Z = 1.0
	} else {
		var anchor r2.Point
		if sections.X == 0 && sections.Y == 0 {
			anchor.X = segmentWidth
			anchor.Y = segmentWidth
		} else if sections.X == 0 && sections.Y == 1 {
			anchor.X = segmentWidth
			anchor.Y = p.Y
		} else if sections.X == 0 && sections.Y == 2 {
			anchor.X = segmentWidth
			anchor.Y = 2 * segmentWidth
		} else if sections.X == 1 && sections.Y == 0 {
			anchor.X = p.X
			anchor.Y = segmentWidth
		} else if sections.X == 1 && sections.Y == 2 {
			anchor.X = p.X
			anchor.Y = 2 * segmentWidth
		} else if sections.X == 2 && sections.Y == 0 {
			anchor.X = 2 * segmentWidth
			anchor.Y = segmentWidth
		} else if sections.X == 2 && sections.Y == 1 {
			anchor.X = 2 * segmentWidth
			anchor.Y = p.Y
		} else if sections.X == 2 && sections.Y == 2 {
			anchor.X = 2 * segmentWidth
			anchor.Y = 2 * segmentWidth
		}
		v := p.Sub(anchor)
		d := v.Normalize()
		norm := v.Norm()
		offset := norm/s.Interval - math.Ceil(norm/s.Interval)
		// y = cos(2 * pi * offset)
		// \delta y  = -sin(2 * pi * offset) * 2 * pi
		localN := r2.Point{
			X: 1,
			Y: -math.Sin(2*math.Pi*offset) * 2 * math.Pi,
		}.Normalize().Ortho()
		n.X = d.X * localN.X
		n.Y = d.Y * localN.X
		n.Z = localN.Y
	}
	return n
}

func (s Scene) Color(p r2.Point) r3.Vector {
	return r3.Vector{
		X: 0.8,
		Y: 0.8,
		Z: 0.8,
	}
}

func (s Scene) Render() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, int(dimen.X), int(dimen.Y)))
	for i := 0; i < int(dimen.X); i++ {
		for j := 0; j < int(dimen.X); j++ {
			p := r2.Point{
				X: float64(i),
				Y: float64(j),
			}
			v := r3.Vector{
				X: p.X,
				Y: p.Y,
				Z: 0,
			}
			n := s.Norm(p)
			d := 0.0
			for _, l := range s.Lights {
				lightDir := l.Pos.Sub(v).Normalize()
				diff := max(0, n.Dot(lightDir))
				diffuse := diff * l.Color
				d += diffuse
			}
			c := s.Color(p)
			img.SetRGBA(i, j, vector2Color(c.Mul(d+ambient)))
		}
	}
	return img
}

func vector2Color(vector r3.Vector) color.RGBA {
	return color.RGBA{
		R: uint8(max(0, min(255, 255.0*vector.X))),
		G: uint8(max(0, min(255, 255.0*vector.Y))),
		B: uint8(max(0, min(255, 255.0*vector.Z))),
		A: 255,
	}
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
