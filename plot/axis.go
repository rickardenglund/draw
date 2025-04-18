package plot

import (
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/rickardenglund/draw/theme"
)

func drawAxisX(target rl.Rectangle, s scale) {
	w := theme.MeaureTextPad(fmt.Sprintf(getFmt(s.l.maxX), s.l.maxX))
	n := int(math.Floor(float64(target.Width / w.X)))
	dx := (s.l.maxX - s.l.minX) / float32(n)
	for i := range n {
		xTick(target, s, s.l.minX+float32(i)*dx, getFmt(s.l.maxX))
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

func xTick(target rl.Rectangle, s scale, x float32, format string) {
	sp := s.transform(rl.NewVector2(x, 0))
	rad := float32(3)
	sp.Y = target.Y + rad
	rl.DrawCircleV(sp, rad, theme.Salmon)
	text := fmt.Sprintf(format, x)

	m := theme.MeaureTextPad(text)
	tp := rl.Vector2Add(sp, rl.NewVector2(-(m.X/2), 3))
	theme.DrawTextPad(text, tp, theme.Charcoal)
}

func drawAxisY(target rl.Rectangle, s scale) {
	w := theme.MeaureTextPad("10")

	rl.DrawLineEx(
		rl.NewVector2(target.X+target.Width, target.Y),
		rl.NewVector2(target.X+target.Width, target.Y+target.Height),
		2,
		theme.Salmon,
	)

	//format := getFmt(s.l.maxY)
	n := int(math.Floor(float64(target.Height / w.Y)))
	dy := (s.l.maxY - s.l.minY) / float32(n)
	for i := range n {
		v := s.l.minY + float32(i)*dy
		tickY(target, s, v, getFmt(v))
	}
}

func getWidest(a, b float32) rl.Vector2 {
	aw := theme.MeaureTextPad(fmt.Sprintf(getFmt(a), a))
	bw := theme.MeaureTextPad(fmt.Sprintf(getFmt(b), b))

	if aw.X > bw.X {
		return aw
	}

	return bw
}

func tickY(target rl.Rectangle, s scale, y float32, format string) {
	tickWidth := float32(3)

	tickPos := s.transform(rl.NewVector2(0, y))
	tickPos.X = target.X + target.Width
	screenY := tickPos.Y

	p0 := rl.Vector2Add(tickPos, rl.NewVector2(-tickWidth, 0))
	p1 := rl.Vector2Add(tickPos, rl.NewVector2(tickWidth, 0))
	rl.DrawLineEx(p0, p1, 1.5, theme.Charcoal)

	text := fmt.Sprintf(format, y)
	ms := theme.MeaureTextPad(text)

	textPos := rl.NewVector2(target.X, screenY-ms.Y/2)
	textPos.X = target.X
	theme.DrawTextPad(text, textPos, theme.Charcoal)
}
