package plot

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/rickardenglund/draw/theme"
)

// Hints that the data goes on outside the plot area: a fade along the top
// or bottom edge where points are above or below it, and a chevron at a side
// with points beyond it.
const (
	fadeDepth  = float32(8)
	hintAlpha  = float32(.35)
	chevronLen = float32(6)
)

// drawEdgeFades shades the top or bottom edge over each run of consecutive
// points (in screen coordinates, sorted by x) that are above or below target.
func drawEdgeFades(ps []rl.Vector2, target rl.Rectangle) {
	left, right := target.X, target.X+target.Width
	top, bottom := target.Y, target.Y+target.Height
	side := func(p rl.Vector2) int {
		switch {
		case p.X < left || p.X > right:
			return 0
		case p.Y < top:
			return -1
		case p.Y > bottom:
			return 1
		}
		return 0
	}
	half := float32(2)
	if len(ps) > 1 {
		half = max(half, (ps[1].X-ps[0].X)/2)
	}
	for i := 0; i < len(ps); {
		sd := side(ps[i])
		if sd == 0 {
			i++
			continue
		}
		j := i
		for j+1 < len(ps) && side(ps[j+1]) == sd {
			j++
		}
		x0 := max(left, ps[i].X-half)
		x1 := min(right, ps[j].X+half)
		edge := rl.ColorAlpha(theme.Charcoal, hintAlpha)
		clear := rl.ColorAlpha(theme.Charcoal, 0)
		if sd < 0 {
			rl.DrawRectangleGradientV(int32(x0), int32(top), int32(x1-x0+1), int32(fadeDepth), edge, clear)
		} else {
			rl.DrawRectangleGradientV(int32(x0), int32(bottom-fadeDepth), int32(x1-x0+1), int32(fadeDepth), clear, edge)
		}
		i = j + 1
	}
}

// drawSideHints draws a faint chevron at the left or right of target when
// points (in screen coordinates, sorted by x) lie beyond that side.
func drawSideHints(ps []rl.Vector2, target rl.Rectangle) {
	if len(ps) == 0 {
		return
	}
	clr := rl.ColorAlpha(theme.Charcoal, hintAlpha)
	midY := target.Y + target.Height/2
	chevron := func(tipX, dir float32) {
		tip := rl.NewVector2(tipX, midY)
		back := tipX + dir*chevronLen
		rl.DrawLineEx(tip, rl.NewVector2(back, midY-chevronLen), 2, clr)
		rl.DrawLineEx(tip, rl.NewVector2(back, midY+chevronLen), 2, clr)
	}
	if ps[0].X < target.X {
		chevron(target.X+4, 1)
	}
	if ps[len(ps)-1].X > target.X+target.Width {
		chevron(target.X+target.Width-4, -1)
	}
}
