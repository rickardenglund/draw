package plot

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type scale struct {
	minX, minY, maxX, maxY float32
	target                 rl.Rectangle
}

func newScale(ps []rl.Vector2, target rl.Rectangle) scale {
	minX, maxX, minY, maxY := minmax(ps)
	return scale{
		target: target,
		minX:   minX,
		minY:   minY,
		maxX:   maxX,
		maxY:   maxY,
	}
}

func minmax(ps []rl.Vector2) (float32, float32, float32, float32) {
	minX, maxX := math.MaxFloat32, -math.MaxFloat32
	minY, maxY := 0.0, 0.0

	for _, p := range ps {
		minX = math.Min(minX, float64(p.X))
		maxX = math.Max(minX, float64(p.X))

		minY = math.Min(minY, float64(p.Y))
		maxY = math.Max(maxY, float64(p.Y))
	}

	return float32(minX), float32(maxX), float32(minY), float32(maxY)
}

func (s scale) transform(pos rl.Vector2) rl.Vector2 {
	sp := pos
	sp.X = inter(sp.X, s.minX, s.maxX, 0, s.target.Width)
	sp.Y = inter(sp.Y, s.minY, s.maxY, s.target.Height/2, -s.target.Height/2)
	sp = rl.Vector2Add(sp, rl.NewVector2(s.target.X, s.target.Y))
	sp.Y += s.target.Height / 2

	return sp
}
