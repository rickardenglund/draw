package plot

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func (p *Plot) handleKeys(target rl.Rectangle, s scale) {
	mp := rl.GetMousePosition()
	if rl.CheckCollisionPointRec(mp, target) {
		scroll := rl.GetMouseWheelMove()
		if scroll != 0 {
			zoom := float32(.1)
			zoom *= scroll
			p.SetLimits(p.getLimits().Zoomed(zoom))
		}
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			p.selection.a = &mp
		}
		if p.selection.IsActive() {
			if rl.IsMouseButtonDown(rl.MouseButtonLeft) {
				p.selection.b = &mp
			}
			if rl.IsMouseButtonReleased(rl.MouseButtonLeft) {
				r := p.selection.GetRect()
				if max(r.Width, r.Height) > 20 {
					p.SetLimits(p.selection.GetLimits(s))
				}
				p.selection.Reset()
			}
		}

		if rl.IsKeyPressed(rl.KeyF) {
			p.SetLimits(minmax(p.cur).Zoomed(.01))
		}
	} else {
		if p.selection.IsActive() {
			p.SetLimits(p.selection.GetLimits(s))
			p.selection.Reset()
		}
	}

}

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
