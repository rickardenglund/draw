package plot

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/rickardenglund/draw/theme"
)

// stemHeadSpacing is how far apart (in pixels) stems must be to get a dot
// on top; closer ones would merge into a bar. Points on the baseline get
// none.
const stemHeadSpacing = 6

// SetStems draws every point as a vertical line from y = 0 instead of
// joining the points, as suits a spectrum. The y range then always
// includes 0.
func (p *Plot[point]) SetStems(on bool) {
	p.stems = on
	p.SetLimits(p.fit(p.cur).CenterZoomed(.01))
}

// xOnly keeps the current y range of a stem plot, so zooming there only
// changes x: the stems stay anchored at 0 with their full heights in view.
func (p *Plot[point]) xOnly(ls limits) limits {
	if p.stems {
		ls.minY, ls.maxY = p.curLimits.minY, p.curLimits.maxY
	}
	return ls
}

// selectionRect is the dragged zoom area; in a stem plot it spans the
// whole data height, as only x is zoomed.
func (p *Plot[point]) selectionRect(target rl.Rectangle) rl.Rectangle {
	r := p.selection.GetRect()
	if p.stems {
		r.Y, r.Height = target.Y, target.Height
	}
	return r
}

// fit is the range showing all of ps.
func (p *Plot[point]) fit(ps []point) limits {
	ls := minmax(ps)
	if p.stems && len(ps) > 0 {
		ls.minY = min(ls.minY, 0)
		ls.maxY = max(ls.maxY, 0)
	}
	return ls
}

// drawStems draws ps, already in screen coordinates and sorted by x, as
// stems from the baseline y = 0. Where several points share a pixel column
// only the one farthest from the baseline is drawn, so peaks are kept
// however many points there are.
func drawStems(ps []rl.Vector2, s scale, target rl.Rectangle) {
	base := clamp(s.transform(rl.NewVector2(0, 0)), target).Y
	stem := func(x, y float32, th float32) {
		top := clamp(rl.NewVector2(x, y), target)
		rl.DrawLineEx(rl.NewVector2(x, base), top, th, theme.Charcoal)
	}
	left, right := target.X, target.X+target.Width

	spacing := float32(0)
	if len(ps) > 1 {
		spacing = (ps[len(ps)-1].X - ps[0].X) / float32(len(ps)-1)
	}
	if len(ps) == 1 || spacing >= 2 {
		th := min(2, spacing/2)
		if len(ps) == 1 {
			th = 2
		}
		for _, sp := range ps {
			if sp.X < left || sp.X > right {
				continue
			}
			stem(sp.X, sp.Y, th)
			if (len(ps) == 1 || spacing >= stemHeadSpacing) && abs(sp.Y-base) >= 1 && rl.CheckCollisionPointRec(sp, target) {
				rl.DrawCircleV(sp, 3, theme.Charcoal)
			}
		}
		return
	}

	col, y, have := 0, float32(0), false
	flush := func() {
		if have {
			stem(float32(col)+.5, y, 1)
		}
	}
	for _, sp := range ps {
		if sp.X < left || sp.X > right {
			continue
		}
		c := int(sp.X)
		if !have || c != col {
			flush()
			col, y, have = c, sp.Y, true
			continue
		}
		if abs(sp.Y-base) > abs(y-base) {
			y = sp.Y
		}
	}
	flush()
}
