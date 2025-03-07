package plot

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type scale struct {
	target rl.Rectangle
	l      limits
}

func newScale(l limits, target rl.Rectangle) scale {
	return scale{
		target: target,
		l:      l,
	}
}

type limits struct {
	minX, maxX float32
	minY, maxY float32
}

func minmax(ps []rl.Vector2) limits {
	minX, maxX := math.MaxFloat32, -math.MaxFloat32
	minY, maxY := math.MaxFloat32, -math.MaxFloat32

	for _, p := range ps {
		minX = min(minX, float64(p.X))
		maxX = max(maxX, float64(p.X))

		minY = min(minY, float64(p.Y))
		maxY = max(maxY, float64(p.Y))
	}

	return limits{
		minX: float32(minX),
		maxX: float32(maxX),
		minY: float32(minY),
		maxY: float32(maxY),
	}
}

func (s scale) transform(pos rl.Vector2) rl.Vector2 {
	sp := pos
	sp.X = inter(sp.X, s.l.minX, s.l.maxX, 0, s.target.Width)
	sp.Y = inter(sp.Y, s.l.minY, s.l.maxY, s.target.Height/2, -s.target.Height/2)
	sp = rl.Vector2Add(sp, rl.NewVector2(s.target.X, s.target.Y))
	sp.Y += s.target.Height / 2

	return sp
}
