package plot

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"

	"draw/theme"
)

func drawAxisX(target rl.Rectangle, s scale) {
	//rl.DrawRectangleRec(target, theme.Muave)
	n := 5
	dx := (s.maxX - s.minX) / float32(n)
	for i := range n {
		xTick(target, s, s.minX+float32(i)*dx, getFmt(s.minX, s.maxX))
	}
}

func getFmt(minV, maxV float32) string {
	switch {
	case maxV < 1:
		return "%.1f"
	case maxV < .1:
		return "%.2f"
	default:
		return "%.0f"

	}
}

func xTick(target rl.Rectangle, s scale, x float32, format string) {
	sp := s.transform(rl.NewVector2(x, 0))
	sp.Y = target.Y
	rl.DrawCircleV(sp, 3, theme.Salmon)
	text := fmt.Sprintf(format, x)
	fontSize := float32(14)
	spacing := float32(1)

	m := rl.MeasureTextEx(theme.Font, text, fontSize, spacing)
	tp := rl.Vector2Add(sp, rl.NewVector2(-(m.X/2), 3))
	rl.DrawTextEx(theme.Font, text, tp, fontSize, spacing, theme.Charcoal)
}
func drawAxisY(target rl.Rectangle, s scale) {
	//rl.DrawRectangleRec(target, theme.Muave)

	rl.DrawLineEx(rl.NewVector2(target.X+target.Width, target.Y), rl.NewVector2(target.X+target.Width, target.Y+target.Height), 2, theme.Salmon)

	format := getFmt(s.minY, s.maxY)
	n := 5
	dy := (s.maxY - s.minY) / float32(n)
	for i := range n {
		tick(target, s, s.minY+float32(i)*dy, format)
	}

}

func tick(target rl.Rectangle, s scale, y float32, format string) {
	tickWidth := target.Width / 4

	p := s.transform(rl.NewVector2(0, y))

	p.X = target.X + target.Width - tickWidth/2
	p2 := rl.Vector2Add(p, rl.NewVector2(tickWidth, 0))
	rl.DrawLineEx(p, p2, 1.5, theme.Charcoal)

	text := fmt.Sprintf(format, y)
	fontSize := float32(14)
	spacing := float32(1)
	m := rl.MeasureTextEx(theme.Font, text, fontSize, spacing)
	p = rl.Vector2Add(p, rl.NewVector2(-m.X-3, -m.Y/2))
	rl.DrawTextEx(theme.Font, text, p, fontSize, spacing, theme.Charcoal)
}
