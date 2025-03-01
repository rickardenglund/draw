package plot

import (
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/rickardenglund/draw/theme"
)

func drawAxisX(target rl.Rectangle, s scale) {
	w := theme.MeaureTextPad(fmt.Sprintf(getFmt(s.l.minX, s.l.maxX), s.l.maxX))
	n := int(math.Floor(float64(target.Width / w.X)))
	dx := (s.l.maxX - s.l.minX) / float32(n)
	for i := range n {
		xTick(target, s, s.l.minX+float32(i)*dx, getFmt(s.l.minX, s.l.maxX))
	}
}

func getFmt(minV, maxV float32) string {
	if maxV < 0 {
		maxV *= -1
	}

	switch {
	case maxV < 10 && maxV >= .1:
		return "%.1f"
	case maxV < .1 && maxV >= .01:
		return "%.3f"
	case maxV < .01:
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
	w := theme.MeaureTextPad(fmt.Sprintf(getFmt(s.l.minY, s.l.maxY), s.l.maxY))
	rl.DrawLineEx(rl.NewVector2(target.X+target.Width, target.Y), rl.NewVector2(target.X+target.Width, target.Y+target.Height), 2, theme.Salmon)

	format := getFmt(s.l.minY, s.l.maxY)
	n := int(math.Floor(float64(target.Height / w.Y)))
	dy := (s.l.maxY - s.l.minY) / float32(n)
	for i := range n {
		tickY(target, s, s.l.minY+float32(i)*dy, format)
	}
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
