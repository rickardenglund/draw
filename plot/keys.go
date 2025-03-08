package plot

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func (p *Plot) handleKeys(target rl.Rectangle, s scale) {
	mp := rl.GetMousePosition()
	active := rl.CheckCollisionPointRec(mp, target)
	if active {
		scroll := rl.GetMouseWheelMove()
		if scroll != 0 {
			zoom := float32(.1)
			zoom *= -scroll
			p.SetLimits(p.getLimits().Zoomed(zoom, s))
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
			p.SetLimits(minmax(p.cur).Zoomed(.01, s))
		}
	} else {
		if p.selection.IsActive() {
			p.SetLimits(p.selection.GetLimits(s))
			p.selection.Reset()
		}
	}

}
