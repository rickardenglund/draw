package plot

import (
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
