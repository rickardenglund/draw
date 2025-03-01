package theme

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

// Dreams
//var (
//	Muave = rl.NewColor(224, 176, 255, 255)
//	Lavender = rl.NewColor(230, 230, 250, 255)
//	Black = rl.Black
//	Orchid = rl.NewColor(218, 112, 214, 255)
//)

var (
	Muave    = rl.NewColor(224, 176, 255, 255)
	Salmon   = rl.NewColor(250, 128, 114, 255)
	Charcoal = rl.NewColor(54, 69, 79, 255)
	Nude     = rl.NewColor(227, 188, 154, 255)
	Font     = rl.GetFontDefault()
	FontSize = float32(20)
	Spacing  = float32(.5)
	pad      = rl.NewVector2(5, 5)
)

func MeaureText(s string) rl.Vector2 {
	return rl.MeasureTextEx(Font, s, FontSize, Spacing)
}

func MeaureTextPad(s string) rl.Vector2 {
	ms := rl.MeasureTextEx(Font, s, FontSize, Spacing)
	return rl.Vector2Add(ms, rl.Vector2Scale(pad, 2))
}

func DrawText(s string, p rl.Vector2, clr rl.Color) {
	rl.DrawTextEx(Font, s, p, FontSize, Spacing, clr)
}

func DrawTextPad(s string, p rl.Vector2, clr rl.Color) {
	p = rl.Vector2Add(p, pad)
	DrawText(s, p, clr)
}
