package widget

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"draw/draw"
	"draw/theme"
)

type Button struct {
	action func()
}

func (b Button) Init() {}

func (b Button) Draw(target rl.Rectangle) {
	mp := rl.GetMousePosition()
	btn := rl.NewRectangle(target.X+10, target.Y+20, 30, 30)

	hover := rl.CheckCollisionPointRec(mp, btn)
	clr := theme.Charcoal
	if hover {
		clr = rl.ColorBrightness(clr, .2)
	}

	rl.DrawRectangleRec(btn, clr)
	rl.DrawRectangleLinesEx(btn, 5, rl.ColorBrightness(theme.Charcoal, -.5))

	if hover && rl.IsMouseButtonReleased(rl.MouseButtonLeft) {
		b.action()
	}
}

func NewButton(action func()) draw.Drawable {
	return &Button{action: action}
}
