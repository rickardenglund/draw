package floating

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/rickardenglund/draw/theme"
)

func TextArea(mp rl.Vector2, msg string) {
	fontSize := float32(20)
	spacing := float32(0)
	ms := rl.MeasureTextEx(theme.Font, msg, fontSize, spacing)

	padding := float32(10)
	paddedRect := rl.NewRectangle(mp.X-(ms.X+2*padding), mp.Y-ms.Y-2*padding, ms.X+2*padding, ms.Y+2*padding)
	if paddedRect.Y < 0 {
		paddedRect.Y = 0
	}
	if paddedRect.X < 0 {
		paddedRect.X = 0
	}

	target := rl.NewVector2(paddedRect.X+padding, paddedRect.Y+padding)

	rl.DrawRectangleRec(paddedRect, rl.ColorAlpha(theme.Nude, .8))
	rl.DrawRectangleLinesEx(paddedRect, 2, rl.ColorAlpha(theme.Charcoal, .8))
	rl.DrawTextEx(theme.Font, msg, target, fontSize, spacing, theme.Charcoal)

}
