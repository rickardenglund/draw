package plot

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type selection struct {
	a, b *rl.Vector2
}

func (s *selection) Reset() {
	s.a = nil
	s.b = nil
}

func (s *selection) IsActive() bool {
	return s.a != nil
}

func (s *selection) GetRect() rl.Rectangle {
	minP := rl.NewVector2(min(s.a.X, s.b.X), min(s.a.Y, s.b.Y))
	maxP := rl.NewVector2(max(s.a.X, s.b.X), max(s.a.Y, s.b.Y))

	size := rl.Vector2Subtract(maxP, minP)

	r := rl.NewRectangle(minP.X, minP.Y, size.X, size.Y)
	return r
}

func (s *selection) GetLimits(sc scale) limits {
	r := s.GetRect()
	dataPA := sc.transformR(rl.NewVector2(r.X, r.Y+r.Height))
	dataPB := sc.transformR(rl.NewVector2(r.X+r.Width, r.Y))
	return limits{
		minX: dataPA.X,
		maxX: dataPB.X,
		minY: dataPA.Y,
		maxY: dataPB.Y,
	}
}
