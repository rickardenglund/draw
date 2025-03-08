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

func (l limits) ZoomedX() limits {
	return limits{
		minX: l.minX * 1.1,
		maxX: l.maxX * .9,
		minY: l.minY * 1.1,
		maxY: l.maxY * .9,
	}
}

func (l limits) ZoomedXOut() limits {
	return limits{
		minX: l.minX * .9,
		maxX: l.maxX * 1.1,
		minY: l.minY * .9,
		maxY: l.maxY * 1.1,
	}
}

func (l limits) CenterZoomed(f float32) limits {
	dx := l.maxX - l.minX
	newDx := dx * (1 + f)
	extraX := newDx - dx

	dy := l.maxY - l.minY
	newDy := dy * (1 + f)
	extraY := newDy - dy

	return limits{
		minX: l.minX - extraX/2,
		maxX: l.maxX + extraX/2,
		minY: l.minY - extraY/2,
		maxY: l.maxY + extraY/2,
	}
}
func (l limits) Zoomed(f float32, s scale) limits {
	mp := rl.GetMousePosition()
	v := s.transformR(mp)

	dx := l.maxX - l.minX
	newDx := dx * (1 + f)
	extraX := newDx - dx
	fx := (v.X - l.minX) / (l.maxX - l.minX)

	dy := l.maxY - l.minY
	newDy := dy * (1 + f)
	extraY := newDy - dy
	fy := (v.Y - l.minY) / (l.maxY - l.minY)

	return limits{
		minX: l.minX - extraX*fx,
		maxX: l.maxX + extraX*(1-fx),
		minY: l.minY - extraY*fy,
		maxY: l.maxY + extraY*(1-fy),
	}
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

func (s scale) transformR(sp rl.Vector2) rl.Vector2 {
	pos := rl.Vector2Subtract(sp, rl.NewVector2(s.target.X, s.target.Y))

	pos.X = inter(pos.X, 0, s.target.Width, s.l.minX, s.l.maxX)
	pos.Y = inter(pos.Y, 0, s.target.Height, s.l.maxY, s.l.minY)

	return pos
}
