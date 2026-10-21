package plot

import (
	"strings"

	"math"

	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/rickardenglund/draw/theme"
)

func drawAxisX(target rl.Rectangle, s scale, ts []tick) {
	for _, t := range ts {
		xTick(target, s, t)
	}
}

func getFmt(v float32) string {
	v = float32(math.Abs(float64(v)))
	switch {
	case v < 10 && v >= .1:
		return "%.1f"
	case v < .1 && v >= .01:
		return "%.3f"
	case v < .01:
		return "%.4f"
	default:
		return "%.0f"

	}
}

func xTick(target rl.Rectangle, s scale, t tick) {
	sp := s.transform(rl.NewVector2(t.v, 0))
	rad := float32(3)
	sp.Y = target.Y + rad
	rl.DrawCircleV(sp, rad, theme.Salmon)

	m := theme.MeaureTextPad(t.label)
	tp := rl.Vector2Add(sp, rl.NewVector2(-(m.X/2), 3))
	// Keep the end labels inside the axis instead of clipped at the edge.
	tp.X = max(target.X-m.X/2, min(tp.X, target.X+target.Width-m.X))
	theme.DrawTextPad(t.label, tp, theme.Charcoal)
}

func drawAxisY(target rl.Rectangle, s scale, ts []tick) {
	rl.DrawLineEx(
		rl.NewVector2(target.X+target.Width, target.Y),
		rl.NewVector2(target.X+target.Width, target.Y+target.Height),
		2,
		theme.Salmon,
	)
	for _, t := range ts {
		tickY(target, s, t)
	}
}

const tickWidth = float32(3)

// tickY draws a tick with its label right-aligned against the axis.
func tickY(target rl.Rectangle, s scale, t tick) {
	tickPos := s.transform(rl.NewVector2(0, t.v))
	tickPos.X = target.X + target.Width

	p0 := rl.Vector2Add(tickPos, rl.NewVector2(-tickWidth, 0))
	p1 := rl.Vector2Add(tickPos, rl.NewVector2(tickWidth, 0))
	rl.DrawLineEx(p0, p1, 1.5, theme.Charcoal)

	ms := theme.MeaureTextPad(t.label)
	textPos := rl.NewVector2(target.X+target.Width-tickWidth-ms.X, tickPos.Y-ms.Y/2)
	theme.DrawTextPad(t.label, textPos, theme.Charcoal)
}

// measureUnit is the length of a rotated unit line; tests replace it as
// text can't be measured without a window.
var measureUnit = func(s string) float32 { return theme.MeaureTextPad(s).X }

// unitLines splits unit at spaces into lines, rotated along an axis of
// length l, that fit it, e.g. "mm/s RMS" on a short axis; ok is false when
// a word alone doesn't fit.
func unitLines(unit string, l float32) (lines []string, ok bool) {
	for _, w := range strings.Fields(unit) {
		if measureUnit(w) > l {
			return nil, false
		}
		if n := len(lines); n > 0 && measureUnit(lines[n-1]+" "+w) <= l {
			lines[n-1] += " " + w
			continue
		}
		lines = append(lines, w)
	}
	return lines, true
}

// drawUnitY draws the unit lines rotated to read upwards, each centred
// along the y axis, from its left edge.
func drawUnitY(target rl.Rectangle, lines []string) {
	lineH := theme.MeaureText("0").Y
	for i, l := range lines {
		ms := theme.MeaureTextPad(l)
		pos := rl.NewVector2(target.X+float32(i)*lineH, target.Y+target.Height/2+ms.X/2)
		rl.DrawTextPro(theme.Font, l, rl.Vector2Add(pos, rl.NewVector2(5, -5)), rl.Vector2Zero(), -90,
			theme.FontSize, theme.Spacing, theme.Charcoal)
	}
}
