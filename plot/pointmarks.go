package plot

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

// pointMarks are x positions marked on the plot, e.g. the values that
// belong to something shown elsewhere.
type pointMarks struct {
	xs     []float32
	active int
	clr    rl.Color
	// activeX is where the active line is drawn, while it moves.
	activeX float32
	// grow is how far the active ring has grown, 0 to 1, and shrink how
	// big the ring at prevX, the previously active mark, still is.
	grow, shrink float32
	prevX        float32
}

// activeMark animates a change of the emphasised mark: its line glides
// from the old x to the new, the old ring shrinks and the new one grows.
type activeMark struct {
	from, to float32
	at       float64
	set      bool
	// prev is the previously active x and prevSize how emphasised its ring
	// was when the change started.
	prev     float32
	prevSize float32
}

const activeMarkDuration = .25

// progress is how far the change has come now, eased, 0 to 1.
func (a activeMark) progress(now float64) float32 {
	f := float32(min(1, (now-a.at)/activeMarkDuration))
	return f * f * (3 - 2*f) // ease in and out
}

// x is the animated x position of the line now.
func (a activeMark) x(now float64) float32 {
	return a.from + (a.to-a.from)*a.progress(now)
}

func (a *activeMark) moveTo(x float32, now float64) {
	if !a.set {
		*a = activeMark{from: x, to: x, at: now - activeMarkDuration, set: true}
		return
	}
	if x != a.to {
		*a = activeMark{
			from: a.x(now), to: x, at: now, set: true,
			prev: a.to, prevSize: a.progress(now),
		}
	}
}

// SetPointMarks marks the data at xs in clr: a faint line across the data
// area and a ring around the point at that x, if there is one. The mark at
// xs[active] is emphasised; when it changes the line glides there and the
// rings change size; -1 for none. nil clears them.
func (p *Plot[point]) SetPointMarks(xs []float32, active int, clr rl.Color) {
	p.pointMarks = pointMarks{xs: xs, active: active, clr: clr}
	if active >= 0 && active < len(xs) {
		now := rl.GetTime()
		p.activeMark.moveTo(xs[active], now)
		f := p.activeMark.progress(now)
		p.pointMarks.activeX = p.activeMark.x(now)
		p.pointMarks.grow = f
		p.pointMarks.prevX = p.activeMark.prev
		p.pointMarks.shrink = p.activeMark.prevSize * (1 - f)
	}
}

// size is how emphasised the ring of the mark at index i is, 0 to 1.
func (m pointMarks) size(i int) float32 {
	switch {
	case i == m.active:
		return m.grow
	case m.xs[i] == m.prevX:
		return m.shrink
	}
	return 0
}

// draw draws the marks for ps, the data in screen coordinates sorted by x.
func (m pointMarks) draw(ps []rl.Vector2, s scale, target rl.Rectangle) {
	const snap = float32(1.5)
	for i, x := range m.xs {
		sx := s.transform(rl.NewVector2(x, 0)).X
		if sx < target.X-snap || sx > target.X+target.Width+snap {
			continue
		}
		if i != m.active {
			rl.DrawLineEx(
				rl.NewVector2(sx, target.Y),
				rl.NewVector2(sx, target.Y+target.Height),
				1, rl.Fade(m.clr, .3))
		}
		e := m.size(i)
		for _, sp := range ps {
			if sp.X-sx > snap {
				break
			}
			if sx-sp.X <= snap && rl.CheckCollisionPointRec(sp, target) {
				rl.DrawRing(sp, 3+2*e, 5+3*e, 0, 360, 24, m.clr)
			}
		}
	}
	if m.active < 0 || m.active >= len(m.xs) {
		return
	}
	sx := s.transform(rl.NewVector2(m.activeX, 0)).X
	if sx < target.X-snap || sx > target.X+target.Width+snap {
		return
	}
	rl.DrawLineEx(
		rl.NewVector2(sx, target.Y),
		rl.NewVector2(sx, target.Y+target.Height),
		2, rl.Fade(m.clr, .8))
}
