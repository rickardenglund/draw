package widget

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/rickardenglund/draw/draw"
	"github.com/rickardenglund/draw/theme"
)

type Button struct {
	action func()
	title  string
}

func (b *Button) Init() {}

func (b *Button) Draw(target rl.Rectangle) {
	fontSize := float32(14)
	spacing := float32(.5)
	ms := rl.MeasureTextEx(theme.Font, b.title, fontSize, spacing)

	mp := rl.GetMousePosition()
	padding := float32(5)

	h := ms.Y + 2*padding
	w := ms.X + 2*padding
	btn := rl.NewRectangle(target.X+target.Width/2-w/2, target.Y+target.Height/2-h/2, w, h)

	hover := rl.CheckCollisionPointRec(mp, btn)
	clr := theme.Charcoal
	if hover {
		clr = rl.ColorBrightness(clr, .2)
	}

	if hover && rl.IsMouseButtonDown(rl.MouseButtonLeft) {
		clr = rl.ColorBrightness(clr, .5)
	}

	rl.DrawRectangleRec(btn, clr)
	rl.DrawRectangleLinesEx(btn, padding/2, rl.ColorBrightness(theme.Charcoal, -.5))
	textPos := rl.NewVector2(btn.X+padding, btn.Y+padding)
	rl.DrawTextEx(theme.Font, b.title, textPos, fontSize, spacing, theme.Nude)

	if hover && rl.IsMouseButtonReleased(rl.MouseButtonLeft) {
		b.action()
	}
}

func NewButton(title string, action func()) draw.Drawable {
	return &Button{action: action, title: title}
}
